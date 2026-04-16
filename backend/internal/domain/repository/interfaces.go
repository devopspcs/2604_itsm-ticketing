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
}

type ProjectActivityLogRepository interface {
	Append(ctx context.Context, log *entity.ProjectActivityLog) error
	ListByProject(ctx context.Context, projectID uuid.UUID, limit int) ([]*entity.ProjectActivityLog, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.ProjectActivityLog, error)
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
