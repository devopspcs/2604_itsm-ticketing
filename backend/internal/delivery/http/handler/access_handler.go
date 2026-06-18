package handler

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/org/itsm/pkg/config"
)

const (
	poseAPIURL   = "https://pose-api.pcsindonesia.co.id/master/user?filter=pcsindonesia.co.id&order=created_at:-1&page_size=-1&page=1"
	poseAPIToken = "51Fr9WJ4nt6bWpvBVpw6piyHbm1VSAsVjA5vfMtczahHOjO7dcbMW2gOhU2Ayr3o"

	pritunlURL      = "https://35.219.125.94"
	pritunlUsername = "devopspcs"
	pritunlPass     = "!*%bXMq|EoYJF0wh"
	pritunlOrgID    = "659b7e2851948a075375d28a"
)

type AccessHandler struct {
	client *http.Client
	cfg    *config.Config
}

func NewAccessHandler(cfg *config.Config) *AccessHandler {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &AccessHandler{
		client: &http.Client{Timeout: 30 * time.Second, Transport: transport},
		cfg:    cfg,
	}
}

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

type ServiceAccessEntry struct {
	Service     string   `json:"service"`
	HasAccess   bool     `json:"has_access"`
	Roles       []string `json:"roles"`
	AccountName string   `json:"account_name"`
	Status      string   `json:"status"`
}

// GetUserAccess checks if a given email has access to registered services
func (h *AccessHandler) GetUserAccess(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, `{"error":"email parameter required"}`, http.StatusBadRequest)
		return
	}
	email = strings.ToLower(email)

	results := []ServiceAccessEntry{}

	// Check POSe access
	poseAccess := h.checkPoseAccess(email)
	results = append(results, poseAccess)

	// Check Pritunl VPN access
	vpnAccess := h.checkPritunlAccess(email)
	results = append(results, vpnAccess)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// GetAllServiceUsers returns all POSe users mapped by email for bulk comparison
func (h *AccessHandler) GetAllServiceUsers(w http.ResponseWriter, r *http.Request) {
	poseUsers, err := h.fetchPoseUsers()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Failed to fetch POSe users: %s"}`, err.Error()), http.StatusBadGateway)
		return
	}

	// Group by email, collect unique roles
	emailMap := map[string]*ServiceAccessEntry{}
	for _, u := range poseUsers {
		e := strings.ToLower(u.Email)
		if !strings.HasSuffix(e, "@pcsindonesia.co.id") {
			continue
		}
		if existing, ok := emailMap[e]; ok {
			// Add role if not already present
			roleExists := false
			for _, r := range existing.Roles {
				if r == u.Role.Name {
					roleExists = true
					break
				}
			}
			if !roleExists && u.Role.Name != "" {
				existing.Roles = append(existing.Roles, u.Role.Name)
			}
		} else {
			roles := []string{}
			if u.Role.Name != "" {
				roles = append(roles, u.Role.Name)
			}
			status := "inactive"
			if u.Status {
				status = "active"
			}
			emailMap[e] = &ServiceAccessEntry{
				Service:     "POSe",
				HasAccess:   true,
				Roles:       roles,
				AccountName: u.Name,
				Status:      status,
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emailMap)
}

func (h *AccessHandler) checkPoseAccess(email string) ServiceAccessEntry {
	poseUsers, err := h.fetchPoseUsers()
	if err != nil {
		return ServiceAccessEntry{
			Service:   "POSe",
			HasAccess: false,
			Status:    "error",
		}
	}

	roles := []string{}
	accountName := ""
	found := false
	for _, u := range poseUsers {
		if strings.ToLower(u.Email) == email {
			found = true
			if u.Name != "" && accountName == "" {
				accountName = u.Name
			}
			// Collect unique roles
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

	return ServiceAccessEntry{
		Service:     "POSe",
		HasAccess:   found,
		Roles:       roles,
		AccountName: accountName,
		Status:      status,
	}
}

func (h *AccessHandler) fetchPoseUsers() ([]poseUserEntry, error) {
	req, err := http.NewRequest("GET", h.cfg.PoseAPIURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+h.cfg.PoseAPIToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "PCS-Pose-App/1.0")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result poseResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}


// Pritunl VPN integration
type pritunlUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Status   bool   `json:"status"`
	Disabled bool   `json:"disabled"`
}

func (h *AccessHandler) checkPritunlAccess(email string) ServiceAccessEntry {
	users, err := h.fetchPritunlUsers()
	if err != nil {
		return ServiceAccessEntry{
			Service:   "VPN (Pritunl)",
			HasAccess: false,
			Status:    "error",
		}
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
			return ServiceAccessEntry{
				Service:     "VPN (Pritunl)",
				HasAccess:   true,
				Roles:       roles,
				AccountName: u.Name,
				Status:      status,
			}
		}
	}

	return ServiceAccessEntry{
		Service:   "VPN (Pritunl)",
		HasAccess: false,
		Roles:     []string{},
		Status:    "no_access",
	}
}

func (h *AccessHandler) fetchPritunlUsers() ([]pritunlUser, error) {
	// Step 1: Login to get session cookie
	loginPayload := fmt.Sprintf(`{"username":"%s","password":"%s"}`, h.cfg.PritunlUsername, h.cfg.PritunlPassword)
	loginReq, _ := http.NewRequest("POST", h.cfg.PritunlURL+"/auth/session", strings.NewReader(loginPayload))
	loginReq.Header.Set("Content-Type", "application/json")

	loginResp, err := h.client.Do(loginReq)
	if err != nil {
		return nil, fmt.Errorf("pritunl login failed: %w", err)
	}
	defer loginResp.Body.Close()

	// Extract session cookie
	var sessionCookie string
	for _, c := range loginResp.Cookies() {
		if c.Name == "session" {
			sessionCookie = c.Value
			break
		}
	}
	if sessionCookie == "" {
		return nil, fmt.Errorf("no session cookie from pritunl")
	}

	// Step 2: Get CSRF token from /state
	stateReq, _ := http.NewRequest("GET", h.cfg.PritunlURL+"/state", nil)
	stateReq.Header.Set("Cookie", "session="+sessionCookie)
	stateResp, err := h.client.Do(stateReq)
	if err != nil {
		return nil, fmt.Errorf("pritunl state failed: %w", err)
	}
	defer stateResp.Body.Close()

	var stateData struct {
		CsrfToken string `json:"csrf_token"`
	}
	if err := json.NewDecoder(stateResp.Body).Decode(&stateData); err != nil {
		return nil, fmt.Errorf("pritunl state parse failed: %w", err)
	}

	// Step 3: Get users with session + CSRF
	usersReq, _ := http.NewRequest("GET", fmt.Sprintf("%s/user/%s", h.cfg.PritunlURL, h.cfg.PritunlOrgID), nil)
	usersReq.Header.Set("Cookie", "session="+sessionCookie)
	usersReq.Header.Set("Csrf-Token", stateData.CsrfToken)

	usersResp, err := h.client.Do(usersReq)
	if err != nil {
		return nil, fmt.Errorf("pritunl users fetch failed: %w", err)
	}
	defer usersResp.Body.Close()

	body, _ := io.ReadAll(usersResp.Body)
	if usersResp.StatusCode != 200 {
		return nil, fmt.Errorf("pritunl users returned %d: %s", usersResp.StatusCode, string(body))
	}

	var users []pritunlUser
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("pritunl users parse failed: %w", err)
	}

	return users, nil
}
