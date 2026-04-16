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

type projectColumnRepository struct {
	db *pgxpool.Pool
}

func NewProjectColumnRepository(db *pgxpool.Pool) repository.ProjectColumnRepository {
	return &projectColumnRepository{db: db}
}

func (r *projectColumnRepository) Create(ctx context.Context, col *entity.ProjectColumn) error {
	query := `INSERT INTO project_columns (id, project_id, name, position, created_at)
		VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.Exec(ctx, query, col.ID, col.ProjectID, col.Name, col.Position, col.CreatedAt)
	return err
}

func (r *projectColumnRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.ProjectColumn, error) {
	query := `SELECT id, project_id, name, position, created_at FROM project_columns WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanProjectColumn(row)
}

func (r *projectColumnRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.ProjectColumn, error) {
	query := `SELECT id, project_id, name, position, created_at FROM project_columns WHERE project_id=$1 ORDER BY position ASC`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []*entity.ProjectColumn
	for rows.Next() {
		c := &entity.ProjectColumn{}
		if err := rows.Scan(&c.ID, &c.ProjectID, &c.Name, &c.Position, &c.CreatedAt); err != nil {
			return nil, err
		}
		cols = append(cols, c)
	}
	return cols, nil
}

func (r *projectColumnRepository) Update(ctx context.Context, col *entity.ProjectColumn) error {
	query := `UPDATE project_columns SET name=$1 WHERE id=$2`
	_, err := r.db.Exec(ctx, query, col.Name, col.ID)
	return err
}

func (r *projectColumnRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM project_columns WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *projectColumnRepository) HasRecords(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM project_records WHERE column_id=$1`
	if err := r.db.QueryRow(ctx, query, id).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *projectColumnRepository) GetMaxPosition(ctx context.Context, projectID uuid.UUID) (int, error) {
	var maxPos *int
	query := `SELECT MAX(position) FROM project_columns WHERE project_id=$1`
	if err := r.db.QueryRow(ctx, query, projectID).Scan(&maxPos); err != nil {
		return -1, err
	}
	if maxPos == nil {
		return -1, nil
	}
	return *maxPos, nil
}

func (r *projectColumnRepository) BulkUpdatePositions(ctx context.Context, projectID uuid.UUID, positions map[uuid.UUID]int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for id, pos := range positions {
		if _, err := tx.Exec(ctx, "UPDATE project_columns SET position=$1 WHERE id=$2 AND project_id=$3", pos, id, projectID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func scanProjectColumn(row pgx.Row) (*entity.ProjectColumn, error) {
	c := &entity.ProjectColumn{}
	err := row.Scan(&c.ID, &c.ProjectID, &c.Name, &c.Position, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return c, nil
}
