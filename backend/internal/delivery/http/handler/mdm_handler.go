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

// GetDeviceStats returns summary stats from all devices (laptop + mobile)
func (h *MDMHandler) GetDeviceStats(w http.ResponseWriter, r *http.Request) {
	// Fetch laptop/desktop devices
	laptopURL := fmt.Sprintf("%s/mdm/device?page=1&limit=1000", mdmBaseURL)
	laptopResp, err := h.client.Get(laptopURL)
	if err != nil {
		http.Error(w, `{"error":"Failed to connect to MDM service"}`, http.StatusBadGateway)
		return
	}
	defer laptopResp.Body.Close()

	var laptopResult struct {
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
	if err := json.NewDecoder(laptopResp.Body).Decode(&laptopResult); err != nil {
		http.Error(w, `{"error":"Failed to parse MDM laptop response"}`, http.StatusBadGateway)
		return
	}

	// Fetch mobile devices
	mobileURL := fmt.Sprintf("%s/mdm/device-mobile?page=1&limit=1000&device_type=smartphone", mdmBaseURL)
	mobileResp, err := h.client.Get(mobileURL)
	if err != nil {
		// Mobile API might not be available, continue with laptop only
		mobileResp = nil
	}

	var mobileCount int
	mobileOS := map[string]int{}
	if mobileResp != nil {
		defer mobileResp.Body.Close()
		var mobileResult struct {
			Data struct {
				Data []struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
					OS   string `json:"os"`
				} `json:"data"`
				TotalItems int `json:"total_items"`
			} `json:"data"`
		}
		if err := json.NewDecoder(mobileResp.Body).Decode(&mobileResult); err == nil {
			mobileCount = mobileResult.Data.TotalItems
			for _, d := range mobileResult.Data.Data {
				mobileOS[d.OS]++
			}
		}
	}

	// Combine stats
	laptopTotal := laptopResult.Data.TotalItems
	onlineCount := 0
	osCount := map[string]int{}
	for _, d := range laptopResult.Data.Data {
		if d.IsOnline {
			onlineCount++
		}
		osCount[d.OS]++
	}
	for os, count := range mobileOS {
		osCount[os] += count
	}

	totalDevices := laptopTotal + mobileCount

	stats := map[string]interface{}{
		"total_devices":  totalDevices,
		"laptop_count":   laptopTotal,
		"mobile_count":   mobileCount,
		"online_count":   onlineCount,
		"offline_count":  laptopTotal - onlineCount,
		"os_breakdown":   osCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
