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

type commentMentionRepository struct {
	db *pgxpool.Pool
}

func NewCommentMentionRepository(db *pgxpool.Pool) repository.CommentMentionRepository {
	return &commentMentionRepository{db: db}
}

func (r *commentMentionRepository) Create(ctx context.Context, mention *entity.CommentMention) error {
	query := `INSERT INTO comment_mentions (id, comment_id, mentioned_user_id, created_at)
		VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query,
		mention.ID, mention.CommentID, mention.MentionedUserID, mention.CreatedAt)
	return err
}

func (r *commentMentionRepository) ListByComment(ctx context.Context, commentID uuid.UUID) ([]*entity.CommentMention, error) {
	query := `SELECT id, comment_id, mentioned_user_id, created_at FROM comment_mentions WHERE comment_id=$1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mentions []*entity.CommentMention
	for rows.Next() {
		mention := &entity.CommentMention{}
		if err := rows.Scan(&mention.ID, &mention.CommentID, &mention.MentionedUserID, &mention.CreatedAt); err != nil {
			return nil, err
		}
		mentions = append(mentions, mention)
	}
	return mentions, nil
}

func (r *commentMentionRepository) DeleteByComment(ctx context.Context, commentID uuid.UUID) error {
	query := `DELETE FROM comment_mentions WHERE comment_id=$1`
	_, err := r.db.Exec(ctx, query, commentID)
	return err
}

func scanCommentMention(row pgx.Row) (*entity.CommentMention, error) {
	mention := &entity.CommentMention{}
	err := row.Scan(&mention.ID, &mention.CommentID, &mention.MentionedUserID, &mention.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}
	return mention, nil
}
