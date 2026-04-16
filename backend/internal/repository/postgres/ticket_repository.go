package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	"github.com/org/itsm/pkg/apperror"
)

type ticketRepository struct {
	db *pgxpool.Pool
}

func NewTicketRepository(db *pgxpool.Pool) repository.TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *entity.Ticket) error {
	query := `INSERT INTO tickets (id, title, description, type, category, priority, status, created_by, assigned_to, created_at, updated_at, resolved_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.db.Exec(ctx, query,
		ticket.ID, ticket.Title, ticket.Description, ticket.Type,
		ticket.Category, ticket.Priority, ticket.Status,
		ticket.CreatedBy, ticket.AssignedTo, ticket.CreatedAt, ticket.UpdatedAt, ticket.ResolvedAt)
	return err
}

func (r *ticketRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Ticket, error) {
	query := `SELECT id, title, description, type, category, priority, status, created_by, assigned_to, created_at, updated_at, resolved_at FROM tickets WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanTicket(row)
}

func (r *ticketRepository) List(ctx context.Context, filter repository.TicketFilter) (*repository.PaginatedTickets, error) {
	args := []interface{}{}
	where := "WHERE 1=1"
	i := 1

	if filter.Status != nil {
		where += fmt.Sprintf(" AND status=$%d", i)
		args = append(args, *filter.Status)
		i++
	}
	if filter.Type != nil {
		where += fmt.Sprintf(" AND type=$%d", i)
		args = append(args, *filter.Type)
		i++
	}
	if filter.Priority != nil {
		where += fmt.Sprintf(" AND priority=$%d", i)
		args = append(args, *filter.Priority)
		i++
	}
	if filter.Category != nil {
		where += fmt.Sprintf(" AND category=$%d", i)
		args = append(args, *filter.Category)
		i++
	}
	if filter.AssignedTo != nil {
		where += fmt.Sprintf(" AND assigned_to=$%d", i)
		args = append(args, *filter.AssignedTo)
		i++
	}
	if filter.CreatedBy != nil {
		where += fmt.Sprintf(" AND created_by=$%d", i)
		args = append(args, *filter.CreatedBy)
		i++
	}
	if filter.DateFrom != nil {
		where += fmt.Sprintf(" AND created_at>=$%d", i)
		args = append(args, *filter.DateFrom)
		i++
	}
	if filter.DateTo != nil {
		where += fmt.Sprintf(" AND created_at<=$%d", i)
		args = append(args, *filter.DateTo)
		i++
	}
	if filter.Search != nil && *filter.Search != "" {
		where += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", i, i)
		args = append(args, "%"+*filter.Search+"%")
		i++
	}

	countQuery := "SELECT COUNT(*) FROM tickets " + where
	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	listArgs := append(args, pageSize, offset)
	listQuery := fmt.Sprintf(
		"SELECT id, title, description, type, category, priority, status, created_by, assigned_to, created_at, updated_at, resolved_at FROM tickets %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		where, i, i+1,
	)

	rows, err := r.db.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*entity.Ticket
	for rows.Next() {
		t, err := scanTicketRow(rows)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return &repository.PaginatedTickets{Tickets: tickets, Total: total, Page: page, PageSize: pageSize}, nil
}

func (r *ticketRepository) Update(ctx context.Context, ticket *entity.Ticket) error {
	query := `UPDATE tickets SET title=$1, description=$2, category=$3, priority=$4, status=$5, assigned_to=$6, updated_at=$7, resolved_at=$8 WHERE id=$9`
	_, err := r.db.Exec(ctx, query, ticket.Title, ticket.Description, ticket.Category, ticket.Priority, ticket.Status, ticket.AssignedTo, time.Now().UTC(), ticket.ResolvedAt, ticket.ID)
	return err
}

func scanTicket(row pgx.Row) (*entity.Ticket, error) {
	t := &entity.Ticket{}
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Type, &t.Category, &t.Priority, &t.Status, &t.CreatedBy, &t.AssignedTo, &t.CreatedAt, &t.UpdatedAt, &t.ResolvedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return t, nil
}

func scanTicketRow(rows pgx.Rows) (*entity.Ticket, error) {
	t := &entity.Ticket{}
	err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Type, &t.Category, &t.Priority, &t.Status, &t.CreatedBy, &t.AssignedTo, &t.CreatedAt, &t.UpdatedAt, &t.ResolvedAt)
	return t, err
}
