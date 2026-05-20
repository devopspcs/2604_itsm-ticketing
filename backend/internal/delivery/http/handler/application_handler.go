package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/org/itsm/internal/delivery/http/middleware"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type ApplicationHandler struct {
	appUC domainUC.ApplicationUseCase
}

func NewApplicationHandler(appUC domainUC.ApplicationUseCase) *ApplicationHandler {
	return &ApplicationHandler{appUC: appUC}
}

// GET /me/apps — get current user's accessible apps
func (h *ApplicationHandler) GetMyApps(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		apperror.WriteError(w, apperror.ErrForbidden)
		return
	}
	apps, err := h.appUC.GetUserApps(r.Context(), claims.UserID)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, apps)
}

// GET /applications — list all apps
func (h *ApplicationHandler) ListApps(w http.ResponseWriter, r *http.Request) {
	apps, err := h.appUC.ListApps(r.Context())
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, apps)
}

// POST /applications — create app (admin only)
func (h *ApplicationHandler) CreateApp(w http.ResponseWriter, r *http.Request) {
	var req domainUC.CreateApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	app, err := h.appUC.CreateApp(r.Context(), req)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, app)
}

// PATCH /applications/{id} — update app (admin only)
func (h *ApplicationHandler) UpdateApp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	app, err := h.appUC.UpdateApp(r.Context(), id, req)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, app)
}

// DELETE /applications/{id} — delete app (admin only)
func (h *ApplicationHandler) DeleteApp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.appUC.DeleteApp(r.Context(), id); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "application deleted"})
}

// GET /applications/{id}/users — list users with access to app
func (h *ApplicationHandler) GetAppUsers(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	users, err := h.appUC.GetAppUsers(r.Context(), id)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, users)
}

// POST /applications/{id}/access — grant access
func (h *ApplicationHandler) GrantAccess(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	claims, ok := middleware.GetClaims(r)
	if !ok {
		apperror.WriteError(w, apperror.ErrForbidden)
		return
	}
	var req domainUC.GrantAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	req.AppID = appID
	if err := h.appUC.GrantAccess(r.Context(), req, claims.UserID); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "access granted"})
}

// DELETE /applications/{id}/access/{userId} — revoke access
func (h *ApplicationHandler) RevokeAccess(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	userID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.appUC.RevokeAccess(r.Context(), userID, appID); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "access revoked"})
}

// POST /applications/{id}/access/bulk — bulk grant access
func (h *ApplicationHandler) BulkGrantAccess(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	claims, ok := middleware.GetClaims(r)
	if !ok {
		apperror.WriteError(w, apperror.ErrForbidden)
		return
	}
	var req domainUC.BulkGrantAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	req.AppID = appID
	if err := h.appUC.BulkGrantAccess(r.Context(), req, claims.UserID); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "bulk access granted"})
}
