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

type projectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) repository.ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *entity.Project) error {
	query := `INSERT INTO projects (id, name, icon_color, created_by, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.db.Exec(ctx, query,
		project.ID, project.Name, project.IconColor,
		project.CreatedBy, project.CreatedAt, project.UpdatedAt)
	return err
}

func (r *projectRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Project, error) {
	query := `SELECT id, name, icon_color, created_by, created_at, updated_at FROM projects WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanProject(row)
}

func (r *projectRepository) List(ctx context.Context, createdBy uuid.UUID) ([]*entity.Project, error) {
	query := `SELECT id, name, icon_color, created_by, created_at, updated_at FROM projects WHERE created_by=$1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*entity.Project
	for rows.Next() {
		p := &entity.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.IconColor, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *projectRepository) Update(ctx context.Context, project *entity.Project) error {
	query := `UPDATE projects SET name=$1, icon_color=$2, updated_at=$3 WHERE id=$4`
	_, err := r.db.Exec(ctx, query, project.Name, project.IconColor, project.UpdatedAt, project.ID)
	return err
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func scanProject(row pgx.Row) (*entity.Project, error) {
	p := &entity.Project{}
	err := row.Scan(&p.ID, &p.Name, &p.IconColor, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return p, nil
}
