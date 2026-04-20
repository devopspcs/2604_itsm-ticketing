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

type releaseRepository struct {
	db *pgxpool.Pool
}

func NewReleaseRepository(db *pgxpool.Pool) repository.ReleaseRepository {
	return &releaseRepository{db: db}
}

func (r *releaseRepository) Create(ctx context.Context, release *entity.Release) error {
	query := `INSERT INTO releases (id, project_id, name, version, description, start_date, release_date, status, created_by, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.db.Exec(ctx, query,
		release.ID, release.ProjectID, release.Name, release.Version, release.Description,
		release.StartDate, release.ReleaseDate, release.Status, release.CreatedBy,
		release.CreatedAt, release.UpdatedAt)
	return err
}

func (r *releaseRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Release, error) {
	query := `SELECT id, project_id, name, version, description, start_date, release_date, status, created_by, created_at, updated_at
		FROM releases WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanRelease(row)
}

func (r *releaseRepository) ListByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.Release, error) {
	query := `SELECT id, project_id, name, version, description, start_date, release_date, status, created_by, created_at, updated_at
		FROM releases WHERE project_id=$1 ORDER BY release_date ASC NULLS LAST`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var releases []*entity.Release
	for rows.Next() {
		rel := &entity.Release{}
		if err := rows.Scan(&rel.ID, &rel.ProjectID, &rel.Name, &rel.Version, &rel.Description,
			&rel.StartDate, &rel.ReleaseDate, &rel.Status, &rel.CreatedBy,
			&rel.CreatedAt, &rel.UpdatedAt); err != nil {
			return nil, err
		}
		releases = append(releases, rel)
	}
	return releases, nil
}

func (r *releaseRepository) Update(ctx context.Context, release *entity.Release) error {
	query := `UPDATE releases SET name=$1, version=$2, description=$3, start_date=$4, release_date=$5, status=$6, updated_at=$7 WHERE id=$8`
	_, err := r.db.Exec(ctx, query,
		release.Name, release.Version, release.Description, release.StartDate,
		release.ReleaseDate, release.Status, release.UpdatedAt, release.ID)
	return err
}

func (r *releaseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM releases WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func scanRelease(row pgx.Row) (*entity.Release, error) {
	rel := &entity.Release{}
	err := row.Scan(&rel.ID, &rel.ProjectID, &rel.Name, &rel.Version, &rel.Description,
		&rel.StartDate, &rel.ReleaseDate, &rel.Status, &rel.CreatedBy,
		&rel.CreatedAt, &rel.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return rel, nil
}
