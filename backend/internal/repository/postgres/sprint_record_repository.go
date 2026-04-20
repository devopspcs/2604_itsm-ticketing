package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	"github.com/org/itsm/pkg/apperror"
)

type sprintRecordRepository struct {
	db *pgxpool.Pool
}

func NewSprintRecordRepository(db *pgxpool.Pool) repository.SprintRecordRepository {
	return &sprintRecordRepository{db: db}
}

func (r *sprintRecordRepository) Create(ctx context.Context, sr *entity.SprintRecord) error {
	query := `INSERT INTO sprint_records (id, sprint_id, record_id, priority, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		sr.ID, sr.SprintID, sr.RecordID, sr.Priority, sr.CreatedAt)
	return err
}

func (r *sprintRecordRepository) ListBySprint(ctx context.Context, sprintID uuid.UUID) ([]*entity.SprintRecord, error) {
	query := `SELECT id, sprint_id, record_id, priority, created_at FROM sprint_records WHERE sprint_id=$1 ORDER BY priority ASC`
	rows, err := r.db.Query(ctx, query, sprintID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*entity.SprintRecord
	for rows.Next() {
		sr := &entity.SprintRecord{}
		if err := rows.Scan(&sr.ID, &sr.SprintID, &sr.RecordID, &sr.Priority, &sr.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, sr)
	}
	return records, nil
}

func (r *sprintRecordRepository) FindByRecord(ctx context.Context, recordID uuid.UUID) (*entity.SprintRecord, error) {
	query := `SELECT id, sprint_id, record_id, priority, created_at FROM sprint_records WHERE record_id=$1`
	row := r.db.QueryRow(ctx, query, recordID)
	return scanSprintRecord(row)
}

func (r *sprintRecordRepository) Update(ctx context.Context, sr *entity.SprintRecord) error {
	query := `UPDATE sprint_records SET priority=$1 WHERE id=$2`
	_, err := r.db.Exec(ctx, query, sr.Priority, sr.ID)
	return err
}

func (r *sprintRecordRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sprint_records WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *sprintRecordRepository) BulkAssign(ctx context.Context, sprintID uuid.UUID, recordIDs []uuid.UUID) error {
	if len(recordIDs) == 0 {
		return nil
	}

	query := `INSERT INTO sprint_records (id, sprint_id, record_id, priority, created_at)
		SELECT gen_random_uuid(), $1, unnest($2::uuid[]), row_number() OVER (ORDER BY unnest), NOW()`
	_, err := r.db.Exec(ctx, query, sprintID, recordIDs)
	return err
}

func (r *sprintRecordRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.SprintRecord, error) {
	query := `SELECT sr.id, sr.sprint_id, sr.record_id, sr.priority, sr.created_at 
              FROM sprint_records sr
              JOIN sprints s ON sr.sprint_id = s.id
              WHERE s.project_id=$1
              ORDER BY sr.priority ASC`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*entity.SprintRecord
	for rows.Next() {
		sr := &entity.SprintRecord{}
		if err := rows.Scan(&sr.ID, &sr.SprintID, &sr.RecordID, &sr.Priority, &sr.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, sr)
	}
	return records, nil
}

func scanSprintRecord(row pgx.Row) (*entity.SprintRecord, error) {
	sr := &entity.SprintRecord{}
	err := row.Scan(&sr.ID, &sr.SprintID, &sr.RecordID, &sr.Priority, &sr.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return sr, nil
}
