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

type componentUseCase struct {
	componentRepo repository.ComponentRepository
	recordRepo    repository.ProjectRecordRepository
	memberRepo    repository.ProjectMemberRepository
}

func NewComponentUseCase(
	componentRepo repository.ComponentRepository,
	recordRepo repository.ProjectRecordRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.ComponentUseCase {
	return &componentUseCase{
		componentRepo: componentRepo,
		recordRepo:    recordRepo,
		memberRepo:    memberRepo,
	}
}

func (uc *componentUseCase) Create(ctx context.Context, projectID uuid.UUID, req domainUC.CreateComponentRequest, requester domainUC.UserClaims) (*entity.Component, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, apperror.ErrValidation
	}
	now := time.Now().UTC()
	comp := &entity.Component{
		ID:          uuid.New(),
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		LeadUserID:  req.LeadUserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.componentRepo.Create(ctx, comp); err != nil {
		return nil, err
	}
	return comp, nil
}

func (uc *componentUseCase) List(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) ([]*domainUC.ComponentWithCount, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}
	components, err := uc.componentRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	var result []*domainUC.ComponentWithCount
	for _, comp := range components {
		count, err := uc.componentRepo.CountRecords(ctx, comp.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, &domainUC.ComponentWithCount{
			Component:   *comp,
			RecordCount: count,
		})
	}
	if result == nil {
		result = []*domainUC.ComponentWithCount{}
	}
	return result, nil
}

func (uc *componentUseCase) Get(ctx context.Context, projectID uuid.UUID, componentID uuid.UUID, requester domainUC.UserClaims) (*domainUC.ComponentWithCount, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}
	comp, err := uc.componentRepo.FindByID(ctx, componentID)
	if err != nil {
		return nil, err
	}
	count, err := uc.componentRepo.CountRecords(ctx, componentID)
	if err != nil {
		return nil, err
	}
	return &domainUC.ComponentWithCount{
		Component:   *comp,
		RecordCount: count,
	}, nil
}

func (uc *componentUseCase) Update(ctx context.Context, projectID uuid.UUID, componentID uuid.UUID, req domainUC.UpdateComponentRequest, requester domainUC.UserClaims) (*entity.Component, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}
	comp, err := uc.componentRepo.FindByID(ctx, componentID)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		comp.Name = *req.Name
	}
	if req.Description != nil {
		comp.Description = *req.Description
	}
	if req.LeadUserID != nil {
		comp.LeadUserID = req.LeadUserID
	}
	comp.UpdatedAt = time.Now().UTC()
	if err := uc.componentRepo.Update(ctx, comp); err != nil {
		return nil, err
	}
	return comp, nil
}

func (uc *componentUseCase) Delete(ctx context.Context, projectID uuid.UUID, componentID uuid.UUID, requester domainUC.UserClaims) error {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return apperror.ErrForbidden
	}
	// Nullify records first
	if err := uc.recordRepo.ClearComponentByComponentID(ctx, componentID); err != nil {
		return err
	}
	return uc.componentRepo.Delete(ctx, componentID)
}

func (uc *componentUseCase) AssignRecord(ctx context.Context, projectID uuid.UUID, componentID uuid.UUID, recordID uuid.UUID, requester domainUC.UserClaims) error {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return apperror.ErrForbidden
	}
	return uc.recordRepo.SetComponentID(ctx, recordID, &componentID)
}

func (uc *componentUseCase) RemoveRecord(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, requester domainUC.UserClaims) error {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return apperror.ErrForbidden
	}
	return uc.recordRepo.SetComponentID(ctx, recordID, nil)
}

func (uc *componentUseCase) ListRecords(ctx context.Context, projectID uuid.UUID, componentID uuid.UUID, requester domainUC.UserClaims) ([]*entity.ProjectRecord, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}
	cid := componentID
	filter := repository.ProjectRecordFilter{
		ProjectID: &projectID,
	}
	// Use ListByProject with a filter — but we need to filter by component_id
	// The existing filter doesn't support component_id, so we use ListByProjectPaginated with a large page
	// Actually, let's use the existing ListByProject and filter in-memory for now
	records, err := uc.recordRepo.ListByProject(ctx, projectID, filter)
	if err != nil {
		return nil, err
	}
	var result []*entity.ProjectRecord
	for _, rec := range records {
		if rec.ComponentID != nil && *rec.ComponentID == cid {
			result = append(result, rec)
		}
	}
	if result == nil {
		result = []*entity.ProjectRecord{}
	}
	return result, nil
}
