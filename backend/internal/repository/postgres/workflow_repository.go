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

type workflowRepository struct {
	db *pgxpool.Pool
}

func NewWorkflowRepository(db *pgxpool.Pool) repository.WorkflowRepository {
	return &workflowRepository{db: db}
}

func (r *workflowRepository) Create(ctx context.Context, workflow *entity.Workflow) error {
	query := `INSERT INTO workflows (id, project_id, name, initial_status, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		workflow.ID, workflow.ProjectID, workflow.Name, workflow.InitialStatus, workflow.CreatedAt)
	return err
}

func (r *workflowRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Workflow, error) {
	query := `SELECT id, project_id, name, initial_status, created_at FROM workflows WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanWorkflow(row)
}

func (r *workflowRepository) FindByProject(ctx context.Context, projectID uuid.UUID) (*entity.Workflow, error) {
	query := `SELECT id, project_id, name, initial_status, created_at FROM workflows WHERE project_id=$1`
	row := r.db.QueryRow(ctx, query, projectID)
	return scanWorkflow(row)
}

func (r *workflowRepository) Update(ctx context.Context, workflow *entity.Workflow) error {
	query := `UPDATE workflows SET name=$1, initial_status=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, workflow.Name, workflow.InitialStatus, workflow.ID)
	return err
}

func scanWorkflow(row pgx.Row) (*entity.Workflow, error) {
	workflow := &entity.Workflow{}
	err := row.Scan(&workflow.ID, &workflow.ProjectID, &workflow.Name, &workflow.InitialStatus, &workflow.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return workflow, nil
}
