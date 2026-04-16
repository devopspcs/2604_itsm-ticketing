package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
)

type projectActivityLogRepository struct {
	db *pgxpool.Pool
}

func NewProjectActivityLogRepository(db *pgxpool.Pool) repository.ProjectActivityLogRepository {
	return &projectActivityLogRepository{db: db}
}

func (r *projectActivityLogRepository) Append(ctx context.Context, log *entity.ProjectActivityLog) error {
	query := `INSERT INTO project_activity_logs (id, project_id, record_id, actor_id, action, detail, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(ctx, query,
		log.ID, log.ProjectID, log.RecordID, log.ActorID, log.Action, log.Detail, log.CreatedAt)
	return err
}

func (r *projectActivityLogRepository) ListByProject(ctx context.Context, projectID uuid.UUID, limit int) ([]*entity.ProjectActivityLog, error) {
	query := `SELECT id, project_id, record_id, actor_id, action, detail, created_at
		FROM project_activity_logs WHERE project_id=$1 ORDER BY created_at DESC LIMIT $2`
	return r.queryLogs(ctx, query, projectID, limit)
}

func (r *projectActivityLogRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]*entity.ProjectActivityLog, error) {
	query := `SELECT id, project_id, record_id, actor_id, action, detail, created_at
		FROM project_activity_logs WHERE actor_id=$1 ORDER BY created_at DESC LIMIT $2`
	return r.queryLogs(ctx, query, userID, limit)
}

func (r *projectActivityLogRepository) queryLogs(ctx context.Context, query string, args ...interface{}) ([]*entity.ProjectActivityLog, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*entity.ProjectActivityLog
	for rows.Next() {
		l := &entity.ProjectActivityLog{}
		if err := rows.Scan(&l.ID, &l.ProjectID, &l.RecordID, &l.ActorID, &l.Action, &l.Detail, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}
