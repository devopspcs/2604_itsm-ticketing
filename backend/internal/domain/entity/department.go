package entity

import (
	"time"

	"github.com/google/uuid"
)

type Division struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Department struct {
	ID         uuid.UUID `json:"id"`
	DivisionID uuid.UUID `json:"division_id"`
	Name       string    `json:"name"`
	Code       string    `json:"code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Team struct {
	ID           uuid.UUID `json:"id"`
	DepartmentID uuid.UUID `json:"department_id"`
	Name         string    `json:"name"`
	Email        *string   `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
