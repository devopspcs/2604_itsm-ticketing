package notification

// EmailSender is the interface for sending email notifications
type EmailSender interface {
	IsConfigured() bool
	SendTicketCreated(toEmail, toName, ticketTitle, ticketID, ticketType, priority string)
	SendTicketAssigned(toEmail, toName, ticketTitle, ticketID, ticketType, priority, category, assignedBy string)
	SendApprovalRequested(toEmail, toName, ticketTitle, ticketID, requestedBy string)
	SendApprovalDecided(toEmail, toName, ticketTitle, ticketID, decision, decidedBy string)
}
