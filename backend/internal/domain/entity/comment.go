package entity

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID         uuid.UUID `json:"id"`
	RecordID   uuid.UUID `json:"record_id" validate:"required"`
	AuthorID   uuid.UUID `json:"author_id" validate:"required"`
	AuthorName string    `json:"author_name"`
	Text       string    `json:"text" validate:"required"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CommentMention struct {
	ID              uuid.UUID `json:"id"`
	CommentID       uuid.UUID `json:"comment_id" validate:"required"`
	MentionedUserID uuid.UUID `json:"mentioned_user_id" validate:"required"`
	CreatedAt       time.Time `json:"created_at"`
}
