package handler

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/itsm/pkg/apperror"
)

type ExternalService struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	URL          string          `json:"url"`
	AuthUsername string          `json:"auth_username"`
	AuthToken    string          `json:"auth_token"`
	AuthPassword string          `json:"auth_password,omitempty"`
	ExtraConfig  json.RawMessage `json:"extra_config"`
	IsActive     bool            `json:"is_active"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

type AccessHandler struct {
	client *http.Client
	db     *pgxpool.Pool
}

func NewAccessHandler(db *pgxpool.Pool) *AccessHandler {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &AccessHandler{
		client: &http.Client{Timeout: 30 * time.Second, Transport: transport},
		db:     db,
	}
}

type ServiceAccessEntry struct {
	Service     string   `json:"service"`
	ServiceType string   `json:"service_type"`
	HasAccess   bool     `json:"has_access"`
	Roles       []string `json:"roles"`
	AccountName string   `json:"account_name"`
	Status      string   `json:"status"`
}

// ==================== CRUD for External Services ====================

func (h *AccessHandler) ListServices(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `SELECT id::text, name, type, url, auth_username, COALESCE(auth_token,''), COALESCE(extra_config,'{}')::text, is_active, created_at::text, updated_at::text FROM external_services ORDER BY name`)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Failed to list services: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var services []ExternalService
	for rows.Next() {
		var s ExternalService
		var extraStr string
		if err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.URL, &s.AuthUsername, &s.AuthToken, &extraStr, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			continue
		}
		s.ExtraConfig = json.RawMessage(extraStr)
		s.AuthPassword = ""
		services = append(services, s)
	}
	if services == nil {
		services = []ExternalService{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (h *AccessHandler) CreateService(w http.ResponseWriter, r *http.Request) {
	var s ExternalService
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}
	if s.Name == "" || s.Type == "" || s.URL == "" {
		http.Error(w, `{"error":"name, type, and url are required"}`, http.StatusBadRequest)
		return
	}
	if s.ExtraConfig == nil {
		s.ExtraConfig = json.RawMessage(`{}`)
	}
	err := h.db.QueryRow(r.Context(),
		`INSERT INTO external_services (name, type, url, auth_username, auth_token, auth_password, extra_config, is_active) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
		s.Name, s.Type, s.URL, s.AuthUsername, s.AuthToken, s.AuthPassword, s.ExtraConfig, true).Scan(&s.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func (h *AccessHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var s ExternalService
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}
	if s.ExtraConfig == nil {
		s.ExtraConfig = json.RawMessage(`{}`)
	}
	_, err := h.db.Exec(r.Context(),
		`UPDATE external_services SET name=$1, type=$2, url=$3, auth_username=$4, auth_token=$5, auth_password=$6, extra_config=$7, is_active=$8, updated_at=NOW() WHERE id=$9`,
		s.Name, s.Type, s.URL, s.AuthUsername, s.AuthToken, s.AuthPassword, s.ExtraConfig, s.IsActive, id)
	if err != nil {
		http.Error(w, `{"error":"Failed to update"}`, http.StatusInternalServerError)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "updated"})
}

func (h *AccessHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.db.Exec(r.Context(), `DELETE FROM external_services WHERE id=$1`, id)
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// ==================== Access Check (reads from cache) ====================

func (h *AccessHandler) GetUserAccess(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(r.URL.Query().Get("email"))
	if email == "" {
		http.Error(w, `{"error":"email parameter required"}`, http.StatusBadRequest)
		return
	}

	// Read from cache
	rows, err := h.db.Query(r.Context(),
		`SELECT service_name, service_type, account_name, roles, status FROM service_access_cache WHERE email = $1`, email)
	if err != nil {
		http.Error(w, `{"error":"Failed to query cache"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	found := map[string]bool{}
	var results []ServiceAccessEntry
	for rows.Next() {
		var name, stype, accountName, status string
		var roles []string
		if err := rows.Scan(&name, &stype, &accountName, &roles, &status); err != nil {
			continue
		}
		found[name] = true
		results = append(results, ServiceAccessEntry{
			Service: name, ServiceType: stype, HasAccess: true,
			Roles: roles, AccountName: accountName, Status: status,
		})
	}

	// Also list services where user has NO access
	svcRows, _ := h.db.Query(r.Context(), `SELECT name, type FROM external_services WHERE is_active = true`)
	if svcRows != nil {
		defer svcRows.Close()
		for svcRows.Next() {
			var name, stype string
			svcRows.Scan(&name, &stype)
			if !found[name] {
				results = append(results, ServiceAccessEntry{
					Service: name, ServiceType: stype, HasAccess: false, Roles: []string{}, Status: "no_access",
				})
			}
		}
	}

	if results == nil {
		results = []ServiceAccessEntry{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// ==================== Sync Status ====================

func (h *AccessHandler) GetSyncStatus(w http.ResponseWriter, r *http.Request) {
	var lastSynced *time.Time
	var status string
	var totalRecords int
	h.db.QueryRow(r.Context(), `SELECT last_synced_at, status, total_records FROM sync_status WHERE id='service_access'`).Scan(&lastSynced, &status, &totalRecords)

	resp := map[string]interface{}{
		"last_synced_at": lastSynced,
		"status":         status,
		"total_records":  totalRecords,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// ==================== Sync (Admin triggers this) ====================

func (h *AccessHandler) SyncAll(w http.ResponseWriter, r *http.Request) {
	// Mark as syncing
	h.db.Exec(r.Context(), `UPDATE sync_status SET status='syncing' WHERE id='service_access'`)

	// Get all active services
	rows, err := h.db.Query(r.Context(), `SELECT id, name, type, url, auth_username, auth_token, auth_password, extra_config FROM external_services WHERE is_active = true`)
	if err != nil {
		http.Error(w, `{"error":"Failed to load services"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var services []ExternalService
	for rows.Next() {
		var s ExternalService
		rows.Scan(&s.ID, &s.Name, &s.Type, &s.URL, &s.AuthUsername, &s.AuthToken, &s.AuthPassword, &s.ExtraConfig)
		services = append(services, s)
	}

	// Clear old cache
	h.db.Exec(r.Context(), `DELETE FROM service_access_cache`)

	totalRecords := 0

	for _, svc := range services {
		var entries []struct {
			Email       string
			AccountName string
			Roles       []string
			Status      string
		}

		switch svc.Type {
		case "pose":
			entries = h.syncPose(svc)
		case "pritunl":
			entries = h.syncPritunl(svc)
		}

		for _, e := range entries {
			h.db.Exec(r.Context(),
				`INSERT INTO service_access_cache (service_id, service_name, service_type, email, account_name, roles, status) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
				svc.ID, svc.Name, svc.Type, strings.ToLower(e.Email), e.AccountName, e.Roles, e.Status)
			totalRecords++
		}
	}

	// Update sync status
	h.db.Exec(r.Context(), `UPDATE sync_status SET last_synced_at=NOW(), status='completed', total_records=$1 WHERE id='service_access'`, totalRecords)

	apperror.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":       "sync completed",
		"total_records": totalRecords,
		"services":      len(services),
	})
}

// ==================== POSe Sync ====================

type poseUserEntry struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Status bool   `json:"status"`
	Role   struct {
		Name string `json:"name"`
	} `json:"Role"`
}

func (h *AccessHandler) syncPose(svc ExternalService) []struct {
	Email       string
	AccountName string
	Roles       []string
	Status      string
} {
	req, _ := http.NewRequest("GET", svc.URL, nil)
	req.Header.Set("Authorization", "Bearer "+svc.AuthToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "PCS-Pose-App/1.0")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Data []poseUserEntry `json:"data"`
	}
	json.Unmarshal(body, &result)

	// Group by email
	emailMap := map[string]*struct {
		Email       string
		AccountName string
		Roles       []string
		Status      string
	}{}

	for _, u := range result.Data {
		e := strings.ToLower(u.Email)
		if !strings.HasSuffix(e, "@pcsindonesia.co.id") {
			continue
		}
		if existing, ok := emailMap[e]; ok {
			roleExists := false
			for _, r := range existing.Roles {
				if r == u.Role.Name {
					roleExists = true
				}
			}
			if !roleExists && u.Role.Name != "" {
				existing.Roles = append(existing.Roles, u.Role.Name)
			}
		} else {
			roles := []string{}
			if u.Role.Name != "" {
				roles = []string{u.Role.Name}
			}
			status := "inactive"
			if u.Status {
				status = "active"
			}
			emailMap[e] = &struct {
				Email       string
				AccountName string
				Roles       []string
				Status      string
			}{Email: e, AccountName: u.Name, Roles: roles, Status: status}
		}
	}

	var entries []struct {
		Email       string
		AccountName string
		Roles       []string
		Status      string
	}
	for _, v := range emailMap {
		entries = append(entries, *v)
	}
	return entries
}

// ==================== Pritunl Sync ====================

type pritunlUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Status   bool   `json:"status"`
	Disabled bool   `json:"disabled"`
}

func (h *AccessHandler) syncPritunl(svc ExternalService) []struct {
	Email       string
	AccountName string
	Roles       []string
	Status      string
} {
	var extra struct {
		OrgID string `json:"org_id"`
	}
	json.Unmarshal(svc.ExtraConfig, &extra)
	if extra.OrgID == "" {
		return nil
	}

	// Login
	loginPayload := fmt.Sprintf(`{"username":"%s","password":"%s"}`, svc.AuthUsername, svc.AuthPassword)
	loginReq, _ := http.NewRequest("POST", svc.URL+"/auth/session", strings.NewReader(loginPayload))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, err := h.client.Do(loginReq)
	if err != nil {
		return nil
	}
	defer loginResp.Body.Close()

	var sessionCookie string
	for _, c := range loginResp.Cookies() {
		if c.Name == "session" {
			sessionCookie = c.Value
		}
	}
	if sessionCookie == "" {
		return nil
	}

	// CSRF
	stateReq, _ := http.NewRequest("GET", svc.URL+"/state", nil)
	stateReq.Header.Set("Cookie", "session="+sessionCookie)
	stateResp, err := h.client.Do(stateReq)
	if err != nil {
		return nil
	}
	defer stateResp.Body.Close()
	var stateData struct {
		CsrfToken string `json:"csrf_token"`
	}
	json.NewDecoder(stateResp.Body).Decode(&stateData)

	// Users
	usersReq, _ := http.NewRequest("GET", fmt.Sprintf("%s/user/%s", svc.URL, extra.OrgID), nil)
	usersReq.Header.Set("Cookie", "session="+sessionCookie)
	usersReq.Header.Set("Csrf-Token", stateData.CsrfToken)
	usersResp, err := h.client.Do(usersReq)
	if err != nil {
		return nil
	}
	defer usersResp.Body.Close()
	body, _ := io.ReadAll(usersResp.Body)

	var users []pritunlUser
	json.Unmarshal(body, &users)

	var entries []struct {
		Email       string
		AccountName string
		Roles       []string
		Status      string
	}
	for _, u := range users {
		if u.Email == "" || u.Email == "None" {
			continue
		}
		status := "active"
		if u.Disabled {
			status = "disabled"
		}
		roles := []string{"VPN User"}
		if u.Status {
			roles = append(roles, "Online")
		}
		entries = append(entries, struct {
			Email       string
			AccountName string
			Roles       []string
			Status      string
		}{Email: strings.ToLower(u.Email), AccountName: u.Name, Roles: roles, Status: status})
	}
	return entries
}
