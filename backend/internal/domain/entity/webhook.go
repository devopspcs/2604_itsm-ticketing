package entity

import (
	"time"

	"github.com/google/uuid"
)

type WebhookEvent string

const (
	EventTicketCreated       WebhookEvent = "ticket.created"
	EventTicketStatusChanged WebhookEvent = "ticket.status_changed"
	EventTicketAssigned      WebhookEvent = "ticket.assigned"
	EventApprovalRequested   WebhookEvent = "approval.requested"
	EventApprovalDecided     WebhookEvent = "approval.decided"
)

type WebhookConfig struct {
	ID        uuid.UUID      `json:"id"`
	URL       string         `json:"url"`
	Events    []WebhookEvent `json:"events"`
	SecretKey string         `json:"secret_key"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type WebhookLog struct {
	ID             uuid.UUID    `json:"id"`
	WebhookID      uuid.UUID    `json:"webhook_id"`
	Event          WebhookEvent `json:"event"`
	Payload        []byte       `json:"payload"`
	ResponseStatus *int         `json:"response_status"`
	Attempt        int          `json:"attempt"`
	SentAt         time.Time    `json:"sent_at"`
}
