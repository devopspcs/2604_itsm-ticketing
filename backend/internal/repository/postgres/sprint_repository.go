package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	"github.com/org/itsm/pkg/apperror"
)

type sprintRepository struct {
	db *pgxpool.Pool
}

func NewSprintRepository(db *pgxpool.Pool) repository.SprintRepository {
	return &sprintRepository{db: db}
}

func (r *sprintRepository) Create(ctx context.Context, sprint *entity.Sprint) error {
	query := `INSERT INTO sprints (id, project_id, name, goal, start_date, end_date, status, actual_start_date, actual_end_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Exec(ctx, query,
		sprint.ID, sprint.ProjectID, sprint.Name, sprint.Goal, sprint.StartDate, sprint.EndDate,
		sprint.Status, sprint.ActualStartDate, sprint.ActualEndDate, sprint.CreatedAt)
	return err
}

func (r *sprintRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Sprint, error) {
	query := `SELECT id, project_id, name, goal, start_date, end_date, status, actual_start_date, actual_end_date, created_at FROM sprints WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanSprint(row)
}

func (r *sprintRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.Sprint, error) {
	query := `SELECT id, project_id, name, goal, start_date, end_date, status, actual_start_date, actual_end_date, created_at FROM sprints WHERE project_id=$1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sprints []*entity.Sprint
	for rows.Next() {
		sprint := &entity.Sprint{}
		if err := rows.Scan(&sprint.ID, &sprint.ProjectID, &sprint.Name, &sprint.Goal, &sprint.StartDate, &sprint.EndDate,
			&sprint.Status, &sprint.ActualStartDate, &sprint.ActualEndDate, &sprint.CreatedAt); err != nil {
			return nil, err
		}
		sprints = append(sprints, sprint)
	}
	return sprints, nil
}

func (r *sprintRepository) ListActive(ctx context.Context, projectID uuid.UUID) ([]*entity.Sprint, error) {
	query := `SELECT id, project_id, name, goal, start_date, end_date, status, actual_start_date, actual_end_date, created_at FROM sprints WHERE project_id=$1 AND status='Active' ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sprints []*entity.Sprint
	for rows.Next() {
		sprint := &entity.Sprint{}
		if err := rows.Scan(&sprint.ID, &sprint.ProjectID, &sprint.Name, &sprint.Goal, &sprint.StartDate, &sprint.EndDate,
			&sprint.Status, &sprint.ActualStartDate, &sprint.ActualEndDate, &sprint.CreatedAt); err != nil {
			return nil, err
		}
		sprints = append(sprints, sprint)
	}
	return sprints, nil
}

func (r *sprintRepository) Update(ctx context.Context, sprint *entity.Sprint) error {
	query := `UPDATE sprints SET name=$1, goal=$2, start_date=$3, end_date=$4, status=$5, actual_start_date=$6, actual_end_date=$7 WHERE id=$8`
	_, err := r.db.Exec(ctx, query,
		sprint.Name, sprint.Goal, sprint.StartDate, sprint.EndDate, sprint.Status,
		sprint.ActualStartDate, sprint.ActualEndDate, sprint.ID)
	return err
}

func (r *sprintRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sprints WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func scanSprint(row pgx.Row) (*entity.Sprint, error) {
	sprint := &entity.Sprint{}
	var startDate, endDate, actualStartDate, actualEndDate *time.Time
	err := row.Scan(&sprint.ID, &sprint.ProjectID, &sprint.Name, &sprint.Goal, &startDate, &endDate,
		&sprint.Status, &actualStartDate, &actualEndDate, &sprint.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	sprint.StartDate = startDate
	sprint.EndDate = endDate
	sprint.ActualStartDate = actualStartDate
	sprint.ActualEndDate = actualEndDate
	return sprint, nil
}
