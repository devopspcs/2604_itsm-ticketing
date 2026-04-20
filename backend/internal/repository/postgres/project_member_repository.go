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

type projectMemberRepository struct {
	db *pgxpool.Pool
}

func NewProjectMemberRepository(db *pgxpool.Pool) repository.ProjectMemberRepository {
	return &projectMemberRepository{db: db}
}

func (r *projectMemberRepository) Add(ctx context.Context, member *entity.ProjectMember) error {
	query := `INSERT INTO project_members (project_id, user_id, role, created_at) VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, member.ProjectID, member.UserID, member.Role, member.CreatedAt)
	return err
}

func (r *projectMemberRepository) Remove(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM project_members WHERE project_id=$1 AND user_id=$2", projectID, userID)
	return err
}

func (r *projectMemberRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.ProjectMember, error) {
	rows, err := r.db.Query(ctx, "SELECT project_id, user_id, role, created_at FROM project_members WHERE project_id=$1 ORDER BY created_at ASC", projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []*entity.ProjectMember
	for rows.Next() {
		m := &entity.ProjectMember{}
		if err := rows.Scan(&m.ProjectID, &m.UserID, &m.Role, &m.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *projectMemberRepository) ListProjectsByUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, "SELECT project_id FROM project_members WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *projectMemberRepository) IsMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM project_members WHERE project_id=$1 AND user_id=$2", projectID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *projectMemberRepository) GetRole(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (entity.ProjectMemberRole, error) {
	var role entity.ProjectMemberRole
	err := r.db.QueryRow(ctx, "SELECT role FROM project_members WHERE project_id=$1 AND user_id=$2", projectID, userID).Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", apperror.ErrForbidden
		}
		return "", err
	}
	return role, nil
}

func (r *projectMemberRepository) GetMember(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*entity.ProjectMember, error) {
	query := `SELECT project_id, user_id, role, created_at FROM project_members WHERE project_id=$1 AND user_id=$2`
	row := r.db.QueryRow(ctx, query, projectID, userID)
	member := &entity.ProjectMember{}
	err := row.Scan(&member.ProjectID, &member.UserID, &member.Role, &member.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return member, nil
}
