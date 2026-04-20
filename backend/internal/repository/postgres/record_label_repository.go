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

type recordLabelRepository struct {
	db *pgxpool.Pool
}

func NewRecordLabelRepository(db *pgxpool.Pool) repository.RecordLabelRepository {
	return &recordLabelRepository{db: db}
}

func (r *recordLabelRepository) Create(ctx context.Context, rl *entity.RecordLabel) error {
	query := `INSERT INTO record_labels (id, record_id, label_id, created_at)
		VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query,
		rl.ID, rl.RecordID, rl.LabelID, rl.CreatedAt)
	return err
}

func (r *recordLabelRepository) ListByRecord(ctx context.Context, recordID uuid.UUID) ([]*entity.RecordLabel, error) {
	query := `SELECT id, record_id, label_id, created_at FROM record_labels WHERE record_id=$1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []*entity.RecordLabel
	for rows.Next() {
		rl := &entity.RecordLabel{}
		if err := rows.Scan(&rl.ID, &rl.RecordID, &rl.LabelID, &rl.CreatedAt); err != nil {
			return nil, err
		}
		labels = append(labels, rl)
	}
	return labels, nil
}

func (r *recordLabelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM record_labels WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *recordLabelRepository) DeleteByRecord(ctx context.Context, recordID uuid.UUID) error {
	query := `DELETE FROM record_labels WHERE record_id=$1`
	_, err := r.db.Exec(ctx, query, recordID)
	return err
}

func scanRecordLabel(row pgx.Row) (*entity.RecordLabel, error) {
	rl := &entity.RecordLabel{}
	err := row.Scan(&rl.ID, &rl.RecordID, &rl.LabelID, &rl.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return rl, nil
}
