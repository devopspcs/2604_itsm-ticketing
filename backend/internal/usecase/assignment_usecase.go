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
	activityRepo     repository.ActivityLogRepository
	notificationRepo repository.NotificationRepository
	webhookUC        domainUC.WebhookUseCase
	emailSvc         notifinfra.EmailSender
}

func NewAssignmentUseCase(
	ticketRepo repository.TicketRepository,
	userRepo repository.UserRepository,
	activityRepo repository.ActivityLogRepository,
	notificationRepo repository.NotificationRepository,
	webhookUC domainUC.WebhookUseCase,
	emailSvc notifinfra.EmailSender,
) domainUC.AssignmentUseCase {
	return &assignmentUseCase{
		ticketRepo:       ticketRepo,
		userRepo:         userRepo,
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
			ticket.Title, ticket.ID.String(),
			string(ticket.Type), string(ticket.Priority), ticket.Category,
			assignerName,
		)
	}

	return nil
}

// checkAssignmentACL validates position-based assignment permissions.
func (uc *assignmentUseCase) checkAssignmentACL(ctx context.Context, requester domainUC.UserClaims, assignee *entity.User) error {
	// Admin can assign to anyone
	if requester.Role == entity.RoleAdmin {
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
