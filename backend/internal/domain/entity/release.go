package entity

import (
	"time"

	"github.com/google/uuid"
)

type Release struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   uuid.UUID  `json:"project_id"`
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	ReleaseDate *time.Time `json:"release_date"`
	Status      string     `json:"status"` // Planning, In Progress, Released, Archived
	CreatedBy   uuid.UUID  `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ReleaseRecord struct {
	ID        uuid.UUID `json:"id"`
	ReleaseID uuid.UUID `json:"release_id"`
	RecordID  uuid.UUID `json:"record_id"`
	CreatedAt time.Time `json:"created_at"`
}
