package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type searchUseCase struct {
	recordRepo      repository.ProjectRecordRepository
	projectRepo     repository.ProjectRepository
	memberRepo      repository.ProjectMemberRepository
}

func NewSearchUseCase(
	recordRepo repository.ProjectRecordRepository,
	projectRepo repository.ProjectRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.SearchUseCase {
	return &searchUseCase{
		recordRepo:      recordRepo,
		projectRepo:     projectRepo,
		memberRepo:      memberRepo,
	}
}

// SearchRecords searches for records with advanced filtering
func (uc *searchUseCase) SearchRecords(ctx context.Context, projectID uuid.UUID, query string, filters domainUC.SearchFilters, requester domainUC.UserClaims) ([]*entity.ProjectRecord, error) {
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
	allRecords, err := uc.recordRepo.ListByProject(ctx, projectID, repository.ProjectRecordFilter{})
	if err != nil {
		return nil, err
	}

	// Apply filters
	var results []*entity.ProjectRecord
	for _, record := range allRecords {
		if uc.matchesFilters(record, query, filters) {
			results = append(results, record)
		}
	}

	if results == nil {
		results = []*entity.ProjectRecord{}
	}

	return results, nil
}

// SaveFilter saves a filter configuration
func (uc *searchUseCase) SaveFilter(ctx context.Context, projectID uuid.UUID, req domainUC.SaveFilterRequest, requester domainUC.UserClaims) error {
	// Verify project exists
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return apperror.ErrNotFound
	}

	// Verify user is a member
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return apperror.ErrForbidden
	}

	// TODO: Implement saved filter storage
	// This would require a SavedFilter repository and entity
	// For now, just validate the input

	if strings.TrimSpace(req.Name) == "" {
		return apperror.ErrValidation
	}

	return nil
}

// ListSavedFilters returns all saved filters for a project
func (uc *searchUseCase) ListSavedFilters(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) ([]*domainUC.SavedFilter, error) {
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

	// TODO: Implement saved filter retrieval
	// This would require a SavedFilter repository
	// For now, return empty list

	return []*domainUC.SavedFilter{}, nil
}

// matchesFilters checks if a record matches all provided filters
func (uc *searchUseCase) matchesFilters(record *entity.ProjectRecord, query string, filters domainUC.SearchFilters) bool {
	// Text search on title and description
	if strings.TrimSpace(query) != "" {
		query = strings.ToLower(query)
		title := strings.ToLower(record.Title)
		description := strings.ToLower(record.Description)

		if !strings.Contains(title, query) && !strings.Contains(description, query) {
			return false
		}
	}

	// Filter by issue type
	if filters.IssueType != nil {
		if record.IssueTypeID == nil || *record.IssueTypeID != *filters.IssueType {
			return false
		}
	}

	// Filter by status
	if filters.Status != nil {
		if record.Status == nil || *record.Status != *filters.Status {
			return false
		}
	}

	// Filter by assignee
	if filters.Assignee != nil {
		if record.AssignedTo == nil || *record.AssignedTo != *filters.Assignee {
			return false
		}
	}

	// Filter by due date range
	if filters.DueDateFrom != nil && record.DueDate != nil {
		if record.DueDate.Before(*filters.DueDateFrom) {
			return false
		}
	}

	if filters.DueDateTo != nil && record.DueDate != nil {
		if record.DueDate.After(*filters.DueDateTo) {
			return false
		}
	}

	// TODO: Filter by label (would need to query record_labels)
	// TODO: Filter by sprint (would need to query sprint_records)
	// TODO: Filter by custom fields (would need to query custom_field_values)

	return true
}
