package entity

import (
	"time"

	"github.com/google/uuid"
)

type ApprovalDecision string

const (
	DecisionApproved ApprovalDecision = "approved"
	DecisionRejected ApprovalDecision = "rejected"
)

type ApprovalConfig struct {
	ID         uuid.UUID  `json:"id"`
	TicketType TicketType `json:"ticket_type"`
	Level      int        `json:"level"`
	ApproverID uuid.UUID  `json:"approver_id"`
}

type Approval struct {
	ID         uuid.UUID         `json:"id"`
	TicketID   uuid.UUID         `json:"ticket_id"`
	ApproverID uuid.UUID         `json:"approver_id"`
	Level      int               `json:"level"`
	Decision   *ApprovalDecision `json:"decision"`
	Comment    string            `json:"comment"`
	DecidedAt  *time.Time        `json:"decided_at"`
}
