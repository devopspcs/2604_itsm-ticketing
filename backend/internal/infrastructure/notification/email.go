package notification

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/org/itsm/pkg/config"
	"github.com/org/itsm/pkg/logger"
)

type EmailService struct {
	host    string
	port    string
	user    string
	pass    string
	from    string
	baseURL string
	log     *logger.Logger
}

func NewEmailService(cfg *config.Config, log *logger.Logger) *EmailService {
	return &EmailService{
		host:    cfg.SMTPHost,
		port:    cfg.SMTPPort,
		user:    cfg.SMTPUser,
		pass:    cfg.SMTPPass,
		from:    cfg.EmailFrom,
		baseURL: cfg.BaseURL,
		log:     log,
	}
}

func (s *EmailService) IsConfigured() bool {
	return s.host != "" && s.user != "" && s.pass != ""
}

func (s *EmailService) SendTicketCreated(toEmail, toName, ticketTitle, ticketID, ticketType, priority string) {
	subject := fmt.Sprintf("[PCS ITSM] New Ticket Created: %s", ticketTitle)
	body := fmt.Sprintf(`
<html>
<body style="font-family: Inter, Arial, sans-serif; background: #f8f9fa; padding: 20px;">
<div style="max-width: 600px; margin: 0 auto; background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
  <div style="background: #00307d; padding: 24px 32px; color: #fff;">
    <h1 style="margin: 0; font-size: 20px;">🆕 New Ticket Created</h1>
  </div>
  <div style="padding: 32px;">
    <p>Hi <strong>%s</strong>,</p>
    <p>A new ticket has been created in the PCS Ticketing System:</p>
    <table style="width: 100%%; border-collapse: collapse; margin: 16px 0;">
      <tr><td style="padding: 8px 0; color: #737783;">Title:</td><td style="padding: 8px 0; font-weight: 600;">%s</td></tr>
      <tr><td style="padding: 8px 0; color: #737783;">Type:</td><td style="padding: 8px 0;">%s</td></tr>
      <tr><td style="padding: 8px 0; color: #737783;">Priority:</td><td style="padding: 8px 0;">%s</td></tr>
    </table>
    <a href="%s/tickets/%s" style="display: inline-block; background: #00307d; color: #fff; padding: 12px 24px; border-radius: 8px; text-decoration: none; font-weight: 600; margin-top: 16px;">View Ticket</a>
  </div>
  <div style="padding: 16px 32px; background: #f3f4f5; color: #737783; font-size: 12px;">
    PCS Ticketing System · IT Service Management
  </div>
</div>
</body>
</html>`, toName, ticketTitle, ticketType, priority, s.baseURL, ticketID)

	go s.send(toEmail, subject, body)
}

func (s *EmailService) SendTicketAssigned(toEmail, toName, ticketTitle, ticketID, ticketType, priority, category, assignedBy string) {
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

	subject := fmt.Sprintf("[PCS ITSM] Ticket Assigned: %s-%s - %s", prefix, shortID, ticketTitle)
	body := fmt.Sprintf(`
<html>
<body style="font-family: Inter, Arial, sans-serif; background: #f8f9fa; padding: 20px;">
<div style="max-width: 600px; margin: 0 auto; background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
  <div style="background: #0045ab; padding: 24px 32px; color: #fff;">
    <h1 style="margin: 0; font-size: 20px;">👤 Ticket Assigned to You</h1>
  </div>
  <div style="padding: 32px;">
    <p>Hi <strong>%s</strong>,</p>
    <p>A ticket has been assigned to you by <strong>%s</strong>.</p>
    <p>Please review and take action on this ticket.</p>
    <table style="width: 100%%; border-collapse: collapse; margin: 20px 0; background: #f8f9fa; border-radius: 8px;">
      <tr><td style="padding: 12px 16px; border-bottom: 1px solid #e2e8f0; color: #737783; font-size: 13px; width: 140px;">Ticket Number</td><td style="padding: 12px 16px; border-bottom: 1px solid #e2e8f0; font-weight: 600;">%s-%s</td></tr>
      <tr><td style="padding: 12px 16px; border-bottom: 1px solid #e2e8f0; color: #737783; font-size: 13px;">Title</td><td style="padding: 12px 16px; border-bottom: 1px solid #e2e8f0; font-weight: 600;">%s</td></tr>
      <tr><td style="padding: 12px 16px; border-bottom: 1px solid #e2e8f0; color: #737783; font-size: 13px;">Type</td><td style="padding: 12px 16px; border-bottom: 1px solid #e2e8f0;">%s</td></tr>
      <tr><td style="padding: 12px 16px; border-bottom: 1px solid #e2e8f0; color: #737783; font-size: 13px;">Category</td><td style="padding: 12px 16px; border-bottom: 1px solid #e2e8f0;">%s</td></tr>
      <tr><td style="padding: 12px 16px; border-bottom: 1px solid #e2e8f0; color: #737783; font-size: 13px;">Priority</td><td style="padding: 12px 16px; border-bottom: 1px solid #e2e8f0; font-weight: 600; color: %s;">%s</td></tr>
      <tr><td style="padding: 12px 16px; color: #737783; font-size: 13px;">Assigned By</td><td style="padding: 12px 16px;">%s</td></tr>
    </table>
    <a href="%s/tickets/%s" style="display: inline-block; background: #00307d; color: #fff; padding: 12px 24px; border-radius: 8px; text-decoration: none; font-weight: 600;">View Ticket Details</a>
  </div>
  <div style="padding: 16px 32px; background: #f3f4f5; color: #737783; font-size: 12px;">
    PCS Ticketing System · IT Service Management
  </div>
</div>
</body>
</html>`,
		toName, assignedBy,
		prefix, shortID, ticketTitle, ticketType, category,
		priorityColor(priority), priority, assignedBy,
		s.baseURL, ticketID)

	go s.send(toEmail, subject, body)
}

func (s *EmailService) SendApprovalRequested(toEmail, toName, ticketTitle, ticketID, requestedBy string) {
	subject := fmt.Sprintf("[PCS ITSM] Approval Required: %s", ticketTitle)
	body := fmt.Sprintf(`
<html>
<body style="font-family: Inter, Arial, sans-serif; background: #f8f9fa; padding: 20px;">
<div style="max-width: 600px; margin: 0 auto; background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
  <div style="background: #f59e0b; padding: 24px 32px; color: #fff;">
    <h1 style="margin: 0; font-size: 20px;">⏳ Approval Required</h1>
  </div>
  <div style="padding: 32px;">
    <p>Hi <strong>%s</strong>,</p>
    <p><strong>%s</strong> has submitted a ticket for your approval:</p>
    <div style="background: #fef3c7; padding: 16px; border-radius: 8px; border-left: 4px solid #f59e0b; margin: 16px 0;">
      <p style="margin: 0; font-weight: 600; font-size: 16px;">%s</p>
    </div>
    <p>Please review and approve or reject this request.</p>
    <a href="%s/tickets/%s" style="display: inline-block; background: #00307d; color: #fff; padding: 12px 24px; border-radius: 8px; text-decoration: none; font-weight: 600; margin-top: 16px;">Review Ticket</a>
  </div>
  <div style="padding: 16px 32px; background: #f3f4f5; color: #737783; font-size: 12px;">
    PCS Ticketing System · IT Service Management
  </div>
</div>
</body>
</html>`, toName, requestedBy, ticketTitle, s.baseURL, ticketID)

	go s.send(toEmail, subject, body)
}

func (s *EmailService) SendApprovalDecided(toEmail, toName, ticketTitle, ticketID, decision, decidedBy string) {
	emoji := "✅"
	color := "#10b981"
	if decision == "rejected" {
		emoji = "❌"
		color = "#ba1a1a"
	}
	subject := fmt.Sprintf("[PCS ITSM] Ticket %s: %s", strings.Title(decision), ticketTitle)
	body := fmt.Sprintf(`
<html>
<body style="font-family: Inter, Arial, sans-serif; background: #f8f9fa; padding: 20px;">
<div style="max-width: 600px; margin: 0 auto; background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
  <div style="background: %s; padding: 24px 32px; color: #fff;">
    <h1 style="margin: 0; font-size: 20px;">%s Ticket %s</h1>
  </div>
  <div style="padding: 32px;">
    <p>Hi <strong>%s</strong>,</p>
    <p>Your ticket has been <strong>%s</strong> by <strong>%s</strong>:</p>
    <div style="background: #f3f4f5; padding: 16px; border-radius: 8px; margin: 16px 0;">
      <p style="margin: 0; font-weight: 600; font-size: 16px;">%s</p>
    </div>
    <a href="%s/tickets/%s" style="display: inline-block; background: #00307d; color: #fff; padding: 12px 24px; border-radius: 8px; text-decoration: none; font-weight: 600; margin-top: 16px;">View Ticket</a>
  </div>
  <div style="padding: 16px 32px; background: #f3f4f5; color: #737783; font-size: 12px;">
    PCS Ticketing System · IT Service Management
  </div>
</div>
</body>
</html>`, color, emoji, strings.Title(decision), toName, decision, decidedBy, ticketTitle, s.baseURL, ticketID)

	go s.send(toEmail, subject, body)
}

func (s *EmailService) send(to, subject, htmlBody string) {
	if !s.IsConfigured() {
		s.log.Warn("email not configured, skipping", "to", to, "subject", subject)
		return
	}

	auth := smtp.PlainAuth("", s.user, s.pass, s.host)

	headers := fmt.Sprintf("From: PCS Ticketing <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		s.from, to, subject)

	msg := []byte(headers + htmlBody)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	if err := smtp.SendMail(addr, auth, s.from, []string{to}, msg); err != nil {
		s.log.Error("failed to send email", "to", to, "error", err)
	} else {
		s.log.Info("email sent", "to", to, "subject", subject)
	}
}

func priorityColor(p string) string {
	switch p {
	case "critical":
		return "#ba1a1a"
	case "high":
		return "#f59e0b"
	case "medium":
		return "#3b82f6"
	default:
		return "#10b981"
	}
}
