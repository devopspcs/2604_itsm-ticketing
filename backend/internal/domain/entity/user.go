package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleUser     Role = "user"
	RoleApprover Role = "approver"
	RoleAdmin    Role = "admin"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	FullName     string     `json:"full_name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         Role       `json:"role"`
	IsActive     bool       `json:"is_active"`
	DepartmentID *uuid.UUID `json:"department_id"`
	DivisionID   *uuid.UUID `json:"division_id"`
	TeamID       *uuid.UUID `json:"team_id"`
	Position     *Position  `json:"position"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
