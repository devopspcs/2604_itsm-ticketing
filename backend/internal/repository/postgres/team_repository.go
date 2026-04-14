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

type teamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) repository.TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(ctx context.Context, team *entity.Team) error {
	query := `INSERT INTO teams (id, division_id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, team.ID, team.DivisionID, team.Name, team.CreatedAt, team.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return apperror.ErrConflict
		}
		return err
	}
	return nil
}

func (r *teamRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
	query := `SELECT id, division_id, name, created_at, updated_at FROM teams WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	t := &entity.Team{}
	err := row.Scan(&t.ID, &t.DivisionID, &t.Name, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *teamRepository) Update(ctx context.Context, team *entity.Team) error {
	query := `UPDATE teams SET division_id=$1, name=$2, updated_at=$3 WHERE id=$4`
	_, err := r.db.Exec(ctx, query, team.DivisionID, team.Name, time.Now().UTC(), team.ID)
	return err
}

func (r *teamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM teams WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *teamRepository) List(ctx context.Context, filter repository.TeamFilter) ([]*entity.Team, error) {
	query := `SELECT id, division_id, name, created_at, updated_at FROM teams WHERE 1=1`
	args := []interface{}{}
	i := 1
	if filter.DivisionID != nil {
		query += ` AND division_id = $` + itoa(i)
		args = append(args, *filter.DivisionID)
		i++
	}
	query += ` ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var teams []*entity.Team
	for rows.Next() {
		t := &entity.Team{}
		if err := rows.Scan(&t.ID, &t.DivisionID, &t.Name, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

func (r *teamRepository) HasUsers(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE team_id = $1`, id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
