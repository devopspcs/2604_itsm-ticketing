package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type WebhookHandler struct {
	webhookUC domainUC.WebhookUseCase
}

func NewWebhookHandler(webhookUC domainUC.WebhookUseCase) *WebhookHandler {
	return &WebhookHandler{webhookUC: webhookUC}
}

func (h *WebhookHandler) List(w http.ResponseWriter, r *http.Request) {
	configs, err := h.webhookUC.ListWebhookConfigs(r.Context())
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, configs)
}

func (h *WebhookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domainUC.CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	config, err := h.webhookUC.CreateWebhookConfig(r.Context(), req)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, config)
}

func (h *WebhookHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.webhookUC.UpdateWebhookConfig(r.Context(), id, req); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "updated"})
}

func (h *WebhookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.webhookUC.DeleteWebhookConfig(r.Context(), id); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}
