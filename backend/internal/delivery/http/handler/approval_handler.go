package handler

import (
	"encoding/json"
	"net/http"

	"github.com/org/itsm/internal/delivery/http/middleware"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type ApprovalHandler struct {
	approvalUC domainUC.ApprovalUseCase
}

func NewApprovalHandler(approvalUC domainUC.ApprovalUseCase) *ApprovalHandler {
	return &ApprovalHandler{approvalUC: approvalUC}
}

func (h *ApprovalHandler) Decide(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	var req domainUC.ApprovalDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.approvalUC.Decide(r.Context(), req, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "decision recorded"})
}
