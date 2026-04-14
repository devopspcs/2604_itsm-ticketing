package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/org/itsm/internal/delivery/http/middleware"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type TicketHandler struct {
	ticketUC     domainUC.TicketUseCase
	approvalUC   domainUC.ApprovalUseCase
	assignmentUC domainUC.AssignmentUseCase
	activityRepo repository.ActivityLogRepository
}

func NewTicketHandler(
	ticketUC domainUC.TicketUseCase,
	approvalUC domainUC.ApprovalUseCase,
	assignmentUC domainUC.AssignmentUseCase,
	activityRepo repository.ActivityLogRepository,
) *TicketHandler {
	return &TicketHandler{
		ticketUC:     ticketUC,
		approvalUC:   approvalUC,
		assignmentUC: assignmentUC,
		activityRepo: activityRepo,
	}
}

func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)

	q := r.URL.Query()
	filter := repository.TicketFilter{Page: 1, PageSize: 20}

	// Parse pagination
	if p := q.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			filter.Page = v
		}
	}
	if ps := q.Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			filter.PageSize = v
		}
	}

	// Parse filters
	if s := q.Get("status"); s != "" {
		st := entity.TicketStatus(s)
		filter.Status = &st
	}
	if t := q.Get("type"); t != "" {
		tt := entity.TicketType(t)
		filter.Type = &tt
	}
	if p := q.Get("priority"); p != "" {
		pp := entity.Priority(p)
		filter.Priority = &pp
	}
	if c := q.Get("category"); c != "" {
		filter.Category = &c
	}
	if s := q.Get("search"); s != "" {
		filter.Search = &s
	}

	result, err := h.ticketUC.ListTickets(r.Context(), filter, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, result)
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	var req domainUC.CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	ticket, err := h.ticketUC.CreateTicket(r.Context(), req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, ticket)
}

func (h *TicketHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	ticket, err := h.ticketUC.GetTicket(r.Context(), id, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	ticket, err := h.ticketUC.UpdateTicket(r.Context(), id, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, ticket)
}

func (h *TicketHandler) Submit(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.approvalUC.SubmitForApproval(r.Context(), id, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "submitted for approval"})
}

func (h *TicketHandler) Assign(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var body struct {
		AssigneeID string `json:"assignee_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	assigneeID, err := uuid.Parse(body.AssigneeID)
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.assignmentUC.AssignTicket(r.Context(), id, assigneeID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "ticket assigned"})
}

func (h *TicketHandler) GetApprovals(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	history, err := h.approvalUC.GetApprovalHistory(r.Context(), id)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, history)
}

func (h *TicketHandler) GetActivities(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	logs, err := h.activityRepo.FindByTicketID(r.Context(), id)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, logs)
}
