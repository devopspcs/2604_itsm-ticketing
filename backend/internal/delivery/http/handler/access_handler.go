package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	poseAPIURL   = "https://pose-api.pcsindonesia.co.id/master/user?filter=pcsindonesia.co.id&order=created_at:-1&page_size=-1&page=1"
	poseAPIToken = "51Fr9WJ4nt6bWpvBVpw6piyHbm1VSAsVjA5vfMtczahHOjO7dcbMW2gOhU2Ayr3o"
)

type AccessHandler struct {
	client *http.Client
}

func NewAccessHandler() *AccessHandler {
	return &AccessHandler{
		client: &http.Client{Timeout: 30 * time.Second},
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

	// Pritunl VPN placeholder (will implement when API token is available)
	results = append(results, ServiceAccessEntry{
		Service:   "VPN (Pritunl)",
		HasAccess: false,
		Roles:     []string{},
		Status:    "not_configured",
	})

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
	req, err := http.NewRequest("GET", poseAPIURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+poseAPIToken)
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
