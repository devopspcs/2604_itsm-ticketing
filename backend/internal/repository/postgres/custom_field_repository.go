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

type customFieldRepository struct {
	db *pgxpool.Pool
}

func NewCustomFieldRepository(db *pgxpool.Pool) repository.CustomFieldRepository {
	return &customFieldRepository{db: db}
}

func (r *customFieldRepository) Create(ctx context.Context, field *entity.CustomField) error {
	query := `INSERT INTO custom_fields (id, project_id, name, field_type, is_required, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query,
		field.ID, field.ProjectID, field.Name, field.FieldType,
		field.IsRequired, field.CreatedAt)
	return err
}

func (r *customFieldRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.CustomField, error) {
	query := `SELECT id, project_id, name, field_type, is_required, created_at FROM custom_fields WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanCustomField(row)
}

func (r *customFieldRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.CustomField, error) {
	query := `SELECT id, project_id, name, field_type, is_required, created_at FROM custom_fields WHERE project_id=$1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []*entity.CustomField
	for rows.Next() {
		field := &entity.CustomField{}
		if err := rows.Scan(&field.ID, &field.ProjectID, &field.Name, &field.FieldType, &field.IsRequired, &field.CreatedAt); err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func (r *customFieldRepository) Update(ctx context.Context, field *entity.CustomField) error {
	query := `UPDATE custom_fields SET name=$1, is_required=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, field.Name, field.IsRequired, field.ID)
	return err
}

func (r *customFieldRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM custom_fields WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func scanCustomField(row pgx.Row) (*entity.CustomField, error) {
	field := &entity.CustomField{}
	err := row.Scan(&field.ID, &field.ProjectID, &field.Name, &field.FieldType, &field.IsRequired, &field.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return field, nil
}
