package usecase

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type attachmentUseCase struct {
	attachmentRepo  repository.AttachmentRepository
	recordRepo      repository.ProjectRecordRepository
	activityRepo    repository.ProjectActivityLogRepository
	memberRepo      repository.ProjectMemberRepository
}

func NewAttachmentUseCase(
	attachmentRepo repository.AttachmentRepository,
	recordRepo repository.ProjectRecordRepository,
	activityRepo repository.ProjectActivityLogRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.AttachmentUseCase {
	return &attachmentUseCase{
		attachmentRepo:  attachmentRepo,
		recordRepo:      recordRepo,
		activityRepo:    activityRepo,
		memberRepo:      memberRepo,
	}
}

// UploadAttachment uploads a file and attaches it to a record
func (uc *attachmentUseCase) UploadAttachment(ctx context.Context, recordID uuid.UUID, file *domainUC.FileUpload, requester domainUC.UserClaims) (*entity.Attachment, error) {
	// Get record
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is a member of the project
	isMember, _ := uc.memberRepo.IsMember(ctx, record.ProjectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	// Validate file size (max 50MB)
	const maxFileSize = 50 * 1024 * 1024 // 50MB
	if file.FileSize > maxFileSize {
		return nil, apperror.ErrValidation
	}

	// Validate file type
	if !uc.isValidFileType(file.FileType) {
		return nil, apperror.ErrValidation
	}

	// Generate file path
	fileID := uuid.New().String()
	filePath := "uploads/attachments/" + fileID + "/" + file.FileName

	// Store file to disk
	dirPath := "uploads/attachments/" + fileID
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, apperror.ErrInternal
	}
	if err := os.WriteFile(filePath, file.Content, 0644); err != nil {
		return nil, apperror.ErrInternal
	}

	// Create attachment record
	attachment := &entity.Attachment{
		ID:         uuid.New(),
		RecordID:   recordID,
		FileName:   file.FileName,
		FileSize:   file.FileSize,
		FileType:   file.FileType,
		FilePath:   filePath,
		UploaderID: requester.UserID,
		CreatedAt:  time.Now().UTC(),
	}

	if err := uc.attachmentRepo.Create(ctx, attachment); err != nil {
		return nil, err
	}

	// Log activity
	uc.logActivity(ctx, record.ProjectID, &recordID, requester.UserID, "attachment_uploaded", "Attachment uploaded: "+file.FileName)

	return attachment, nil
}

// ListAttachments returns all attachments for a record
func (uc *attachmentUseCase) ListAttachments(ctx context.Context, recordID uuid.UUID, requester domainUC.UserClaims) ([]*entity.Attachment, error) {
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

	attachments, err := uc.attachmentRepo.ListByRecord(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if attachments == nil {
		attachments = []*entity.Attachment{}
	}

	return attachments, nil
}

// GetAttachment returns a single attachment by ID
func (uc *attachmentUseCase) GetAttachment(ctx context.Context, attachmentID uuid.UUID, requester domainUC.UserClaims) (*entity.Attachment, error) {
	attachment, err := uc.attachmentRepo.FindByID(ctx, attachmentID)
	if err != nil {
		return nil, err
	}
	if attachment == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is a member of the project
	record, err := uc.recordRepo.FindByID(ctx, attachment.RecordID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, apperror.ErrNotFound
	}

	isMember, _ := uc.memberRepo.IsMember(ctx, record.ProjectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	return attachment, nil
}

// DeleteAttachment deletes an attachment
func (uc *attachmentUseCase) DeleteAttachment(ctx context.Context, attachmentID uuid.UUID, requester domainUC.UserClaims) error {
	// Get attachment
	attachment, err := uc.attachmentRepo.FindByID(ctx, attachmentID)
	if err != nil {
		return err
	}
	if attachment == nil {
		return apperror.ErrNotFound
	}

	// Verify user is the uploader
	if attachment.UploaderID != requester.UserID {
		return apperror.ErrForbidden
	}

	// Get record
	record, err := uc.recordRepo.FindByID(ctx, attachment.RecordID)
	if err != nil {
		return err
	}
	if record == nil {
		return apperror.ErrNotFound
	}

	// TODO: Delete file from storage

	// Delete attachment record
	if err := uc.attachmentRepo.Delete(ctx, attachmentID); err != nil {
		return err
	}

	// Log activity
	uc.logActivity(ctx, record.ProjectID, &attachment.RecordID, requester.UserID, "attachment_deleted", "Attachment deleted: "+attachment.FileName)

	return nil
}

// isValidFileType checks if a file type is allowed
func (uc *attachmentUseCase) isValidFileType(fileType string) bool {
	// Normalize file type
	fileType = strings.ToLower(strings.TrimSpace(fileType))

	// Allowed file types
	allowedTypes := map[string]bool{
		// Images
		"image/jpeg":      true,
		"image/jpg":       true,
		"image/png":       true,
		"image/gif":       true,
		"image/webp":      true,

		// Documents
		"application/pdf":                                   true,
		"application/msword":                               true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"application/vnd.ms-excel":                         true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       true,

		// Archives
		"application/zip":       true,
		"application/x-rar-compressed": true,
		"application/x-7z-compressed":  true,

		// Text
		"text/plain":       true,
		"text/csv":         true,
		"text/html":        true,
		"application/json": true,
	}

	return allowedTypes[fileType]
}

// logActivity logs an activity to the project activity log
func (uc *attachmentUseCase) logActivity(ctx context.Context, projectID uuid.UUID, recordID *uuid.UUID, actorID uuid.UUID, action, detail string) {
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
