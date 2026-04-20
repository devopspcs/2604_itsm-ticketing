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

type sprintUseCase struct {
	sprintRepo      repository.SprintRepository
	sprintRecordRepo repository.SprintRecordRepository
	projectRepo     repository.ProjectRepository
	recordRepo      repository.ProjectRecordRepository
	activityRepo    repository.ProjectActivityLogRepository
	memberRepo      repository.ProjectMemberRepository
}

func NewSprintUseCase(
	sprintRepo repository.SprintRepository,
	sprintRecordRepo repository.SprintRecordRepository,
	projectRepo repository.ProjectRepository,
	recordRepo repository.ProjectRecordRepository,
	activityRepo repository.ProjectActivityLogRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.SprintUseCase {
	return &sprintUseCase{
		sprintRepo:       sprintRepo,
		sprintRecordRepo: sprintRecordRepo,
		projectRepo:      projectRepo,
		recordRepo:       recordRepo,
		activityRepo:     activityRepo,
		memberRepo:       memberRepo,
	}
}

// CreateSprint creates a new sprint for a project
func (uc *sprintUseCase) CreateSprint(ctx context.Context, projectID uuid.UUID, req domainUC.CreateSprintRequest, requester domainUC.UserClaims) (*entity.Sprint, error) {
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

	// Validate dates
	if req.EndDate == nil {
		return nil, apperror.ErrValidation
	}

	sprint := &entity.Sprint{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      req.Name,
		Goal:      req.Goal,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    "Planned",
		CreatedAt: time.Now().UTC(),
	}

	if err := uc.sprintRepo.Create(ctx, sprint); err != nil {
		return nil, err
	}

	// Log activity
	uc.logActivity(ctx, projectID, nil, requester.UserID, "sprint_created", "Sprint created: "+sprint.Name)

	return sprint, nil
}

// ListSprints returns all sprints for a project
func (uc *sprintUseCase) ListSprints(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) ([]*entity.Sprint, error) {
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

	sprints, err := uc.sprintRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if sprints == nil {
		sprints = []*entity.Sprint{}
	}

	return sprints, nil
}

// GetActiveSprint returns the currently active sprint for a project
func (uc *sprintUseCase) GetActiveSprint(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) (*entity.Sprint, error) {
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

	sprints, err := uc.sprintRepo.ListActive(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if len(sprints) == 0 {
		return nil, apperror.ErrNotFound
	}

	return sprints[0], nil
}

// StartSprint transitions a sprint from Planned to Active
func (uc *sprintUseCase) StartSprint(ctx context.Context, sprintID uuid.UUID, requester domainUC.UserClaims) (*entity.Sprint, error) {
	// Get sprint
	sprint, err := uc.sprintRepo.FindByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is project owner
	member, err := uc.memberRepo.GetMember(ctx, sprint.ProjectID, requester.UserID)
	if err != nil || member == nil || member.Role != entity.ProjectRoleOwner {
		return nil, apperror.ErrForbidden
	}

	// Validate status transition
	if sprint.Status != "Planned" {
		return nil, apperror.ErrValidation
	}

	// Update sprint
	now := time.Now().UTC()
	sprint.Status = "Active"
	sprint.ActualStartDate = &now

	if err := uc.sprintRepo.Update(ctx, sprint); err != nil {
		return nil, err
	}

	// Log activity
	uc.logActivity(ctx, sprint.ProjectID, nil, requester.UserID, "sprint_started", "Sprint started: "+sprint.Name)

	return sprint, nil
}

// CompleteSprint transitions a sprint from Active to Completed
func (uc *sprintUseCase) CompleteSprint(ctx context.Context, sprintID uuid.UUID, requester domainUC.UserClaims) (*domainUC.SprintMetrics, error) {
	// Get sprint
	sprint, err := uc.sprintRepo.FindByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is project owner
	member, err := uc.memberRepo.GetMember(ctx, sprint.ProjectID, requester.UserID)
	if err != nil || member == nil || member.Role != entity.ProjectRoleOwner {
		return nil, apperror.ErrForbidden
	}

	// Validate status transition
	if sprint.Status != "Active" {
		return nil, apperror.ErrValidation
	}

	// Get all records in sprint
	sprintRecords, err := uc.sprintRecordRepo.ListBySprint(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	// Calculate metrics
	totalRecords := len(sprintRecords)
	completedRecords := 0

	// Move incomplete records to backlog
	for _, sr := range sprintRecords {
		record, err := uc.recordRepo.FindByID(ctx, sr.RecordID)
		if err != nil {
			continue
		}

		// Check if record is completed (assuming "Done" status means completed)
		if record.Status != nil && *record.Status == "Done" {
			completedRecords++
		} else {
			// Remove from sprint (move to backlog)
			if err := uc.sprintRecordRepo.Delete(ctx, sr.ID); err != nil {
				continue
			}
		}
	}

	// Calculate completion percentage
	completionPercent := 0.0
	if totalRecords > 0 {
		completionPercent = (float64(completedRecords) / float64(totalRecords)) * 100
	}

	// Update sprint
	now := time.Now().UTC()
	sprint.Status = "Completed"
	sprint.ActualEndDate = &now

	if err := uc.sprintRepo.Update(ctx, sprint); err != nil {
		return nil, err
	}

	// Log activity
	uc.logActivity(ctx, sprint.ProjectID, nil, requester.UserID, "sprint_completed", "Sprint completed: "+sprint.Name)

	metrics := &domainUC.SprintMetrics{
		TotalRecords:      totalRecords,
		CompletedRecords:  completedRecords,
		CompletionPercent: completionPercent,
		Velocity:          completedRecords,
	}

	return metrics, nil
}

// AssignRecordToSprint assigns a record to a sprint
func (uc *sprintUseCase) AssignRecordToSprint(ctx context.Context, recordID uuid.UUID, sprintID uuid.UUID, requester domainUC.UserClaims) error {
	// Get record
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return err
	}
	if record == nil {
		return apperror.ErrNotFound
	}

	// Get sprint
	sprint, err := uc.sprintRepo.FindByID(ctx, sprintID)
	if err != nil {
		return err
	}
	if sprint == nil {
		return apperror.ErrNotFound
	}

	// Verify sprint belongs to same project as record
	if sprint.ProjectID != record.ProjectID {
		return apperror.ErrValidation
	}

	// Check if already assigned
	existing, err := uc.sprintRecordRepo.FindByRecord(ctx, recordID)
	if err == nil && existing != nil {
		// Already assigned, update priority
		existing.Priority = 0
		return uc.sprintRecordRepo.Update(ctx, existing)
	}

	// Create sprint record
	sprintRecord := &entity.SprintRecord{
		ID:        uuid.New(),
		SprintID:  sprintID,
		RecordID:  recordID,
		Priority:  0,
		CreatedAt: time.Now().UTC(),
	}

	if err := uc.sprintRecordRepo.Create(ctx, sprintRecord); err != nil {
		return err
	}

	// Log activity
	uc.logActivity(ctx, record.ProjectID, &recordID, requester.UserID, "record_assigned_to_sprint", "Record assigned to sprint: "+sprint.Name)

	return nil
}

// RemoveRecordFromSprint removes a record from a sprint
func (uc *sprintUseCase) RemoveRecordFromSprint(ctx context.Context, recordID uuid.UUID, requester domainUC.UserClaims) error {
	// Get record
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return err
	}
	if record == nil {
		return apperror.ErrNotFound
	}

	// Find sprint record
	sprintRecord, err := uc.sprintRecordRepo.FindByRecord(ctx, recordID)
	if err != nil {
		return err
	}
	if sprintRecord == nil {
		return apperror.ErrNotFound
	}

	// Delete sprint record
	if err := uc.sprintRecordRepo.Delete(ctx, sprintRecord.ID); err != nil {
		return err
	}

	// Log activity
	uc.logActivity(ctx, record.ProjectID, &recordID, requester.UserID, "record_removed_from_sprint", "Record removed from sprint")

	return nil
}

// GetSprintRecords returns all records in a sprint
func (uc *sprintUseCase) GetSprintRecords(ctx context.Context, sprintID uuid.UUID, requester domainUC.UserClaims) ([]*entity.ProjectRecord, error) {
	// Get sprint
	sprint, err := uc.sprintRepo.FindByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}
	if sprint == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is a member
	isMember, _ := uc.memberRepo.IsMember(ctx, sprint.ProjectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	// Get sprint records
	sprintRecords, err := uc.sprintRecordRepo.ListBySprint(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	// Get full record details
	var records []*entity.ProjectRecord
	for _, sr := range sprintRecords {
		record, err := uc.recordRepo.FindByID(ctx, sr.RecordID)
		if err == nil && record != nil {
			records = append(records, record)
		}
	}

	if records == nil {
		records = []*entity.ProjectRecord{}
	}

	return records, nil
}

// logActivity logs an activity to the project activity log
func (uc *sprintUseCase) logActivity(ctx context.Context, projectID uuid.UUID, recordID *uuid.UUID, actorID uuid.UUID, action, detail string) {
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
