package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type labelUseCase struct {
	labelRepo       repository.LabelRepository
	recordLabelRepo repository.RecordLabelRepository
	projectRepo     repository.ProjectRepository
	recordRepo      repository.ProjectRecordRepository
	activityRepo    repository.ProjectActivityLogRepository
	memberRepo      repository.ProjectMemberRepository
}

func NewLabelUseCase(
	labelRepo repository.LabelRepository,
	recordLabelRepo repository.RecordLabelRepository,
	projectRepo repository.ProjectRepository,
	recordRepo repository.ProjectRecordRepository,
	activityRepo repository.ProjectActivityLogRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.LabelUseCase {
	return &labelUseCase{
		labelRepo:       labelRepo,
		recordLabelRepo: recordLabelRepo,
		projectRepo:     projectRepo,
		recordRepo:      recordRepo,
		activityRepo:    activityRepo,
		memberRepo:      memberRepo,
	}
}

// CreateLabel creates a new label for a project
func (uc *labelUseCase) CreateLabel(ctx context.Context, projectID uuid.UUID, req domainUC.CreateLabelRequest, requester domainUC.UserClaims) (*entity.Label, error) {
	// Verify project exists
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is a project member
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	// Validate label name
	if strings.TrimSpace(req.Name) == "" {
		return nil, apperror.ErrValidation
	}

	// Check for duplicate label name in project
	existingLabels, err := uc.labelRepo.ListByProject(ctx, projectID)
	if err == nil {
		for _, label := range existingLabels {
			if strings.EqualFold(label.Name, req.Name) {
				return nil, apperror.ErrValidation
			}
		}
	}

	// Create label
	label := &entity.Label{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      req.Name,
		Color:     req.Color,
		CreatedAt: time.Now().UTC(),
	}

	if err := uc.labelRepo.Create(ctx, label); err != nil {
		return nil, err
	}

	// Log activity
	uc.logActivity(ctx, projectID, nil, requester.UserID, "label_created", "Label created: "+req.Name)

	return label, nil
}

// ListLabels returns all labels for a project
func (uc *labelUseCase) ListLabels(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) ([]*entity.Label, error) {
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

	labels, err := uc.labelRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if labels == nil {
		labels = []*entity.Label{}
	}

	return labels, nil
}

// AddLabelToRecord adds a label to a record
func (uc *labelUseCase) AddLabelToRecord(ctx context.Context, recordID uuid.UUID, labelID uuid.UUID, requester domainUC.UserClaims) error {
	// Get record
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return err
	}
	if record == nil {
		return apperror.ErrNotFound
	}

	// Verify user is a member
	isMember, _ := uc.memberRepo.IsMember(ctx, record.ProjectID, requester.UserID)
	if !isMember {
		return apperror.ErrForbidden
	}

	// Get label
	label, err := uc.labelRepo.FindByID(ctx, labelID)
	if err != nil {
		return err
	}
	if label == nil || label.ProjectID != record.ProjectID {
		return apperror.ErrNotFound
	}

	// Check if label already added
	existingLabels, err := uc.recordLabelRepo.ListByRecord(ctx, recordID)
	if err == nil {
		for _, rl := range existingLabels {
			if rl.LabelID == labelID {
				// Already added
				return nil
			}
		}
	}

	// Create record label
	recordLabel := &entity.RecordLabel{
		ID:        uuid.New(),
		RecordID:  recordID,
		LabelID:   labelID,
		CreatedAt: time.Now().UTC(),
	}

	if err := uc.recordLabelRepo.Create(ctx, recordLabel); err != nil {
		return err
	}

	// Log activity
	uc.logActivity(ctx, record.ProjectID, &recordID, requester.UserID, "label_added", "Label added: "+label.Name)

	return nil
}

// RemoveLabelFromRecord removes a label from a record
func (uc *labelUseCase) RemoveLabelFromRecord(ctx context.Context, recordID uuid.UUID, labelID uuid.UUID, requester domainUC.UserClaims) error {
	// Get record
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return err
	}
	if record == nil {
		return apperror.ErrNotFound
	}

	// Verify user is a member
	isMember, _ := uc.memberRepo.IsMember(ctx, record.ProjectID, requester.UserID)
	if !isMember {
		return apperror.ErrForbidden
	}

	// Get label
	label, err := uc.labelRepo.FindByID(ctx, labelID)
	if err != nil {
		return err
	}
	if label == nil {
		return apperror.ErrNotFound
	}

	// Find and delete record label
	recordLabels, err := uc.recordLabelRepo.ListByRecord(ctx, recordID)
	if err != nil {
		return err
	}

	for _, rl := range recordLabels {
		if rl.LabelID == labelID {
			if err := uc.recordLabelRepo.Delete(ctx, rl.ID); err != nil {
				return err
			}
			break
		}
	}

	// Log activity
	uc.logActivity(ctx, record.ProjectID, &recordID, requester.UserID, "label_removed", "Label removed: "+label.Name)

	return nil
}

// GetRecordLabels returns all labels for a record
func (uc *labelUseCase) GetRecordLabels(ctx context.Context, recordID uuid.UUID, requester domainUC.UserClaims) ([]*entity.Label, error) {
	// Get record
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is a member
	isMember, _ := uc.memberRepo.IsMember(ctx, record.ProjectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	// Get record labels
	recordLabels, err := uc.recordLabelRepo.ListByRecord(ctx, recordID)
	if err != nil {
		return nil, err
	}

	// Get label details
	var labels []*entity.Label
	for _, rl := range recordLabels {
		label, err := uc.labelRepo.FindByID(ctx, rl.LabelID)
		if err == nil && label != nil {
			labels = append(labels, label)
		}
	}

	if labels == nil {
		labels = []*entity.Label{}
	}

	return labels, nil
}

// DeleteLabel deletes a label from a project
func (uc *labelUseCase) DeleteLabel(ctx context.Context, labelID uuid.UUID, requester domainUC.UserClaims) error {
	// Get label
	label, err := uc.labelRepo.FindByID(ctx, labelID)
	if err != nil {
		return err
	}
	if label == nil {
		return apperror.ErrNotFound
	}

	// Verify user is project owner
	member, err := uc.memberRepo.GetMember(ctx, label.ProjectID, requester.UserID)
	if err != nil || member == nil || member.Role != entity.ProjectRoleOwner {
		return apperror.ErrForbidden
	}

	// Delete label (cascade will handle record_labels)
	if err := uc.labelRepo.Delete(ctx, labelID); err != nil {
		return err
	}

	// Log activity
	uc.logActivity(ctx, label.ProjectID, nil, requester.UserID, "label_deleted", "Label deleted: "+label.Name)

	return nil
}

// logActivity logs an activity to the project activity log
func (uc *labelUseCase) logActivity(ctx context.Context, projectID uuid.UUID, recordID *uuid.UUID, actorID uuid.UUID, action, detail string) {
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
