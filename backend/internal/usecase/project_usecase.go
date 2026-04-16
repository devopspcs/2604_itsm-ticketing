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

type projectBoardUseCase struct {
	projectRepo     repository.ProjectRepository
	columnRepo      repository.ProjectColumnRepository
	recordRepo      repository.ProjectRecordRepository
	activityLogRepo repository.ProjectActivityLogRepository
}

func NewProjectBoardUseCase(
	projectRepo repository.ProjectRepository,
	columnRepo repository.ProjectColumnRepository,
	recordRepo repository.ProjectRecordRepository,
	activityLogRepo repository.ProjectActivityLogRepository,
) domainUC.ProjectBoardUseCase {
	return &projectBoardUseCase{
		projectRepo:     projectRepo,
		columnRepo:      columnRepo,
		recordRepo:      recordRepo,
		activityLogRepo: activityLogRepo,
	}
}

// --- Project CRUD ---

func (uc *projectBoardUseCase) CreateProject(ctx context.Context, req domainUC.CreateProjectRequest, requester domainUC.UserClaims) (*entity.Project, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, apperror.ErrValidation
	}
	now := time.Now().UTC()
	iconColor := req.IconColor
	if iconColor == "" {
		iconColor = "#3b82f6"
	}
	project := &entity.Project{
		ID:        uuid.New(),
		Name:      req.Name,
		IconColor: iconColor,
		CreatedBy: requester.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}
	uc.logActivity(ctx, project.ID, nil, requester.UserID, "project_created", "Project created: "+project.Name)
	return project, nil
}

func (uc *projectBoardUseCase) GetProject(ctx context.Context, id uuid.UUID, requester domainUC.UserClaims) (*entity.Project, error) {
	project, err := uc.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if project.CreatedBy != requester.UserID {
		return nil, apperror.ErrForbidden
	}
	return project, nil
}

func (uc *projectBoardUseCase) ListProjects(ctx context.Context, requester domainUC.UserClaims) ([]*entity.Project, error) {
	return uc.projectRepo.List(ctx, requester.UserID)
}

func (uc *projectBoardUseCase) UpdateProject(ctx context.Context, id uuid.UUID, req domainUC.UpdateProjectRequest, requester domainUC.UserClaims) (*entity.Project, error) {
	project, err := uc.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if project.CreatedBy != requester.UserID {
		return nil, apperror.ErrForbidden
	}
	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.IconColor != nil {
		project.IconColor = *req.IconColor
	}
	project.UpdatedAt = time.Now().UTC()
	if err := uc.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}
	uc.logActivity(ctx, project.ID, nil, requester.UserID, "project_updated", "Project updated: "+project.Name)
	return project, nil
}

func (uc *projectBoardUseCase) DeleteProject(ctx context.Context, id uuid.UUID, requester domainUC.UserClaims) error {
	project, err := uc.projectRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if project.CreatedBy != requester.UserID {
		return apperror.ErrForbidden
	}
	return uc.projectRepo.Delete(ctx, id)
}

// --- Column CRUD ---

func (uc *projectBoardUseCase) CreateColumn(ctx context.Context, projectID uuid.UUID, req domainUC.CreateColumnRequest, requester domainUC.UserClaims) (*entity.ProjectColumn, error) {
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.CreatedBy != requester.UserID {
		return nil, apperror.ErrForbidden
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, apperror.ErrValidation
	}
	maxPos, err := uc.columnRepo.GetMaxPosition(ctx, projectID)
	if err != nil {
		return nil, err
	}
	col := &entity.ProjectColumn{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      req.Name,
		Position:  maxPos + 1,
		CreatedAt: time.Now().UTC(),
	}
	if err := uc.columnRepo.Create(ctx, col); err != nil {
		return nil, err
	}
	uc.logActivity(ctx, projectID, nil, requester.UserID, "column_created", "Column created: "+col.Name)
	return col, nil
}

func (uc *projectBoardUseCase) UpdateColumn(ctx context.Context, projectID uuid.UUID, columnID uuid.UUID, req domainUC.UpdateColumnRequest, requester domainUC.UserClaims) (*entity.ProjectColumn, error) {
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.CreatedBy != requester.UserID {
		return nil, apperror.ErrForbidden
	}
	col, err := uc.columnRepo.FindByID(ctx, columnID)
	if err != nil {
		return nil, err
	}
	col.Name = req.Name
	if err := uc.columnRepo.Update(ctx, col); err != nil {
		return nil, err
	}
	return col, nil
}

func (uc *projectBoardUseCase) DeleteColumn(ctx context.Context, projectID uuid.UUID, columnID uuid.UUID, requester domainUC.UserClaims) error {
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project.CreatedBy != requester.UserID {
		return apperror.ErrForbidden
	}
	hasRecords, err := uc.columnRepo.HasRecords(ctx, columnID)
	if err != nil {
		return err
	}
	if hasRecords {
		return apperror.New("COLUMN_HAS_RECORDS", "Column has records, please move or delete them first", 409)
	}
	col, err := uc.columnRepo.FindByID(ctx, columnID)
	if err != nil {
		return err
	}
	if err := uc.columnRepo.Delete(ctx, columnID); err != nil {
		return err
	}
	// Adjust positions of remaining columns
	cols, err := uc.columnRepo.ListByProject(ctx, projectID)
	if err != nil {
		return err
	}
	positions := make(map[uuid.UUID]int)
	for i, c := range cols {
		positions[c.ID] = i
	}
	if len(positions) > 0 {
		_ = uc.columnRepo.BulkUpdatePositions(ctx, projectID, positions)
	}
	uc.logActivity(ctx, projectID, nil, requester.UserID, "column_deleted", "Column deleted: "+col.Name)
	return nil
}

func (uc *projectBoardUseCase) ReorderColumns(ctx context.Context, projectID uuid.UUID, req domainUC.ReorderColumnsRequest, requester domainUC.UserClaims) error {
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project.CreatedBy != requester.UserID {
		return apperror.ErrForbidden
	}
	positions := make(map[uuid.UUID]int)
	for i, id := range req.ColumnIDs {
		positions[id] = i
	}
	return uc.columnRepo.BulkUpdatePositions(ctx, projectID, positions)
}

// --- Record CRUD ---

func (uc *projectBoardUseCase) CreateRecord(ctx context.Context, projectID uuid.UUID, req domainUC.CreateRecordRequest, requester domainUC.UserClaims) (*entity.ProjectRecord, error) {
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.CreatedBy != requester.UserID {
		return nil, apperror.ErrForbidden
	}
	if strings.TrimSpace(req.Title) == "" {
		return nil, apperror.ErrValidation
	}
	// Verify column exists
	_, err = uc.columnRepo.FindByID(ctx, req.ColumnID)
	if err != nil {
		return nil, err
	}
	maxPos, err := uc.recordRepo.GetMaxPosition(ctx, req.ColumnID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	record := &entity.ProjectRecord{
		ID:          uuid.New(),
		ColumnID:    req.ColumnID,
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		AssignedTo:  req.AssignedTo,
		DueDate:     req.DueDate,
		Position:    maxPos + 1,
		IsCompleted: false,
		CreatedBy:   requester.UserID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := uc.recordRepo.Create(ctx, record); err != nil {
		return nil, err
	}
	recID := record.ID
	uc.logActivity(ctx, projectID, &recID, requester.UserID, "record_created", "Record created: "+record.Title)
	return record, nil
}

func (uc *projectBoardUseCase) GetRecord(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, requester domainUC.UserClaims) (*entity.ProjectRecord, error) {
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.CreatedBy != requester.UserID {
		return nil, apperror.ErrForbidden
	}
	return uc.recordRepo.FindByID(ctx, recordID)
}

func (uc *projectBoardUseCase) UpdateRecord(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, req domainUC.UpdateRecordRequest, requester domainUC.UserClaims) (*entity.ProjectRecord, error) {
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.CreatedBy != requester.UserID {
		return nil, apperror.ErrForbidden
	}
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		record.Title = *req.Title
	}
	if req.Description != nil {
		record.Description = *req.Description
	}
	if req.AssignedTo != nil {
		record.AssignedTo = req.AssignedTo
	}
	if req.DueDate != nil {
		record.DueDate = req.DueDate
	}
	record.UpdatedAt = time.Now().UTC()
	if err := uc.recordRepo.Update(ctx, record); err != nil {
		return nil, err
	}
	uc.logActivity(ctx, projectID, &recordID, requester.UserID, "record_updated", "Record updated: "+record.Title)
	return record, nil
}

func (uc *projectBoardUseCase) DeleteRecord(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, requester domainUC.UserClaims) error {
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project.CreatedBy != requester.UserID {
		return apperror.ErrForbidden
	}
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return err
	}
	columnID := record.ColumnID
	if err := uc.recordRepo.Delete(ctx, recordID); err != nil {
		return err
	}
	// Adjust positions of remaining records in the column
	records, err := uc.recordRepo.ListByColumn(ctx, columnID)
	if err != nil {
		return err
	}
	positions := make(map[uuid.UUID]int)
	for i, r := range records {
		positions[r.ID] = i
	}
	if len(positions) > 0 {
		_ = uc.recordRepo.BulkUpdatePositions(ctx, columnID, positions)
	}
	uc.logActivity(ctx, projectID, &recordID, requester.UserID, "record_deleted", "Record deleted: "+record.Title)
	return nil
}

func (uc *projectBoardUseCase) MoveRecord(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, req domainUC.MoveRecordRequest, requester domainUC.UserClaims) error {
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project.CreatedBy != requester.UserID {
		return apperror.ErrForbidden
	}
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return err
	}
	// Verify target column exists
	_, err = uc.columnRepo.FindByID(ctx, req.TargetColumnID)
	if err != nil {
		return err
	}

	sourceColumnID := record.ColumnID
	record.ColumnID = req.TargetColumnID
	record.Position = req.Position
	record.UpdatedAt = time.Now().UTC()
	if err := uc.recordRepo.Update(ctx, record); err != nil {
		return err
	}

	// Adjust positions in source column (if different from target)
	if sourceColumnID != req.TargetColumnID {
		srcRecords, err := uc.recordRepo.ListByColumn(ctx, sourceColumnID)
		if err == nil && len(srcRecords) > 0 {
			positions := make(map[uuid.UUID]int)
			for i, r := range srcRecords {
				positions[r.ID] = i
			}
			_ = uc.recordRepo.BulkUpdatePositions(ctx, sourceColumnID, positions)
		}
	}

	// Adjust positions in target column
	tgtRecords, err := uc.recordRepo.ListByColumn(ctx, req.TargetColumnID)
	if err == nil && len(tgtRecords) > 0 {
		positions := make(map[uuid.UUID]int)
		pos := 0
		for _, r := range tgtRecords {
			if r.ID == recordID {
				positions[r.ID] = req.Position
			} else {
				if pos == req.Position {
					pos++
				}
				positions[r.ID] = pos
				pos++
			}
		}
		_ = uc.recordRepo.BulkUpdatePositions(ctx, req.TargetColumnID, positions)
	}

	uc.logActivity(ctx, projectID, &recordID, requester.UserID, "record_moved", "Record moved: "+record.Title)
	return nil
}

func (uc *projectBoardUseCase) CompleteRecord(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, requester domainUC.UserClaims) (*entity.ProjectRecord, error) {
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.CreatedBy != requester.UserID {
		return nil, apperror.ErrForbidden
	}
	record, err := uc.recordRepo.FindByID(ctx, recordID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	record.IsCompleted = true
	record.CompletedAt = &now
	record.UpdatedAt = now

	// Move to last column (highest position)
	cols, err := uc.columnRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if len(cols) > 0 {
		lastCol := cols[len(cols)-1]
		sourceColumnID := record.ColumnID
		record.ColumnID = lastCol.ID

		maxPos, err := uc.recordRepo.GetMaxPosition(ctx, lastCol.ID)
		if err != nil {
			return nil, err
		}
		record.Position = maxPos + 1

		if err := uc.recordRepo.Update(ctx, record); err != nil {
			return nil, err
		}

		// Adjust source column positions if moved
		if sourceColumnID != lastCol.ID {
			srcRecords, err := uc.recordRepo.ListByColumn(ctx, sourceColumnID)
			if err == nil && len(srcRecords) > 0 {
				positions := make(map[uuid.UUID]int)
				for i, r := range srcRecords {
					positions[r.ID] = i
				}
				_ = uc.recordRepo.BulkUpdatePositions(ctx, sourceColumnID, positions)
			}
		}
	} else {
		if err := uc.recordRepo.Update(ctx, record); err != nil {
			return nil, err
		}
	}

	uc.logActivity(ctx, projectID, &recordID, requester.UserID, "record_completed", "Record completed: "+record.Title)
	return record, nil
}

// --- Views ---

func (uc *projectBoardUseCase) GetHome(ctx context.Context, requester domainUC.UserClaims) (*domainUC.ProjectHomeData, error) {
	overdueCount, err := uc.recordRepo.CountOverdue(ctx, requester.UserID)
	if err != nil {
		return nil, err
	}
	activities, err := uc.activityLogRepo.ListByUser(ctx, requester.UserID, 10)
	if err != nil {
		return nil, err
	}
	return &domainUC.ProjectHomeData{
		OverdueCount:     overdueCount,
		RecentActivities: activities,
	}, nil
}

func (uc *projectBoardUseCase) GetCalendar(ctx context.Context, month int, year int, requester domainUC.UserClaims) ([]*entity.ProjectRecord, error) {
	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1)
	return uc.recordRepo.ListByDueDateRange(ctx, requester.UserID, from, to)
}

func (uc *projectBoardUseCase) GetActivities(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) ([]*entity.ProjectActivityLog, error) {
	project, err := uc.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.CreatedBy != requester.UserID {
		return nil, apperror.ErrForbidden
	}
	return uc.activityLogRepo.ListByProject(ctx, projectID, 50)
}

// --- Helpers ---

func (uc *projectBoardUseCase) logActivity(ctx context.Context, projectID uuid.UUID, recordID *uuid.UUID, actorID uuid.UUID, action, detail string) {
	_ = uc.activityLogRepo.Append(ctx, &entity.ProjectActivityLog{
		ID:        uuid.New(),
		ProjectID: projectID,
		RecordID:  recordID,
		ActorID:   actorID,
		Action:    action,
		Detail:    detail,
		CreatedAt: time.Now().UTC(),
	})
}
