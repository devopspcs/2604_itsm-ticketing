package entity

import (
	"time"

	"github.com/google/uuid"
)

type Sprint struct {
	ID              uuid.UUID  `json:"id"`
	ProjectID       uuid.UUID  `json:"project_id" validate:"required"`
	Name            string     `json:"name" validate:"required"`
	Goal            string     `json:"goal"`
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date" validate:"required"`
	Status          string     `json:"status" validate:"required,oneof=Planned Active Completed"`
	ActualStartDate *time.Time `json:"actual_start_date"`
	ActualEndDate   *time.Time `json:"actual_end_date"`
	CreatedAt       time.Time  `json:"created_at"`
}

type SprintRecord struct {
	ID        uuid.UUID `json:"id"`
	SprintID  uuid.UUID `json:"sprint_id" validate:"required"`
	RecordID  uuid.UUID `json:"record_id" validate:"required"`
	Priority  int       `json:"priority" validate:"required,min=0"`
	CreatedAt time.Time `json:"created_at"`
}
