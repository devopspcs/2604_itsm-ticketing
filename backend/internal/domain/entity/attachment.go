package entity

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	ID           uuid.UUID `json:"id"`
	RecordID     uuid.UUID `json:"record_id" validate:"required"`
	FileName     string    `json:"file_name" validate:"required"`
	FileSize     int64     `json:"file_size" validate:"required,min=1"`
	FileType     string    `json:"file_type"`
	FilePath     string    `json:"file_path" validate:"required"`
	UploaderID   uuid.UUID `json:"uploader_id" validate:"required"`
	UploaderName string    `json:"uploader_name"`
	CreatedAt    time.Time `json:"created_at"`
}
