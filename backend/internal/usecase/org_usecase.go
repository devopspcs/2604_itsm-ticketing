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

type orgUseCase struct {
	divRepo  repository.DivisionRepository
	deptRepo repository.DepartmentRepository
	teamRepo repository.TeamRepository
}

func NewOrgUseCase(
	deptRepo repository.DepartmentRepository,
	divRepo repository.DivisionRepository,
	teamRepo repository.TeamRepository,
) domainUC.OrgUseCase {
	return &orgUseCase{
		deptRepo: deptRepo,
		divRepo:  divRepo,
		teamRepo: teamRepo,
	}
}

// --- Division (top level) ---

func (uc *orgUseCase) CreateDivision(ctx context.Context, req domainUC.CreateDivisionRequest) (*entity.Division, error) {
	if req.Name == "" || req.Code == "" {
		return nil, apperror.ErrValidation
	}
	now := time.Now().UTC()
	div := &entity.Division{
		ID:        uuid.New(),
		Name:      req.Name,
		Code:      req.Code,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.divRepo.Create(ctx, div); err != nil {
		return nil, err
	}
	return div, nil
}

func (uc *orgUseCase) UpdateDivision(ctx context.Context, id uuid.UUID, req domainUC.UpdateDivisionRequest) (*entity.Division, error) {
	div, err := uc.divRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name == "" || req.Code == "" {
		return nil, apperror.ErrValidation
	}
	div.Name = req.Name
	div.Code = req.Code
	if err := uc.divRepo.Update(ctx, div); err != nil {
		return nil, err
	}
	div.UpdatedAt = time.Now().UTC()
	return div, nil
}

func (uc *orgUseCase) DeleteDivision(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.divRepo.FindByID(ctx, id); err != nil {
		return err
	}
	has, err := uc.divRepo.HasDepartments(ctx, id)
	if err != nil {
		return err
	}
	if has {
		return apperror.ErrHasDependencies
	}
	return uc.divRepo.Delete(ctx, id)
}

func (uc *orgUseCase) ListDivisions(ctx context.Context) ([]*entity.Division, error) {
	return uc.divRepo.List(ctx)
}

// --- Department (child of Division) ---

func (uc *orgUseCase) CreateDepartment(ctx context.Context, req domainUC.CreateDepartmentRequest) (*entity.Department, error) {
	if req.Name == "" || req.Code == "" {
		return nil, apperror.ErrValidation
	}
	// Validate division exists
	if _, err := uc.divRepo.FindByID(ctx, req.DivisionID); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	dept := &entity.Department{
		ID:         uuid.New(),
		DivisionID: req.DivisionID,
		Name:       req.Name,
		Code:       req.Code,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := uc.deptRepo.Create(ctx, dept); err != nil {
		return nil, err
	}
	return dept, nil
}

func (uc *orgUseCase) UpdateDepartment(ctx context.Context, id uuid.UUID, req domainUC.UpdateDepartmentRequest) (*entity.Department, error) {
	dept, err := uc.deptRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name == "" || req.Code == "" {
		return nil, apperror.ErrValidation
	}
	if _, err := uc.divRepo.FindByID(ctx, req.DivisionID); err != nil {
		return nil, err
	}
	dept.DivisionID = req.DivisionID
	dept.Name = req.Name
	dept.Code = req.Code
	if err := uc.deptRepo.Update(ctx, dept); err != nil {
		return nil, err
	}
	dept.UpdatedAt = time.Now().UTC()
	return dept, nil
}

func (uc *orgUseCase) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.deptRepo.FindByID(ctx, id); err != nil {
		return err
	}
	has, err := uc.deptRepo.HasTeamsOrUsers(ctx, id)
	if err != nil {
		return err
	}
	if has {
		return apperror.ErrHasDependencies
	}
	return uc.deptRepo.Delete(ctx, id)
}

func (uc *orgUseCase) ListDepartments(ctx context.Context, divisionID *uuid.UUID) ([]*entity.Department, error) {
	return uc.deptRepo.List(ctx, repository.DepartmentFilter{DivisionID: divisionID})
}

// --- Team (child of Department) ---

func (uc *orgUseCase) CreateTeam(ctx context.Context, req domainUC.CreateTeamRequest) (*entity.Team, error) {
	if req.Name == "" {
		return nil, apperror.ErrValidation
	}
	if _, err := uc.deptRepo.FindByID(ctx, req.DepartmentID); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	team := &entity.Team{
		ID:           uuid.New(),
		DepartmentID: req.DepartmentID,
		Name:         req.Name,
		Email:        req.Email,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := uc.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}
	return team, nil
}

func (uc *orgUseCase) UpdateTeam(ctx context.Context, id uuid.UUID, req domainUC.UpdateTeamRequest) (*entity.Team, error) {
	team, err := uc.teamRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name == "" {
		return nil, apperror.ErrValidation
	}
	if _, err := uc.deptRepo.FindByID(ctx, req.DepartmentID); err != nil {
		return nil, err
	}
	team.DepartmentID = req.DepartmentID
	team.Name = req.Name
	team.Email = req.Email
	if err := uc.teamRepo.Update(ctx, team); err != nil {
		return nil, err
	}
	team.UpdatedAt = time.Now().UTC()
	return team, nil
}

func (uc *orgUseCase) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	if _, err := uc.teamRepo.FindByID(ctx, id); err != nil {
		return err
	}
	has, err := uc.teamRepo.HasUsers(ctx, id)
	if err != nil {
		return err
	}
	if has {
		return apperror.ErrHasDependencies
	}
	return uc.teamRepo.Delete(ctx, id)
}

func (uc *orgUseCase) ListTeams(ctx context.Context, departmentID *uuid.UUID) ([]*entity.Team, error) {
	return uc.teamRepo.List(ctx, repository.TeamFilter{DepartmentID: departmentID})
}
