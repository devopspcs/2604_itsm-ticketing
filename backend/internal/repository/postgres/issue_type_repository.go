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

type issueTypeRepository struct {
	db *pgxpool.Pool
}

func NewIssueTypeRepository(db *pgxpool.Pool) repository.IssueTypeRepository {
	return &issueTypeRepository{db: db}
}

func (r *issueTypeRepository) Create(ctx context.Context, issueType *entity.IssueType) error {
	query := `INSERT INTO issue_types (id, name, icon, description, created_at)
		VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query,
		issueType.ID, issueType.Name, issueType.Icon,
		issueType.Description, issueType.CreatedAt)
	return err
}

func (r *issueTypeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.IssueType, error) {
	query := `SELECT id, name, icon, description, created_at FROM issue_types WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanIssueType(row)
}

func (r *issueTypeRepository) List(ctx context.Context) ([]*entity.IssueType, error) {
	query := `SELECT id, name, icon, description, created_at FROM issue_types ORDER BY name ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issueTypes []*entity.IssueType
	for rows.Next() {
		it := &entity.IssueType{}
		if err := rows.Scan(&it.ID, &it.Name, &it.Icon, &it.Description, &it.CreatedAt); err != nil {
			return nil, err
		}
		issueTypes = append(issueTypes, it)
	}
	return issueTypes, nil
}

func scanIssueType(row pgx.Row) (*entity.IssueType, error) {
	it := &entity.IssueType{}
	err := row.Scan(&it.ID, &it.Name, &it.Icon, &it.Description, &it.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return it, nil
}
