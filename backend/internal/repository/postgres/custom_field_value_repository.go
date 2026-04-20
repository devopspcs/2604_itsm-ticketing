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

type customFieldValueRepository struct {
	db *pgxpool.Pool
}

func NewCustomFieldValueRepository(db *pgxpool.Pool) repository.CustomFieldValueRepository {
	return &customFieldValueRepository{db: db}
}

func (r *customFieldValueRepository) Set(ctx context.Context, value *entity.CustomFieldValue) error {
	// First try to find existing value
	existing, err := r.GetByRecordAndField(ctx, value.RecordID, value.FieldID)
	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return err
	}

	if existing != nil {
		// Update existing value
		query := `UPDATE custom_field_values SET value=$1, updated_at=$2 WHERE record_id=$3 AND field_id=$4`
		_, err := r.db.Exec(ctx, query, value.Value, value.UpdatedAt, value.RecordID, value.FieldID)
		return err
	}

	// Insert new value
	query := `INSERT INTO custom_field_values (id, record_id, field_id, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = r.db.Exec(ctx, query,
		value.ID, value.RecordID, value.FieldID, value.Value, value.CreatedAt, value.UpdatedAt)
	return err
}

func (r *customFieldValueRepository) GetByRecord(ctx context.Context, recordID uuid.UUID) ([]*entity.CustomFieldValue, error) {
	query := `SELECT id, record_id, field_id, value, created_at, updated_at FROM custom_field_values WHERE record_id=$1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []*entity.CustomFieldValue
	for rows.Next() {
		val := &entity.CustomFieldValue{}
		if err := rows.Scan(&val.ID, &val.RecordID, &val.FieldID, &val.Value, &val.CreatedAt, &val.UpdatedAt); err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return values, nil
}

func (r *customFieldValueRepository) GetByRecordAndField(ctx context.Context, recordID uuid.UUID, fieldID uuid.UUID) (*entity.CustomFieldValue, error) {
	query := `SELECT id, record_id, field_id, value, created_at, updated_at FROM custom_field_values WHERE record_id=$1 AND field_id=$2`
	row := r.db.QueryRow(ctx, query, recordID, fieldID)
	return scanCustomFieldValue(row)
}

func (r *customFieldValueRepository) DeleteByRecord(ctx context.Context, recordID uuid.UUID) error {
	query := `DELETE FROM custom_field_values WHERE record_id=$1`
	_, err := r.db.Exec(ctx, query, recordID)
	return err
}

func (r *customFieldValueRepository) Update(ctx context.Context, value *entity.CustomFieldValue) error {
	query := `UPDATE custom_field_values SET value=$1, updated_at=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, value.Value, value.UpdatedAt, value.ID)
	return err
}

func scanCustomFieldValue(row pgx.Row) (*entity.CustomFieldValue, error) {
	val := &entity.CustomFieldValue{}
	err := row.Scan(&val.ID, &val.RecordID, &val.FieldID, &val.Value, &val.CreatedAt, &val.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return val, nil
}
