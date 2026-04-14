package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
)

type approvalRepository struct {
	db *pgxpool.Pool
}

func NewApprovalRepository(db *pgxpool.Pool) repository.ApprovalRepository {
	return &approvalRepository{db: db}
}

func (r *approvalRepository) SaveConfig(ctx context.Context, config *entity.ApprovalConfig) error {
	query := `INSERT INTO approval_configs (id, ticket_type, level, approver_id) VALUES ($1,$2,$3,$4)
		ON CONFLICT (ticket_type, level, approver_id) DO NOTHING`
	_, err := r.db.Exec(ctx, query, config.ID, config.TicketType, config.Level, config.ApproverID)
	return err
}

func (r *approvalRepository) FindConfigsByTicketType(ctx context.Context, ticketType entity.TicketType) ([]*entity.ApprovalConfig, error) {
	rows, err := r.db.Query(ctx, `SELECT id, ticket_type, level, approver_id FROM approval_configs WHERE ticket_type=$1 ORDER BY level ASC`, ticketType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var configs []*entity.ApprovalConfig
	for rows.Next() {
		c := &entity.ApprovalConfig{}
		if err := rows.Scan(&c.ID, &c.TicketType, &c.Level, &c.ApproverID); err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, nil
}

func (r *approvalRepository) SaveDecision(ctx context.Context, approval *entity.Approval) error {
	query := `INSERT INTO approvals (id, ticket_id, approver_id, level, decision, comment, decided_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.db.Exec(ctx, query, approval.ID, approval.TicketID, approval.ApproverID, approval.Level, approval.Decision, approval.Comment, approval.DecidedAt)
	return err
}

func (r *approvalRepository) FindDecisionsByTicketID(ctx context.Context, ticketID uuid.UUID) ([]*entity.Approval, error) {
	rows, err := r.db.Query(ctx, `SELECT id, ticket_id, approver_id, level, decision, comment, decided_at FROM approvals WHERE ticket_id=$1 ORDER BY level ASC`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var approvals []*entity.Approval
	for rows.Next() {
		a := &entity.Approval{}
		if err := rows.Scan(&a.ID, &a.TicketID, &a.ApproverID, &a.Level, &a.Decision, &a.Comment, &a.DecidedAt); err != nil {
			return nil, err
		}
		approvals = append(approvals, a)
	}
	return approvals, nil
}

func (r *approvalRepository) FindPendingLevel(ctx context.Context, ticketID uuid.UUID) (int, error) {
	var maxLevel int
	err := r.db.QueryRow(ctx, `SELECT COALESCE(MAX(level), 0) FROM approvals WHERE ticket_id=$1 AND decision IS NOT NULL`, ticketID).Scan(&maxLevel)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}
	return maxLevel + 1, nil
}
