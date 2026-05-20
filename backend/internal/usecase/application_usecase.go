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

type applicationUseCase struct {
	appRepo    repository.ApplicationRepository
	accessRepo repository.UserAppAccessRepository
	userRepo   repository.UserRepository
}

func NewApplicationUseCase(
	appRepo repository.ApplicationRepository,
	accessRepo repository.UserAppAccessRepository,
	userRepo repository.UserRepository,
) domainUC.ApplicationUseCase {
	return &applicationUseCase{
		appRepo:    appRepo,
		accessRepo: accessRepo,
		userRepo:   userRepo,
	}
}

func (uc *applicationUseCase) CreateApp(ctx context.Context, req domainUC.CreateApplicationRequest) (*entity.Application, error) {
	if req.Name == "" || req.Code == "" {
		return nil, apperror.ErrValidation
	}
	now := time.Now().UTC()
	app := &entity.Application{
		ID:          uuid.New(),
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if app.Icon == "" {
		app.Icon = "apps"
	}
	if app.Color == "" {
		app.Color = "#1976d2"
	}
	if err := uc.appRepo.Create(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

func (uc *applicationUseCase) UpdateApp(ctx context.Context, id uuid.UUID, req domainUC.UpdateApplicationRequest) (*entity.Application, error) {
	app, err := uc.appRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		app.Name = *req.Name
	}
	if req.Description != nil {
		app.Description = *req.Description
	}
	if req.Icon != nil {
		app.Icon = *req.Icon
	}
	if req.Color != nil {
		app.Color = *req.Color
	}
	if req.IsActive != nil {
		app.IsActive = *req.IsActive
	}
	app.UpdatedAt = time.Now().UTC()
	if err := uc.appRepo.Update(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

func (uc *applicationUseCase) DeleteApp(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.appRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return uc.appRepo.Delete(ctx, id)
}

func (uc *applicationUseCase) ListApps(ctx context.Context) ([]*entity.Application, error) {
	return uc.appRepo.List(ctx)
}

func (uc *applicationUseCase) GetApp(ctx context.Context, id uuid.UUID) (*entity.Application, error) {
	return uc.appRepo.FindByID(ctx, id)
}

func (uc *applicationUseCase) GrantAccess(ctx context.Context, req domainUC.GrantAccessRequest, grantedBy uuid.UUID) error {
	// Verify app exists
	if _, err := uc.appRepo.FindByID(ctx, req.AppID); err != nil {
		return err
	}
	// Verify user exists
	if _, err := uc.userRepo.FindByID(ctx, req.UserID); err != nil {
		return err
	}
	access := &entity.UserAppAccess{
		ID:        uuid.New(),
		UserID:    req.UserID,
		AppID:     req.AppID,
		Role:      req.Role,
		GrantedAt: time.Now().UTC(),
		GrantedBy: &grantedBy,
	}
	return uc.accessRepo.Grant(ctx, access)
}

func (uc *applicationUseCase) RevokeAccess(ctx context.Context, userID, appID uuid.UUID) error {
	return uc.accessRepo.Revoke(ctx, userID, appID)
}

func (uc *applicationUseCase) GetUserApps(ctx context.Context, userID uuid.UUID) ([]*domainUC.AppWithAccess, error) {
	accessList, err := uc.accessRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []*domainUC.AppWithAccess
	for _, a := range accessList {
		app, err := uc.appRepo.FindByID(ctx, a.AppID)
		if err != nil {
			continue // skip if app not found
		}
		if !app.IsActive {
			continue // skip inactive apps
		}
		result = append(result, &domainUC.AppWithAccess{
			Application: *app,
			Role:        a.Role,
		})
	}
	return result, nil
}

func (uc *applicationUseCase) GetAppUsers(ctx context.Context, appID uuid.UUID) ([]*domainUC.UserWithAppAccess, error) {
	// Verify app exists
	if _, err := uc.appRepo.FindByID(ctx, appID); err != nil {
		return nil, err
	}
	accessList, err := uc.accessRepo.ListByApp(ctx, appID)
	if err != nil {
		return nil, err
	}
	var result []*domainUC.UserWithAppAccess
	for _, a := range accessList {
		user, err := uc.userRepo.FindByID(ctx, a.UserID)
		if err != nil {
			continue
		}
		result = append(result, &domainUC.UserWithAppAccess{
			User: *user,
			Role: a.Role,
		})
	}
	return result, nil
}

func (uc *applicationUseCase) BulkGrantAccess(ctx context.Context, req domainUC.BulkGrantAccessRequest, grantedBy uuid.UUID) error {
	// Verify app exists
	if _, err := uc.appRepo.FindByID(ctx, req.AppID); err != nil {
		return err
	}
	return uc.accessRepo.BulkGrant(ctx, req.UserIDs, req.AppID, req.Role, grantedBy)
}

func (uc *applicationUseCase) HasAccess(ctx context.Context, userID, appID uuid.UUID) (bool, error) {
	return uc.accessRepo.HasAccess(ctx, userID, appID)
}
