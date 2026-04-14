package handler

import (
	"net/http"

	"github.com/org/itsm/internal/delivery/http/middleware"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type DashboardHandler struct {
	dashboardUC domainUC.DashboardUseCase
}

func NewDashboardHandler(dashboardUC domainUC.DashboardUseCase) *DashboardHandler {
	return &DashboardHandler{dashboardUC: dashboardUC}
}

func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	filter := domainUC.DashboardFilter{}
	if from := r.URL.Query().Get("date_from"); from != "" {
		filter.DateFrom = &from
	}
	if to := r.URL.Query().Get("date_to"); to != "" {
		filter.DateTo = &to
	}
	stats, err := h.dashboardUC.GetStats(r.Context(), filter, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, stats)
}
