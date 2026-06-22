package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PublicHandler struct {
	db *pgxpool.Pool
}

func NewPublicHandler(db *pgxpool.Pool) *PublicHandler {
	return &PublicHandler{db: db}
}

type ServiceOverview struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
	Type    string `json:"type"`
}

type PublicOverviewResponse struct {
	TotalStaff   int               `json:"total_staff"`
	TotalDevices int               `json:"total_devices"`
	TotalTickets int               `json:"total_tickets"`
	TotalVPN     int               `json:"total_vpn"`
	Services     []ServiceOverview `json:"services"`
}

func (h *PublicHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	var totalStaff int
	h.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM users WHERE is_active = true`).Scan(&totalStaff)

	var totalTickets int
	h.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM tickets`).Scan(&totalTickets)

	var totalDevices int
	h.db.QueryRow(r.Context(), `SELECT total_records FROM sync_status WHERE id='service_access'`).Scan(&totalDevices)

	var totalVPN int
	h.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM external_services WHERE type='pritunl' AND is_active=true`).Scan(&totalVPN)

	// Build services list
	services := []ServiceOverview{
		{Name: "ITSM Ticketing", Version: "2.0.0", Status: "operational", Type: "platform"},
		{Name: "POSe", Version: "3.2.1", Status: "operational", Type: "pos"},
	}

	// Add VPN services
	rows, _ := h.db.Query(r.Context(), `SELECT name FROM external_services WHERE type='pritunl' AND is_active=true ORDER BY name`)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			rows.Scan(&name)
			services = append(services, ServiceOverview{Name: name, Version: "1.0", Status: "operational", Type: "vpn"})
		}
	}

	resp := PublicOverviewResponse{
		TotalStaff:   totalStaff,
		TotalDevices: totalDevices,
		TotalTickets: totalTickets,
		TotalVPN:     totalVPN,
		Services:     services,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
