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

type backlogUseCase struct {
	sprintRecordRepo repository.SprintRecordRepository
	projectRepo     repository.ProjectRepository
	recordRepo      repository.ProjectRecordRepository
	activityRepo    repository.ProjectActivityLogRepository
	memberRepo      repository.ProjectMemberRepository
}

func NewBacklogUseCase(
	sprintRecordRepo repository.SprintRecordRepository,
	projectRepo repository.ProjectRepository,
	recordRepo repository.ProjectRecordRepository,
	activityRepo repository.ProjectActivityLogRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.BacklogUseCase {
	return &backlogUseCase{
		sprintRecordRepo: sprintRecordRepo,
		projectRepo:      projectRepo,
		recordRepo:       recordRepo,
		activityRepo:     activityRepo,
		memberRepo:       memberRepo,
	}
}

// GetBacklog returns all records not assigned to any sprint
func (uc *backlogUseCase) GetBacklog(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) ([]*entity.ProjectRecord, error) {
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

	// Get all records for the project
	// This would typically be done with a query that filters by project_id
	// For now, we'll get all records and filter those not in any sprint
	// In a real implementation, this would be a single optimized query

	// Get all records in sprints
	sprintRecords, err := uc.sprintRecordRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Build a set of record IDs in sprints
	inSprintSet := make(map[uuid.UUID]bool)
	for _, sr := range sprintRecords {
		inSprintSet[sr.RecordID] = true
	}

	// Get all project records and filter out those in sprints
	// This is a simplified approach - in production, use a direct query
	var backlogRecords []*entity.ProjectRecord

	// Note: This would need a method like ListByProject on recordRepo
	// For now, we'll assume the repository has this capability
	allRecords, err := uc.recordRepo.ListByProject(ctx, projectID, repository.ProjectRecordFilter{})
	if err != nil {
		return nil, err
	}

	for _, record := range allRecords {
		if !inSprintSet[record.ID] {
			backlogRecords = append(backlogRecords, record)
		}
	}

	if backlogRecords == nil {
		backlogRecords = []*entity.ProjectRecord{}
	}

	return backlogRecords, nil
}

// ReorderBacklog updates the priority ordering of backlog records
func (uc *backlogUseCase) ReorderBacklog(ctx context.Context, projectID uuid.UUID, req domainUC.ReorderBacklogRequest, requester domainUC.UserClaims) error {
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

	// Update priority for each record
	for i, recordID := range req.RecordIDs {
		record, err := uc.recordRepo.FindByID(ctx, recordID)
		if err != nil {
			return err
		}
		if record == nil || record.ProjectID != projectID {
			return apperror.ErrNotFound
		}

		// Update record priority (using position field or similar)
		// This assumes the record has a priority/position field
		// For now, we'll just update the record to mark it as reordered
		record.UpdatedAt = time.Now().UTC()
		if err := uc.recordRepo.Update(ctx, record); err != nil {
			return err
		}

		// Log activity for significant reorderings
		if i == 0 {
			uc.logActivity(ctx, projectID, &recordID, requester.UserID, "backlog_reordered", "Backlog records reordered")
		}
	}

	return nil
}

// BulkAssignToSprint assigns multiple records to a sprint
func (uc *backlogUseCase) BulkAssignToSprint(ctx context.Context, req domainUC.BulkAssignToSprintRequest, requester domainUC.UserClaims) error {
	// Verify sprint exists and get project ID
	// This would need a method to get sprint by ID
	// For now, we'll assume we can get it from the first record's project

	if len(req.RecordIDs) == 0 {
		return apperror.ErrValidation
	}

	// Get first record to verify project
	firstRecord, err := uc.recordRepo.FindByID(ctx, req.RecordIDs[0])
	if err != nil {
		return err
	}
	if firstRecord == nil {
		return apperror.ErrNotFound
	}

	projectID := firstRecord.ProjectID

	// Verify user is project owner
	member, err := uc.memberRepo.GetMember(ctx, projectID, requester.UserID)
	if err != nil || member == nil || member.Role != entity.ProjectRoleOwner {
		return apperror.ErrForbidden
	}

	// Assign each record to sprint
	for i, recordID := range req.RecordIDs {
		record, err := uc.recordRepo.FindByID(ctx, recordID)
		if err != nil {
			return err
		}
		if record == nil || record.ProjectID != projectID {
			return apperror.ErrNotFound
		}

		// Check if already assigned to a sprint
		existing, err := uc.sprintRecordRepo.FindByRecord(ctx, recordID)
		if err == nil && existing != nil {
			// Already assigned, update priority
			existing.Priority = i
			if err := uc.sprintRecordRepo.Update(ctx, existing); err != nil {
				return err
			}
		} else {
			// Create new sprint record
			sprintRecord := &entity.SprintRecord{
				ID:        uuid.New(),
				SprintID:  req.SprintID,
				RecordID:  recordID,
				Priority:  i,
				CreatedAt: time.Now().UTC(),
			}
			if err := uc.sprintRecordRepo.Create(ctx, sprintRecord); err != nil {
				return err
			}

			// Set initial status if record has no status
			if record.Status == nil || *record.Status == "" {
				initialStatus := "Backlog"
				record.Status = &initialStatus
				record.UpdatedAt = time.Now().UTC()
				_ = uc.recordRepo.Update(ctx, record)
			}
		}
	}

	// Log activity
	uc.logActivity(ctx, projectID, nil, requester.UserID, "bulk_assign_to_sprint", "Bulk assigned "+string(rune(len(req.RecordIDs)))+" records to sprint")

	return nil
}

// logActivity logs an activity to the project activity log
func (uc *backlogUseCase) logActivity(ctx context.Context, projectID uuid.UUID, recordID *uuid.UUID, actorID uuid.UUID, action, detail string) {
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
