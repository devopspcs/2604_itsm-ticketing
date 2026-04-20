package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
)

type ProjectRecordFilter struct {
	ColumnID    *uuid.UUID
	ProjectID   *uuid.UUID
	AssignedTo  *uuid.UUID
	Search      *string
	DueDateFrom *time.Time
	DueDateTo   *time.Time
}

type ProjectRepository interface {
	Create(ctx context.Context, project *entity.Project) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Project, error)
	List(ctx context.Context, createdBy uuid.UUID) ([]*entity.Project, error)
	Update(ctx context.Context, project *entity.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProjectColumnRepository interface {
	Create(ctx context.Context, col *entity.ProjectColumn) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ProjectColumn, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.ProjectColumn, error)
	Update(ctx context.Context, col *entity.ProjectColumn) error
	Delete(ctx context.Context, id uuid.UUID) error
	HasRecords(ctx context.Context, id uuid.UUID) (bool, error)
	GetMaxPosition(ctx context.Context, projectID uuid.UUID) (int, error)
	BulkUpdatePositions(ctx context.Context, projectID uuid.UUID, positions map[uuid.UUID]int) error
}

type ProjectRecordRepository interface {
	Create(ctx context.Context, record *entity.ProjectRecord) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ProjectRecord, error)
	ListByColumn(ctx context.Context, columnID uuid.UUID) ([]*entity.ProjectRecord, error)
	ListByProject(ctx context.Context, projectID uuid.UUID, filter ProjectRecordFilter) ([]*entity.ProjectRecord, error)
	Update(ctx context.Context, record *entity.ProjectRecord) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetMaxPosition(ctx context.Context, columnID uuid.UUID) (int, error)
	BulkUpdatePositions(ctx context.Context, columnID uuid.UUID, positions map[uuid.UUID]int) error
	ListByDueDateRange(ctx context.Context, createdBy uuid.UUID, from time.Time, to time.Time) ([]*entity.ProjectRecord, error)
	CountOverdue(ctx context.Context, createdBy uuid.UUID) (int, error)
	SetAssignees(ctx context.Context, recordID uuid.UUID, userIDs []uuid.UUID) error
	GetAssignees(ctx context.Context, recordID uuid.UUID) ([]uuid.UUID, error)
	CountByProject(ctx context.Context, projectID uuid.UUID) (int, error)
	CountByProjectAndStatus(ctx context.Context, projectID uuid.UUID, isCompleted bool) (int, error)
	CountByProjectGroupedByStatus(ctx context.Context, projectID uuid.UUID) (map[string]int, error)
	ListByProjectPaginated(ctx context.Context, projectID uuid.UUID, filter IssuesFilter) (*PaginatedRecords, error)
	SetComponentID(ctx context.Context, recordID uuid.UUID, componentID *uuid.UUID) error
	ClearComponentByComponentID(ctx context.Context, componentID uuid.UUID) error
}

type ProjectActivityLogRepository interface {
	Append(ctx context.Context, log *entity.ProjectActivityLog) error
	ListByProject(ctx context.Context, projectID uuid.UUID, limit int) ([]*entity.ProjectActivityLog, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.ProjectActivityLog, error)
	ListByProjectPaginated(ctx context.Context, projectID uuid.UUID, filter ActivityLogFilter) (*PaginatedActivityLogs, error)
}

type ProjectMemberRepository interface {
	Add(ctx context.Context, member *entity.ProjectMember) error
	Remove(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.ProjectMember, error)
	ListProjectsByUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	IsMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (bool, error)
	GetRole(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (entity.ProjectMemberRole, error)
	GetMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*entity.ProjectMember, error)
}

type UserFilter struct {
	Role     *entity.Role
	IsActive *bool
}

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	List(ctx context.Context, filter UserFilter) ([]*entity.User, error)
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *entity.RefreshToken) error
	FindByTokenHash(ctx context.Context, hash string) (*entity.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type TicketFilter struct {
	Status     *entity.TicketStatus
	Type       *entity.TicketType
	Priority   *entity.Priority
	Category   *string
	AssignedTo *uuid.UUID
	CreatedBy  *uuid.UUID
	DateFrom   *time.Time
	DateTo     *time.Time
	Search     *string
	Page       int
	PageSize   int
}

type PaginatedTickets struct {
	Tickets  []*entity.Ticket `json:"tickets"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *entity.Ticket) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Ticket, error)
	List(ctx context.Context, filter TicketFilter) (*PaginatedTickets, error)
	Update(ctx context.Context, ticket *entity.Ticket) error
}

type ApprovalRepository interface {
	SaveConfig(ctx context.Context, config *entity.ApprovalConfig) error
	FindConfigsByTicketType(ctx context.Context, ticketType entity.TicketType) ([]*entity.ApprovalConfig, error)
	SaveDecision(ctx context.Context, approval *entity.Approval) error
	FindDecisionsByTicketID(ctx context.Context, ticketID uuid.UUID) ([]*entity.Approval, error)
	FindPendingLevel(ctx context.Context, ticketID uuid.UUID) (int, error)
}

type ActivityLogRepository interface {
	Append(ctx context.Context, log *entity.ActivityLog) error
	FindByTicketID(ctx context.Context, ticketID uuid.UUID) ([]*entity.ActivityLog, error)
}

type NotificationRepository interface {
	Create(ctx context.Context, notif *entity.Notification) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error)
	MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type WebhookRepository interface {
	Create(ctx context.Context, config *entity.WebhookConfig) error
	FindAll(ctx context.Context) ([]*entity.WebhookConfig, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.WebhookConfig, error)
	Update(ctx context.Context, config *entity.WebhookConfig) error
	Delete(ctx context.Context, id uuid.UUID) error
	SaveLog(ctx context.Context, log *entity.WebhookLog) error
}

type DepartmentRepository interface {
	Create(ctx context.Context, dept *entity.Department) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Department, error)
	Update(ctx context.Context, dept *entity.Department) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*entity.Department, error)
	HasDivisions(ctx context.Context, id uuid.UUID) (bool, error)
}

type DivisionFilter struct {
	DepartmentID *uuid.UUID
}

type DivisionRepository interface {
	Create(ctx context.Context, div *entity.Division) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Division, error)
	Update(ctx context.Context, div *entity.Division) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter DivisionFilter) ([]*entity.Division, error)
	HasTeamsOrUsers(ctx context.Context, id uuid.UUID) (bool, error)
}

type TeamFilter struct {
	DivisionID *uuid.UUID
}

type TeamRepository interface {
	Create(ctx context.Context, team *entity.Team) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Team, error)
	Update(ctx context.Context, team *entity.Team) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter TeamFilter) ([]*entity.Team, error)
	HasUsers(ctx context.Context, id uuid.UUID) (bool, error)
}


// --- Jira-like Features ---

type IssueTypeRepository interface {
	Create(ctx context.Context, issueType *entity.IssueType) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.IssueType, error)
	List(ctx context.Context) ([]*entity.IssueType, error)
}

type IssueTypeSchemeRepository interface {
	Create(ctx context.Context, scheme *entity.IssueTypeScheme) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.IssueTypeScheme, error)
	FindByProject(ctx context.Context, projectID uuid.UUID) (*entity.IssueTypeScheme, error)
	Update(ctx context.Context, scheme *entity.IssueTypeScheme) error
	ListItems(ctx context.Context, schemeID uuid.UUID) ([]*entity.IssueTypeSchemeItem, error)
	AddItem(ctx context.Context, item *entity.IssueTypeSchemeItem) error
	RemoveItem(ctx context.Context, schemeID uuid.UUID, issueTypeID uuid.UUID) error
}

type CustomFieldRepository interface {
	Create(ctx context.Context, field *entity.CustomField) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.CustomField, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.CustomField, error)
	Update(ctx context.Context, field *entity.CustomField) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CustomFieldOptionRepository interface {
	Create(ctx context.Context, option *entity.CustomFieldOption) error
	ListByField(ctx context.Context, fieldID uuid.UUID) ([]*entity.CustomFieldOption, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CustomFieldValueRepository interface {
	Set(ctx context.Context, value *entity.CustomFieldValue) error
	Update(ctx context.Context, value *entity.CustomFieldValue) error
	GetByRecord(ctx context.Context, recordID uuid.UUID) ([]*entity.CustomFieldValue, error)
	GetByRecordAndField(ctx context.Context, recordID uuid.UUID, fieldID uuid.UUID) (*entity.CustomFieldValue, error)
	DeleteByRecord(ctx context.Context, recordID uuid.UUID) error
}

type WorkflowRepository interface {
	Create(ctx context.Context, workflow *entity.Workflow) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Workflow, error)
	FindByProject(ctx context.Context, projectID uuid.UUID) (*entity.Workflow, error)
	Update(ctx context.Context, workflow *entity.Workflow) error
}

type WorkflowStatusRepository interface {
	Create(ctx context.Context, status *entity.WorkflowStatus) error
	ListByWorkflow(ctx context.Context, workflowID uuid.UUID) ([]*entity.WorkflowStatus, error)
	FindByName(ctx context.Context, workflowID uuid.UUID, statusName string) (*entity.WorkflowStatus, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.WorkflowStatus, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type WorkflowTransitionRepository interface {
	Create(ctx context.Context, transition *entity.WorkflowTransition) error
	ListByWorkflow(ctx context.Context, workflowID uuid.UUID) ([]*entity.WorkflowTransition, error)
	ListFromStatus(ctx context.Context, fromStatusID uuid.UUID) ([]*entity.WorkflowTransition, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SprintRepository interface {
	Create(ctx context.Context, sprint *entity.Sprint) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Sprint, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.Sprint, error)
	ListActive(ctx context.Context, projectID uuid.UUID) ([]*entity.Sprint, error)
	Update(ctx context.Context, sprint *entity.Sprint) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SprintRecordRepository interface {
	Create(ctx context.Context, sr *entity.SprintRecord) error
	ListBySprint(ctx context.Context, sprintID uuid.UUID) ([]*entity.SprintRecord, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.SprintRecord, error)
	FindByRecord(ctx context.Context, recordID uuid.UUID) (*entity.SprintRecord, error)
	Update(ctx context.Context, sr *entity.SprintRecord) error
	Delete(ctx context.Context, id uuid.UUID) error
	BulkAssign(ctx context.Context, sprintID uuid.UUID, recordIDs []uuid.UUID) error
}

type CommentRepository interface {
	Create(ctx context.Context, comment *entity.Comment) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error)
	ListByRecord(ctx context.Context, recordID uuid.UUID) ([]*entity.Comment, error)
	Update(ctx context.Context, comment *entity.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CommentMentionRepository interface {
	Create(ctx context.Context, mention *entity.CommentMention) error
	ListByComment(ctx context.Context, commentID uuid.UUID) ([]*entity.CommentMention, error)
	DeleteByComment(ctx context.Context, commentID uuid.UUID) error
}

type AttachmentRepository interface {
	Create(ctx context.Context, attachment *entity.Attachment) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Attachment, error)
	ListByRecord(ctx context.Context, recordID uuid.UUID) ([]*entity.Attachment, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByRecord(ctx context.Context, recordID uuid.UUID) error
}

type LabelRepository interface {
	Create(ctx context.Context, label *entity.Label) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Label, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.Label, error)
	Update(ctx context.Context, label *entity.Label) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type RecordLabelRepository interface {
	Create(ctx context.Context, rl *entity.RecordLabel) error
	ListByRecord(ctx context.Context, recordID uuid.UUID) ([]*entity.RecordLabel, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByRecord(ctx context.Context, recordID uuid.UUID) error
}

// --- Project Board Features ---

type IssuesFilter struct {
	Search     *string
	StatusID   *uuid.UUID
	AssigneeID *uuid.UUID
	IssueType  *uuid.UUID
	LabelID    *uuid.UUID
	Page       int
	PageSize   int
}

type PaginatedRecords struct {
	Records  []*entity.ProjectRecord `json:"records"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
}

type ActivityLogFilter struct {
	ActionType *string
	Page       int
	PageSize   int
}

type PaginatedActivityLogs struct {
	Logs     []*entity.ProjectActivityLog `json:"logs"`
	Total    int64                        `json:"total"`
	Page     int                          `json:"page"`
	PageSize int                          `json:"page_size"`
}

type ReleaseRepository interface {
	Create(ctx context.Context, release *entity.Release) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Release, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.Release, error)
	Update(ctx context.Context, release *entity.Release) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ReleaseRecordRepository interface {
	Create(ctx context.Context, rr *entity.ReleaseRecord) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByRelease(ctx context.Context, releaseID uuid.UUID) error
	ListByRelease(ctx context.Context, releaseID uuid.UUID) ([]*entity.ReleaseRecord, error)
	FindByReleaseAndRecord(ctx context.Context, releaseID, recordID uuid.UUID) (*entity.ReleaseRecord, error)
	CountByRelease(ctx context.Context, releaseID uuid.UUID) (total int, completed int, err error)
}

type ComponentRepository interface {
	Create(ctx context.Context, component *entity.Component) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Component, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.Component, error)
	Update(ctx context.Context, component *entity.Component) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountRecords(ctx context.Context, componentID uuid.UUID) (int, error)
}
