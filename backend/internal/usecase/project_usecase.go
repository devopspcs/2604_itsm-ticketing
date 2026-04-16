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
	memberRepo      repository.ProjectMemberRepository
}

func NewProjectBoardUseCase(
	projectRepo repository.ProjectRepository,
	columnRepo repository.ProjectColumnRepository,
	recordRepo repository.ProjectRecordRepository,
	activityLogRepo repository.ProjectActivityLogRepository,
	memberRepo repository.ProjectMemberRepository,
) domainUC.ProjectBoardUseCase {
	return &projectBoardUseCase{
		projectRepo:     projectRepo,
		columnRepo:      columnRepo,
		recordRepo:      recordRepo,
		activityLogRepo: activityLogRepo,
		memberRepo:      memberRepo,
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
	// Auto-add creator as owner
	_ = uc.memberRepo.Add(ctx, &entity.ProjectMember{
		ProjectID: project.ID,
		UserID:    requester.UserID,
		Role:      entity.ProjectRoleOwner,
		CreatedAt: now,
	})
	uc.logActivity(ctx, project.ID, nil, requester.UserID, "project_created", "Project created: "+project.Name)
	return project, nil
}

func (uc *projectBoardUseCase) GetProject(ctx context.Context, id uuid.UUID, requester domainUC.UserClaims) (*domainUC.ProjectDetailResponse, error) {
	project, err := uc.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	isMember, _ := uc.memberRepo.IsMember(ctx, id, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}

	cols, err := uc.columnRepo.ListByProject(ctx, id)
	if err != nil {
		return nil, err
	}

	var columnsWithRecords []domainUC.ProjectColumnWithRecords
	for _, col := range cols {
		records, err := uc.recordRepo.ListByColumn(ctx, col.ID)
		if err != nil {
			return nil, err
		}
		if records == nil {
			records = []*entity.ProjectRecord{}
		}
		// Load assignees for each record
		for _, rec := range records {
			assignees, err := uc.recordRepo.GetAssignees(ctx, rec.ID)
			if err == nil {
				rec.Assignees = assignees
			}
			if rec.Assignees == nil {
				rec.Assignees = []uuid.UUID{}
			}
		}
		columnsWithRecords = append(columnsWithRecords, domainUC.ProjectColumnWithRecords{
			ProjectColumn: *col,
			Records:       records,
		})
	}
	if columnsWithRecords == nil {
		columnsWithRecords = []domainUC.ProjectColumnWithRecords{}
	}

	return &domainUC.ProjectDetailResponse{
		Project: *project,
		Columns: columnsWithRecords,
	}, nil
}

func (uc *projectBoardUseCase) ListProjects(ctx context.Context, requester domainUC.UserClaims) ([]*entity.Project, error) {
	projectIDs, err := uc.memberRepo.ListProjectsByUser(ctx, requester.UserID)
	if err != nil {
		return nil, err
	}
	var projects []*entity.Project
	for _, pid := range projectIDs {
		p, err := uc.projectRepo.FindByID(ctx, pid)
		if err != nil {
			continue
		}
		projects = append(projects, p)
	}
	if projects == nil {
		projects = []*entity.Project{}
	}
	return projects, nil
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
	// Handle multi-assign
	if req.Assignees != nil {
		if err := uc.recordRepo.SetAssignees(ctx, recordID, req.Assignees); err != nil {
			return nil, err
		}
		record.Assignees = req.Assignees
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

// --- Members ---

func (uc *projectBoardUseCase) AddComment(ctx context.Context, projectID uuid.UUID, recordID uuid.UUID, text string, requester domainUC.UserClaims) error {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return apperror.ErrForbidden
	}
	uc.logActivity(ctx, projectID, &recordID, requester.UserID, "comment", text)
	return nil
}

func (uc *projectBoardUseCase) InviteMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, requester domainUC.UserClaims) error {
	// Only owner can invite
	role, err := uc.memberRepo.GetRole(ctx, projectID, requester.UserID)
	if err != nil {
		return apperror.ErrForbidden
	}
	if role != entity.ProjectRoleOwner {
		return apperror.ErrForbidden
	}
	member := &entity.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      entity.ProjectRoleMember,
		CreatedAt: time.Now().UTC(),
	}
	if err := uc.memberRepo.Add(ctx, member); err != nil {
		return err
	}
	uc.logActivity(ctx, projectID, nil, requester.UserID, "member_invited", "Member invited")
	return nil
}

func (uc *projectBoardUseCase) RemoveMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, requester domainUC.UserClaims) error {
	role, err := uc.memberRepo.GetRole(ctx, projectID, requester.UserID)
	if err != nil {
		return apperror.ErrForbidden
	}
	if role != entity.ProjectRoleOwner {
		return apperror.ErrForbidden
	}
	// Cannot remove owner
	targetRole, _ := uc.memberRepo.GetRole(ctx, projectID, userID)
	if targetRole == entity.ProjectRoleOwner {
		return apperror.New("CANNOT_REMOVE_OWNER", "Cannot remove project owner", 400)
	}
	if err := uc.memberRepo.Remove(ctx, projectID, userID); err != nil {
		return err
	}
	uc.logActivity(ctx, projectID, nil, requester.UserID, "member_removed", "Member removed")
	return nil
}

func (uc *projectBoardUseCase) ListMembers(ctx context.Context, projectID uuid.UUID, requester domainUC.UserClaims) ([]*entity.ProjectMember, error) {
	isMember, _ := uc.memberRepo.IsMember(ctx, projectID, requester.UserID)
	if !isMember {
		return nil, apperror.ErrForbidden
	}
	return uc.memberRepo.ListByProject(ctx, projectID)
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
