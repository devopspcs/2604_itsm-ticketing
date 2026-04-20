package entity

import (
	"time"

	"github.com/google/uuid"
)

type Label struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type RecordLabel struct {
	ID        uuid.UUID `json:"id"`
	RecordID  uuid.UUID `json:"record_id" validate:"required"`
	LabelID   uuid.UUID `json:"label_id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}
