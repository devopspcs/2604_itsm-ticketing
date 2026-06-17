package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const mdmBaseURL = "https://psops.pcsindonesia.com"

type MDMHandler struct {
	client *http.Client
}

func NewMDMHandler() *MDMHandler {
	return &MDMHandler{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

// ListDevices proxies the MDM device list API
func (h *MDMHandler) ListDevices(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "20"
	}
	deviceType := r.URL.Query().Get("device_type")
	category := r.URL.Query().Get("category") // "laptop" or "mobile"

	var url string
	if category == "mobile" {
		url = fmt.Sprintf("%s/mdm/device-mobile?page=%s&limit=%s", mdmBaseURL, page, limit)
		if deviceType != "" {
			url += "&device_type=" + deviceType
		}
	} else {
		url = fmt.Sprintf("%s/mdm/device?page=%s&limit=%s", mdmBaseURL, page, limit)
		if deviceType != "" {
			url += "&device_type=" + deviceType
		}
	}

	resp, err := h.client.Get(url)
	if err != nil {
		http.Error(w, `{"error":"Failed to connect to MDM service"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, `{"error":"Failed to read MDM response"}`, http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// GetDevice proxies a single device detail
func (h *MDMHandler) GetDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")
	if deviceID == "" {
		http.Error(w, `{"error":"Missing device ID"}`, http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("%s/mdm/device/%s", mdmBaseURL, deviceID)

	resp, err := h.client.Get(url)
	if err != nil {
		http.Error(w, `{"error":"Failed to connect to MDM service"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, `{"error":"Failed to read MDM response"}`, http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// GetDeviceStats returns summary stats from all devices
func (h *MDMHandler) GetDeviceStats(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/mdm/device?page=1&limit=1000", mdmBaseURL)

	resp, err := h.client.Get(url)
	if err != nil {
		http.Error(w, `{"error":"Failed to connect to MDM service"}`, http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Data []struct {
				ID       int    `json:"id"`
				Name     string `json:"name"`
				OS       string `json:"os"`
				IsOnline bool   `json:"is_online"`
			} `json:"data"`
			TotalItems int `json:"total_items"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		http.Error(w, `{"error":"Failed to parse MDM response"}`, http.StatusBadGateway)
		return
	}

	totalDevices := result.Data.TotalItems
	onlineCount := 0
	osCount := map[string]int{}
	for _, d := range result.Data.Data {
		if d.IsOnline {
			onlineCount++
		}
		osCount[d.OS]++
	}

	stats := map[string]interface{}{
		"total_devices": totalDevices,
		"online_count":  onlineCount,
		"offline_count": totalDevices - onlineCount,
		"os_breakdown":  osCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
