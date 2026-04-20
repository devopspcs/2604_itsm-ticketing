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

type issueTypeSchemeRepository struct {
	db *pgxpool.Pool
}

func NewIssueTypeSchemeRepository(db *pgxpool.Pool) repository.IssueTypeSchemeRepository {
	return &issueTypeSchemeRepository{db: db}
}

func (r *issueTypeSchemeRepository) Create(ctx context.Context, scheme *entity.IssueTypeScheme) error {
	query := `INSERT INTO issue_type_schemes (id, project_id, name, created_at)
		VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query,
		scheme.ID, scheme.ProjectID, scheme.Name, scheme.CreatedAt)
	return err
}

func (r *issueTypeSchemeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.IssueTypeScheme, error) {
	query := `SELECT id, project_id, name, created_at FROM issue_type_schemes WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanIssueTypeScheme(row)
}

func (r *issueTypeSchemeRepository) FindByProject(ctx context.Context, projectID uuid.UUID) (*entity.IssueTypeScheme, error) {
	query := `SELECT id, project_id, name, created_at FROM issue_type_schemes WHERE project_id=$1`
	row := r.db.QueryRow(ctx, query, projectID)
	return scanIssueTypeScheme(row)
}

func (r *issueTypeSchemeRepository) Update(ctx context.Context, scheme *entity.IssueTypeScheme) error {
	query := `UPDATE issue_type_schemes SET name=$1 WHERE id=$2`
	_, err := r.db.Exec(ctx, query, scheme.Name, scheme.ID)
	return err
}

func (r *issueTypeSchemeRepository) ListItems(ctx context.Context, schemeID uuid.UUID) ([]*entity.IssueTypeSchemeItem, error) {
	query := `SELECT id, scheme_id, issue_type_id, created_at FROM issue_type_scheme_items WHERE scheme_id=$1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, schemeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*entity.IssueTypeSchemeItem
	for rows.Next() {
		item := &entity.IssueTypeSchemeItem{}
		if err := rows.Scan(&item.ID, &item.SchemeID, &item.IssueTypeID, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *issueTypeSchemeRepository) AddItem(ctx context.Context, item *entity.IssueTypeSchemeItem) error {
	query := `INSERT INTO issue_type_scheme_items (id, scheme_id, issue_type_id, created_at)
		VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query,
		item.ID, item.SchemeID, item.IssueTypeID, item.CreatedAt)
	return err
}

func (r *issueTypeSchemeRepository) RemoveItem(ctx context.Context, schemeID uuid.UUID, issueTypeID uuid.UUID) error {
	query := `DELETE FROM issue_type_scheme_items WHERE scheme_id=$1 AND issue_type_id=$2`
	_, err := r.db.Exec(ctx, query, schemeID, issueTypeID)
	return err
}

func scanIssueTypeScheme(row pgx.Row) (*entity.IssueTypeScheme, error) {
	scheme := &entity.IssueTypeScheme{}
	err := row.Scan(&scheme.ID, &scheme.ProjectID, &scheme.Name, &scheme.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return scheme, nil
}
