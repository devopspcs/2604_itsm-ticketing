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

type workflowUseCase struct {
	workflowRepo    repository.WorkflowRepository
	statusRepo      repository.WorkflowStatusRepository
	transitionRepo  repository.WorkflowTransitionRepository
	projectRepo     repository.ProjectRepository
	recordRepo      repository.ProjectRecordRepository
	activityRepo    repository.ProjectActivityLogRepository
	memberRepo      repository.ProjectMemberRepository
}

func NewWorkflowUseCase(
	workflowRepo repository.WorkflowRepository,
	statusRepo repository.WorkflowStatusRepository,
	transitionRepo repository.WorkflowTransitionRepository,
	projectRepo repository.ProjectRepository,
	recordRepo repository.ProjectRecordRepository,
	activityRepo repository.ProjectActivityLogRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.WorkflowUseCase {
	return &workflowUseCase{
		workflowRepo:    workflowRepo,
		statusRepo:      statusRepo,
		transitionRepo:  transitionRepo,
		projectRepo:     projectRepo,
		recordRepo:      recordRepo,
		activityRepo:    activityRepo,
		memberRepo:      memberRepo,
	}
}

// GetWorkflow returns the workflow for a project
func (uc *workflowUseCase) GetWorkflow(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) (*entity.Workflow, error) {
	// Verify project exists
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is a member
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	workflow, err := uc.workflowRepo.FindByProject(ctx, projectID)
	if err != nil && err != apperror.ErrNotFound {
		return nil, err
	}

	// Auto-create default workflow if none exists
	if workflow == nil {
		now := time.Now().UTC()
		workflow = &entity.Workflow{
			ID:            uuid.New(),
			ProjectID:     projectID,
			Name:          "Default Workflow",
			InitialStatus: "Backlog",
			CreatedAt:     now,
		}
		if err := uc.workflowRepo.Create(ctx, workflow); err != nil {
			return nil, err
		}

		// Create default statuses
		defaultStatuses := []string{"Backlog", "To Do", "In Progress", "In Review", "Done"}
		for i, name := range defaultStatuses {
			status := &entity.WorkflowStatus{
				ID:          uuid.New(),
				WorkflowID:  workflow.ID,
				StatusName:  name,
				StatusOrder: i,
				CreatedAt:   now,
			}
			_ = uc.statusRepo.Create(ctx, status)
		}
	}

	return workflow, nil
}

// CreateWorkflow creates a new workflow for a project
func (uc *workflowUseCase) CreateWorkflow(ctx context.Context, projectID uuid.UUID, req domainUC.CreateWorkflowRequest, requester domainUC.UserClaims) (*entity.Workflow, error) {
	// Verify project exists
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is project owner
	member, err := uc.memberRepo.GetMember(ctx, projectID, requester.UserID)
	if err != nil || member == nil || member.Role != entity.ProjectRoleOwner {
		return nil, apperror.ErrForbidden
	}

	// Validate statuses
	if len(req.Statuses) == 0 {
		return nil, apperror.ErrValidation
	}

	// Validate initial status is in statuses list
	initialStatusFound := false
	for _, status := range req.Statuses {
		if status == req.InitialStatus {
			initialStatusFound = true
			break
		}
	}
	if !initialStatusFound {
		return nil, apperror.ErrValidation
	}

	// Create workflow
	workflow := &entity.Workflow{
		ID:             uuid.New(),
		ProjectID:      projectID,
		Name:           req.Name,
		InitialStatus:  req.InitialStatus,
		CreatedAt:      time.Now().UTC(),
	}

	if err := uc.workflowRepo.Create(ctx, workflow); err != nil {
		return nil, err
	}

	// Create statuses
	for i, statusName := range req.Statuses {
		status := &entity.WorkflowStatus{
			ID:         uuid.New(),
			WorkflowID: workflow.ID,
			StatusName: statusName,
			StatusOrder: i,
			CreatedAt:  time.Now().UTC(),
		}
		if err := uc.statusRepo.Create(ctx, status); err != nil {
			return nil, err
		}
	}

	// Log activity
	uc.logActivity(ctx, projectID, nil, requester.UserID, "workflow_created", "Workflow created: "+workflow.Name)

	return workflow, nil
}

// UpdateWorkflow updates a workflow's statuses and transitions
func (uc *workflowUseCase) UpdateWorkflow(ctx context.Context, projectID uuid.UUID, workflowID uuid.UUID, req domainUC.UpdateWorkflowRequest, requester domainUC.UserClaims) error {
	// Verify project exists
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return apperror.ErrNotFound
	}

	// Verify user is project owner
	member, err := uc.memberRepo.GetMember(ctx, projectID, requester.UserID)
	if err != nil || member == nil || member.Role != entity.ProjectRoleOwner {
		return apperror.ErrForbidden
	}

	// Get workflow
	workflow, err := uc.workflowRepo.FindByID(ctx, workflowID)
	if err != nil {
		return err
	}
	if workflow == nil || workflow.ProjectID != projectID {
		return apperror.ErrNotFound
	}

	// Update statuses if provided
	if len(req.Statuses) > 0 {
		// Get existing statuses
		existingStatuses, err := uc.statusRepo.ListByWorkflow(ctx, workflowID)
		if err != nil {
			return err
		}

		// Delete statuses not in new list
		for _, existingStatus := range existingStatuses {
			found := false
			for _, newStatusName := range req.Statuses {
				if existingStatus.StatusName == newStatusName {
					found = true
					break
				}
			}
			if !found {
				if err := uc.statusRepo.Delete(ctx, existingStatus.ID); err != nil {
					return err
				}
			}
		}

		// Add new statuses
		for i, statusName := range req.Statuses {
			found := false
			for _, existingStatus := range existingStatuses {
				if existingStatus.StatusName == statusName {
					found = true
					break
				}
			}
			if !found {
				status := &entity.WorkflowStatus{
					ID:         uuid.New(),
					WorkflowID: workflowID,
					StatusName: statusName,
					StatusOrder: i,
					CreatedAt:  time.Now().UTC(),
				}
				if err := uc.statusRepo.Create(ctx, status); err != nil {
					return err
				}
			}
		}
	}

	// Update transitions if provided
	if len(req.Transitions) > 0 {
		// Get existing transitions
		existingTransitions, err := uc.transitionRepo.ListByWorkflow(ctx, workflowID)
		if err != nil {
			return err
		}

		// Delete existing transitions
		for _, trans := range existingTransitions {
			if err := uc.transitionRepo.Delete(ctx, trans.ID); err != nil {
				return err
			}
		}

		// Create new transitions
		for _, transDef := range req.Transitions {
			// Get status IDs
			fromStatus, err := uc.statusRepo.FindByName(ctx, workflowID, transDef.FromStatus)
			if err != nil || fromStatus == nil {
				return apperror.ErrValidation
			}

			toStatus, err := uc.statusRepo.FindByName(ctx, workflowID, transDef.ToStatus)
			if err != nil || toStatus == nil {
				return apperror.ErrValidation
			}

			transition := &entity.WorkflowTransition{
				ID:            uuid.New(),
				WorkflowID:    workflowID,
				FromStatusID:  fromStatus.ID,
				ToStatusID:    toStatus.ID,
				ValidationRule: "",
				CreatedAt:     time.Now().UTC(),
			}
			if err := uc.transitionRepo.Create(ctx, transition); err != nil {
				return err
			}
		}
	}

	// Log activity
	uc.logActivity(ctx, projectID, nil, requester.UserID, "workflow_updated", "Workflow updated: "+workflow.Name)

	return nil
}

// ListStatuses returns all statuses in a workflow
func (uc *workflowUseCase) ListStatuses(ctx context.Context, workflowID uuid.UUID, requester domainUC.UserClaims) ([]*entity.WorkflowStatus, error) {
	// Get workflow
	workflow, err := uc.workflowRepo.FindByID(ctx, workflowID)
	if err != nil {
		return nil, err
	}
	if workflow == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is a member of the project
	isMember, _ := uc.memberRepo.IsMember(ctx, workflow.ProjectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	statuses, err := uc.statusRepo.ListByWorkflow(ctx, workflowID)
	if err != nil {
		return nil, err
	}
	if statuses == nil {
		statuses = []*entity.WorkflowStatus{}
	}

	return statuses, nil
}

// CanTransition checks if a record can transition to a target status
func (uc *workflowUseCase) CanTransition(ctx context.Context, recordID uuid.UUID, toStatusID uuid.UUID, requester domainUC.UserClaims) (bool, error) {
	// Get record
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return false, err
	}
	if record == nil {
		return false, apperror.ErrNotFound
	}

	// If no status set, can transition to any status
	if record.Status == nil || *record.Status == "" {
		return true, nil
	}

	// Get workflow for the project
	workflow, err := uc.workflowRepo.FindByProject(ctx, record.ProjectID)
	if err != nil {
		return false, err
	}
	if workflow == nil {
		return false, apperror.ErrNotFound
	}

	// Verify target status belongs to this workflow
	targetStatus, err := uc.statusRepo.FindByID(ctx, toStatusID)
	if err != nil || targetStatus == nil {
		return false, apperror.ErrNotFound
	}
	if targetStatus.WorkflowID != workflow.ID {
		return false, nil
	}

	// Try to find current status - it could be stored as UUID or name
	currentStatusID, parseErr := uuid.Parse(*record.Status)
	if parseErr != nil {
		// Status stored as name, look up by name
		currentStatus, err := uc.statusRepo.FindByName(ctx, workflow.ID, *record.Status)
		if err != nil || currentStatus == nil {
			// Current status not found, allow transition
			return true, nil
		}
		currentStatusID = currentStatus.ID
	}

	// Same status, no transition needed
	if currentStatusID == toStatusID {
		return false, nil
	}

	// Check if explicit transitions exist
	transitions, err := uc.transitionRepo.ListFromStatus(ctx, currentStatusID)
	if err != nil {
		return false, err
	}

	// If no transitions defined, allow any transition (open workflow)
	if len(transitions) == 0 {
		return true, nil
	}

	for _, trans := range transitions {
		if trans.ToStatusID == toStatusID {
			return true, nil
		}
	}

	return false, nil
}

// TransitionRecord moves a record to a new status
func (uc *workflowUseCase) TransitionRecord(ctx context.Context, recordID uuid.UUID, req domainUC.TransitionRecordRequest, requester domainUC.UserClaims) error {
	// Get record
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return err
	}
	if record == nil {
		return apperror.ErrNotFound
	}

	// Verify membership
	isMember, _ := uc.memberRepo.IsMember(ctx, record.ProjectID, requester.UserID)
	if !isMember {
		return apperror.ErrForbidden
	}

	// Get target status
	targetStatus, err := uc.statusRepo.FindByID(ctx, req.ToStatusID)
	if err != nil {
		return apperror.ErrNotFound
	}

	// Update record status - store status ID for board column matching
	statusIDStr := targetStatus.ID.String()
	record.Status = &statusIDStr
	record.UpdatedAt = time.Now().UTC()

	if err := uc.recordRepo.Update(ctx, record); err != nil {
		return err
	}

	// Log activity
	detail := "Record status changed to " + targetStatus.StatusName
	uc.logActivity(ctx, record.ProjectID, &recordID, requester.UserID, "record_status_changed", detail)

	return nil
}

// logActivity logs an activity to the project activity log
func (uc *workflowUseCase) logActivity(ctx context.Context, projectID uuid.UUID, recordID *uuid.UUID, actorID uuid.UUID, action, detail string) {
	activity := &entity.ProjectActivityLog{
		ID:        uuid.New(),
		ProjectID: projectID,
		RecordID:  recordID,
		ActorID:   actorID,
		Action:    action,
		Detail:    detail,
		CreatedAt: time.Now().UTC(),
	}
	_ = uc.activityRepo.Append(ctx, activity)
}
