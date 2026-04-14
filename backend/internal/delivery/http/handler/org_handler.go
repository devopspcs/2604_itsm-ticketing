package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type OrgHandler struct {
	orgUC domainUC.OrgUseCase
}

func NewOrgHandler(orgUC domainUC.OrgUseCase) *OrgHandler {
	return &OrgHandler{orgUC: orgUC}
}

// --- Department ---

func (h *OrgHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	depts, err := h.orgUC.ListDepartments(r.Context())
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, depts)
}

func (h *OrgHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var req domainUC.CreateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	dept, err := h.orgUC.CreateDepartment(r.Context(), req)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, dept)
}

func (h *OrgHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	dept, err := h.orgUC.UpdateDepartment(r.Context(), id, req)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, dept)
}

func (h *OrgHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.orgUC.DeleteDepartment(r.Context(), id); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "department deleted"})
}

// --- Division ---

func (h *OrgHandler) ListDivisions(w http.ResponseWriter, r *http.Request) {
	var deptID *uuid.UUID
	if d := r.URL.Query().Get("department_id"); d != "" {
		if id, err := uuid.Parse(d); err == nil {
			deptID = &id
		}
	}
	divs, err := h.orgUC.ListDivisions(r.Context(), deptID)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, divs)
}

func (h *OrgHandler) CreateDivision(w http.ResponseWriter, r *http.Request) {
	var req domainUC.CreateDivisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	div, err := h.orgUC.CreateDivision(r.Context(), req)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, div)
}

func (h *OrgHandler) UpdateDivision(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateDivisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	div, err := h.orgUC.UpdateDivision(r.Context(), id, req)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, div)
}

func (h *OrgHandler) DeleteDivision(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.orgUC.DeleteDivision(r.Context(), id); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "division deleted"})
}

// --- Team ---

func (h *OrgHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	var divID *uuid.UUID
	if d := r.URL.Query().Get("division_id"); d != "" {
		if id, err := uuid.Parse(d); err == nil {
			divID = &id
		}
	}
	teams, err := h.orgUC.ListTeams(r.Context(), divID)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, teams)
}

func (h *OrgHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req domainUC.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	team, err := h.orgUC.CreateTeam(r.Context(), req)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, team)
}

func (h *OrgHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	team, err := h.orgUC.UpdateTeam(r.Context(), id, req)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, team)
}

func (h *OrgHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.orgUC.DeleteTeam(r.Context(), id); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "team deleted"})
}
