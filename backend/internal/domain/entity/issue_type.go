package entity

import (
	"time"

	"github.com/google/uuid"
)

type IssueType struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type IssueTypeScheme struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type IssueTypeSchemeItem struct {
	ID          uuid.UUID `json:"id"`
	SchemeID    uuid.UUID `json:"scheme_id" validate:"required"`
	IssueTypeID uuid.UUID `json:"issue_type_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
}
