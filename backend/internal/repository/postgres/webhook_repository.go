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

type webhookRepository struct {
	db *pgxpool.Pool
}

func NewWebhookRepository(db *pgxpool.Pool) repository.WebhookRepository {
	return &webhookRepository{db: db}
}

func (r *webhookRepository) Create(ctx context.Context, config *entity.WebhookConfig) error {
	query := `INSERT INTO webhook_configs (id, url, events, secret_key, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(ctx, query, config.ID, config.URL, config.Events, config.SecretKey, config.IsActive, config.CreatedAt, config.UpdatedAt)
	return err
}

func (r *webhookRepository) FindAll(ctx context.Context) ([]*entity.WebhookConfig, error) {
	rows, err := r.db.Query(ctx, `SELECT id, url, events, secret_key, is_active, created_at, updated_at FROM webhook_configs ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var configs []*entity.WebhookConfig
	for rows.Next() {
		c := &entity.WebhookConfig{}
		if err := rows.Scan(&c.ID, &c.URL, &c.Events, &c.SecretKey, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, nil
}

func (r *webhookRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.WebhookConfig, error) {
	row := r.db.QueryRow(ctx, `SELECT id, url, events, secret_key, is_active, created_at, updated_at FROM webhook_configs WHERE id=$1`, id)
	c := &entity.WebhookConfig{}
	err := row.Scan(&c.ID, &c.URL, &c.Events, &c.SecretKey, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *webhookRepository) Update(ctx context.Context, config *entity.WebhookConfig) error {
	_, err := r.db.Exec(ctx, `UPDATE webhook_configs SET url=$1, events=$2, secret_key=$3, is_active=$4, updated_at=NOW() WHERE id=$5`,
		config.URL, config.Events, config.SecretKey, config.IsActive, config.ID)
	return err
}

func (r *webhookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM webhook_configs WHERE id=$1`, id)
	return err
}

func (r *webhookRepository) SaveLog(ctx context.Context, log *entity.WebhookLog) error {
	query := `INSERT INTO webhook_logs (id, webhook_id, event, payload, response_status, attempt, sent_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(ctx, query, log.ID, log.WebhookID, log.Event, log.Payload, log.ResponseStatus, log.Attempt, log.SentAt)
	return err
}
