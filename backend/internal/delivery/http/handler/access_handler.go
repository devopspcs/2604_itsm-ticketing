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
	rows, err := h.db.Query(r.Context(), `SELECT id, name, type, url, auth_username, auth_token, extra_config, is_active, created_at, updated_at FROM external_services ORDER BY name`)
	if err != nil {
		http.Error(w, `{"error":"Failed to list services"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var services []ExternalService
	for rows.Next() {
		var s ExternalService
		if err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.URL, &s.AuthUsername, &s.AuthToken, &s.ExtraConfig, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
			continue
		}
		s.AuthPassword = "" // Don't expose password in list
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
		`INSERT INTO external_services (name, type, url, auth_username, auth_token, auth_password, extra_config, is_active) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		s.Name, s.Type, s.URL, s.AuthUsername, s.AuthToken, s.AuthPassword, s.ExtraConfig, true).Scan(&s.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Failed to create service: %s"}`, err.Error()), http.StatusInternalServerError)
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
		http.Error(w, `{"error":"Failed to update service"}`, http.StatusInternalServerError)
		return
	}

	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "updated"})
}

func (h *AccessHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := h.db.Exec(r.Context(), `DELETE FROM external_services WHERE id=$1`, id)
	if err != nil {
		http.Error(w, `{"error":"Failed to delete service"}`, http.StatusInternalServerError)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// ==================== Access Check ====================

func (h *AccessHandler) GetUserAccess(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, `{"error":"email parameter required"}`, http.StatusBadRequest)
		return
	}
	email = strings.ToLower(email)

	// Get all active services from DB
	rows, err := h.db.Query(r.Context(), `SELECT id, name, type, url, auth_username, auth_token, auth_password, extra_config FROM external_services WHERE is_active = true`)
	if err != nil {
		http.Error(w, `{"error":"Failed to load services"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []ServiceAccessEntry
	for rows.Next() {
		var s ExternalService
		if err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.URL, &s.AuthUsername, &s.AuthToken, &s.AuthPassword, &s.ExtraConfig); err != nil {
			continue
		}

		var entry ServiceAccessEntry
		switch s.Type {
		case "pose":
			entry = h.checkPoseAccess(email, s)
		case "pritunl":
			entry = h.checkPritunlAccess(email, s)
		default:
			entry = ServiceAccessEntry{
				Service:     s.Name,
				ServiceType: s.Type,
				HasAccess:   false,
				Status:      "unsupported_type",
			}
		}
		results = append(results, entry)
	}

	if results == nil {
		results = []ServiceAccessEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// ==================== POSe Integration ====================

type poseUserEntry struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	Status bool   `json:"status"`
	Role   struct {
		Name string `json:"name"`
	} `json:"Role"`
}

type poseResponse struct {
	Count int             `json:"count"`
	Data  []poseUserEntry `json:"data"`
}

func (h *AccessHandler) checkPoseAccess(email string, svc ExternalService) ServiceAccessEntry {
	users, err := h.fetchPoseUsers(svc)
	if err != nil {
		return ServiceAccessEntry{Service: svc.Name, ServiceType: svc.Type, HasAccess: false, Status: "error"}
	}

	roles := []string{}
	accountName := ""
	found := false
	for _, u := range users {
		if strings.ToLower(u.Email) == email {
			found = true
			if u.Name != "" && accountName == "" {
				accountName = u.Name
			}
			roleExists := false
			for _, r := range roles {
				if r == u.Role.Name {
					roleExists = true
					break
				}
			}
			if !roleExists && u.Role.Name != "" {
				roles = append(roles, u.Role.Name)
			}
		}
	}

	status := "no_access"
	if found {
		status = "active"
	}
	return ServiceAccessEntry{Service: svc.Name, ServiceType: svc.Type, HasAccess: found, Roles: roles, AccountName: accountName, Status: status}
}

func (h *AccessHandler) fetchPoseUsers(svc ExternalService) ([]poseUserEntry, error) {
	req, err := http.NewRequest("GET", svc.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+svc.AuthToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "PCS-Pose-App/1.0")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result poseResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ==================== Pritunl VPN Integration ====================

type pritunlUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Status   bool   `json:"status"`
	Disabled bool   `json:"disabled"`
}

func (h *AccessHandler) checkPritunlAccess(email string, svc ExternalService) ServiceAccessEntry {
	users, err := h.fetchPritunlUsers(svc)
	if err != nil {
		return ServiceAccessEntry{Service: svc.Name, ServiceType: svc.Type, HasAccess: false, Status: "error"}
	}

	for _, u := range users {
		if strings.ToLower(u.Email) == email {
			status := "active"
			if u.Disabled {
				status = "disabled"
			}
			roles := []string{"VPN User"}
			if u.Status {
				roles = append(roles, "Online")
			}
			return ServiceAccessEntry{Service: svc.Name, ServiceType: svc.Type, HasAccess: true, Roles: roles, AccountName: u.Name, Status: status}
		}
	}

	return ServiceAccessEntry{Service: svc.Name, ServiceType: svc.Type, HasAccess: false, Roles: []string{}, Status: "no_access"}
}

func (h *AccessHandler) fetchPritunlUsers(svc ExternalService) ([]pritunlUser, error) {
	// Parse org_id from extra_config
	var extra struct {
		OrgID string `json:"org_id"`
	}
	json.Unmarshal(svc.ExtraConfig, &extra)
	if extra.OrgID == "" {
		return nil, fmt.Errorf("org_id not configured")
	}

	// Step 1: Login
	loginPayload := fmt.Sprintf(`{"username":"%s","password":"%s"}`, svc.AuthUsername, svc.AuthPassword)
	loginReq, _ := http.NewRequest("POST", svc.URL+"/auth/session", strings.NewReader(loginPayload))
	loginReq.Header.Set("Content-Type", "application/json")

	loginResp, err := h.client.Do(loginReq)
	if err != nil {
		return nil, fmt.Errorf("pritunl login failed: %w", err)
	}
	defer loginResp.Body.Close()

	var sessionCookie string
	for _, c := range loginResp.Cookies() {
		if c.Name == "session" {
			sessionCookie = c.Value
			break
		}
	}
	if sessionCookie == "" {
		return nil, fmt.Errorf("no session cookie")
	}

	// Step 2: Get CSRF token
	stateReq, _ := http.NewRequest("GET", svc.URL+"/state", nil)
	stateReq.Header.Set("Cookie", "session="+sessionCookie)
	stateResp, err := h.client.Do(stateReq)
	if err != nil {
		return nil, fmt.Errorf("state failed: %w", err)
	}
	defer stateResp.Body.Close()

	var stateData struct {
		CsrfToken string `json:"csrf_token"`
	}
	json.NewDecoder(stateResp.Body).Decode(&stateData)

	// Step 3: Get users
	usersReq, _ := http.NewRequest("GET", fmt.Sprintf("%s/user/%s", svc.URL, extra.OrgID), nil)
	usersReq.Header.Set("Cookie", "session="+sessionCookie)
	usersReq.Header.Set("Csrf-Token", stateData.CsrfToken)

	usersResp, err := h.client.Do(usersReq)
	if err != nil {
		return nil, fmt.Errorf("users fetch failed: %w", err)
	}
	defer usersResp.Body.Close()

	body, _ := io.ReadAll(usersResp.Body)
	if usersResp.StatusCode != 200 {
		return nil, fmt.Errorf("users returned %d", usersResp.StatusCode)
	}

	var users []pritunlUser
	json.Unmarshal(body, &users)
	return users, nil
}
