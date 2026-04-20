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

type customFieldUseCase struct {
	fieldRepo       repository.CustomFieldRepository
	optionRepo      repository.CustomFieldOptionRepository
	valueRepo       repository.CustomFieldValueRepository
	projectRepo     repository.ProjectRepository
	recordRepo      repository.ProjectRecordRepository
	activityRepo    repository.ProjectActivityLogRepository
	memberRepo      repository.ProjectMemberRepository
}

func NewCustomFieldUseCase(
	fieldRepo repository.CustomFieldRepository,
	optionRepo repository.CustomFieldOptionRepository,
	valueRepo repository.CustomFieldValueRepository,
	projectRepo repository.ProjectRepository,
	recordRepo repository.ProjectRecordRepository,
	activityRepo repository.ProjectActivityLogRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.CustomFieldUseCase {
	return &customFieldUseCase{
		fieldRepo:       fieldRepo,
		optionRepo:      optionRepo,
		valueRepo:       valueRepo,
		projectRepo:     projectRepo,
		recordRepo:      recordRepo,
		activityRepo:    activityRepo,
		memberRepo:      memberRepo,
	}
}

// CreateCustomField creates a new custom field for a project
func (uc *customFieldUseCase) CreateCustomField(ctx context.Context, projectID uuid.UUID, req domainUC.CreateCustomFieldRequest, requester domainUC.UserClaims) (*entity.CustomField, error) {
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

	// Validate field type
	validTypes := map[string]bool{
		"text":        true,
		"textarea":    true,
		"dropdown":    true,
		"multiselect": true,
		"date":        true,
		"number":      true,
		"checkbox":    true,
	}
	if !validTypes[req.FieldType] {
		return nil, apperror.ErrValidation
	}

	// Validate options for dropdown/multiselect
	if (req.FieldType == "dropdown" || req.FieldType == "multiselect") && len(req.Options) == 0 {
		return nil, apperror.ErrValidation
	}

	// Create field
	field := &entity.CustomField{
		ID:         uuid.New(),
		ProjectID:  projectID,
		Name:       req.Name,
		FieldType:  req.FieldType,
		IsRequired: req.IsRequired,
		CreatedAt:  time.Now().UTC(),
	}

	if err := uc.fieldRepo.Create(ctx, field); err != nil {
		return nil, err
	}

	// Create options if provided
	for i, optionValue := range req.Options {
		option := &entity.CustomFieldOption{
			ID:          uuid.New(),
			FieldID:     field.ID,
			OptionValue: optionValue,
			OptionOrder: i,
			CreatedAt:   time.Now().UTC(),
		}
		if err := uc.optionRepo.Create(ctx, option); err != nil {
			return nil, err
		}
	}

	// Log activity
	uc.logActivity(ctx, projectID, nil, requester.UserID, "custom_field_created", "Custom field created: "+field.Name)

	return field, nil
}

// ListCustomFields returns all custom fields for a project
func (uc *customFieldUseCase) ListCustomFields(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) ([]*entity.CustomField, error) {
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

	fields, err := uc.fieldRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if fields == nil {
		fields = []*entity.CustomField{}
	}

	return fields, nil
}

// UpdateCustomField updates a custom field's name and options
func (uc *customFieldUseCase) UpdateCustomField(ctx context.Context, projectID uuid.UUID, fieldID uuid.UUID, req domainUC.UpdateCustomFieldRequest, requester domainUC.UserClaims) (*entity.CustomField, error) {
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

	// Get field
	field, err := uc.fieldRepo.FindByID(ctx, fieldID)
	if err != nil {
		return nil, err
	}
	if field == nil || field.ProjectID != projectID {
		return nil, apperror.ErrNotFound
	}

	// Update name if provided
	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		field.Name = *req.Name
	}

	// Update field
	if err := uc.fieldRepo.Update(ctx, field); err != nil {
		return nil, err
	}

	// Update options if provided
	if len(req.Options) > 0 {
		// Delete existing options
		existingOptions, err := uc.optionRepo.ListByField(ctx, fieldID)
		if err == nil {
			for _, opt := range existingOptions {
				_ = uc.optionRepo.Delete(ctx, opt.ID)
			}
		}

		// Create new options
		for i, optionValue := range req.Options {
			option := &entity.CustomFieldOption{
				ID:          uuid.New(),
				FieldID:     field.ID,
				OptionValue: optionValue,
				OptionOrder: i,
				CreatedAt:   time.Now().UTC(),
			}
			if err := uc.optionRepo.Create(ctx, option); err != nil {
				return nil, err
			}
		}
	}

	// Log activity
	uc.logActivity(ctx, projectID, nil, requester.UserID, "custom_field_updated", "Custom field updated: "+field.Name)

	return field, nil
}

// DeleteCustomField deletes a custom field and all its values
func (uc *customFieldUseCase) DeleteCustomField(ctx context.Context, projectID uuid.UUID, fieldID uuid.UUID, requester domainUC.UserClaims) error {
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

	// Get field
	field, err := uc.fieldRepo.FindByID(ctx, fieldID)
	if err != nil {
		return err
	}
	if field == nil || field.ProjectID != projectID {
		return apperror.ErrNotFound
	}

	// Delete field (cascade will handle options and values)
	if err := uc.fieldRepo.Delete(ctx, fieldID); err != nil {
		return err
	}

	// Log activity
	uc.logActivity(ctx, projectID, nil, requester.UserID, "custom_field_deleted", "Custom field deleted: "+field.Name)

	return nil
}

// SetFieldValue sets the value of a custom field for a record
func (uc *customFieldUseCase) SetFieldValue(ctx context.Context, recordID uuid.UUID, fieldID uuid.UUID, value string, requester domainUC.UserClaims) error {
	// Get record
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return err
	}
	if record == nil {
		return apperror.ErrNotFound
	}

	// Get field
	field, err := uc.fieldRepo.FindByID(ctx, fieldID)
	if err != nil {
		return err
	}
	if field == nil {
		return apperror.ErrNotFound
	}

	// Validate field value based on type
	if err := uc.validateFieldValue(ctx, field, value); err != nil {
		return err
	}

	// Check if value already exists
	existingValue, err := uc.valueRepo.GetByRecordAndField(ctx, recordID, fieldID)
	if err == nil && existingValue != nil {
		// Update existing value
		existingValue.Value = value
		existingValue.UpdatedAt = time.Now().UTC()
		return uc.valueRepo.Update(ctx, existingValue)
	}

	// Create new value
	fieldValue := &entity.CustomFieldValue{
		ID:        uuid.New(),
		RecordID:  recordID,
		FieldID:   fieldID,
		Value:     value,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return uc.valueRepo.Set(ctx, fieldValue)
}

// GetFieldValues returns all custom field values for a record
func (uc *customFieldUseCase) GetFieldValues(ctx context.Context, recordID uuid.UUID, requester domainUC.UserClaims) ([]*entity.CustomFieldValue, error) {
	// Get record
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, apperror.ErrNotFound
	}

	values, err := uc.valueRepo.GetByRecord(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if values == nil {
		values = []*entity.CustomFieldValue{}
	}

	return values, nil
}

// validateFieldValue validates a field value based on field type
func (uc *customFieldUseCase) validateFieldValue(ctx context.Context, field *entity.CustomField, value string) error {
	// Check required fields
	if field.IsRequired && strings.TrimSpace(value) == "" {
		return apperror.ErrValidation
	}

	// Type-specific validation
	switch field.FieldType {
	case "number":
		// Basic number validation - could be enhanced with regex
		if strings.TrimSpace(value) != "" {
			// Validate it's a valid number format
			// This is a simplified check - in production, use proper parsing
			if !isValidNumber(value) {
				return apperror.ErrValidation
			}
		}
	case "date":
		// Basic date validation
		if strings.TrimSpace(value) != "" {
			if !isValidDate(value) {
				return apperror.ErrValidation
			}
		}
	case "dropdown", "multiselect":
		// Validate against available options
		if strings.TrimSpace(value) != "" {
			options, err := uc.optionRepo.ListByField(ctx, field.ID)
			if err != nil {
				return err
			}
			found := false
			for _, opt := range options {
				if opt.OptionValue == value {
					found = true
					break
				}
			}
			if !found {
				return apperror.ErrValidation
			}
		}
	}

	return nil
}

// isValidNumber checks if a string is a valid number
func isValidNumber(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	// Simple check - could be enhanced with regex or strconv
	for i, ch := range s {
		if i == 0 && (ch == '-' || ch == '+') {
			continue
		}
		if ch == '.' {
			continue
		}
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// isValidDate checks if a string is a valid date (YYYY-MM-DD format)
func isValidDate(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) != 10 {
		return false
	}
	// Basic format check
	if s[4] != '-' || s[7] != '-' {
		return false
	}
	// Could be enhanced with actual date parsing
	return true
}

// logActivity logs an activity to the project activity log
func (uc *customFieldUseCase) logActivity(ctx context.Context, projectID uuid.UUID, recordID *uuid.UUID, actorID uuid.UUID, action, detail string) {
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
