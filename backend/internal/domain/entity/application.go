package entity

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserAppAccess struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	AppID     uuid.UUID  `json:"app_id"`
	Role      string     `json:"role"`
	GrantedAt time.Time  `json:"granted_at"`
	GrantedBy *uuid.UUID `json:"granted_by"`
}
