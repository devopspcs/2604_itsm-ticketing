package entity

import (
	"time"

	"github.com/google/uuid"
)

type Component struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   uuid.UUID  `json:"project_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	LeadUserID  *uuid.UUID `json:"lead_user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
