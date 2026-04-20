package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
)

type UserClaims struct {
	UserID uuid.UUID
	Role   entity.Role
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthUseCase interface {
	Login(ctx context.Context, req LoginRequest) (*TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type CreateUserRequest struct {
	FullName     string           `json:"full_name" validate:"required"`
	Email        string           `json:"email" validate:"required,email"`
	Password     string           `json:"password" validate:"required,min=8"`
	Role         entity.Role      `json:"role" validate:"required,oneof=user approver admin"`
	DepartmentID *uuid.UUID       `json:"department_id"`
	DivisionID   *uuid.UUID       `json:"division_id"`
	TeamID       *uuid.UUID       `json:"team_id"`
	Position     *entity.Position `json:"position"`
}

type UpdateUserOrgRequest struct {
	DepartmentID *uuid.UUID       `json:"department_id"`
	DivisionID   *uuid.UUID       `json:"division_id"`
	TeamID       *uuid.UUID       `json:"team_id"`
	Position     *entity.Position `json:"position"`
}

type UserUseCase interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*entity.User, error)
	UpdateUserRole(ctx context.Context, userID uuid.UUID, role entity.Role) error
	DeactivateUser(ctx context.Context, userID uuid.UUID) error
	ActivateUser(ctx context.Context, userID uuid.UUID) error
	GetUsers(ctx context.Context) ([]*entity.User, error)
	UpdateUserOrg(ctx context.Context, userID uuid.UUID, req UpdateUserOrgRequest) (*entity.User, error)
}

type CreateTicketRequest struct {
	Title       string            `json:"title" validate:"required"`
	Description string            `json:"description" validate:"required"`
	Type        entity.TicketType `json:"type" validate:"required,oneof=change_request incident helpdesk_request"`
	Category    string            `json:"category"`
	Priority    entity.Priority   `json:"priority" validate:"required,oneof=low medium high critical"`
}

type UpdateTicketRequest struct {
	Title       *string              `json:"title"`
	Description *string              `json:"description"`
	Category    *string              `json:"category"`
	Priority    *entity.Priority     `json:"priority"`
	Status      *entity.TicketStatus `json:"status"`
}

type TicketUseCase interface {
	CreateTicket(ctx context.Context, req CreateTicketRequest, requester UserClaims) (*entity.Ticket, error)
	GetTicket(ctx context.Context, id uuid.UUID, requester UserClaims) (*entity.Ticket, error)
	ListTickets(ctx context.Context, filter repository.TicketFilter, requester UserClaims) (*repository.PaginatedTickets, error)
	UpdateTicket(ctx context.Context, id uuid.UUID, req UpdateTicketRequest, requester UserClaims) (*entity.Ticket, error)
}

type ApprovalDecisionRequest struct {
	TicketID uuid.UUID               `json:"ticket_id"`
	Decision entity.ApprovalDecision `json:"decision" validate:"required,oneof=approved rejected"`
	Comment  string                  `json:"comment"`
}

type ApprovalUseCase interface {
	SubmitForApproval(ctx context.Context, ticketID uuid.UUID, requester UserClaims) error
	Decide(ctx context.Context, req ApprovalDecisionRequest, approver UserClaims) error
	GetApprovalHistory(ctx context.Context, ticketID uuid.UUID) ([]*entity.Approval, error)
}

type AssignmentUseCase interface {
	AssignTicket(ctx context.Context, ticketID uuid.UUID, assigneeID uuid.UUID, requester UserClaims) error
}

type DashboardFilter struct {
	DateFrom *string `json:"date_from"`
	DateTo   *string `json:"date_to"`
}

type DashboardStats struct {
	TotalTickets  int64                         `json:"total_tickets"`
	ByStatus      map[entity.TicketStatus]int64 `json:"by_status"`
	ByType        map[entity.TicketType]int64   `json:"by_type"`
	ByPriority    map[entity.Priority]int64     `json:"by_priority"`
	RecentTickets      []*entity.Ticket              `json:"recent_tickets"`
	SLAComplianceRate  float64                       `json:"sla_compliance_rate"`
	AvgResolutionHours float64                       `json:"avg_resolution_hours"`
	OnTimeCount        int64                         `json:"on_time_count"`
	BreachedCount      int64                         `json:"breached_count"`
}

type DashboardUseCase interface {
	GetStats(ctx context.Context, filter DashboardFilter, requester UserClaims) (*DashboardStats, error)
}

type NotificationUseCase interface {
	GetNotifications(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error)
	MarkAsRead(ctx context.Context, notifID uuid.UUID, userID uuid.UUID) error
}

type CreateWebhookRequest struct {
	URL    string                `json:"url" validate:"required,url"`
	Events []entity.WebhookEvent `json:"events" validate:"required,min=1"`
	Secret string                `json:"secret" validate:"required"`
}

type WebhookUseCase interface {
	CreateWebhookConfig(ctx context.Context, req CreateWebhookRequest) (*entity.WebhookConfig, error)
	UpdateWebhookConfig(ctx context.Context, id uuid.UUID, req CreateWebhookRequest) error
	DeleteWebhookConfig(ctx context.Context, id uuid.UUID) error
	ListWebhookConfigs(ctx context.Context) ([]*entity.WebhookConfig, error)
	Dispatch(ctx context.Context, event entity.WebhookEvent, payload interface{}) error
}

type CreateDepartmentRequest struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
}

type UpdateDepartmentRequest struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
}

type CreateDivisionRequest struct {
	DepartmentID uuid.UUID `json:"department_id" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	Code         string    `json:"code" validate:"required"`
}

type UpdateDivisionRequest struct {
	DepartmentID uuid.UUID `json:"department_id" validate:"required"`
	Name         string    `json:"name" validate:"required"`
	Code         string    `json:"code" validate:"required"`
}

type CreateTeamRequest struct {
	DivisionID uuid.UUID `json:"division_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
}

type UpdateTeamRequest struct {
	DivisionID uuid.UUID `json:"division_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
}

type OrgUseCase interface {
	CreateDepartment(ctx context.Context, req CreateDepartmentRequest) (*entity.Department, error)
	UpdateDepartment(ctx context.Context, id uuid.UUID, req UpdateDepartmentRequest) (*entity.Department, error)
	DeleteDepartment(ctx context.Context, id uuid.UUID) error
	ListDepartments(ctx context.Context) ([]*entity.Department, error)

	CreateDivision(ctx context.Context, req CreateDivisionRequest) (*entity.Division, error)
	UpdateDivision(ctx context.Context, id uuid.UUID, req UpdateDivisionRequest) (*entity.Division, error)
	DeleteDivision(ctx context.Context, id uuid.UUID) error
	ListDivisions(ctx context.Context, departmentID *uuid.UUID) ([]*entity.Division, error)

	CreateTeam(ctx context.Context, req CreateTeamRequest) (*entity.Team, error)
	UpdateTeam(ctx context.Context, id uuid.UUID, req UpdateTeamRequest) (*entity.Team, error)
	DeleteTeam(ctx context.Context, id uuid.UUID) error
	ListTeams(ctx context.Context, divisionID *uuid.UUID) ([]*entity.Team, error)
}

// --- Project Board ---

type CreateProjectRequest struct {
	Name      string `json:"name" validate:"required"`
	IconColor string `json:"icon_color"`
}

type UpdateProjectRequest struct {
	Name      *string `json:"name"`
	IconColor *string `json:"icon_color"`
}

type CreateColumnRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateColumnRequest struct {
	Name string `json:"name" validate:"required"`
}

type ReorderColumnsRequest struct {
	ColumnIDs []uuid.UUID `json:"column_ids" validate:"required"`
}

type CreateRecordRequest struct {
	ColumnID    uuid.UUID  `json:"column_id" validate:"required"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description"`
	AssignedTo  *uuid.UUID `json:"assigned_to"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateRecordRequest struct {
	Title       *string      `json:"title"`
	Description *string      `json:"description"`
	AssignedTo  *uuid.UUID   `json:"assigned_to"`
	Assignees   []uuid.UUID  `json:"assignees"`
	DueDate     *time.Time   `json:"due_date"`
}

type MoveRecordRequest struct {
	TargetColumnID uuid.UUID `json:"target_column_id" validate:"required"`
	Position       int       `json:"position" validate:"min=0"`
}

type ProjectHomeData struct {
	OverdueCount     int                          `json:"overdue_count"`
	RecentActivities []*entity.ProjectActivityLog `json:"recent_activities"`
}

type ProjectColumnWithRecords struct {
	entity.ProjectColumn
	Records []*entity.ProjectRecord `json:"records"`
}

type ProjectDetailResponse struct {
	entity.Project
	Columns []ProjectColumnWithRecords `json:"columns"`
}

type ProjectBoardUseCase interface {
	// Project CRUD
	CreateProject(ctx context.Context, req CreateProjectRequest, requester UserClaims) (*entity.Project, error)
	GetProject(ctx context.Context, id uuid.UUID, requester UserClaims) (*ProjectDetailResponse, error)
	ListProjects(ctx context.Context, requester UserClaims) ([]*entity.Project, error)
	UpdateProject(ctx context.Context, id uuid.UUID, req UpdateProjectRequest, requester UserClaims) (*entity.Project, error)
	DeleteProject(ctx context.Context, id uuid.UUID, requester UserClaims) error

	// Column CRUD
	CreateColumn(ctx context.Context, projectID uuid.UUID, req CreateColumnRequest, requester UserClaims) (*entity.ProjectColumn, error)
	UpdateColumn(ctx context.Context, projectID uuid.UUID, columnID uuid.UUID, req UpdateColumnRequest, requester UserClaims) (*entity.ProjectColumn, error)
	DeleteColumn(ctx context.Context, projectID uuid.UUID, columnID uuid.UUID, requester UserClaims) error
	ReorderColumns(ctx context.Context, projectID uuid.UUID, req ReorderColumnsRequest, requester UserClaims) error

	// Record CRUD
	CreateRecord(ctx context.Context, projectID uuid.UUID, req CreateRecordRequest, requester UserClaims) (*entity.ProjectRecord, error)
	GetRecord(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, requester UserClaims) (*entity.ProjectRecord, error)
	UpdateRecord(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, req UpdateRecordRequest, requester UserClaims) (*entity.ProjectRecord, error)
	DeleteRecord(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, requester UserClaims) error
	MoveRecord(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, req MoveRecordRequest, requester UserClaims) error
	CompleteRecord(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, requester UserClaims) (*entity.ProjectRecord, error)

	// Views
	GetHome(ctx context.Context, requester UserClaims) (*ProjectHomeData, error)
	GetCalendar(ctx context.Context, month int, year int, requester UserClaims) ([]*entity.ProjectRecord, error)
	GetActivities(ctx context.Context, projectID uuid.UUID, requester UserClaims) ([]*entity.ProjectActivityLog, error)

	// Members
	InviteMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, requester UserClaims) error
	RemoveMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, requester UserClaims) error
	ListMembers(ctx context.Context, projectID uuid.UUID, requester UserClaims) ([]*entity.ProjectMember, error)

	// Comments
	AddComment(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, text string, requester UserClaims) error
}


// --- Jira-like Features ---

type CreateCustomFieldRequest struct {
	Name       string   `json:"name" validate:"required"`
	FieldType  string   `json:"field_type" validate:"required,oneof=text textarea dropdown multiselect date number checkbox"`
	IsRequired bool     `json:"is_required"`
	Options    []string `json:"options"`
}

type UpdateCustomFieldRequest struct {
	Name    *string  `json:"name"`
	Options []string `json:"options"`
}

type CreateWorkflowRequest struct {
	Name          string   `json:"name" validate:"required"`
	InitialStatus string   `json:"initial_status" validate:"required"`
	Statuses      []string `json:"statuses" validate:"required,min=1"`
}

type UpdateWorkflowRequest struct {
	Statuses    []string       `json:"statuses"`
	Transitions []TransitionDef `json:"transitions"`
}

type TransitionDef struct {
	FromStatus string `json:"from_status" validate:"required"`
	ToStatus   string `json:"to_status" validate:"required"`
}

type CreateSprintRequest struct {
	Name      string     `json:"name" validate:"required"`
	Goal      string     `json:"goal"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date" validate:"required"`
}

type SprintMetrics struct {
	TotalRecords      int     `json:"total_records"`
	CompletedRecords  int     `json:"completed_records"`
	CompletionPercent float64 `json:"completion_percent"`
	Velocity          int     `json:"velocity"`
}

type SearchFilters struct {
	IssueType       *uuid.UUID
	Status          *string
	Assignee        *uuid.UUID
	Label           *uuid.UUID
	Sprint          *uuid.UUID
	DueDateFrom     *time.Time
	DueDateTo       *time.Time
	CustomFields    map[uuid.UUID]string
}

type SavedFilter struct {
	ID      uuid.UUID    `json:"id"`
	Name    string       `json:"name"`
	Filters SearchFilters `json:"filters"`
}

type FileUpload struct {
	FileName string
	FileSize int64
	FileType string
	Content  []byte
}

type CreateIssueTypeSchemeRequest struct {
	Name         string      `json:"name" validate:"required"`
	IssueTypeIDs []uuid.UUID `json:"issue_type_ids" validate:"required,min=1"`
}

type TransitionRecordRequest struct {
	ToStatusID uuid.UUID `json:"to_status_id" validate:"required"`
}

type ReorderBacklogRequest struct {
	RecordIDs []uuid.UUID `json:"record_ids" validate:"required"`
}

type BulkAssignToSprintRequest struct {
	SprintID  uuid.UUID   `json:"sprint_id" validate:"required"`
	RecordIDs []uuid.UUID `json:"record_ids" validate:"required"`
}

type AddCommentRequest struct {
	Text string `json:"text" validate:"required"`
}

type UpdateCommentRequest struct {
	Text string `json:"text" validate:"required"`
}

type CreateLabelRequest struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color" validate:"required"`
}

type BulkChangeStatusRequest struct {
	RecordIDs []uuid.UUID `json:"record_ids" validate:"required"`
	StatusID  uuid.UUID   `json:"status_id" validate:"required"`
}

type BulkAssignToRequest struct {
	RecordIDs []uuid.UUID `json:"record_ids" validate:"required"`
	AssigneeID uuid.UUID  `json:"assignee_id" validate:"required"`
}

type BulkAddLabelRequest struct {
	RecordIDs []uuid.UUID `json:"record_ids" validate:"required"`
	LabelID   uuid.UUID   `json:"label_id" validate:"required"`
}

type BulkDeleteRequest struct {
	RecordIDs []uuid.UUID `json:"record_ids" validate:"required"`
}

type SaveFilterRequest struct {
	Name    string        `json:"name" validate:"required"`
	Filters SearchFilters `json:"filters"`
}

type IssueTypeUseCase interface {
	ListIssueTypes(ctx context.Context) ([]*entity.IssueType, error)
	GetIssueTypeScheme(ctx context.Context, projectID uuid.UUID, requester UserClaims) (*entity.IssueTypeScheme, error)
	CreateIssueTypeScheme(ctx context.Context, projectID uuid.UUID, req CreateIssueTypeSchemeRequest, requester UserClaims) (*entity.IssueTypeScheme, error)
	UpdateIssueTypeScheme(ctx context.Context, schemeID uuid.UUID, issueTypeIDs []uuid.UUID, requester UserClaims) error
}

type CustomFieldUseCase interface {
	CreateCustomField(ctx context.Context, projectID uuid.UUID, req CreateCustomFieldRequest, requester UserClaims) (*entity.CustomField, error)
	ListCustomFields(ctx context.Context, projectID uuid.UUID, requester UserClaims) ([]*entity.CustomField, error)
	UpdateCustomField(ctx context.Context, projectID uuid.UUID, fieldID uuid.UUID, req UpdateCustomFieldRequest, requester UserClaims) (*entity.CustomField, error)
	DeleteCustomField(ctx context.Context, projectID uuid.UUID, fieldID uuid.UUID, requester UserClaims) error
	SetFieldValue(ctx context.Context, recordID uuid.UUID, fieldID uuid.UUID, value string, requester UserClaims) error
	GetFieldValues(ctx context.Context, recordID uuid.UUID, requester UserClaims) ([]*entity.CustomFieldValue, error)
}

type WorkflowUseCase interface {
	GetWorkflow(ctx context.Context, projectID uuid.UUID, requester UserClaims) (*entity.Workflow, error)
	CreateWorkflow(ctx context.Context, projectID uuid.UUID, req CreateWorkflowRequest, requester UserClaims) (*entity.Workflow, error)
	UpdateWorkflow(ctx context.Context, projectID uuid.UUID, workflowID uuid.UUID, req UpdateWorkflowRequest, requester UserClaims) error
	ListStatuses(ctx context.Context, workflowID uuid.UUID, requester UserClaims) ([]*entity.WorkflowStatus, error)
	CanTransition(ctx context.Context, recordID uuid.UUID, toStatusID uuid.UUID, requester UserClaims) (bool, error)
	TransitionRecord(ctx context.Context, recordID uuid.UUID, req TransitionRecordRequest, requester UserClaims) error
}

type SprintUseCase interface {
	CreateSprint(ctx context.Context, projectID uuid.UUID, req CreateSprintRequest, requester UserClaims) (*entity.Sprint, error)
	ListSprints(ctx context.Context, projectID uuid.UUID, requester UserClaims) ([]*entity.Sprint, error)
	GetActiveSprint(ctx context.Context, projectID uuid.UUID, requester UserClaims) (*entity.Sprint, error)
	StartSprint(ctx context.Context, sprintID uuid.UUID, requester UserClaims) (*entity.Sprint, error)
	CompleteSprint(ctx context.Context, sprintID uuid.UUID, requester UserClaims) (*SprintMetrics, error)
	AssignRecordToSprint(ctx context.Context, recordID uuid.UUID, sprintID uuid.UUID, requester UserClaims) error
	RemoveRecordFromSprint(ctx context.Context, recordID uuid.UUID, requester UserClaims) error
	GetSprintRecords(ctx context.Context, sprintID uuid.UUID, requester UserClaims) ([]*entity.ProjectRecord, error)
}

type BacklogUseCase interface {
	GetBacklog(ctx context.Context, projectID uuid.UUID, requester UserClaims) ([]*entity.ProjectRecord, error)
	ReorderBacklog(ctx context.Context, projectID uuid.UUID, req ReorderBacklogRequest, requester UserClaims) error
	BulkAssignToSprint(ctx context.Context, req BulkAssignToSprintRequest, requester UserClaims) error
}

type CommentUseCase interface {
	AddComment(ctx context.Context, recordID uuid.UUID, req AddCommentRequest, requester UserClaims) (*entity.Comment, error)
	ListComments(ctx context.Context, recordID uuid.UUID, requester UserClaims) ([]*entity.Comment, error)
	UpdateComment(ctx context.Context, commentID uuid.UUID, req UpdateCommentRequest, requester UserClaims) (*entity.Comment, error)
	DeleteComment(ctx context.Context, commentID uuid.UUID, requester UserClaims) error
	ParseMentions(text string) []string
}

type AttachmentUseCase interface {
	UploadAttachment(ctx context.Context, recordID uuid.UUID, file *FileUpload, requester UserClaims) (*entity.Attachment, error)
	ListAttachments(ctx context.Context, recordID uuid.UUID, requester UserClaims) ([]*entity.Attachment, error)
	GetAttachment(ctx context.Context, attachmentID uuid.UUID, requester UserClaims) (*entity.Attachment, error)
	DeleteAttachment(ctx context.Context, attachmentID uuid.UUID, requester UserClaims) error
}

type LabelUseCase interface {
	CreateLabel(ctx context.Context, projectID uuid.UUID, req CreateLabelRequest, requester UserClaims) (*entity.Label, error)
	ListLabels(ctx context.Context, projectID uuid.UUID, requester UserClaims) ([]*entity.Label, error)
	AddLabelToRecord(ctx context.Context, recordID uuid.UUID, labelID uuid.UUID, requester UserClaims) error
	RemoveLabelFromRecord(ctx context.Context, recordID uuid.UUID, labelID uuid.UUID, requester UserClaims) error
	GetRecordLabels(ctx context.Context, recordID uuid.UUID, requester UserClaims) ([]*entity.Label, error)
	DeleteLabel(ctx context.Context, labelID uuid.UUID, requester UserClaims) error
}

type BulkOperationUseCase interface {
	BulkChangeStatus(ctx context.Context, req BulkChangeStatusRequest, requester UserClaims) error
	BulkAssignTo(ctx context.Context, req BulkAssignToRequest, requester UserClaims) error
	BulkAddLabel(ctx context.Context, req BulkAddLabelRequest, requester UserClaims) error
	BulkDelete(ctx context.Context, req BulkDeleteRequest, requester UserClaims) error
}

type SearchUseCase interface {
	SearchRecords(ctx context.Context, projectID uuid.UUID, query string, filters SearchFilters, requester UserClaims) ([]*entity.ProjectRecord, error)
	SaveFilter(ctx context.Context, projectID uuid.UUID, req SaveFilterRequest, requester UserClaims) error
	ListSavedFilters(ctx context.Context, projectID uuid.UUID, requester UserClaims) ([]*SavedFilter, error)
}

// --- Project Board Features ---

// Reports
type ReportsSummary struct {
	TotalRecords   int            `json:"total_records"`
	CompletedCount int            `json:"completed_count"`
	OpenCount      int            `json:"open_count"`
	ByStatus       map[string]int `json:"by_status"`
}

type VelocityDataPoint struct {
	SprintName     string `json:"sprint_name"`
	TotalRecords   int    `json:"total_records"`
	CompletedCount int    `json:"completed_count"`
}

type BurndownData struct {
	SprintName string     `json:"sprint_name"`
	StartDate  *time.Time `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	TotalCount int        `json:"total_count"`
	DoneCount  int        `json:"done_count"`
	HasActive  bool       `json:"has_active"`
}

type ReportsUseCase interface {
	GetSummary(ctx context.Context, projectID uuid.UUID, requester UserClaims) (*ReportsSummary, error)
	GetVelocity(ctx context.Context, projectID uuid.UUID, requester UserClaims) ([]*VelocityDataPoint, error)
	GetBurndown(ctx context.Context, projectID uuid.UUID, requester UserClaims) (*BurndownData, error)
}

// Releases
type CreateReleaseRequest struct {
	Name        string     `json:"name" validate:"required"`
	Version     string     `json:"version" validate:"required"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	ReleaseDate *time.Time `json:"release_date"`
	Status      string     `json:"status"`
}

type UpdateReleaseRequest struct {
	Name        *string    `json:"name"`
	Version     *string    `json:"version"`
	Description *string    `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	ReleaseDate *time.Time `json:"release_date"`
	Status      *string    `json:"status"`
}

type ReleaseWithProgress struct {
	entity.Release
	TotalRecords    int     `json:"total_records"`
	CompletedCount  int     `json:"completed_count"`
	ProgressPercent float64 `json:"progress_percent"`
}

type ReleaseUseCase interface {
	Create(ctx context.Context, projectID uuid.UUID, req CreateReleaseRequest, requester UserClaims) (*entity.Release, error)
	List(ctx context.Context, projectID uuid.UUID, requester UserClaims) ([]*ReleaseWithProgress, error)
	Get(ctx context.Context, projectID uuid.UUID, releaseID uuid.UUID, requester UserClaims) (*ReleaseWithProgress, error)
	Update(ctx context.Context, projectID uuid.UUID, releaseID uuid.UUID, req UpdateReleaseRequest, requester UserClaims) (*entity.Release, error)
	Delete(ctx context.Context, projectID uuid.UUID, releaseID uuid.UUID, requester UserClaims) error
	AssignRecord(ctx context.Context, projectID uuid.UUID, releaseID uuid.UUID, recordID uuid.UUID, requester UserClaims) error
	RemoveRecord(ctx context.Context, projectID uuid.UUID, releaseID uuid.UUID, recordID uuid.UUID, requester UserClaims) error
}

// Components
type CreateComponentRequest struct {
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description"`
	LeadUserID  *uuid.UUID `json:"lead_user_id"`
}

type UpdateComponentRequest struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	LeadUserID  *uuid.UUID `json:"lead_user_id"`
}

type ComponentWithCount struct {
	entity.Component
	RecordCount int `json:"record_count"`
}

type ComponentUseCase interface {
	Create(ctx context.Context, projectID uuid.UUID, req CreateComponentRequest, requester UserClaims) (*entity.Component, error)
	List(ctx context.Context, projectID uuid.UUID, requester UserClaims) ([]*ComponentWithCount, error)
	Get(ctx context.Context, projectID uuid.UUID, componentID uuid.UUID, requester UserClaims) (*ComponentWithCount, error)
	Update(ctx context.Context, projectID uuid.UUID, componentID uuid.UUID, req UpdateComponentRequest, requester UserClaims) (*entity.Component, error)
	Delete(ctx context.Context, projectID uuid.UUID, componentID uuid.UUID, requester UserClaims) error
	AssignRecord(ctx context.Context, projectID uuid.UUID, componentID uuid.UUID, recordID uuid.UUID, requester UserClaims) error
	RemoveRecord(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, requester UserClaims) error
	ListRecords(ctx context.Context, projectID uuid.UUID, componentID uuid.UUID, requester UserClaims) ([]*entity.ProjectRecord, error)
}
