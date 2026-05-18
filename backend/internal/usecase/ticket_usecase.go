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
	teamRepo         repository.TeamRepository
	webhookUC        domainUC.WebhookUseCase
}

func NewTicketUseCase(
	ticketRepo repository.TicketRepository,
	activityRepo repository.ActivityLogRepository,
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
	webhookUC domainUC.WebhookUseCase,
) domainUC.TicketUseCase {
	return &ticketUseCase{
		ticketRepo:       ticketRepo,
		activityRepo:     activityRepo,
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		teamRepo:         teamRepo,
		webhookUC:        webhookUC,
	}
}

func (uc *ticketUseCase) CreateTicket(ctx context.Context, req domainUC.CreateTicketRequest, requester domainUC.UserClaims) (*entity.Ticket, error) {
	// Role-based ticket type enforcement
	switch requester.Role {
	case entity.RoleUser:
		// Users can only create incident and request tickets
		if req.Type != entity.TypeIncident && req.Type != entity.TypeRequest {
			return nil, apperror.ErrForbidden.WithDetails(map[string]interface{}{
				"message": "Users can only create incident or request tickets",
			})
		}
	case entity.RoleAgent:
		// Agents can create all types including change_request
	case entity.RoleApprover, entity.RoleAdmin:
		// Approvers and admins can create all types
	default:
		return nil, apperror.ErrForbidden
	}

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

	// Notify agents and admins about new tickets from users
	go uc.notifyAgentsAndAdmins(ctx, ticket)

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
	// Users can only see their own tickets
	if requester.Role == entity.RoleUser {
		if ticket.CreatedBy != requester.UserID {
			return nil, apperror.ErrForbidden
		}
	}
	// Approvers can only see tickets from teams in their division or unassigned
	if requester.Role == entity.RoleApprover {
		approver, err := uc.userRepo.FindByID(ctx, requester.UserID)
		if err == nil && approver.DivisionID != nil {
			if ticket.AssignedTeamID != nil {
				// Check if ticket's team is in approver's division
				teams, _ := uc.teamRepo.List(ctx, repository.TeamFilter{DivisionID: approver.DivisionID})
				found := false
				for _, t := range teams {
					if t.ID == *ticket.AssignedTeamID {
						found = true
						break
					}
				}
				if !found {
					return nil, apperror.ErrForbidden
				}
			}
			// unassigned tickets are visible to all approvers
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

	// Agent sees all tickets (they handle them)
	if requester.Role == entity.RoleAgent {
		return uc.ticketRepo.List(ctx, filter)
	}

	// Approver sees tickets assigned to teams in their division + unassigned tickets
	if requester.Role == entity.RoleApprover {
		approver, err := uc.userRepo.FindByID(ctx, requester.UserID)
		if err != nil || approver.DivisionID == nil {
			// Approver without division: see all (backward compatible)
			return uc.ticketRepo.List(ctx, filter)
		}
		// Get all teams in the approver's division
		teams, err := uc.teamRepo.List(ctx, repository.TeamFilter{DivisionID: approver.DivisionID})
		if err != nil || len(teams) == 0 {
			return uc.ticketRepo.List(ctx, filter)
		}
		var teamIDs []uuid.UUID
		for _, t := range teams {
			teamIDs = append(teamIDs, t.ID)
		}
		filter.AssignedTeamIDsOrUnassigned = teamIDs
		return uc.ticketRepo.List(ctx, filter)
	}

	// Check if user has org assignment for position-based visibility
	user, err := uc.userRepo.FindByID(ctx, requester.UserID)
	if err == nil && user.Position != nil {
		return uc.listTicketsByPosition(ctx, filter, user)
	}

	// Regular user: only own tickets
	filter.CreatedBy = &requester.UserID
	return uc.ticketRepo.List(ctx, filter)
}

// listTicketsByPosition applies position-based visibility filtering.
func (uc *ticketUseCase) listTicketsByPosition(ctx context.Context, filter repository.TicketFilter, user *entity.User) (*repository.PaginatedTickets, error) {
	switch *user.Position {
	case entity.PositionStaff:
		// Staff: only own tickets
		filter.CreatedBy = &user.ID
		return uc.ticketRepo.List(ctx, filter)

	case entity.PositionLeader:
		// Leader: tickets from all team members
		if user.TeamID == nil {
			filter.CreatedBy = &user.ID
			return uc.ticketRepo.List(ctx, filter)
		}
		return uc.listTicketsByTeam(ctx, filter, *user.TeamID)

	case entity.PositionManager:
		// Manager: tickets from all division members
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

	// Users can only update their own tickets (title, description only while still open)
	if requester.Role == entity.RoleUser {
		if ticket.CreatedBy != requester.UserID {
			return nil, apperror.ErrForbidden
		}
		// Users cannot change status
		if req.Status != nil {
			return nil, apperror.ErrForbidden.WithDetails(map[string]interface{}{
				"message": "Users cannot change ticket status",
			})
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

func (uc *ticketUseCase) DeleteTicket(ctx context.Context, id uuid.UUID, requester domainUC.UserClaims) error {
	ticket, err := uc.ticketRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Only admin and agent can delete tickets
	if requester.Role != entity.RoleAdmin && requester.Role != entity.RoleAgent {
		return apperror.ErrForbidden
	}

	if err := uc.ticketRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Log the deletion with ticket ID in the activity log
	ticketTitle := ticket.Title
	_ = uc.activityRepo.Append(ctx, &entity.ActivityLog{
		ID:        uuid.New(),
		TicketID:  id,
		ActorID:   requester.UserID,
		Action:    entity.ActionTicketDeleted,
		OldValue:  &ticketTitle,
		CreatedAt: time.Now().UTC(),
	})

	return nil
}

func (uc *ticketUseCase) notifyAgentsAndAdmins(ctx context.Context, ticket *entity.Ticket) {
	users, err := uc.userRepo.List(ctx, repository.UserFilter{})
	if err != nil {
		return
	}
	for _, u := range users {
		if u.Role == entity.RoleAdmin || u.Role == entity.RoleAgent {
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
