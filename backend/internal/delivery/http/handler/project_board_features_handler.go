package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/org/itsm/internal/delivery/http/middleware"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type ProjectBoardFeaturesHandler struct {
	reportsUC   domainUC.ReportsUseCase
	releaseUC   domainUC.ReleaseUseCase
	componentUC domainUC.ComponentUseCase
	recordRepo  repository.ProjectRecordRepository
	actLogRepo  repository.ProjectActivityLogRepository
	memberRepo  repository.ProjectMemberRepository
}

func NewProjectBoardFeaturesHandler(
	reportsUC domainUC.ReportsUseCase,
	releaseUC domainUC.ReleaseUseCase,
	componentUC domainUC.ComponentUseCase,
	recordRepo repository.ProjectRecordRepository,
	actLogRepo repository.ProjectActivityLogRepository,
	memberRepo repository.ProjectMemberRepository,
) *ProjectBoardFeaturesHandler {
	return &ProjectBoardFeaturesHandler{
		reportsUC:   reportsUC,
		releaseUC:   releaseUC,
		componentUC: componentUC,
		recordRepo:  recordRepo,
		actLogRepo:  actLogRepo,
		memberRepo:  memberRepo,
	}
}

// --- Reports ---

func (h *ProjectBoardFeaturesHandler) GetReportsSummary(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	summary, err := h.reportsUC.GetSummary(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, summary)
}

func (h *ProjectBoardFeaturesHandler) GetReportsVelocity(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	data, err := h.reportsUC.GetVelocity(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, data)
}

func (h *ProjectBoardFeaturesHandler) GetReportsBurndown(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	data, err := h.reportsUC.GetBurndown(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, data)
}

// --- Releases ---

func (h *ProjectBoardFeaturesHandler) CreateRelease(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.CreateReleaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	release, err := h.releaseUC.Create(r.Context(), projectID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, release)
}

func (h *ProjectBoardFeaturesHandler) ListReleases(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	releases, err := h.releaseUC.List(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, releases)
}

func (h *ProjectBoardFeaturesHandler) GetRelease(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	releaseID, err := uuid.Parse(chi.URLParam(r, "releaseId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	release, err := h.releaseUC.Get(r.Context(), projectID, releaseID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, release)
}

func (h *ProjectBoardFeaturesHandler) UpdateRelease(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	releaseID, err := uuid.Parse(chi.URLParam(r, "releaseId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateReleaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	release, err := h.releaseUC.Update(r.Context(), projectID, releaseID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, release)
}

func (h *ProjectBoardFeaturesHandler) DeleteRelease(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	releaseID, err := uuid.Parse(chi.URLParam(r, "releaseId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.releaseUC.Delete(r.Context(), projectID, releaseID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectBoardFeaturesHandler) AssignRecordToRelease(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	releaseID, err := uuid.Parse(chi.URLParam(r, "releaseId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.releaseUC.AssignRecord(r.Context(), projectID, releaseID, recordID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, map[string]string{"message": "record assigned to release"})
}

func (h *ProjectBoardFeaturesHandler) RemoveRecordFromRelease(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	releaseID, err := uuid.Parse(chi.URLParam(r, "releaseId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.releaseUC.RemoveRecord(r.Context(), projectID, releaseID, recordID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Components ---

func (h *ProjectBoardFeaturesHandler) CreateComponent(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.CreateComponentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	comp, err := h.componentUC.Create(r.Context(), projectID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, comp)
}

func (h *ProjectBoardFeaturesHandler) ListComponents(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	components, err := h.componentUC.List(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, components)
}

func (h *ProjectBoardFeaturesHandler) GetComponent(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	componentID, err := uuid.Parse(chi.URLParam(r, "componentId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	comp, err := h.componentUC.Get(r.Context(), projectID, componentID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, comp)
}

func (h *ProjectBoardFeaturesHandler) UpdateComponent(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	componentID, err := uuid.Parse(chi.URLParam(r, "componentId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateComponentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	comp, err := h.componentUC.Update(r.Context(), projectID, componentID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, comp)
}

func (h *ProjectBoardFeaturesHandler) DeleteComponent(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	componentID, err := uuid.Parse(chi.URLParam(r, "componentId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.componentUC.Delete(r.Context(), projectID, componentID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectBoardFeaturesHandler) AssignRecordToComponent(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	componentID, err := uuid.Parse(chi.URLParam(r, "componentId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.componentUC.AssignRecord(r.Context(), projectID, componentID, recordID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, map[string]string{"message": "record assigned to component"})
}

func (h *ProjectBoardFeaturesHandler) RemoveRecordFromComponent(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	_, err = uuid.Parse(chi.URLParam(r, "componentId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.componentUC.RemoveRecord(r.Context(), projectID, recordID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectBoardFeaturesHandler) ListComponentRecords(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	componentID, err := uuid.Parse(chi.URLParam(r, "componentId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	records, err := h.componentUC.ListRecords(r.Context(), projectID, componentID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, records)
}

// --- Issues (direct repository call) ---

func (h *ProjectBoardFeaturesHandler) ListIssues(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	// Check membership
	isMember, _ := h.memberRepo.IsMember(r.Context(), projectID, claims.UserID)
	if !isMember {
		apperror.WriteError(w, apperror.ErrForbidden)
		return
	}

	q := r.URL.Query()
	filter := repository.IssuesFilter{
		Page:     1,
		PageSize: 20,
	}
	if v := q.Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p >= 1 {
			filter.Page = p
		}
	}
	if v := q.Get("page_size"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil && ps >= 1 && ps <= 100 {
			filter.PageSize = ps
		}
	}
	if v := q.Get("search"); v != "" {
		filter.Search = &v
	}
	if v := q.Get("status_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.StatusID = &id
		}
	}
	if v := q.Get("assignee_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.AssigneeID = &id
		}
	}
	if v := q.Get("issue_type"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.IssueType = &id
		}
	}
	if v := q.Get("label_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.LabelID = &id
		}
	}

	result, err := h.recordRepo.ListByProjectPaginated(r.Context(), projectID, filter)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, result)
}

// --- Activity Log (direct repository call) ---

func (h *ProjectBoardFeaturesHandler) ListActivityLog(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	// Check membership
	isMember, _ := h.memberRepo.IsMember(r.Context(), projectID, claims.UserID)
	if !isMember {
		apperror.WriteError(w, apperror.ErrForbidden)
		return
	}

	q := r.URL.Query()
	filter := repository.ActivityLogFilter{
		Page:     1,
		PageSize: 20,
	}
	if v := q.Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p >= 1 {
			filter.Page = p
		}
	}
	if v := q.Get("page_size"); v != "" {
		if ps, err := strconv.Atoi(v); err == nil && ps >= 1 && ps <= 100 {
			filter.PageSize = ps
		}
	}
	if v := q.Get("action_type"); v != "" {
		filter.ActionType = &v
	}

	result, err := h.actLogRepo.ListByProjectPaginated(r.Context(), projectID, filter)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, result)
}
