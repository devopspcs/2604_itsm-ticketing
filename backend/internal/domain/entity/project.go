package entity

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IconColor string    `json:"icon_color"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectColumn struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}

type ProjectRecord struct {
	ID          uuid.UUID    `json:"id"`
	ColumnID    uuid.UUID    `json:"column_id"`
	ProjectID   uuid.UUID    `json:"project_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	AssignedTo  *uuid.UUID   `json:"assigned_to"`
	Assignees   []uuid.UUID  `json:"assignees"`
	DueDate     *time.Time   `json:"due_date"`
	Position    int          `json:"position"`
	IsCompleted bool         `json:"is_completed"`
	CompletedAt *time.Time   `json:"completed_at"`
	CreatedBy   uuid.UUID    `json:"created_by"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type ProjectActivityLog struct {
	ID        uuid.UUID  `json:"id"`
	ProjectID uuid.UUID  `json:"project_id"`
	RecordID  *uuid.UUID `json:"record_id"`
	ActorID   uuid.UUID  `json:"actor_id"`
	Action    string     `json:"action"`
	Detail    string     `json:"detail"`
	CreatedAt time.Time  `json:"created_at"`
}

type ProjectMemberRole string

const (
	ProjectRoleOwner  ProjectMemberRole = "owner"
	ProjectRoleMember ProjectMemberRole = "member"
)

type ProjectMember struct {
	ProjectID uuid.UUID         `json:"project_id"`
	UserID    uuid.UUID         `json:"user_id"`
	Role      ProjectMemberRole `json:"role"`
	CreatedAt time.Time         `json:"created_at"`
}
