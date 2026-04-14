package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	"github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
	"github.com/org/itsm/pkg/password"
)

type userUseCase struct {
	userRepo repository.UserRepository
	divRepo  repository.DivisionRepository
	teamRepo repository.TeamRepository
}

func NewUserUseCase(
	userRepo repository.UserRepository,
	divRepo repository.DivisionRepository,
	teamRepo repository.TeamRepository,
) usecase.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		divRepo:  divRepo,
		teamRepo: teamRepo,
	}
}

// CreateUser hashes the password, checks for duplicate email, validates org hierarchy, and saves the new user.
func (u *userUseCase) CreateUser(ctx context.Context, req usecase.CreateUserRequest) (*entity.User, error) {
	existing, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, apperror.ErrConflict
	}

	hashed, err := password.Hash(req.Password)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	// Validate org hierarchy if provided
	if err := u.validateOrgHierarchy(ctx, req.DepartmentID, req.DivisionID, req.TeamID); err != nil {
		return nil, err
	}

	now := time.Now()
	user := &entity.User{
		ID:           uuid.New(),
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: hashed,
		Role:         req.Role,
		IsActive:     true,
		DepartmentID: req.DepartmentID,
		DivisionID:   req.DivisionID,
		TeamID:       req.TeamID,
		Position:     req.Position,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, apperror.ErrInternal
	}

	return user, nil
}

// UpdateUserRole updates the role of an existing user.
func (u *userUseCase) UpdateUserRole(ctx context.Context, userID uuid.UUID, role entity.Role) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return apperror.ErrNotFound
	}

	user.Role = role
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(ctx, user); err != nil {
		return apperror.ErrInternal
	}

	return nil
}

// DeactivateUser sets a user's IsActive flag to false.
func (u *userUseCase) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return apperror.ErrNotFound
	}

	user.IsActive = false
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(ctx, user); err != nil {
		return apperror.ErrInternal
	}

	return nil
}

// ActivateUser sets a user's IsActive flag to true.
func (u *userUseCase) ActivateUser(ctx context.Context, userID uuid.UUID) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return apperror.ErrNotFound
	}

	user.IsActive = true
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(ctx, user); err != nil {
		return apperror.ErrInternal
	}

	return nil
}

// GetUsers returns a list of users.
func (u *userUseCase) GetUsers(ctx context.Context) ([]*entity.User, error) {
	users, err := u.userRepo.List(ctx, repository.UserFilter{})
	if err != nil {
		return nil, apperror.ErrInternal
	}

	return users, nil
}

// UpdateUserOrg updates the org assignment for a user with hierarchy validation.
func (u *userUseCase) UpdateUserOrg(ctx context.Context, userID uuid.UUID, req usecase.UpdateUserOrgRequest) (*entity.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.ErrNotFound
	}

	if err := u.validateOrgHierarchy(ctx, req.DepartmentID, req.DivisionID, req.TeamID); err != nil {
		return nil, err
	}

	user.DepartmentID = req.DepartmentID
	user.DivisionID = req.DivisionID
	user.TeamID = req.TeamID
	user.Position = req.Position
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, apperror.ErrInternal
	}

	return user, nil
}

// validateOrgHierarchy checks that division belongs to department and team belongs to division.
func (u *userUseCase) validateOrgHierarchy(ctx context.Context, deptID, divID, teamID *uuid.UUID) error {
	if divID != nil {
		div, err := u.divRepo.FindByID(ctx, *divID)
		if err != nil {
			return apperror.ErrNotFound
		}
		if deptID == nil || div.DepartmentID != *deptID {
			return apperror.ErrInvalidHierarchy
		}
	}
	if teamID != nil {
		team, err := u.teamRepo.FindByID(ctx, *teamID)
		if err != nil {
			return apperror.ErrNotFound
		}
		if divID == nil || team.DivisionID != *divID {
			return apperror.ErrInvalidHierarchy
		}
	}
	return nil
}
