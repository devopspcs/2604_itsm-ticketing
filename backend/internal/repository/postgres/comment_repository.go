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

type commentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) repository.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *entity.Comment) error {
	query := `INSERT INTO comments (id, record_id, author_id, text, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query,
		comment.ID, comment.RecordID, comment.AuthorID, comment.Text, comment.CreatedAt, comment.UpdatedAt)
	return err
}

func (r *commentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Comment, error) {
	query := `SELECT id, record_id, author_id, text, created_at, updated_at FROM comments WHERE id=$1`
	row := r.db.QueryRow(ctx, query, id)
	return scanComment(row)
}

func (r *commentRepository) ListByRecord(ctx context.Context, recordID uuid.UUID) ([]*entity.Comment, error) {
	query := `SELECT c.id, c.record_id, c.author_id, COALESCE(u.full_name, u.email, '') as author_name, c.text, c.created_at, c.updated_at
		FROM comments c LEFT JOIN users u ON c.author_id = u.id
		WHERE c.record_id=$1 ORDER BY c.created_at ASC`
	rows, err := r.db.Query(ctx, query, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		comment := &entity.Comment{}
		if err := rows.Scan(&comment.ID, &comment.RecordID, &comment.AuthorID, &comment.AuthorName, &comment.Text, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *commentRepository) Update(ctx context.Context, comment *entity.Comment) error {
	query := `UPDATE comments SET text=$1, updated_at=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, comment.Text, comment.UpdatedAt, comment.ID)
	return err
}

func (r *commentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM comments WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func scanComment(row pgx.Row) (*entity.Comment, error) {
	comment := &entity.Comment{}
	err := row.Scan(&comment.ID, &comment.RecordID, &comment.AuthorID, &comment.Text, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return comment, nil
}
