package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type ticketUseCase struct {
	ticketRepo       repository.TicketRepository
	activityRepo     repository.ActivityLogRepository
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
	webhookUC        domainUC.WebhookUseCase
}

func NewTicketUseCase(
	ticketRepo repository.TicketRepository,
	activityRepo repository.ActivityLogRepository,
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
	webhookUC domainUC.WebhookUseCase,
) domainUC.TicketUseCase {
	return &ticketUseCase{
		ticketRepo:       ticketRepo,
		activityRepo:     activityRepo,
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		webhookUC:        webhookUC,
	}
}

func (uc *ticketUseCase) CreateTicket(ctx context.Context, req domainUC.CreateTicketRequest, requester domainUC.UserClaims) (*entity.Ticket, error) {
	now := time.Now().UTC()
	ticket := &entity.Ticket{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Category:    req.Category,
		Priority:    req.Priority,
		Status:      entity.StatusOpen,
		CreatedBy:   requester.UserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.ticketRepo.Create(ctx, ticket); err != nil {
		return nil, err
	}

	// Log activity
	_ = uc.activityRepo.Append(ctx, &entity.ActivityLog{
		ID:        uuid.New(),
		TicketID:  ticket.ID,
		ActorID:   requester.UserID,
		Action:    entity.ActionTicketCreated,
		CreatedAt: now,
	})

	// Notify admins/approvers
	go uc.notifyAdminsAndApprovers(ctx, ticket)

	// Dispatch webhook
	if uc.webhookUC != nil {
		go uc.webhookUC.Dispatch(context.Background(), entity.EventTicketCreated, ticket)
	}

	return ticket, nil
}

func (uc *ticketUseCase) GetTicket(ctx context.Context, id uuid.UUID, requester domainUC.UserClaims) (*entity.Ticket, error) {
	ticket, err := uc.ticketRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if requester.Role == entity.RoleUser {
		isCreator := ticket.CreatedBy == requester.UserID
		isAssignee := ticket.AssignedTo != nil && *ticket.AssignedTo == requester.UserID
		if !isCreator && !isAssignee {
			return nil, apperror.ErrForbidden
		}
	}
	return ticket, nil
}

func (uc *ticketUseCase) ListTickets(ctx context.Context, filter repository.TicketFilter, requester domainUC.UserClaims) (*repository.PaginatedTickets, error) {
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	// Admin sees everything
	if requester.Role == entity.RoleAdmin {
		return uc.ticketRepo.List(ctx, filter)
	}

	// Check if user has org assignment for position-based visibility
	user, err := uc.userRepo.FindByID(ctx, requester.UserID)
	if err == nil && user.Position != nil {
		return uc.listTicketsByPosition(ctx, filter, user)
	}

	// Fallback to role-based visibility
	if requester.Role != entity.RoleUser {
		// Approver sees all
		return uc.ticketRepo.List(ctx, filter)
	}

	// Regular user: own tickets + assigned tickets
	return uc.listOwnAndAssignedTickets(ctx, filter, requester)
}

// listTicketsByPosition applies position-based visibility filtering.
func (uc *ticketUseCase) listTicketsByPosition(ctx context.Context, filter repository.TicketFilter, user *entity.User) (*repository.PaginatedTickets, error) {
	switch *user.Position {
	case entity.PositionStaff:
		// Staff: only own tickets
		filter.CreatedBy = &user.ID
		return uc.ticketRepo.List(ctx, filter)

	case entity.PositionLeader, entity.PositionManager:
		// Leader/Manager: tickets from all team members
		if user.TeamID == nil {
			filter.CreatedBy = &user.ID
			return uc.ticketRepo.List(ctx, filter)
		}
		return uc.listTicketsByTeam(ctx, filter, *user.TeamID)

	case entity.PositionDivisionManager:
		// Division Manager: tickets from all division members
		if user.DivisionID == nil {
			filter.CreatedBy = &user.ID
			return uc.ticketRepo.List(ctx, filter)
		}
		return uc.listTicketsByDivision(ctx, filter, *user.DivisionID)

	default:
		filter.CreatedBy = &user.ID
		return uc.ticketRepo.List(ctx, filter)
	}
}

func (uc *ticketUseCase) listTicketsByTeam(ctx context.Context, filter repository.TicketFilter, teamID uuid.UUID) (*repository.PaginatedTickets, error) {
	// Get all users in the team
	allUsers, err := uc.userRepo.List(ctx, repository.UserFilter{})
	if err != nil {
		return nil, err
	}

	var teamUserIDs []uuid.UUID
	for _, u := range allUsers {
		if u.TeamID != nil && *u.TeamID == teamID {
			teamUserIDs = append(teamUserIDs, u.ID)
		}
	}

	return uc.listTicketsByUserIDs(ctx, filter, teamUserIDs)
}

func (uc *ticketUseCase) listTicketsByDivision(ctx context.Context, filter repository.TicketFilter, divisionID uuid.UUID) (*repository.PaginatedTickets, error) {
	allUsers, err := uc.userRepo.List(ctx, repository.UserFilter{})
	if err != nil {
		return nil, err
	}

	var divUserIDs []uuid.UUID
	for _, u := range allUsers {
		if u.DivisionID != nil && *u.DivisionID == divisionID {
			divUserIDs = append(divUserIDs, u.ID)
		}
	}

	return uc.listTicketsByUserIDs(ctx, filter, divUserIDs)
}

func (uc *ticketUseCase) listTicketsByUserIDs(ctx context.Context, filter repository.TicketFilter, userIDs []uuid.UUID) (*repository.PaginatedTickets, error) {
	if len(userIDs) == 0 {
		return &repository.PaginatedTickets{
			Tickets:  []*entity.Ticket{},
			Total:    0,
			Page:     filter.Page,
			PageSize: filter.PageSize,
		}, nil
	}

	// Fetch tickets for each user and merge
	seen := map[uuid.UUID]bool{}
	var merged []*entity.Ticket
	for _, uid := range userIDs {
		uid := uid
		f := filter
		f.CreatedBy = &uid
		f.PageSize = 1000
		f.Page = 1
		result, err := uc.ticketRepo.List(ctx, f)
		if err != nil {
			return nil, err
		}
		for _, t := range result.Tickets {
			if !seen[t.ID] {
				seen[t.ID] = true
				merged = append(merged, t)
			}
		}
	}

	// Apply pagination
	total := int64(len(merged))
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(merged) {
		start = len(merged)
	}
	if end > len(merged) {
		end = len(merged)
	}

	return &repository.PaginatedTickets{
		Tickets:  merged[start:end],
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (uc *ticketUseCase) listOwnAndAssignedTickets(ctx context.Context, filter repository.TicketFilter, requester domainUC.UserClaims) (*repository.PaginatedTickets, error) {
	createdFilter := filter
	createdFilter.CreatedBy = &requester.UserID
	createdFilter.AssignedTo = nil
	createdFilter.PageSize = 1000
	createdFilter.Page = 1

	assignedFilter := filter
	assignedFilter.AssignedTo = &requester.UserID
	assignedFilter.CreatedBy = nil
	assignedFilter.PageSize = 1000
	assignedFilter.Page = 1

	createdResult, err := uc.ticketRepo.List(ctx, createdFilter)
	if err != nil {
		return nil, err
	}
	assignedResult, err := uc.ticketRepo.List(ctx, assignedFilter)
	if err != nil {
		return nil, err
	}

	seen := map[uuid.UUID]bool{}
	var merged []*entity.Ticket
	for _, t := range append(createdResult.Tickets, assignedResult.Tickets...) {
		if !seen[t.ID] {
			seen[t.ID] = true
			merged = append(merged, t)
		}
	}

	total := int64(len(merged))
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(merged) {
		start = len(merged)
	}
	if end > len(merged) {
		end = len(merged)
	}

	return &repository.PaginatedTickets{
		Tickets:  merged[start:end],
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (uc *ticketUseCase) UpdateTicket(ctx context.Context, id uuid.UUID, req domainUC.UpdateTicketRequest, requester domainUC.UserClaims) (*entity.Ticket, error) {
	ticket, err := uc.ticketRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if requester.Role == entity.RoleUser {
		isCreator := ticket.CreatedBy == requester.UserID
		isAssignee := ticket.AssignedTo != nil && *ticket.AssignedTo == requester.UserID
		if !isCreator && !isAssignee {
			return nil, apperror.ErrForbidden
		}
	}

	if req.Title != nil {
		ticket.Title = *req.Title
	}
	if req.Description != nil {
		ticket.Description = *req.Description
	}
	if req.Category != nil {
		ticket.Category = *req.Category
	}
	if req.Priority != nil {
		ticket.Priority = *req.Priority
	}
	if req.Status != nil {
		oldStatus := ticket.Status
		ticket.Status = *req.Status

		now := time.Now().UTC()
		if *req.Status == entity.StatusDone {
			ticket.ResolvedAt = &now
		} else if oldStatus == entity.StatusDone {
			ticket.ResolvedAt = nil
		}
	}
	ticket.UpdatedAt = time.Now().UTC()

	if err := uc.ticketRepo.Update(ctx, ticket); err != nil {
		return nil, err
	}

	_ = uc.activityRepo.Append(ctx, &entity.ActivityLog{
		ID:        uuid.New(),
		TicketID:  ticket.ID,
		ActorID:   requester.UserID,
		Action:    entity.ActionFieldUpdated,
		CreatedAt: time.Now().UTC(),
	})

	return ticket, nil
}

func (uc *ticketUseCase) notifyAdminsAndApprovers(ctx context.Context, ticket *entity.Ticket) {
	users, err := uc.userRepo.List(ctx, repository.UserFilter{})
	if err != nil {
		return
	}
	for _, u := range users {
		if u.Role == entity.RoleAdmin || u.Role == entity.RoleApprover {
			_ = uc.notificationRepo.Create(ctx, &entity.Notification{
				ID:        uuid.New(),
				UserID:    u.ID,
				TicketID:  ticket.ID,
				Message:   "New ticket created: " + ticket.Title,
				IsRead:    false,
				CreatedAt: time.Now().UTC(),
			})
		}
	}
}
