package entity

import (
	"time"

	"github.com/google/uuid"
)

type CustomField struct {
	ID         uuid.UUID `json:"id"`
	ProjectID  uuid.UUID `json:"project_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	FieldType  string    `json:"field_type" validate:"required,oneof=text textarea dropdown multiselect date number checkbox"`
	IsRequired bool      `json:"is_required"`
	CreatedAt  time.Time `json:"created_at"`
}

type CustomFieldOption struct {
	ID          uuid.UUID `json:"id"`
	FieldID     uuid.UUID `json:"field_id" validate:"required"`
	OptionValue string    `json:"option_value" validate:"required"`
	OptionOrder int       `json:"option_order" validate:"required,min=0"`
	CreatedAt   time.Time `json:"created_at"`
}

type CustomFieldValue struct {
	ID        uuid.UUID `json:"id"`
	RecordID  uuid.UUID `json:"record_id" validate:"required"`
	FieldID   uuid.UUID `json:"field_id" validate:"required"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
