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

type labelRepository struct {
	db *pgxpool.Pool
}

func NewLabelRepository(db *pgxpool.Pool) repository.LabelRepository {
	return &labelRepository{db: db}
}

func (r *labelRepository) Create(ctx context.Context, label *entity.Label) error {
	query := `INSERT INTO labels (id, project_id, name, color, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		label.ID, label.ProjectID, label.Name, label.Color, label.CreatedAt)
	return err
}

func (r *labelRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Label, error) {
	query := `SELECT id, project_id, name, color, created_at FROM labels WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanLabel(row)
}

func (r *labelRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.Label, error) {
	query := `SELECT id, project_id, name, color, created_at FROM labels WHERE project_id=$1 ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []*entity.Label
	for rows.Next() {
		label := &entity.Label{}
		if err := rows.Scan(&label.ID, &label.ProjectID, &label.Name, &label.Color, &label.CreatedAt); err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}
	return labels, nil
}

func (r *labelRepository) Update(ctx context.Context, label *entity.Label) error {
	query := `UPDATE labels SET name=$1, color=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, label.Name, label.Color, label.ID)
	return err
}

func (r *labelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM labels WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func scanLabel(row pgx.Row) (*entity.Label, error) {
	label := &entity.Label{}
	err := row.Scan(&label.ID, &label.ProjectID, &label.Name, &label.Color, &label.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return label, nil
}
