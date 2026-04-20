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

type releaseRecordRepository struct {
	db *pgxpool.Pool
}

func NewReleaseRecordRepository(db *pgxpool.Pool) repository.ReleaseRecordRepository {
	return &releaseRecordRepository{db: db}
}

func (r *releaseRecordRepository) Create(ctx context.Context, rr *entity.ReleaseRecord) error {
	query := `INSERT INTO release_records (id, release_id, record_id, created_at)
		VALUES ($1,$2,$3,$4)`
	_, err := r.db.Exec(ctx, query, rr.ID, rr.ReleaseID, rr.RecordID, rr.CreatedAt)
	return err
}

func (r *releaseRecordRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM release_records WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *releaseRecordRepository) DeleteByRelease(ctx context.Context, releaseID uuid.UUID) error {
	query := `DELETE FROM release_records WHERE release_id=$1`
	_, err := r.db.Exec(ctx, query, releaseID)
	return err
}

func (r *releaseRecordRepository) ListByRelease(ctx context.Context, releaseID uuid.UUID) ([]*entity.ReleaseRecord, error) {
	query := `SELECT id, release_id, record_id, created_at FROM release_records WHERE release_id=$1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, releaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*entity.ReleaseRecord
	for rows.Next() {
		rr := &entity.ReleaseRecord{}
		if err := rows.Scan(&rr.ID, &rr.ReleaseID, &rr.RecordID, &rr.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, rr)
	}
	return records, nil
}

func (r *releaseRecordRepository) FindByReleaseAndRecord(ctx context.Context, releaseID, recordID uuid.UUID) (*entity.ReleaseRecord, error) {
	query := `SELECT id, release_id, record_id, created_at FROM release_records WHERE release_id=$1 AND record_id=$2`
	row := r.db.QueryRow(ctx, query, releaseID, recordID)
	rr := &entity.ReleaseRecord{}
	err := row.Scan(&rr.ID, &rr.ReleaseID, &rr.RecordID, &rr.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return rr, nil
}

func (r *releaseRecordRepository) CountByRelease(ctx context.Context, releaseID uuid.UUID) (total int, completed int, err error) {
	query := `SELECT
		COUNT(*) AS total,
		COUNT(*) FILTER (WHERE pr.is_completed = true) AS completed
		FROM release_records rr
		JOIN project_records pr ON pr.id = rr.record_id
		WHERE rr.release_id = $1`
	err = r.db.QueryRow(ctx, query, releaseID).Scan(&total, &completed)
	return
}
