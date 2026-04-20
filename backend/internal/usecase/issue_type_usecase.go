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

type issueTypeUseCase struct {
	issueTypeRepo   repository.IssueTypeRepository
	schemeRepo      repository.IssueTypeSchemeRepository
	projectRepo     repository.ProjectRepository
	activityRepo    repository.ProjectActivityLogRepository
	memberRepo      repository.ProjectMemberRepository
}

func NewIssueTypeUseCase(
	issueTypeRepo repository.IssueTypeRepository,
	schemeRepo repository.IssueTypeSchemeRepository,
	projectRepo repository.ProjectRepository,
	activityRepo repository.ProjectActivityLogRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.IssueTypeUseCase {
	return &issueTypeUseCase{
		issueTypeRepo:   issueTypeRepo,
		schemeRepo:      schemeRepo,
		projectRepo:     projectRepo,
		activityRepo:    activityRepo,
		memberRepo:      memberRepo,
	}
}

// ListIssueTypes returns all available issue types
func (uc *issueTypeUseCase) ListIssueTypes(ctx context.Context) ([]*entity.IssueType, error) {
	return uc.issueTypeRepo.List(ctx)
}

// GetIssueTypeScheme returns the issue type scheme for a project
func (uc *issueTypeUseCase) GetIssueTypeScheme(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) (*entity.IssueTypeScheme, error) {
	// Verify project exists
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is a member of the project
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	scheme, err := uc.schemeRepo.FindByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if scheme == nil {
		return nil, apperror.ErrNotFound
	}

	return scheme, nil
}

// CreateIssueTypeScheme creates a new issue type scheme for a project
func (uc *issueTypeUseCase) CreateIssueTypeScheme(ctx context.Context, projectID uuid.UUID, req domainUC.CreateIssueTypeSchemeRequest, requester domainUC.UserClaims) (*entity.IssueTypeScheme, error) {
	// Verify project exists
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, apperror.ErrNotFound
	}

	// Verify user is project owner
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	// Check if user is owner
	member, err := uc.memberRepo.GetMember(ctx, projectID, requester.UserID)
	if err != nil || member == nil || member.Role != entity.ProjectRoleOwner {
		return nil, apperror.ErrForbidden
	}

	// Validate issue types exist
	for _, issueTypeID := range req.IssueTypeIDs {
		issueType, err := uc.issueTypeRepo.FindByID(ctx, issueTypeID)
		if err != nil || issueType == nil {
			return nil, apperror.ErrValidation
		}
	}

	// Create scheme
	scheme := &entity.IssueTypeScheme{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      req.Name,
		CreatedAt: time.Now().UTC(),
	}

	if err := uc.schemeRepo.Create(ctx, scheme); err != nil {
		return nil, err
	}

	// Add issue types to scheme
	for _, issueTypeID := range req.IssueTypeIDs {
		item := &entity.IssueTypeSchemeItem{
			ID:         uuid.New(),
			SchemeID:   scheme.ID,
			IssueTypeID: issueTypeID,
			CreatedAt:  time.Now().UTC(),
		}
		if err := uc.schemeRepo.AddItem(ctx, item); err != nil {
			return nil, err
		}
	}

	// Log activity
	uc.logActivity(ctx, projectID, nil, requester.UserID, "issue_type_scheme_created", "Issue type scheme created: "+req.Name)

	return scheme, nil
}

// UpdateIssueTypeScheme updates the issue types in a scheme
func (uc *issueTypeUseCase) UpdateIssueTypeScheme(ctx context.Context, schemeID uuid.UUID, issueTypeIDs []uuid.UUID, requester domainUC.UserClaims) error {
	// Get the scheme
	scheme, err := uc.schemeRepo.FindByID(ctx, schemeID)
	if err != nil {
		return err
	}
	if scheme == nil {
		return apperror.ErrNotFound
	}

	// Verify user is project owner
	member, err := uc.memberRepo.GetMember(ctx, scheme.ProjectID, requester.UserID)
	if err != nil || member == nil || member.Role != entity.ProjectRoleOwner {
		return apperror.ErrForbidden
	}

	// Validate issue types exist
	for _, issueTypeID := range issueTypeIDs {
		issueType, err := uc.issueTypeRepo.FindByID(ctx, issueTypeID)
		if err != nil || issueType == nil {
			return apperror.ErrValidation
		}
	}

	// Get current items
	currentItems, err := uc.schemeRepo.ListItems(ctx, schemeID)
	if err != nil {
		return err
	}

	// Remove items not in new list
	for _, item := range currentItems {
		found := false
		for _, newID := range issueTypeIDs {
			if item.IssueTypeID == newID {
				found = true
				break
			}
		}
		if !found {
			if err := uc.schemeRepo.RemoveItem(ctx, schemeID, item.IssueTypeID); err != nil {
				return err
			}
		}
	}

	// Add new items
	for _, issueTypeID := range issueTypeIDs {
		found := false
		for _, item := range currentItems {
			if item.IssueTypeID == issueTypeID {
				found = true
				break
			}
		}
		if !found {
			newItem := &entity.IssueTypeSchemeItem{
				ID:         uuid.New(),
				SchemeID:   schemeID,
				IssueTypeID: issueTypeID,
				CreatedAt:  time.Now().UTC(),
			}
			if err := uc.schemeRepo.AddItem(ctx, newItem); err != nil {
				return err
			}
		}
	}

	// Log activity
	uc.logActivity(ctx, scheme.ProjectID, nil, requester.UserID, "issue_type_scheme_updated", "Issue type scheme updated: "+scheme.Name)

	return nil
}

// logActivity logs an activity to the project activity log
func (uc *issueTypeUseCase) logActivity(ctx context.Context, projectID uuid.UUID, recordID *uuid.UUID, actorID uuid.UUID, action, detail string) {
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
