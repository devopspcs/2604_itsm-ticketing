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
	query := `INSERT INTO project_records (id, column_id, project_id, title, description, assigned_to, due_date, position, is_completed, completed_at, created_by, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := r.db.Exec(ctx, query,
		record.ID, record.ColumnID, record.ProjectID, record.Title, record.Description,
		record.AssignedTo, record.DueDate, record.Position, record.IsCompleted, record.CompletedAt,
		record.CreatedBy, record.CreatedAt, record.UpdatedAt)
	return err
}

func (r *projectRecordRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.ProjectRecord, error) {
	query := `SELECT id, column_id, project_id, title, description, assigned_to, due_date, position, is_completed, completed_at, created_by, created_at, updated_at FROM project_records WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanProjectRecord(row)
}

func (r *projectRecordRepository) ListByColumn(ctx context.Context, columnID uuid.UUID) ([]*entity.ProjectRecord, error) {
	query := `SELECT id, column_id, project_id, title, description, assigned_to, due_date, position, is_completed, completed_at, created_by, created_at, updated_at FROM project_records WHERE column_id=$1 ORDER BY position ASC`
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

	query := "SELECT id, column_id, project_id, title, description, assigned_to, due_date, position, is_completed, completed_at, created_by, created_at, updated_at FROM project_records " + where + " ORDER BY position ASC"
	return r.queryRecords(ctx, query, args...)
}

func (r *projectRecordRepository) Update(ctx context.Context, record *entity.ProjectRecord) error {
	query := `UPDATE project_records SET column_id=$1, title=$2, description=$3, assigned_to=$4, due_date=$5, position=$6, is_completed=$7, completed_at=$8, updated_at=$9 WHERE id=$10`
	_, err := r.db.Exec(ctx, query,
		record.ColumnID, record.Title, record.Description, record.AssignedTo,
		record.DueDate, record.Position, record.IsCompleted, record.CompletedAt,
		record.UpdatedAt, record.ID)
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
	query := `SELECT id, column_id, project_id, title, description, assigned_to, due_date, position, is_completed, completed_at, created_by, created_at, updated_at
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
			&rec.CreatedBy, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
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
		&rec.CreatedBy, &rec.CreatedAt, &rec.UpdatedAt)
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
