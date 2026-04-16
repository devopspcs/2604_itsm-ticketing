package usecase

import (
	"context"

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
