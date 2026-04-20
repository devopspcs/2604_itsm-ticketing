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

type workflowStatusRepository struct {
	db *pgxpool.Pool
}

func NewWorkflowStatusRepository(db *pgxpool.Pool) repository.WorkflowStatusRepository {
	return &workflowStatusRepository{db: db}
}

func (r *workflowStatusRepository) Create(ctx context.Context, status *entity.WorkflowStatus) error {
	query := `INSERT INTO workflow_statuses (id, workflow_id, status_name, status_order, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		status.ID, status.WorkflowID, status.StatusName, status.StatusOrder, status.CreatedAt)
	return err
}

func (r *workflowStatusRepository) ListByWorkflow(ctx context.Context, workflowID uuid.UUID) ([]*entity.WorkflowStatus, error) {
	query := `SELECT id, workflow_id, status_name, status_order, created_at FROM workflow_statuses WHERE workflow_id=$1 ORDER BY status_order ASC`
	rows, err := r.db.Query(ctx, query, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []*entity.WorkflowStatus
	for rows.Next() {
		status := &entity.WorkflowStatus{}
		if err := rows.Scan(&status.ID, &status.WorkflowID, &status.StatusName, &status.StatusOrder, &status.CreatedAt); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

func (r *workflowStatusRepository) FindByName(ctx context.Context, workflowID uuid.UUID, statusName string) (*entity.WorkflowStatus, error) {
	query := `SELECT id, workflow_id, status_name, status_order, created_at FROM workflow_statuses WHERE workflow_id=$1 AND status_name=$2`
	row := r.db.QueryRow(ctx, query, workflowID, statusName)
	return scanWorkflowStatus(row)
}

func (r *workflowStatusRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM workflow_statuses WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *workflowStatusRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.WorkflowStatus, error) {
	query := `SELECT id, workflow_id, status_name, status_order, created_at FROM workflow_statuses WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanWorkflowStatus(row)
}

func scanWorkflowStatus(row pgx.Row) (*entity.WorkflowStatus, error) {
	status := &entity.WorkflowStatus{}
	err := row.Scan(&status.ID, &status.WorkflowID, &status.StatusName, &status.StatusOrder, &status.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return status, nil
}
