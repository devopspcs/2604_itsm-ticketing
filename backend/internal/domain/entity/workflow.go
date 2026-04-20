package entity

import (
	"time"

	"github.com/google/uuid"
)

type Workflow struct {
	ID            uuid.UUID `json:"id"`
	ProjectID     uuid.UUID `json:"project_id" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	InitialStatus string    `json:"initial_status" validate:"required"`
	CreatedAt     time.Time `json:"created_at"`
}

type WorkflowStatus struct {
	ID         uuid.UUID `json:"id"`
	WorkflowID uuid.UUID `json:"workflow_id" validate:"required"`
	StatusName string    `json:"status_name" validate:"required"`
	StatusOrder int      `json:"status_order" validate:"required,min=0"`
	CreatedAt  time.Time `json:"created_at"`
}

type WorkflowTransition struct {
	ID             uuid.UUID `json:"id"`
	WorkflowID     uuid.UUID `json:"workflow_id" validate:"required"`
	FromStatusID   uuid.UUID `json:"from_status_id" validate:"required"`
	ToStatusID     uuid.UUID `json:"to_status_id" validate:"required"`
	ValidationRule string    `json:"validation_rule"`
	CreatedAt      time.Time `json:"created_at"`
}
