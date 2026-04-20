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

type componentRepository struct {
	db *pgxpool.Pool
}

func NewComponentRepository(db *pgxpool.Pool) repository.ComponentRepository {
	return &componentRepository{db: db}
}

func (r *componentRepository) Create(ctx context.Context, component *entity.Component) error {
	query := `INSERT INTO components (id, project_id, name, description, lead_user_id, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(ctx, query,
		component.ID, component.ProjectID, component.Name, component.Description,
		component.LeadUserID, component.CreatedAt, component.UpdatedAt)
	return err
}

func (r *componentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Component, error) {
	query := `SELECT id, project_id, name, description, lead_user_id, created_at, updated_at
		FROM components WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanComponent(row)
}

func (r *componentRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.Component, error) {
	query := `SELECT id, project_id, name, description, lead_user_id, created_at, updated_at
		FROM components WHERE project_id=$1 ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var components []*entity.Component
	for rows.Next() {
		c := &entity.Component{}
		if err := rows.Scan(&c.ID, &c.ProjectID, &c.Name, &c.Description,
			&c.LeadUserID, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		components = append(components, c)
	}
	return components, nil
}

func (r *componentRepository) Update(ctx context.Context, component *entity.Component) error {
	query := `UPDATE components SET name=$1, description=$2, lead_user_id=$3, updated_at=$4 WHERE id=$5`
	_, err := r.db.Exec(ctx, query,
		component.Name, component.Description, component.LeadUserID, component.UpdatedAt, component.ID)
	return err
}

func (r *componentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM components WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *componentRepository) CountRecords(ctx context.Context, componentID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM project_records WHERE component_id=$1`
	if err := r.db.QueryRow(ctx, query, componentID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func scanComponent(row pgx.Row) (*entity.Component, error) {
	c := &entity.Component{}
	err := row.Scan(&c.ID, &c.ProjectID, &c.Name, &c.Description,
		&c.LeadUserID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return c, nil
}
