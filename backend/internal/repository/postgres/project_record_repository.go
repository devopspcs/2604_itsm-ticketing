package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	"github.com/org/itsm/pkg/apperror"
)

type projectRecordRepository struct {
	db *pgxpool.Pool
}

func NewProjectRecordRepository(db *pgxpool.Pool) repository.ProjectRecordRepository {
	return &projectRecordRepository{db: db}
}

func (r *projectRecordRepository) Create(ctx context.Context, record *entity.ProjectRecord) error {
	query := `INSERT INTO project_records (id, column_id, project_id, title, description, assigned_to, due_date, position, is_completed, completed_at, created_by, created_at, updated_at, issue_type_id, status, parent_record_id, component_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)`
	_, err := r.db.Exec(ctx, query,
		record.ID, record.ColumnID, record.ProjectID, record.Title, record.Description,
		record.AssignedTo, record.DueDate, record.Position, record.IsCompleted, record.CompletedAt,
		record.CreatedBy, record.CreatedAt, record.UpdatedAt, record.IssueTypeID, record.Status, record.ParentRecordID, record.ComponentID)
	return err
}

func (r *projectRecordRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.ProjectRecord, error) {
	query := `SELECT id, column_id, project_id, title, description, assigned_to, due_date, position, is_completed, completed_at, created_by, created_at, updated_at, issue_type_id, status, parent_record_id, component_id FROM project_records WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanProjectRecord(row)
}

func (r *projectRecordRepository) ListByColumn(ctx context.Context, columnID uuid.UUID) ([]*entity.ProjectRecord, error) {
	query := `SELECT id, column_id, project_id, title, description, assigned_to, due_date, position, is_completed, completed_at, created_by, created_at, updated_at, issue_type_id, status, parent_record_id, component_id FROM project_records WHERE column_id=$1 ORDER BY position ASC`
	return r.queryRecords(ctx, query, columnID)
}

func (r *projectRecordRepository) ListByProject(ctx context.Context, projectID uuid.UUID, filter repository.ProjectRecordFilter) ([]*entity.ProjectRecord, error) {
	args := []interface{}{projectID}
	where := "WHERE project_id=$1"
	i := 2

	if filter.ColumnID != nil {
		where += fmt.Sprintf(" AND column_id=$%d", i)
		args = append(args, *filter.ColumnID)
		i++
	}
	if filter.AssignedTo != nil {
		where += fmt.Sprintf(" AND assigned_to=$%d", i)
		args = append(args, *filter.AssignedTo)
		i++
	}
	if filter.Search != nil && *filter.Search != "" {
		where += fmt.Sprintf(" AND title ILIKE $%d", i)
		args = append(args, "%"+*filter.Search+"%")
		i++
	}
	if filter.DueDateFrom != nil {
		where += fmt.Sprintf(" AND due_date>=$%d", i)
		args = append(args, *filter.DueDateFrom)
		i++
	}
	if filter.DueDateTo != nil {
		where += fmt.Sprintf(" AND due_date<=$%d", i)
		args = append(args, *filter.DueDateTo)
		i++
	}

	query := "SELECT id, column_id, project_id, title, description, assigned_to, due_date, position, is_completed, completed_at, created_by, created_at, updated_at, issue_type_id, status, parent_record_id, component_id FROM project_records " + where + " ORDER BY position ASC"
	return r.queryRecords(ctx, query, args...)
}

func (r *projectRecordRepository) Update(ctx context.Context, record *entity.ProjectRecord) error {
	query := `UPDATE project_records SET column_id=$1, title=$2, description=$3, assigned_to=$4, due_date=$5, position=$6, is_completed=$7, completed_at=$8, updated_at=$9, issue_type_id=$10, status=$11, parent_record_id=$12, component_id=$13 WHERE id=$14`
	_, err := r.db.Exec(ctx, query,
		record.ColumnID, record.Title, record.Description, record.AssignedTo,
		record.DueDate, record.Position, record.IsCompleted, record.CompletedAt,
		record.UpdatedAt, record.IssueTypeID, record.Status, record.ParentRecordID, record.ComponentID, record.ID)
	return err
}

func (r *projectRecordRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM project_records WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *projectRecordRepository) GetMaxPosition(ctx context.Context, columnID uuid.UUID) (int, error) {
	var maxPos *int
	query := `SELECT MAX(position) FROM project_records WHERE column_id=$1`
	if err := r.db.QueryRow(ctx, query, columnID).Scan(&maxPos); err != nil {
		return -1, err
	}
	if maxPos == nil {
		return -1, nil
	}
	return *maxPos, nil
}

func (r *projectRecordRepository) BulkUpdatePositions(ctx context.Context, columnID uuid.UUID, positions map[uuid.UUID]int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for id, pos := range positions {
		if _, err := tx.Exec(ctx, "UPDATE project_records SET position=$1 WHERE id=$2 AND column_id=$3", pos, id, columnID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *projectRecordRepository) ListByDueDateRange(ctx context.Context, createdBy uuid.UUID, from time.Time, to time.Time) ([]*entity.ProjectRecord, error) {
	query := `SELECT id, column_id, project_id, title, description, assigned_to, due_date, position, is_completed, completed_at, created_by, created_at, updated_at, issue_type_id, status, parent_record_id, component_id
		FROM project_records WHERE created_by=$1 AND due_date>=$2 AND due_date<=$3 ORDER BY due_date ASC`
	return r.queryRecords(ctx, query, createdBy, from, to)
}

func (r *projectRecordRepository) CountOverdue(ctx context.Context, createdBy uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM project_records WHERE created_by=$1 AND due_date < $2 AND is_completed=false`
	if err := r.db.QueryRow(ctx, query, createdBy, time.Now().UTC()).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *projectRecordRepository) queryRecords(ctx context.Context, query string, args ...interface{}) ([]*entity.ProjectRecord, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*entity.ProjectRecord
	for rows.Next() {
		rec := &entity.ProjectRecord{}
		if err := rows.Scan(&rec.ID, &rec.ColumnID, &rec.ProjectID, &rec.Title, &rec.Description,
			&rec.AssignedTo, &rec.DueDate, &rec.Position, &rec.IsCompleted, &rec.CompletedAt,
			&rec.CreatedBy, &rec.CreatedAt, &rec.UpdatedAt, &rec.IssueTypeID, &rec.Status, &rec.ParentRecordID, &rec.ComponentID); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

func scanProjectRecord(row pgx.Row) (*entity.ProjectRecord, error) {
	rec := &entity.ProjectRecord{}
	err := row.Scan(&rec.ID, &rec.ColumnID, &rec.ProjectID, &rec.Title, &rec.Description,
		&rec.AssignedTo, &rec.DueDate, &rec.Position, &rec.IsCompleted, &rec.CompletedAt,
		&rec.CreatedBy, &rec.CreatedAt, &rec.UpdatedAt, &rec.IssueTypeID, &rec.Status, &rec.ParentRecordID, &rec.ComponentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return rec, nil
}

func (r *projectRecordRepository) SetAssignees(ctx context.Context, recordID uuid.UUID, userIDs []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete existing assignees
	if _, err := tx.Exec(ctx, "DELETE FROM project_record_assignees WHERE record_id=$1", recordID); err != nil {
		return err
	}

	// Insert new assignees
	for _, uid := range userIDs {
		if _, err := tx.Exec(ctx, "INSERT INTO project_record_assignees (record_id, user_id) VALUES ($1, $2)", recordID, uid); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *projectRecordRepository) GetAssignees(ctx context.Context, recordID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, "SELECT user_id FROM project_record_assignees WHERE record_id=$1 ORDER BY created_at ASC", recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if ids == nil {
		ids = []uuid.UUID{}
	}
	return ids, nil
}

func (r *projectRecordRepository) CountByProject(ctx context.Context, projectID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM project_records WHERE project_id=$1`
	if err := r.db.QueryRow(ctx, query, projectID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *projectRecordRepository) CountByProjectAndStatus(ctx context.Context, projectID uuid.UUID, isCompleted bool) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM project_records WHERE project_id=$1 AND is_completed=$2`
	if err := r.db.QueryRow(ctx, query, projectID, isCompleted).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *projectRecordRepository) CountByProjectGroupedByStatus(ctx context.Context, projectID uuid.UUID) (map[string]int, error) {
	query := `SELECT COALESCE(ws.status_name, 'No Status') AS sname, COUNT(*) AS cnt
		FROM project_records pr
		LEFT JOIN workflow_statuses ws ON ws.id::text = pr.status
		WHERE pr.project_id=$1
		GROUP BY sname`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var name string
		var cnt int
		if err := rows.Scan(&name, &cnt); err != nil {
			return nil, err
		}
		result[name] = cnt
	}
	return result, nil
}

func (r *projectRecordRepository) ListByProjectPaginated(ctx context.Context, projectID uuid.UUID, filter repository.IssuesFilter) (*repository.PaginatedRecords, error) {
	args := []interface{}{projectID}
	where := "WHERE pr.project_id=$1"
	i := 2

	if filter.Search != nil && *filter.Search != "" {
		where += fmt.Sprintf(" AND (pr.title ILIKE $%d OR pr.description ILIKE $%d)", i, i)
		args = append(args, "%"+*filter.Search+"%")
		i++
	}
	if filter.StatusID != nil {
		where += fmt.Sprintf(" AND pr.status=$%d", i)
		args = append(args, filter.StatusID.String())
		i++
	}
	if filter.AssigneeID != nil {
		where += fmt.Sprintf(" AND pr.assigned_to=$%d", i)
		args = append(args, *filter.AssigneeID)
		i++
	}
	if filter.IssueType != nil {
		where += fmt.Sprintf(" AND pr.issue_type_id=$%d", i)
		args = append(args, *filter.IssueType)
		i++
	}
	if filter.LabelID != nil {
		where += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM record_labels rl WHERE rl.record_id=pr.id AND rl.label_id=$%d)", i)
		args = append(args, *filter.LabelID)
		i++
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM project_records pr " + where
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	// Paginated query
	offset := (filter.Page - 1) * filter.PageSize
	dataQuery := fmt.Sprintf("SELECT pr.id, pr.column_id, pr.project_id, pr.title, pr.description, pr.assigned_to, pr.due_date, pr.position, pr.is_completed, pr.completed_at, pr.created_by, pr.created_at, pr.updated_at, pr.issue_type_id, pr.status, pr.parent_record_id, pr.component_id FROM project_records pr %s ORDER BY pr.created_at DESC LIMIT $%d OFFSET $%d", where, i, i+1)
	args = append(args, filter.PageSize, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*entity.ProjectRecord
	for rows.Next() {
		rec := &entity.ProjectRecord{}
		if err := rows.Scan(&rec.ID, &rec.ColumnID, &rec.ProjectID, &rec.Title, &rec.Description,
			&rec.AssignedTo, &rec.DueDate, &rec.Position, &rec.IsCompleted, &rec.CompletedAt,
			&rec.CreatedBy, &rec.CreatedAt, &rec.UpdatedAt, &rec.IssueTypeID, &rec.Status, &rec.ParentRecordID, &rec.ComponentID); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return &repository.PaginatedRecords{
		Records:  records,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

func (r *projectRecordRepository) SetComponentID(ctx context.Context, recordID uuid.UUID, componentID *uuid.UUID) error {
	query := `UPDATE project_records SET component_id=$1, updated_at=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, componentID, time.Now().UTC(), recordID)
	return err
}

func (r *projectRecordRepository) ClearComponentByComponentID(ctx context.Context, componentID uuid.UUID) error {
	query := `UPDATE project_records SET component_id=NULL, updated_at=$1 WHERE component_id=$2`
	_, err := r.db.Exec(ctx, query, time.Now().UTC(), componentID)
	return err
}
