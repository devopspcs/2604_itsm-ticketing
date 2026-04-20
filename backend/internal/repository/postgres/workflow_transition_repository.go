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

type workflowTransitionRepository struct {
	db *pgxpool.Pool
}

func NewWorkflowTransitionRepository(db *pgxpool.Pool) repository.WorkflowTransitionRepository {
	return &workflowTransitionRepository{db: db}
}

func (r *workflowTransitionRepository) Create(ctx context.Context, transition *entity.WorkflowTransition) error {
	query := `INSERT INTO workflow_transitions (id, workflow_id, from_status_id, to_status_id, validation_rule, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query,
		transition.ID, transition.WorkflowID, transition.FromStatusID, transition.ToStatusID,
		transition.ValidationRule, transition.CreatedAt)
	return err
}

func (r *workflowTransitionRepository) ListByWorkflow(ctx context.Context, workflowID uuid.UUID) ([]*entity.WorkflowTransition, error) {
	query := `SELECT id, workflow_id, from_status_id, to_status_id, validation_rule, created_at FROM workflow_transitions WHERE workflow_id=$1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, workflowID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transitions []*entity.WorkflowTransition
	for rows.Next() {
		transition := &entity.WorkflowTransition{}
		if err := rows.Scan(&transition.ID, &transition.WorkflowID, &transition.FromStatusID, &transition.ToStatusID, &transition.ValidationRule, &transition.CreatedAt); err != nil {
			return nil, err
		}
		transitions = append(transitions, transition)
	}
	return transitions, nil
}

func (r *workflowTransitionRepository) ListFromStatus(ctx context.Context, fromStatusID uuid.UUID) ([]*entity.WorkflowTransition, error) {
	query := `SELECT id, workflow_id, from_status_id, to_status_id, validation_rule, created_at FROM workflow_transitions WHERE from_status_id=$1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, fromStatusID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transitions []*entity.WorkflowTransition
	for rows.Next() {
		transition := &entity.WorkflowTransition{}
		if err := rows.Scan(&transition.ID, &transition.WorkflowID, &transition.FromStatusID, &transition.ToStatusID, &transition.ValidationRule, &transition.CreatedAt); err != nil {
			return nil, err
		}
		transitions = append(transitions, transition)
	}
	return transitions, nil
}

func (r *workflowTransitionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM workflow_transitions WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func scanWorkflowTransition(row pgx.Row) (*entity.WorkflowTransition, error) {
	transition := &entity.WorkflowTransition{}
	err := row.Scan(&transition.ID, &transition.WorkflowID, &transition.FromStatusID, &transition.ToStatusID, &transition.ValidationRule, &transition.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return transition, nil
}
