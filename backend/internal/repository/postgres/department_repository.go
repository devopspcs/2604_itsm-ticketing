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
	query := `INSERT INTO departments (id, name, code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, dept.ID, dept.Name, dept.Code, dept.CreatedAt, dept.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return apperror.ErrConflict
		}
		return err
	}
	return nil
}

func (r *departmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Department, error) {
	query := `SELECT id, name, code, created_at, updated_at FROM departments WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	d := &entity.Department{}
	err := row.Scan(&d.ID, &d.Name, &d.Code, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return d, nil
}

func (r *departmentRepository) Update(ctx context.Context, dept *entity.Department) error {
	query := `UPDATE departments SET name=$1, code=$2, updated_at=$3 WHERE id=$4`
	_, err := r.db.Exec(ctx, query, dept.Name, dept.Code, time.Now().UTC(), dept.ID)
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

func (r *departmentRepository) List(ctx context.Context) ([]*entity.Department, error) {
	query := `SELECT id, name, code, created_at, updated_at FROM departments ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var depts []*entity.Department
	for rows.Next() {
		d := &entity.Department{}
		if err := rows.Scan(&d.ID, &d.Name, &d.Code, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		depts = append(depts, d)
	}
	return depts, nil
}

func (r *departmentRepository) HasDivisions(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM divisions WHERE department_id = $1`, id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
