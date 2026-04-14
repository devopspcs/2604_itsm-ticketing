package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
)

type Dispatcher struct {
	webhookRepo repository.WebhookRepository
	userRepo    repository.UserRepository
	client      *http.Client
	baseURL     string
}

func NewDispatcher(webhookRepo repository.WebhookRepository, userRepo repository.UserRepository, baseURL string) *Dispatcher {
	return &Dispatcher{
		webhookRepo: webhookRepo,
		userRepo:    userRepo,
		client:      &http.Client{Timeout: 10 * time.Second},
		baseURL:     baseURL,
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, event entity.WebhookEvent, payload interface{}, configs []*entity.WebhookConfig) error {
	// Enrich payload with user names before marshaling
	enriched := d.enrichPayload(ctx, payload)
	body, err := json.Marshal(enriched)
	if err != nil {
		return err
	}

	for _, config := range configs {
		go d.sendWithRetry(ctx, config, event, body)
	}
	return nil
}

func (d *Dispatcher) sendWithRetry(ctx context.Context, config *entity.WebhookConfig, event entity.WebhookEvent, body []byte) {
	delays := []time.Duration{0, time.Second, 2 * time.Second, 4 * time.Second}
	for attempt, delay := range delays {
		if delay > 0 {
			time.Sleep(delay)
		}
		status, err := d.send(config, event, body)
		logEntry := &entity.WebhookLog{
			ID:        uuid.New(),
			WebhookID: config.ID,
			Event:     event,
			Payload:   body,
			Attempt:   attempt + 1,
			SentAt:    time.Now().UTC(),
		}
		if status > 0 {
			logEntry.ResponseStatus = &status
		}
		_ = d.webhookRepo.SaveLog(ctx, logEntry)
		if err == nil && status >= 200 && status < 300 {
			return
		}
	}
}

func (d *Dispatcher) send(config *entity.WebhookConfig, event entity.WebhookEvent, body []byte) (int, error) {
	var sendBody []byte

	// Detect Google Chat webhook URL and format accordingly
	if isGoogleChatURL(config.URL) {
		sendBody = formatGoogleChat(event, body, d.baseURL)
	} else if isSlackURL(config.URL) {
		sendBody = formatSlack(event, body)
	} else {
		sendBody = body
	}

	req, err := http.NewRequest(http.MethodPost, config.URL, bytes.NewReader(sendBody))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ITSM-Event", string(event))
	req.Header.Set("X-ITSM-Signature", computeHMAC(sendBody, config.SecretKey))

	resp, err := d.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func isGoogleChatURL(url string) bool {
	return len(url) > 0 && (contains(url, "chat.googleapis.com") || contains(url, "chat.google.com"))
}

func isSlackURL(url string) bool {
	return len(url) > 0 && contains(url, "hooks.slack.com")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func formatGoogleChat(event entity.WebhookEvent, rawPayload []byte, baseURL string) []byte {
	var data map[string]interface{}
	json.Unmarshal(rawPayload, &data)

	title := getStr(data, "title")
	priority := getStr(data, "priority")
	ticketType := getStr(data, "type")
	category := getStr(data, "category")
	ticketID := getStr(data, "id")
	createdBy := getStr(data, "created_by")
	assignedTo := getStr(data, "assigned_to")

	shortID := ticketID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}

	prefix := "REQ"
	switch ticketType {
	case "incident":
		prefix = "INC"
	case "change_request":
		prefix = "CHG"
	}

	// Event header
	eventLabel := ""
	switch event {
	case entity.EventTicketCreated:
		eventLabel = "New Ticket Has Been Created"
	case entity.EventTicketStatusChanged:
		eventLabel = "Ticket Status Changed"
	case entity.EventTicketAssigned:
		eventLabel = "Ticket Has Been Assigned"
	case entity.EventApprovalRequested:
		eventLabel = "Approval Request Submitted"
	case entity.EventApprovalDecided:
		eventLabel = "Approval Decision Made"
	default:
		eventLabel = string(event)
	}

	text := fmt.Sprintf("*PCS Ticketing System*\n%s\n\n", eventLabel)
	text += fmt.Sprintf("*Ticket Number:* %s-%s\n\n", prefix, shortID)
	text += fmt.Sprintf("*Title:* %s\n\n", title)

	if category != "" {
		text += fmt.Sprintf("*Category:* %s\n\n", category)
	}
	if ticketType != "" {
		text += fmt.Sprintf("*Type:* %s\n\n", ticketType)
	}
	if priority != "" {
		text += fmt.Sprintf("*Priority:* %s\n\n", priority)
	}
	if createdBy != "" {
		text += fmt.Sprintf("*Requestor:* %s\n\n", createdBy)
	}
	if assignedTo != "" && assignedTo != "<nil>" {
		text += fmt.Sprintf("*Assigned To:* %s\n\n", assignedTo)
	}

	text += fmt.Sprintf("*Date & Time:* %s\n\n", time.Now().UTC().Format("January 02 2006 15:04:05"))

	if ticketID != "" && baseURL != "" {
		text += fmt.Sprintf("<%s/tickets/%s|Go to Ticket Details>", baseURL, ticketID)
	}

	gcPayload := map[string]string{"text": text}
	b, _ := json.Marshal(gcPayload)
	return b
}

func formatSlack(event entity.WebhookEvent, rawPayload []byte) []byte {
	var data map[string]interface{}
	json.Unmarshal(rawPayload, &data)

	title := getStr(data, "title")
	status := getStr(data, "status")

	text := fmt.Sprintf("*ITSM %s*\nTitle: %s\nStatus: %s", string(event), title, status)
	slackPayload := map[string]string{"text": text}
	b, _ := json.Marshal(slackPayload)
	return b
}

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// enrichPayload resolves UUID fields to user names for better webhook readability
func (d *Dispatcher) enrichPayload(ctx context.Context, payload interface{}) interface{} {
	if d.userRepo == nil {
		return payload
	}

	// Marshal to map, resolve UUIDs, return enriched map
	raw, err := json.Marshal(payload)
	if err != nil {
		return payload
	}

	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return payload
	}

	// Resolve created_by
	if createdBy, ok := m["created_by"].(string); ok && createdBy != "" {
		if name := d.resolveUserName(ctx, createdBy); name != "" {
			m["created_by"] = name
			m["created_by_id"] = createdBy
		}
	}

	// Resolve assigned_to
	if assignedTo, ok := m["assigned_to"].(string); ok && assignedTo != "" {
		if name := d.resolveUserName(ctx, assignedTo); name != "" {
			m["assigned_to"] = name
			m["assigned_to_id"] = assignedTo
		}
	}

	return m
}

func (d *Dispatcher) resolveUserName(ctx context.Context, uuidStr string) string {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return ""
	}
	user, err := d.userRepo.FindByID(ctx, id)
	if err != nil {
		return ""
	}
	return user.FullName
}

func computeHMAC(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return fmt.Sprintf("sha256=%x", mac.Sum(nil))
}
