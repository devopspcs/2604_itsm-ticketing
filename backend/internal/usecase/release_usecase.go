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

type releaseUseCase struct {
	releaseRepo       repository.ReleaseRepository
	releaseRecordRepo repository.ReleaseRecordRepository
	memberRepo        repository.ProjectMemberRepository
}

func NewReleaseUseCase(
	releaseRepo repository.ReleaseRepository,
	releaseRecordRepo repository.ReleaseRecordRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.ReleaseUseCase {
	return &releaseUseCase{
		releaseRepo:       releaseRepo,
		releaseRecordRepo: releaseRecordRepo,
		memberRepo:        memberRepo,
	}
}

func (uc *releaseUseCase) Create(ctx context.Context, projectID uuid.UUID, req domainUC.CreateReleaseRequest, requester domainUC.UserClaims) (*entity.Release, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Version) == "" {
		return nil, apperror.ErrValidation
	}
	status := req.Status
	if status == "" {
		status = "Planning"
	}
	now := time.Now().UTC()
	release := &entity.Release{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        req.Name,
		Version:     req.Version,
		Description: req.Description,
		StartDate:   req.StartDate,
		ReleaseDate: req.ReleaseDate,
		Status:      status,
		CreatedBy:   requester.UserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.releaseRepo.Create(ctx, release); err != nil {
		return nil, err
	}
	return release, nil
}

func (uc *releaseUseCase) List(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) ([]*domainUC.ReleaseWithProgress, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}
	releases, err := uc.releaseRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	var result []*domainUC.ReleaseWithProgress
	for _, rel := range releases {
		total, completed, err := uc.releaseRecordRepo.CountByRelease(ctx, rel.ID)
		if err != nil {
			return nil, err
		}
		progress := 0.0
		if total > 0 {
			progress = (float64(completed) / float64(total)) * 100
		}
		result = append(result, &domainUC.ReleaseWithProgress{
			Release:         *rel,
			TotalRecords:    total,
			CompletedCount:  completed,
			ProgressPercent: progress,
		})
	}
	if result == nil {
		result = []*domainUC.ReleaseWithProgress{}
	}
	return result, nil
}

func (uc *releaseUseCase) Get(ctx context.Context, projectID uuid.UUID, releaseID uuid.UUID, requester domainUC.UserClaims) (*domainUC.ReleaseWithProgress, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}
	rel, err := uc.releaseRepo.FindByID(ctx, releaseID)
	if err != nil {
		return nil, err
	}
	total, completed, err := uc.releaseRecordRepo.CountByRelease(ctx, releaseID)
	if err != nil {
		return nil, err
	}
	progress := 0.0
	if total > 0 {
		progress = (float64(completed) / float64(total)) * 100
	}
	return &domainUC.ReleaseWithProgress{
		Release:         *rel,
		TotalRecords:    total,
		CompletedCount:  completed,
		ProgressPercent: progress,
	}, nil
}

func (uc *releaseUseCase) Update(ctx context.Context, projectID uuid.UUID, releaseID uuid.UUID, req domainUC.UpdateReleaseRequest, requester domainUC.UserClaims) (*entity.Release, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}
	rel, err := uc.releaseRepo.FindByID(ctx, releaseID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		rel.Name = *req.Name
	}
	if req.Version != nil {
		rel.Version = *req.Version
	}
	if req.Description != nil {
		rel.Description = *req.Description
	}
	if req.StartDate != nil {
		rel.StartDate = req.StartDate
	}
	if req.ReleaseDate != nil {
		rel.ReleaseDate = req.ReleaseDate
	}
	if req.Status != nil {
		rel.Status = *req.Status
	}
	rel.UpdatedAt = time.Now().UTC()
	if err := uc.releaseRepo.Update(ctx, rel); err != nil {
		return nil, err
	}
	return rel, nil
}

func (uc *releaseUseCase) Delete(ctx context.Context, projectID uuid.UUID, releaseID uuid.UUID, requester domainUC.UserClaims) error {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return apperror.ErrForbidden
	}
	// Cascade: delete all release_records first
	if err := uc.releaseRecordRepo.DeleteByRelease(ctx, releaseID); err != nil {
		return err
	}
	return uc.releaseRepo.Delete(ctx, releaseID)
}

func (uc *releaseUseCase) AssignRecord(ctx context.Context, projectID uuid.UUID, releaseID uuid.UUID, recordID uuid.UUID, requester domainUC.UserClaims) error {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return apperror.ErrForbidden
	}
	// Check if already assigned
	existing, err := uc.releaseRecordRepo.FindByReleaseAndRecord(ctx, releaseID, recordID)
	if err == nil && existing != nil {
		return apperror.New("DUPLICATE_ASSOCIATION", "Record is already assigned to this release", 409)
	}
	rr := &entity.ReleaseRecord{
		ID:        uuid.New(),
		ReleaseID: releaseID,
		RecordID:  recordID,
		CreatedAt: time.Now().UTC(),
	}
	return uc.releaseRecordRepo.Create(ctx, rr)
}

func (uc *releaseUseCase) RemoveRecord(ctx context.Context, projectID uuid.UUID, releaseID uuid.UUID, recordID uuid.UUID, requester domainUC.UserClaims) error {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return apperror.ErrForbidden
	}
	rr, err := uc.releaseRecordRepo.FindByReleaseAndRecord(ctx, releaseID, recordID)
	if err != nil {
		return err
	}
	return uc.releaseRecordRepo.Delete(ctx, rr.ID)
}
