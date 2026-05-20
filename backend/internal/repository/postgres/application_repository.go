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

// --- ApplicationRepository ---

type applicationRepository struct {
	db *pgxpool.Pool
}

func NewApplicationRepository(db *pgxpool.Pool) repository.ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(ctx context.Context, app *entity.Application) error {
	query := `INSERT INTO applications (id, name, code, description, icon, color, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(ctx, query, app.ID, app.Name, app.Code, app.Description, app.Icon, app.Color, app.IsActive, app.CreatedAt, app.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return apperror.ErrConflict
		}
		return apperror.ErrInternal
	}
	return nil
}

func (r *applicationRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Application, error) {
	query := `SELECT id, name, code, description, icon, color, is_active, created_at, updated_at
		FROM applications WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	app := &entity.Application{}
	err := row.Scan(&app.ID, &app.Name, &app.Code, &app.Description, &app.Icon, &app.Color, &app.IsActive, &app.CreatedAt, &app.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, apperror.ErrInternal
	}
	return app, nil
}

func (r *applicationRepository) FindByCode(ctx context.Context, code string) (*entity.Application, error) {
	query := `SELECT id, name, code, description, icon, color, is_active, created_at, updated_at
		FROM applications WHERE code = $1`
	row := r.db.QueryRow(ctx, query, code)
	app := &entity.Application{}
	err := row.Scan(&app.ID, &app.Name, &app.Code, &app.Description, &app.Icon, &app.Color, &app.IsActive, &app.CreatedAt, &app.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, apperror.ErrInternal
	}
	return app, nil
}

func (r *applicationRepository) Update(ctx context.Context, app *entity.Application) error {
	query := `UPDATE applications SET name=$1, description=$2, icon=$3, color=$4, is_active=$5, updated_at=$6
		WHERE id=$7`
	tag, err := r.db.Exec(ctx, query, app.Name, app.Description, app.Icon, app.Color, app.IsActive, app.UpdatedAt, app.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return apperror.ErrConflict
		}
		return apperror.ErrInternal
	}
	if tag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (r *applicationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM applications WHERE id=$1`, id)
	if err != nil {
		return apperror.ErrInternal
	}
	if tag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (r *applicationRepository) List(ctx context.Context) ([]*entity.Application, error) {
	query := `SELECT id, name, code, description, icon, color, is_active, created_at, updated_at
		FROM applications ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	defer rows.Close()

	var apps []*entity.Application
	for rows.Next() {
		app := &entity.Application{}
		if err := rows.Scan(&app.ID, &app.Name, &app.Code, &app.Description, &app.Icon, &app.Color, &app.IsActive, &app.CreatedAt, &app.UpdatedAt); err != nil {
			return nil, apperror.ErrInternal
		}
		apps = append(apps, app)
	}
	return apps, nil
}

// --- UserAppAccessRepository ---

type userAppAccessRepository struct {
	db *pgxpool.Pool
}

func NewUserAppAccessRepository(db *pgxpool.Pool) repository.UserAppAccessRepository {
	return &userAppAccessRepository{db: db}
}

func (r *userAppAccessRepository) Grant(ctx context.Context, access *entity.UserAppAccess) error {
	query := `INSERT INTO user_app_access (id, user_id, app_id, role, granted_at, granted_by)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, access.ID, access.UserID, access.AppID, access.Role, access.GrantedAt, access.GrantedBy)
	if err != nil {
		if isUniqueViolation(err) {
			return apperror.ErrConflict
		}
		return apperror.ErrInternal
	}
	return nil
}

func (r *userAppAccessRepository) Revoke(ctx context.Context, userID, appID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM user_app_access WHERE user_id=$1 AND app_id=$2`, userID, appID)
	if err != nil {
		return apperror.ErrInternal
	}
	if tag.RowsAffected() == 0 {
		return apperror.ErrNotFound
	}
	return nil
}

func (r *userAppAccessRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*entity.UserAppAccess, error) {
	query := `SELECT id, user_id, app_id, role, granted_at, granted_by
		FROM user_app_access WHERE user_id=$1 ORDER BY granted_at ASC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	defer rows.Close()

	var list []*entity.UserAppAccess
	for rows.Next() {
		a := &entity.UserAppAccess{}
		if err := rows.Scan(&a.ID, &a.UserID, &a.AppID, &a.Role, &a.GrantedAt, &a.GrantedBy); err != nil {
			return nil, apperror.ErrInternal
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *userAppAccessRepository) ListByApp(ctx context.Context, appID uuid.UUID) ([]*entity.UserAppAccess, error) {
	query := `SELECT id, user_id, app_id, role, granted_at, granted_by
		FROM user_app_access WHERE app_id=$1 ORDER BY granted_at ASC`
	rows, err := r.db.Query(ctx, query, appID)
	if err != nil {
		return nil, apperror.ErrInternal
	}
	defer rows.Close()

	var list []*entity.UserAppAccess
	for rows.Next() {
		a := &entity.UserAppAccess{}
		if err := rows.Scan(&a.ID, &a.UserID, &a.AppID, &a.Role, &a.GrantedAt, &a.GrantedBy); err != nil {
			return nil, apperror.ErrInternal
		}
		list = append(list, a)
	}
	return list, nil
}

func (r *userAppAccessRepository) HasAccess(ctx context.Context, userID, appID uuid.UUID) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM user_app_access WHERE user_id=$1 AND app_id=$2`, userID, appID).Scan(&count)
	if err != nil {
		return false, apperror.ErrInternal
	}
	return count > 0, nil
}

func (r *userAppAccessRepository) BulkGrant(ctx context.Context, userIDs []uuid.UUID, appID uuid.UUID, role string, grantedBy uuid.UUID) error {
	batch := &pgx.Batch{}
	for _, uid := range userIDs {
		id := uuid.New()
		batch.Queue(
			`INSERT INTO user_app_access (id, user_id, app_id, role, granted_at, granted_by)
			VALUES ($1, $2, $3, $4, NOW(), $5)
			ON CONFLICT (user_id, app_id) DO UPDATE SET role=$3, granted_by=$5`,
			id, uid, appID, role, grantedBy,
		)
	}
	br := r.db.SendBatch(ctx, batch)
	defer br.Close()
	for range userIDs {
		if _, err := br.Exec(); err != nil {
			return apperror.ErrInternal
		}
	}
	return nil
}
