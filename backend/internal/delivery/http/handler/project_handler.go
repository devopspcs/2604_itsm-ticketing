package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/org/itsm/internal/delivery/http/middleware"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type ProjectHandler struct {
	projectUC domainUC.ProjectBoardUseCase
}

func NewProjectHandler(projectUC domainUC.ProjectBoardUseCase) *ProjectHandler {
	return &ProjectHandler{projectUC: projectUC}
}

// --- Project ---

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	var req domainUC.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	project, err := h.projectUC.CreateProject(r.Context(), req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projects, err := h.projectUC.ListProjects(r.Context(), claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	project, err := h.projectUC.GetProject(r.Context(), id, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	project, err := h.projectUC.UpdateProject(r.Context(), id, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.projectUC.DeleteProject(r.Context(), id, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Views ---

func (h *ProjectHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	data, err := h.projectUC.GetHome(r.Context(), claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, data)
}

func (h *ProjectHandler) GetCalendar(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	q := r.URL.Query()
	month, _ := strconv.Atoi(q.Get("month"))
	year, _ := strconv.Atoi(q.Get("year"))
	if month < 1 || month > 12 || year < 1 {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	records, err := h.projectUC.GetCalendar(r.Context(), month, year, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, records)
}

func (h *ProjectHandler) GetActivities(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	activities, err := h.projectUC.GetActivities(r.Context(), id, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, activities)
}

// --- Columns ---

func (h *ProjectHandler) CreateColumn(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.CreateColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	col, err := h.projectUC.CreateColumn(r.Context(), projectID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, col)
}

func (h *ProjectHandler) UpdateColumn(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	columnID, err := uuid.Parse(chi.URLParam(r, "columnId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	col, err := h.projectUC.UpdateColumn(r.Context(), projectID, columnID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, col)
}

func (h *ProjectHandler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	columnID, err := uuid.Parse(chi.URLParam(r, "columnId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.projectUC.DeleteColumn(r.Context(), projectID, columnID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) ReorderColumns(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.ReorderColumnsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.projectUC.ReorderColumns(r.Context(), projectID, req, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "columns reordered"})
}

// --- Records ---

func (h *ProjectHandler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.CreateRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	record, err := h.projectUC.CreateRecord(r.Context(), projectID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, record)
}

func (h *ProjectHandler) GetRecord(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	record, err := h.projectUC.GetRecord(r.Context(), projectID, recordID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, record)
}

func (h *ProjectHandler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	record, err := h.projectUC.UpdateRecord(r.Context(), projectID, recordID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, record)
}

func (h *ProjectHandler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.projectUC.DeleteRecord(r.Context(), projectID, recordID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) MoveRecord(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.MoveRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.projectUC.MoveRecord(r.Context(), projectID, recordID, req, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "record moved"})
}

func (h *ProjectHandler) CompleteRecord(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	record, err := h.projectUC.CompleteRecord(r.Context(), projectID, recordID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, record)
}

// --- Members ---

func (h *ProjectHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	members, err := h.projectUC.ListMembers(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, members)
}

func (h *ProjectHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var body struct {
		UserID uuid.UUID `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.projectUC.InviteMember(r.Context(), projectID, body.UserID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, map[string]string{"message": "member invited"})
}

func (h *ProjectHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	memberID, err := uuid.Parse(chi.URLParam(r, "memberId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.projectUC.RemoveMember(r.Context(), projectID, memberID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var body struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Text == "" {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.projectUC.AddComment(r.Context(), projectID, recordID, body.Text, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, map[string]string{"message": "comment added"})
}
