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

type divisionRepository struct {
	db *pgxpool.Pool
}

func NewDivisionRepository(db *pgxpool.Pool) repository.DivisionRepository {
	return &divisionRepository{db: db}
}

func (r *divisionRepository) Create(ctx context.Context, div *entity.Division) error {
	query := `INSERT INTO divisions (id, department_id, name, code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, div.ID, div.DepartmentID, div.Name, div.Code, div.CreatedAt, div.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return apperror.ErrConflict
		}
		return err
	}
	return nil
}

func (r *divisionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Division, error) {
	query := `SELECT id, department_id, name, code, created_at, updated_at FROM divisions WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	d := &entity.Division{}
	err := row.Scan(&d.ID, &d.DepartmentID, &d.Name, &d.Code, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return d, nil
}

func (r *divisionRepository) Update(ctx context.Context, div *entity.Division) error {
	query := `UPDATE divisions SET department_id=$1, name=$2, code=$3, updated_at=$4 WHERE id=$5`
	_, err := r.db.Exec(ctx, query, div.DepartmentID, div.Name, div.Code, time.Now().UTC(), div.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return apperror.ErrConflict
		}
		return err
	}
	return nil
}

func (r *divisionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM divisions WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *divisionRepository) List(ctx context.Context, filter repository.DivisionFilter) ([]*entity.Division, error) {
	query := `SELECT id, department_id, name, code, created_at, updated_at FROM divisions WHERE 1=1`
	args := []interface{}{}
	i := 1
	if filter.DepartmentID != nil {
		query += ` AND department_id = $` + itoa(i)
		args = append(args, *filter.DepartmentID)
		i++
	}
	query += ` ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var divs []*entity.Division
	for rows.Next() {
		d := &entity.Division{}
		if err := rows.Scan(&d.ID, &d.DepartmentID, &d.Name, &d.Code, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		divs = append(divs, d)
	}
	return divs, nil
}

func (r *divisionRepository) HasTeamsOrUsers(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM teams WHERE division_id = $1`, id).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE division_id = $1`, id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
