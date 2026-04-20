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

type bulkOperationUseCase struct {
	recordRepo      repository.ProjectRecordRepository
	statusRepo      repository.WorkflowStatusRepository
	labelRepo       repository.LabelRepository
	recordLabelRepo repository.RecordLabelRepository
	projectRepo     repository.ProjectRepository
	activityRepo    repository.ProjectActivityLogRepository
	memberRepo      repository.ProjectMemberRepository
}

func NewBulkOperationUseCase(
	recordRepo repository.ProjectRecordRepository,
	statusRepo repository.WorkflowStatusRepository,
	labelRepo repository.LabelRepository,
	recordLabelRepo repository.RecordLabelRepository,
	projectRepo repository.ProjectRepository,
	activityRepo repository.ProjectActivityLogRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.BulkOperationUseCase {
	return &bulkOperationUseCase{
		recordRepo:      recordRepo,
		statusRepo:      statusRepo,
		labelRepo:       labelRepo,
		recordLabelRepo: recordLabelRepo,
		projectRepo:     projectRepo,
		activityRepo:    activityRepo,
		memberRepo:      memberRepo,
	}
}

// BulkChangeStatus changes the status of multiple records
func (uc *bulkOperationUseCase) BulkChangeStatus(ctx context.Context, req domainUC.BulkChangeStatusRequest, requester domainUC.UserClaims) error {
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

	// Get target status
	targetStatus, err := uc.statusRepo.FindByID(ctx, req.StatusID)
	if err != nil {
		return err
	}
	if targetStatus == nil {
		return apperror.ErrNotFound
	}

	// Update each record
	for _, recordID := range req.RecordIDs {
		record, err := uc.recordRepo.FindByID(ctx, recordID)
		if err != nil {
			return err
		}
		if record == nil || record.ProjectID != projectID {
			return apperror.ErrNotFound
		}

		// Update status - store status ID for consistency with board
		statusIDStr := targetStatus.ID.String()
		record.Status = &statusIDStr
		record.UpdatedAt = time.Now().UTC()

		if err := uc.recordRepo.Update(ctx, record); err != nil {
			return err
		}

		// Log activity for each record
		detail := "Status changed to " + targetStatus.StatusName
		uc.logActivity(ctx, projectID, &recordID, requester.UserID, "record_status_changed", detail)
	}

	return nil
}

// BulkAssignTo assigns multiple records to a user
func (uc *bulkOperationUseCase) BulkAssignTo(ctx context.Context, req domainUC.BulkAssignToRequest, requester domainUC.UserClaims) error {
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

	// Update each record
	for _, recordID := range req.RecordIDs {
		record, err := uc.recordRepo.FindByID(ctx, recordID)
		if err != nil {
			return err
		}
		if record == nil || record.ProjectID != projectID {
			return apperror.ErrNotFound
		}

		// Update assignee
		record.AssignedTo = &req.AssigneeID
		record.UpdatedAt = time.Now().UTC()

		if err := uc.recordRepo.Update(ctx, record); err != nil {
			return err
		}

		// Log activity
		uc.logActivity(ctx, projectID, &recordID, requester.UserID, "record_assigned", "Record assigned to user")
	}

	return nil
}

// BulkAddLabel adds a label to multiple records
func (uc *bulkOperationUseCase) BulkAddLabel(ctx context.Context, req domainUC.BulkAddLabelRequest, requester domainUC.UserClaims) error {
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

	// Get label
	label, err := uc.labelRepo.FindByID(ctx, req.LabelID)
	if err != nil {
		return err
	}
	if label == nil || label.ProjectID != projectID {
		return apperror.ErrNotFound
	}

	// Add label to each record
	for _, recordID := range req.RecordIDs {
		record, err := uc.recordRepo.FindByID(ctx, recordID)
		if err != nil {
			return err
		}
		if record == nil || record.ProjectID != projectID {
			return apperror.ErrNotFound
		}

		// Check if label already added
		existingLabels, err := uc.recordLabelRepo.ListByRecord(ctx, recordID)
		if err == nil {
			alreadyAdded := false
			for _, rl := range existingLabels {
				if rl.LabelID == req.LabelID {
					alreadyAdded = true
					break
				}
			}
			if alreadyAdded {
				continue
			}
		}

		// Create record label
		recordLabel := &entity.RecordLabel{
			ID:        uuid.New(),
			RecordID:  recordID,
			LabelID:   req.LabelID,
			CreatedAt: time.Now().UTC(),
		}

		if err := uc.recordLabelRepo.Create(ctx, recordLabel); err != nil {
			return err
		}

		// Log activity
		uc.logActivity(ctx, projectID, &recordID, requester.UserID, "label_added", "Label added: "+label.Name)
	}

	return nil
}

// BulkDelete deletes multiple records
func (uc *bulkOperationUseCase) BulkDelete(ctx context.Context, req domainUC.BulkDeleteRequest, requester domainUC.UserClaims) error {
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

	// Delete each record
	for _, recordID := range req.RecordIDs {
		record, err := uc.recordRepo.FindByID(ctx, recordID)
		if err != nil {
			return err
		}
		if record == nil || record.ProjectID != projectID {
			return apperror.ErrNotFound
		}

		// Delete record (cascade will handle related records)
		if err := uc.recordRepo.Delete(ctx, recordID); err != nil {
			return err
		}

		// Log activity
		uc.logActivity(ctx, projectID, &recordID, requester.UserID, "record_deleted", "Record deleted: "+record.Title)
	}

	return nil
}

// logActivity logs an activity to the project activity log
func (uc *bulkOperationUseCase) logActivity(ctx context.Context, projectID uuid.UUID, recordID *uuid.UUID, actorID uuid.UUID, action, detail string) {
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
