package usecase

import (
	"context"
	"sort"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
)

type reportsUseCase struct {
	recordRepo       repository.ProjectRecordRepository
	memberRepo       repository.ProjectMemberRepository
	sprintRepo       repository.SprintRepository
	sprintRecordRepo repository.SprintRecordRepository
}

func NewReportsUseCase(
	recordRepo repository.ProjectRecordRepository,
	memberRepo repository.ProjectMemberRepository,
	sprintRepo repository.SprintRepository,
	sprintRecordRepo repository.SprintRecordRepository,
) domainUC.ReportsUseCase {
	return &reportsUseCase{
		recordRepo:       recordRepo,
		memberRepo:       memberRepo,
		sprintRepo:       sprintRepo,
		sprintRecordRepo: sprintRecordRepo,
	}
}

func (uc *reportsUseCase) GetSummary(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) (*domainUC.ReportsSummary, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	total, err := uc.recordRepo.CountByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	completed, err := uc.recordRepo.CountByProjectAndStatus(ctx, projectID, true)
	if err != nil {
		return nil, err
	}
	open := total - completed

	byStatus, err := uc.recordRepo.CountByProjectGroupedByStatus(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if byStatus == nil {
		byStatus = make(map[string]int)
	}

	return &domainUC.ReportsSummary{
		TotalRecords:   total,
		CompletedCount: completed,
		OpenCount:      open,
		ByStatus:       byStatus,
	}, nil
}

func (uc *reportsUseCase) GetVelocity(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) ([]*domainUC.VelocityDataPoint, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	sprints, err := uc.sprintRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Filter to completed sprints and sort by start date ASC
	var completedSprints []*entity.Sprint
	for _, s := range sprints {
		if s.Status == "Completed" {
			completedSprints = append(completedSprints, s)
		}
	}
	sort.Slice(completedSprints, func(i, j int) bool {
		si := completedSprints[i].StartDate
		sj := completedSprints[j].StartDate
		if si == nil && sj == nil {
			return false
		}
		if si == nil {
			return true
		}
		if sj == nil {
			return false
		}
		return si.Before(*sj)
	})

	var result []*domainUC.VelocityDataPoint
	for _, s := range completedSprints {
		sprintRecords, err := uc.sprintRecordRepo.ListBySprint(ctx, s.ID)
		if err != nil {
			return nil, err
		}
		totalRecords := len(sprintRecords)
		completedCount := 0
		for _, sr := range sprintRecords {
			rec, err := uc.recordRepo.FindByID(ctx, sr.RecordID)
			if err != nil {
				continue
			}
			if rec.IsCompleted {
				completedCount++
			}
		}
		result = append(result, &domainUC.VelocityDataPoint{
			SprintName:     s.Name,
			TotalRecords:   totalRecords,
			CompletedCount: completedCount,
		})
	}

	if result == nil {
		result = []*domainUC.VelocityDataPoint{}
	}
	return result, nil
}

func (uc *reportsUseCase) GetBurndown(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) (*domainUC.BurndownData, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	activeSprints, err := uc.sprintRepo.ListActive(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if len(activeSprints) == 0 {
		return &domainUC.BurndownData{HasActive: false}, nil
	}

	active := activeSprints[0]
	sprintRecords, err := uc.sprintRecordRepo.ListBySprint(ctx, active.ID)
	if err != nil {
		return nil, err
	}

	totalCount := len(sprintRecords)
	doneCount := 0
	for _, sr := range sprintRecords {
		rec, err := uc.recordRepo.FindByID(ctx, sr.RecordID)
		if err != nil {
			continue
		}
		if rec.IsCompleted {
			doneCount++
		}
	}

	return &domainUC.BurndownData{
		SprintName: active.Name,
		StartDate:  active.StartDate,
		EndDate:    active.EndDate,
		TotalCount: totalCount,
		DoneCount:  doneCount,
		HasActive:  true,
	}, nil
}
