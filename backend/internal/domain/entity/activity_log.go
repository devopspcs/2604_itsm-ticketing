package entity

import (
	"time"

	"github.com/google/uuid"
)

type ActivityAction string

const (
	ActionTicketCreated     ActivityAction = "ticket_created"
	ActionStatusChanged     ActivityAction = "status_changed"
	ActionAssigned          ActivityAction = "assigned"
	ActionReassigned        ActivityAction = "reassigned"
	ActionApprovalRequested ActivityAction = "approval_requested"
	ActionApprovalDecided   ActivityAction = "approval_decided"
	ActionFieldUpdated      ActivityAction = "field_updated"
)

type ActivityLog struct {
	ID        uuid.UUID      `json:"id"`
	TicketID  uuid.UUID      `json:"ticket_id"`
	ActorID   uuid.UUID      `json:"actor_id"`
	Action    ActivityAction `json:"action"`
	OldValue  *string        `json:"old_value"`
	NewValue  *string        `json:"new_value"`
	CreatedAt time.Time      `json:"created_at"`
}
