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

type JiraHandler struct {
	issueTypeUC    domainUC.IssueTypeUseCase
	customFieldUC  domainUC.CustomFieldUseCase
	workflowUC     domainUC.WorkflowUseCase
	sprintUC       domainUC.SprintUseCase
	backlogUC      domainUC.BacklogUseCase
	commentUC      domainUC.CommentUseCase
	attachmentUC   domainUC.AttachmentUseCase
	labelUC        domainUC.LabelUseCase
	bulkOpUC       domainUC.BulkOperationUseCase
	searchUC       domainUC.SearchUseCase
}

func NewJiraHandler(
	issueTypeUC domainUC.IssueTypeUseCase,
	customFieldUC domainUC.CustomFieldUseCase,
	workflowUC domainUC.WorkflowUseCase,
	sprintUC domainUC.SprintUseCase,
	backlogUC domainUC.BacklogUseCase,
	commentUC domainUC.CommentUseCase,
	attachmentUC domainUC.AttachmentUseCase,
	labelUC domainUC.LabelUseCase,
	bulkOpUC domainUC.BulkOperationUseCase,
	searchUC domainUC.SearchUseCase,
) *JiraHandler {
	return &JiraHandler{
		issueTypeUC:   issueTypeUC,
		customFieldUC: customFieldUC,
		workflowUC:    workflowUC,
		sprintUC:      sprintUC,
		backlogUC:     backlogUC,
		commentUC:     commentUC,
		attachmentUC:  attachmentUC,
		labelUC:       labelUC,
		bulkOpUC:      bulkOpUC,
		searchUC:      searchUC,
	}
}

// Issue Type Handlers

func (h *JiraHandler) ListIssueTypes(w http.ResponseWriter, r *http.Request) {
	issueTypes, err := h.issueTypeUC.ListIssueTypes(r.Context())
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, issueTypes)
}

func (h *JiraHandler) GetIssueTypeScheme(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	scheme, err := h.issueTypeUC.GetIssueTypeScheme(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, scheme)
}

func (h *JiraHandler) CreateIssueTypeScheme(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.CreateIssueTypeSchemeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	scheme, err := h.issueTypeUC.CreateIssueTypeScheme(r.Context(), projectID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, scheme)
}

// Custom Field Handlers

func (h *JiraHandler) CreateCustomField(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.CreateCustomFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	field, err := h.customFieldUC.CreateCustomField(r.Context(), projectID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, field)
}

func (h *JiraHandler) ListCustomFields(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	fields, err := h.customFieldUC.ListCustomFields(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, fields)
}

func (h *JiraHandler) UpdateCustomField(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	fieldID, err := uuid.Parse(chi.URLParam(r, "fieldId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateCustomFieldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	field, err := h.customFieldUC.UpdateCustomField(r.Context(), projectID, fieldID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, field)
}

func (h *JiraHandler) DeleteCustomField(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	fieldID, err := uuid.Parse(chi.URLParam(r, "fieldId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.customFieldUC.DeleteCustomField(r.Context(), projectID, fieldID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Workflow Handlers

func (h *JiraHandler) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	workflow, err := h.workflowUC.GetWorkflow(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, workflow)
}

func (h *JiraHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.CreateWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	workflow, err := h.workflowUC.CreateWorkflow(r.Context(), projectID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, workflow)
}

func (h *JiraHandler) UpdateWorkflow(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	workflowID, err := uuid.Parse(chi.URLParam(r, "workflowId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.workflowUC.UpdateWorkflow(r.Context(), projectID, workflowID, req, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JiraHandler) TransitionRecord(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.TransitionRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.workflowUC.TransitionRecord(r.Context(), recordID, req, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Sprint Handlers

func (h *JiraHandler) CreateSprint(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.CreateSprintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	sprint, err := h.sprintUC.CreateSprint(r.Context(), projectID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, sprint)
}

func (h *JiraHandler) ListSprints(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	sprints, err := h.sprintUC.ListSprints(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, sprints)
}

func (h *JiraHandler) GetActiveSprint(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	sprint, err := h.sprintUC.GetActiveSprint(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, sprint)
}

func (h *JiraHandler) StartSprint(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	sprintID, err := uuid.Parse(chi.URLParam(r, "sprintId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	sprint, err := h.sprintUC.StartSprint(r.Context(), sprintID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, sprint)
}

func (h *JiraHandler) CompleteSprint(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	sprintID, err := uuid.Parse(chi.URLParam(r, "sprintId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	metrics, err := h.sprintUC.CompleteSprint(r.Context(), sprintID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, metrics)
}

func (h *JiraHandler) GetSprintRecords(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	sprintID, err := uuid.Parse(chi.URLParam(r, "sprintId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	records, err := h.sprintUC.GetSprintRecords(r.Context(), sprintID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, records)
}

// Backlog Handlers

func (h *JiraHandler) GetBacklog(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	records, err := h.backlogUC.GetBacklog(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, records)
}

func (h *JiraHandler) ReorderBacklog(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.ReorderBacklogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.backlogUC.ReorderBacklog(r.Context(), projectID, req, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JiraHandler) BulkAssignToSprint(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	var req domainUC.BulkAssignToSprintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.backlogUC.BulkAssignToSprint(r.Context(), req, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Comment Handlers

func (h *JiraHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.AddCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	comment, err := h.commentUC.AddComment(r.Context(), recordID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, comment)
}

func (h *JiraHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	comments, err := h.commentUC.ListComments(r.Context(), recordID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, comments)
}

func (h *JiraHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	commentID, err := uuid.Parse(chi.URLParam(r, "commentId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	comment, err := h.commentUC.UpdateComment(r.Context(), commentID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, comment)
}

func (h *JiraHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	commentID, err := uuid.Parse(chi.URLParam(r, "commentId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.commentUC.DeleteComment(r.Context(), commentID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Attachment Handlers

func (h *JiraHandler) UploadAttachment(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := r.ParseMultipartForm(50 * 1024 * 1024); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	defer file.Close()

	// Read file content
	fileContent := make([]byte, handler.Size)
	if _, err := file.Read(fileContent); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}

	fileUpload := &domainUC.FileUpload{
		FileName: handler.Filename,
		FileSize: handler.Size,
		FileType: handler.Header.Get("Content-Type"),
		Content:  fileContent,
	}

	attachment, err := h.attachmentUC.UploadAttachment(r.Context(), recordID, fileUpload, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, attachment)
}

func (h *JiraHandler) ListAttachments(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	attachments, err := h.attachmentUC.ListAttachments(r.Context(), recordID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, attachments)
}

func (h *JiraHandler) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	attachmentID, err := uuid.Parse(chi.URLParam(r, "attachmentId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.attachmentUC.DeleteAttachment(r.Context(), attachmentID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JiraHandler) DownloadAttachment(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	attachmentID, err := uuid.Parse(chi.URLParam(r, "attachmentId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}

	attachment, err := h.attachmentUC.GetAttachment(r.Context(), attachmentID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}

	http.ServeFile(w, r, attachment.FilePath)
}

// Label Handlers

func (h *JiraHandler) CreateLabel(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.CreateLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	label, err := h.labelUC.CreateLabel(r.Context(), projectID, req, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusCreated, label)
}

func (h *JiraHandler) ListLabels(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	labels, err := h.labelUC.ListLabels(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, labels)
}

func (h *JiraHandler) AddLabelToRecord(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	labelID, err := uuid.Parse(chi.URLParam(r, "labelId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.labelUC.AddLabelToRecord(r.Context(), recordID, labelID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JiraHandler) GetRecordLabels(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	labels, err := h.labelUC.GetRecordLabels(r.Context(), recordID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, labels)
}

func (h *JiraHandler) RemoveLabelFromRecord(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	recordID, err := uuid.Parse(chi.URLParam(r, "recordId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	labelID, err := uuid.Parse(chi.URLParam(r, "labelId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.labelUC.RemoveLabelFromRecord(r.Context(), recordID, labelID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JiraHandler) DeleteLabel(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	labelID, err := uuid.Parse(chi.URLParam(r, "labelId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.labelUC.DeleteLabel(r.Context(), labelID, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Bulk Operation Handlers

func (h *JiraHandler) BulkChangeStatus(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	var req domainUC.BulkChangeStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.bulkOpUC.BulkChangeStatus(r.Context(), req, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JiraHandler) BulkAssignTo(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	var req domainUC.BulkAssignToRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.bulkOpUC.BulkAssignTo(r.Context(), req, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JiraHandler) BulkAddLabel(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	var req domainUC.BulkAddLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.bulkOpUC.BulkAddLabel(r.Context(), req, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *JiraHandler) BulkDelete(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	var req domainUC.BulkDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.bulkOpUC.BulkDelete(r.Context(), req, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Search Handlers

func (h *JiraHandler) SearchRecords(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	query := r.URL.Query().Get("q")
	var filters domainUC.SearchFilters
	if err := r.ParseForm(); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	records, err := h.searchUC.SearchRecords(r.Context(), projectID, query, filters, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, records)
}

func (h *JiraHandler) SaveFilter(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	var req domainUC.SaveFilterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	if err := h.searchUC.SaveFilter(r.Context(), projectID, req, claims); err != nil {
		apperror.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *JiraHandler) ListSavedFilters(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	filters, err := h.searchUC.ListSavedFilters(r.Context(), projectID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, filters)
}

func (h *JiraHandler) ListWorkflowStatuses(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetClaims(r)
	workflowID, err := uuid.Parse(chi.URLParam(r, "workflowId"))
	if err != nil {
		apperror.WriteError(w, apperror.ErrValidation)
		return
	}
	statuses, err := h.workflowUC.ListStatuses(r.Context(), workflowID, claims)
	if err != nil {
		apperror.WriteError(w, err)
		return
	}
	apperror.WriteJSON(w, http.StatusOK, statuses)
}
