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

type approvalUseCase struct {
	ticketRepo       repository.TicketRepository
	approvalRepo     repository.ApprovalRepository
	activityRepo     repository.ActivityLogRepository
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
	webhookUC        domainUC.WebhookUseCase
}

func NewApprovalUseCase(
	ticketRepo repository.TicketRepository,
	approvalRepo repository.ApprovalRepository,
	activityRepo repository.ActivityLogRepository,
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
	webhookUC domainUC.WebhookUseCase,
) domainUC.ApprovalUseCase {
	return &approvalUseCase{
		ticketRepo:       ticketRepo,
		approvalRepo:     approvalRepo,
		activityRepo:     activityRepo,
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		webhookUC:        webhookUC,
	}
}

func (uc *approvalUseCase) SubmitForApproval(ctx context.Context, ticketID uuid.UUID, requester domainUC.UserClaims) error {
	ticket, err := uc.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return err
	}
	if requester.Role == entity.RoleUser && ticket.CreatedBy != requester.UserID {
		return apperror.ErrForbidden
	}

	ticket.Status = entity.StatusPendingApproval
	ticket.UpdatedAt = time.Now().UTC()
	if err := uc.ticketRepo.Update(ctx, ticket); err != nil {
		return err
	}

	_ = uc.activityRepo.Append(ctx, &entity.ActivityLog{
		ID:        uuid.New(),
		TicketID:  ticket.ID,
		ActorID:   requester.UserID,
		Action:    entity.ActionApprovalRequested,
		CreatedAt: time.Now().UTC(),
	})

	// Dispatch webhook
	if uc.webhookUC != nil {
		go uc.webhookUC.Dispatch(context.Background(), entity.EventApprovalRequested, ticket)
	}

	return nil
}

func (uc *approvalUseCase) Decide(ctx context.Context, req domainUC.ApprovalDecisionRequest, approver domainUC.UserClaims) error {
	// Check role-based permission first
	if approver.Role != entity.RoleApprover && approver.Role != entity.RoleAdmin {
		// Check if user has org-based approval permission
		approverUser, err := uc.userRepo.FindByID(ctx, approver.UserID)
		if err != nil || approverUser.Position == nil {
			return apperror.ErrForbidden
		}
		pos := *approverUser.Position
		if pos != entity.PositionManager && pos != entity.PositionDivisionManager {
			return apperror.ErrForbidden
		}
	}

	ticket, err := uc.ticketRepo.FindByID(ctx, req.TicketID)
	if err != nil {
		return err
	}

	// Org-based scope validation: if approver has org position, validate scope
	approverUser, _ := uc.userRepo.FindByID(ctx, approver.UserID)
	if approverUser != nil && approverUser.Position != nil && approver.Role != entity.RoleAdmin {
		creator, _ := uc.userRepo.FindByID(ctx, ticket.CreatedBy)
		if creator != nil {
			if err := uc.validateApprovalScope(approverUser, creator); err != nil {
				return err
			}
		}
	}

	// Get configs to know total levels
	configs, err := uc.approvalRepo.FindConfigsByTicketType(ctx, ticket.Type)
	if err != nil {
		return err
	}

	// Get current pending level
	pendingLevel, err := uc.approvalRepo.FindPendingLevel(ctx, ticket.ID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	decision := req.Decision
	approval := &entity.Approval{
		ID:         uuid.New(),
		TicketID:   ticket.ID,
		ApproverID: approver.UserID,
		Level:      pendingLevel,
		Decision:   &decision,
		Comment:    req.Comment,
		DecidedAt:  &now,
	}

	if err := uc.approvalRepo.SaveDecision(ctx, approval); err != nil {
		return err
	}

	_ = uc.activityRepo.Append(ctx, &entity.ActivityLog{
		ID:        uuid.New(),
		TicketID:  ticket.ID,
		ActorID:   approver.UserID,
		Action:    entity.ActionApprovalDecided,
		NewValue:  strPtr(string(decision)),
		CreatedAt: now,
	})

	if decision == entity.DecisionRejected {
		ticket.Status = entity.StatusRejected
	} else {
		// Check if all levels approved
		maxLevel := 0
		for _, c := range configs {
			if c.Level > maxLevel {
				maxLevel = c.Level
			}
		}
		if pendingLevel >= maxLevel || len(configs) == 0 {
			ticket.Status = entity.StatusApproved
		}
		// else stays pending_approval for next level
	}

	ticket.UpdatedAt = now
	if err := uc.ticketRepo.Update(ctx, ticket); err != nil {
		return err
	}

	// Notify ticket creator
	_ = uc.notificationRepo.Create(ctx, &entity.Notification{
		ID:        uuid.New(),
		UserID:    ticket.CreatedBy,
		TicketID:  ticket.ID,
		Message:   "Your ticket approval decision: " + string(decision),
		IsRead:    false,
		CreatedAt: now,
	})

	// Dispatch webhook
	if uc.webhookUC != nil {
		go uc.webhookUC.Dispatch(context.Background(), entity.EventApprovalDecided, ticket)
	}

	return nil
}

func (uc *approvalUseCase) GetApprovalHistory(ctx context.Context, ticketID uuid.UUID) ([]*entity.Approval, error) {
	return uc.approvalRepo.FindDecisionsByTicketID(ctx, ticketID)
}

// validateApprovalScope checks that the approver has org-based authority over the ticket creator.
func (uc *approvalUseCase) validateApprovalScope(approver, creator *entity.User) error {
	if approver.Position == nil {
		return apperror.ErrForbidden
	}
	switch *approver.Position {
	case entity.PositionManager:
		// Manager can approve tickets from users in the same team
		if approver.TeamID == nil || creator.TeamID == nil || *approver.TeamID != *creator.TeamID {
			return apperror.ErrForbidden
		}
	case entity.PositionDivisionManager:
		// Division Manager can approve tickets from users in the same division
		if approver.DivisionID == nil || creator.DivisionID == nil || *approver.DivisionID != *creator.DivisionID {
			return apperror.ErrForbidden
		}
	default:
		return apperror.ErrForbidden
	}
	return nil
}

// ResolveApprovers finds org-based approvers for a ticket creator.
// Returns Manager (level 1) and Division Manager (level 2) if found.
func (uc *approvalUseCase) ResolveApprovers(ctx context.Context, creatorID uuid.UUID) ([]struct {
	UserID uuid.UUID
	Level  int
}, error) {
	creator, err := uc.userRepo.FindByID(ctx, creatorID)
	if err != nil || creator.TeamID == nil {
		return nil, nil // fallback to existing approval
	}

	var approvers []struct {
		UserID uuid.UUID
		Level  int
	}

	// Find Manager in the same team
	users, err := uc.userRepo.List(ctx, repository.UserFilter{})
	if err != nil {
		return nil, nil
	}

	for _, u := range users {
		if u.Position == nil || !u.IsActive {
			continue
		}
		if *u.Position == entity.PositionManager && u.TeamID != nil && *u.TeamID == *creator.TeamID && u.ID != creatorID {
			approvers = append(approvers, struct {
				UserID uuid.UUID
				Level  int
			}{UserID: u.ID, Level: 1})
			break
		}
	}

	// Find Division Manager in the same division
	if creator.DivisionID != nil {
		for _, u := range users {
			if u.Position == nil || !u.IsActive {
				continue
			}
			if *u.Position == entity.PositionDivisionManager && u.DivisionID != nil && *u.DivisionID == *creator.DivisionID && u.ID != creatorID {
				approvers = append(approvers, struct {
					UserID uuid.UUID
					Level  int
				}{UserID: u.ID, Level: 2})
				break
			}
		}
	}

	return approvers, nil
}

func strPtr(s string) *string { return &s }
