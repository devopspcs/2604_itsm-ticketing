package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	notifinfra "github.com/org/itsm/internal/infrastructure/notification"
	"github.com/org/itsm/pkg/apperror"
)

type assignmentUseCase struct {
	ticketRepo       repository.TicketRepository
	userRepo         repository.UserRepository
	teamRepo         repository.TeamRepository
	activityRepo     repository.ActivityLogRepository
	notificationRepo repository.NotificationRepository
	webhookUC        domainUC.WebhookUseCase
	emailSvc         notifinfra.EmailSender
}

func NewAssignmentUseCase(
	ticketRepo repository.TicketRepository,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
	activityRepo repository.ActivityLogRepository,
	notificationRepo repository.NotificationRepository,
	webhookUC domainUC.WebhookUseCase,
	emailSvc notifinfra.EmailSender,
) domainUC.AssignmentUseCase {
	return &assignmentUseCase{
		ticketRepo:       ticketRepo,
		userRepo:         userRepo,
		teamRepo:         teamRepo,
		activityRepo:     activityRepo,
		notificationRepo: notificationRepo,
		webhookUC:        webhookUC,
		emailSvc:         emailSvc,
	}
}

func (uc *assignmentUseCase) AssignTicket(ctx context.Context, ticketID uuid.UUID, assigneeID uuid.UUID, requester domainUC.UserClaims) error {
	assignee, err := uc.userRepo.FindByID(ctx, assigneeID)
	if err != nil {
		return apperror.ErrNotFound
	}
	if !assignee.IsActive {
		return apperror.ErrValidation.WithDetails(map[string]interface{}{"assignee": "user is inactive"})
	}

	// Position-based ACL check
	if err := uc.checkAssignmentACL(ctx, requester, assignee); err != nil {
		return err
	}

	ticket, err := uc.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return err
	}

	oldAssignee := ticket.AssignedTo
	ticket.AssignedTo = &assigneeID
	ticket.Status = entity.StatusInProgress
	ticket.UpdatedAt = time.Now().UTC()

	if err := uc.ticketRepo.Update(ctx, ticket); err != nil {
		return err
	}

	action := entity.ActionAssigned
	if oldAssignee != nil {
		action = entity.ActionReassigned
	}

	_ = uc.activityRepo.Append(ctx, &entity.ActivityLog{
		ID:        uuid.New(),
		TicketID:  ticket.ID,
		ActorID:   requester.UserID,
		Action:    action,
		NewValue:  strPtr(assigneeID.String()),
		CreatedAt: time.Now().UTC(),
	})

	_ = uc.notificationRepo.Create(ctx, &entity.Notification{
		ID:        uuid.New(),
		UserID:    assigneeID,
		TicketID:  ticket.ID,
		Message:   "You have been assigned to ticket: " + ticket.Title,
		IsRead:    false,
		CreatedAt: time.Now().UTC(),
	})

	// Dispatch webhook
	if uc.webhookUC != nil {
		go uc.webhookUC.Dispatch(context.Background(), entity.EventTicketAssigned, ticket)
	}

	// Send email notification
	if uc.emailSvc != nil && uc.emailSvc.IsConfigured() {
		reqUser, _ := uc.userRepo.FindByID(ctx, requester.UserID)
		assignerName := "System"
		if reqUser != nil {
			assignerName = reqUser.FullName
		}
		uc.emailSvc.SendTicketAssigned(
			assignee.Email, assignee.FullName,
			ticket.Title, ticket.ID.String(), ticket.TicketNumber,
			string(ticket.Type), string(ticket.Priority), ticket.Category,
			assignerName,
		)
	}

	return nil
}

// checkAssignmentACL validates position-based assignment permissions.
func (uc *assignmentUseCase) checkAssignmentACL(ctx context.Context, requester domainUC.UserClaims, assignee *entity.User) error {
	// Admin, Agent, and Approver can assign to anyone
	if requester.Role == entity.RoleAdmin || requester.Role == entity.RoleAgent || requester.Role == entity.RoleApprover {
		return nil
	}

	// Get requester's full user data for org info
	reqUser, err := uc.userRepo.FindByID(ctx, requester.UserID)
	if err != nil {
		// If we can't find the user, fall back to allowing (backward compat)
		return nil
	}

	// If requester has no org position, allow (backward compatibility)
	if reqUser.Position == nil {
		return nil
	}

	switch *reqUser.Position {
	case entity.PositionLeader:
		// Leader can only assign to Staff in the same team
		if reqUser.TeamID == nil || assignee.TeamID == nil || *reqUser.TeamID != *assignee.TeamID {
			return apperror.ErrForbidden
		}
		if assignee.Position == nil || *assignee.Position != entity.PositionStaff {
			return apperror.ErrForbidden
		}
		return nil

	case entity.PositionManager:
		// Manager can assign to any member of the same team
		if reqUser.TeamID == nil || assignee.TeamID == nil || *reqUser.TeamID != *assignee.TeamID {
			return apperror.ErrForbidden
		}
		return nil

	case entity.PositionDivisionManager:
		// Division Manager can assign to any member of the same division
		if reqUser.DivisionID == nil || assignee.DivisionID == nil || *reqUser.DivisionID != *assignee.DivisionID {
			return apperror.ErrForbidden
		}
		return nil

	case entity.PositionStaff:
		// Staff cannot assign
		return apperror.ErrForbidden

	default:
		return nil
	}
}

// AssignTicketToTeam assigns a ticket to a team and sends email notification to the team email.
func (uc *assignmentUseCase) AssignTicketToTeam(ctx context.Context, ticketID uuid.UUID, teamID uuid.UUID, requester domainUC.UserClaims) error {
	// Verify team exists
	team, err := uc.teamRepo.FindByID(ctx, teamID)
	if err != nil {
		return apperror.ErrNotFound.WithDetails(map[string]interface{}{"team": "team not found"})
	}

	// Check ACL: admin, agent, approver can assign to team; staff cannot
	if requester.Role != entity.RoleAdmin && requester.Role != entity.RoleAgent && requester.Role != entity.RoleApprover {
		reqUser, err := uc.userRepo.FindByID(ctx, requester.UserID)
		if err != nil {
			return nil // backward compat
		}
		if reqUser.Position != nil && *reqUser.Position == entity.PositionStaff {
			return apperror.ErrForbidden
		}
	}

	// Get ticket
	ticket, err := uc.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return err
	}

	// Set team assignment
	ticket.AssignedTeamID = &teamID
	ticket.UpdatedAt = time.Now().UTC()

	// If ticket is still open, move to in_progress
	if ticket.Status == entity.StatusOpen {
		ticket.Status = entity.StatusInProgress
	}

	if err := uc.ticketRepo.Update(ctx, ticket); err != nil {
		return err
	}

	// Log activity
	_ = uc.activityRepo.Append(ctx, &entity.ActivityLog{
		ID:        uuid.New(),
		TicketID:  ticket.ID,
		ActorID:   requester.UserID,
		Action:    entity.ActionAssignedToTeam,
		NewValue:  strPtr(teamID.String()),
		CreatedAt: time.Now().UTC(),
	})

	// Notify all team members via in-app notification
	allUsers, _ := uc.userRepo.List(ctx, repository.UserFilter{})
	for _, u := range allUsers {
		if u.TeamID != nil && *u.TeamID == teamID {
			_ = uc.notificationRepo.Create(ctx, &entity.Notification{
				ID:        uuid.New(),
				UserID:    u.ID,
				TicketID:  ticket.ID,
				Message:   "Ticket assigned to your team: " + ticket.Title,
				IsRead:    false,
				CreatedAt: time.Now().UTC(),
			})
		}
	}

	// Dispatch webhook
	if uc.webhookUC != nil {
		go uc.webhookUC.Dispatch(context.Background(), entity.EventTicketAssigned, ticket)
	}

	// Send email to team email
	if uc.emailSvc != nil && uc.emailSvc.IsConfigured() && team.Email != nil && *team.Email != "" {
		reqUser, _ := uc.userRepo.FindByID(ctx, requester.UserID)
		assignerName := "System"
		if reqUser != nil {
			assignerName = reqUser.FullName
		}
		uc.emailSvc.SendTicketAssignedToTeam(
			*team.Email, team.Name,
			ticket.Title, ticket.ID.String(), ticket.TicketNumber,
			string(ticket.Type), string(ticket.Priority), ticket.Category,
			assignerName,
		)
	}

	return nil
}
