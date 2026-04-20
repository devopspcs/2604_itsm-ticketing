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

type customFieldOptionRepository struct {
	db *pgxpool.Pool
}

func NewCustomFieldOptionRepository(db *pgxpool.Pool) repository.CustomFieldOptionRepository {
	return &customFieldOptionRepository{db: db}
}

func (r *customFieldOptionRepository) Create(ctx context.Context, option *entity.CustomFieldOption) error {
	query := `INSERT INTO custom_field_options (id, field_id, option_value, option_order, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		option.ID, option.FieldID, option.OptionValue, option.OptionOrder, option.CreatedAt)
	return err
}

func (r *customFieldOptionRepository) ListByField(ctx context.Context, fieldID uuid.UUID) ([]*entity.CustomFieldOption, error) {
	query := `SELECT id, field_id, option_value, option_order, created_at FROM custom_field_options WHERE field_id=$1 ORDER BY option_order ASC`
	rows, err := r.db.Query(ctx, query, fieldID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []*entity.CustomFieldOption
	for rows.Next() {
		option := &entity.CustomFieldOption{}
		if err := rows.Scan(&option.ID, &option.FieldID, &option.OptionValue, &option.OptionOrder, &option.CreatedAt); err != nil {
			return nil, err
		}
		options = append(options, option)
	}
	return options, nil
}

func (r *customFieldOptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM custom_field_options WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func scanCustomFieldOption(row pgx.Row) (*entity.CustomFieldOption, error) {
	option := &entity.CustomFieldOption{}
	err := row.Scan(&option.ID, &option.FieldID, &option.OptionValue, &option.OptionOrder, &option.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return option, nil
}
