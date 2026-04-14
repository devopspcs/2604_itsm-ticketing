package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
)

type activityLogRepository struct {
	db *pgxpool.Pool
}

func NewActivityLogRepository(db *pgxpool.Pool) repository.ActivityLogRepository {
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Append(ctx context.Context, log *entity.ActivityLog) error {
	query := `INSERT INTO activity_logs (id, ticket_id, actor_id, action, old_value, new_value, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(ctx, query, log.ID, log.TicketID, log.ActorID, log.Action, log.OldValue, log.NewValue, log.CreatedAt)
	return err
}

func (r *activityLogRepository) FindByTicketID(ctx context.Context, ticketID uuid.UUID) ([]*entity.ActivityLog, error) {
	rows, err := r.db.Query(ctx, `SELECT id, ticket_id, actor_id, action, old_value, new_value, created_at FROM activity_logs WHERE ticket_id=$1 ORDER BY created_at ASC`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []*entity.ActivityLog
	for rows.Next() {
		l := &entity.ActivityLog{}
		if err := rows.Scan(&l.ID, &l.TicketID, &l.ActorID, &l.Action, &l.OldValue, &l.NewValue, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
