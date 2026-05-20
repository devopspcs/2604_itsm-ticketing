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

type departmentRepository struct {
	db *pgxpool.Pool
}

func NewDepartmentRepository(db *pgxpool.Pool) repository.DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(ctx context.Context, dept *entity.Department) error {
	query := `INSERT INTO departments (id, division_id, name, code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, dept.ID, dept.DivisionID, dept.Name, dept.Code, dept.CreatedAt, dept.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return apperror.ErrConflict
		}
		return err
	}
	return nil
}

func (r *departmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Department, error) {
	query := `SELECT id, division_id, name, code, created_at, updated_at FROM departments WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	d := &entity.Department{}
	err := row.Scan(&d.ID, &d.DivisionID, &d.Name, &d.Code, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return d, nil
}

func (r *departmentRepository) Update(ctx context.Context, dept *entity.Department) error {
	query := `UPDATE departments SET division_id=$1, name=$2, code=$3, updated_at=$4 WHERE id=$5`
	_, err := r.db.Exec(ctx, query, dept.DivisionID, dept.Name, dept.Code, time.Now().UTC(), dept.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return apperror.ErrConflict
		}
		return err
	}
	return nil
}

func (r *departmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM departments WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *departmentRepository) List(ctx context.Context, filter repository.DepartmentFilter) ([]*entity.Department, error) {
	query := `SELECT id, division_id, name, code, created_at, updated_at FROM departments WHERE 1=1`
	args := []interface{}{}
	i := 1
	if filter.DivisionID != nil {
		query += ` AND division_id = $` + itoa(i)
		args = append(args, *filter.DivisionID)
		i++
	}
	_ = i
	query += ` ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var depts []*entity.Department
	for rows.Next() {
		d := &entity.Department{}
		if err := rows.Scan(&d.ID, &d.DivisionID, &d.Name, &d.Code, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		depts = append(depts, d)
	}
	return depts, nil
}

func (r *departmentRepository) HasTeamsOrUsers(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM teams WHERE department_id = $1`, id).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE department_id = $1`, id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
