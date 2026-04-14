package entity

import (
	"time"

	"github.com/google/uuid"
)

type TicketType string

const (
	TypeChangeRequest   TicketType = "change_request"
	TypeIncident        TicketType = "incident"
	TypeHelpdeskRequest TicketType = "helpdesk_request"
)

type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

type TicketStatus string

const (
	StatusOpen            TicketStatus = "open"
	StatusInProgress      TicketStatus = "in_progress"
	StatusPendingApproval TicketStatus = "pending_approval"
	StatusApproved        TicketStatus = "approved"
	StatusRejected        TicketStatus = "rejected"
	StatusDone            TicketStatus = "done"
)

type Ticket struct {
	ID          uuid.UUID    `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Type        TicketType   `json:"type"`
	Category    string       `json:"category"`
	Priority    Priority     `json:"priority"`
	Status      TicketStatus `json:"status"`
	CreatedBy   uuid.UUID    `json:"created_by"`
	AssignedTo  *uuid.UUID   `json:"assigned_to"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
