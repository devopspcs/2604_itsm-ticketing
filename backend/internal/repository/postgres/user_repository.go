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

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (id, full_name, email, password, role, is_active, department_id, division_id, team_id, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.Exec(ctx, query,
		user.ID, user.FullName, user.Email, user.PasswordHash,
		user.Role, user.IsActive, user.DepartmentID, user.DivisionID,
		user.TeamID, user.Position, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return apperror.ErrConflict
		}
		return err
	}
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, full_name, email, password, role, is_active, department_id, division_id, team_id, position, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)
	return scanUser(row)
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `SELECT id, full_name, email, password, role, is_active, department_id, division_id, team_id, position, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return scanUser(row)
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET full_name=$1, email=$2, role=$3, is_active=$4, department_id=$5, division_id=$6, team_id=$7, position=$8, updated_at=$9 WHERE id=$10`
	_, err := r.db.Exec(ctx, query, user.FullName, user.Email, user.Role, user.IsActive,
		user.DepartmentID, user.DivisionID, user.TeamID, user.Position,
		time.Now().UTC(), user.ID)
	return err
}

func (r *userRepository) List(ctx context.Context, filter repository.UserFilter) ([]*entity.User, error) {
	query := `SELECT id, full_name, email, password, role, is_active, department_id, division_id, team_id, position, created_at, updated_at FROM users WHERE 1=1`
	args := []interface{}{}
	i := 1
	if filter.Role != nil {
		query += ` AND role = $` + itoa(i)
		args = append(args, *filter.Role)
		i++
	}
	if filter.IsActive != nil {
		query += ` AND is_active = $` + itoa(i)
		args = append(args, *filter.IsActive)
		i++
	}
	query += ` ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*entity.User
	for rows.Next() {
		u := &entity.User{}
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive,
			&u.DepartmentID, &u.DivisionID, &u.TeamID, &u.Position,
			&u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func scanUser(row pgx.Row) (*entity.User, error) {
	u := &entity.User{}
	err := row.Scan(&u.ID, &u.FullName, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive,
		&u.DepartmentID, &u.DivisionID, &u.TeamID, &u.Position,
		&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return u, nil
}
