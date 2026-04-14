package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
)

type notificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) repository.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notif *entity.Notification) error {
	query := `INSERT INTO notifications (id, user_id, ticket_id, message, is_read, created_at) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.db.Exec(ctx, query, notif.ID, notif.UserID, notif.TicketID, notif.Message, notif.IsRead, notif.CreatedAt)
	return err
}

func (r *notificationRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error) {
	rows, err := r.db.Query(ctx, `SELECT id, user_id, ticket_id, message, is_read, created_at FROM notifications WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notifs []*entity.Notification
	for rows.Next() {
		n := &entity.Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.TicketID, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE notifications SET is_read=true WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}
