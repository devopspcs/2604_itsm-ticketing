package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/org/itsm/internal/delivery/http/handler"
	"github.com/org/itsm/internal/delivery/http/middleware"
	"github.com/org/itsm/internal/domain/entity"
	jwtpkg "github.com/org/itsm/pkg/jwt"
)

type Handlers struct {
	Auth         *handler.AuthHandler
	User         *handler.UserHandler
	Ticket       *handler.TicketHandler
	Approval     *handler.ApprovalHandler
	Dashboard    *handler.DashboardHandler
	Notification *handler.NotificationHandler
	Webhook      *handler.WebhookHandler
	Attachment   *handler.AttachmentHandler
	Org          *handler.OrgHandler
	SSO          *handler.SSOHandler
}

func NewRouter(h *Handlers, jwtManager *jwtpkg.Manager, db interface{ Ping() error }) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, req *http.Request) {
		status := "ok"
		if err := db.Ping(); err != nil {
			status = "db_error"
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"` + status + `"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes
		r.Post("/auth/login", h.Auth.Login)
		r.Post("/auth/refresh", h.Auth.Refresh)

		// SSO routes (public)
		r.Get("/auth/sso/login-url", h.SSO.GetLoginURL)
		r.Get("/auth/sso/redirect", h.SSO.Redirect)
		r.Get("/auth/sso/callback", h.SSO.Callback)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(jwtManager))

			r.Post("/auth/logout", h.Auth.Logout)

			// Notifications
			r.Get("/notifications", h.Notification.List)
			r.Patch("/notifications/{id}/read", h.Notification.MarkRead)

			// Tickets
			r.Get("/tickets", h.Ticket.List)
			r.Post("/tickets", h.Ticket.Create)
			r.Get("/tickets/{id}", h.Ticket.Get)
			r.Patch("/tickets/{id}", h.Ticket.Update)
			r.Post("/tickets/{id}/submit", h.Ticket.Submit)
			r.Post("/tickets/{id}/notes", h.Attachment.AddNote)
			r.Get("/tickets/{id}/notes", h.Attachment.ListNotes)
			r.Get("/tickets/{id}/notes/{noteId}/image", h.Attachment.GetNoteImage)
			r.Get("/tickets/{id}/approvals", h.Ticket.GetApprovals)
			r.Get("/tickets/{id}/activities", h.Ticket.GetActivities)
			r.Post("/tickets/{id}/attachments", h.Attachment.Upload)
			r.Get("/tickets/{id}/attachments", h.Attachment.List)
			r.Get("/tickets/{id}/attachments/{attachmentId}", h.Attachment.Download)
			r.Delete("/tickets/{id}/attachments/{attachmentId}", h.Attachment.Delete)

			// Assign — all authenticated users can assign
			r.Post("/tickets/{id}/assign", h.Ticket.Assign)

			// Approvals — approver or admin only
			r.With(middleware.RequireRole(entity.RoleAdmin, entity.RoleApprover)).
				Post("/approvals/decide", h.Approval.Decide)

			// Dashboard — accessible to all authenticated users
			r.Get("/dashboard", h.Dashboard.GetStats)

			// User list for assign dropdown — all authenticated users can see
			r.Get("/users/list", h.User.List)

			// Org structure — list endpoints (authenticated)
			r.Get("/departments", h.Org.ListDepartments)
			r.Get("/divisions", h.Org.ListDivisions)
			r.Get("/teams", h.Org.ListTeams)

			// Admin-only routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole(entity.RoleAdmin))

				r.Get("/users", h.User.List)
				r.Post("/users", h.User.Create)
				r.Patch("/users/{id}/role", h.User.UpdateRole)
				r.Patch("/users/{id}/deactivate", h.User.Deactivate)
				r.Patch("/users/{id}/activate", h.User.Activate)
				r.Patch("/users/{id}/org", h.User.UpdateUserOrg)

				// Org structure — CUD endpoints (admin only)
				r.Post("/departments", h.Org.CreateDepartment)
				r.Patch("/departments/{id}", h.Org.UpdateDepartment)
				r.Delete("/departments/{id}", h.Org.DeleteDepartment)

				r.Post("/divisions", h.Org.CreateDivision)
				r.Patch("/divisions/{id}", h.Org.UpdateDivision)
				r.Delete("/divisions/{id}", h.Org.DeleteDivision)

				r.Post("/teams", h.Org.CreateTeam)
				r.Patch("/teams/{id}", h.Org.UpdateTeam)
				r.Delete("/teams/{id}", h.Org.DeleteTeam)

				r.Get("/webhooks", h.Webhook.List)
				r.Post("/webhooks", h.Webhook.Create)
				r.Patch("/webhooks/{id}", h.Webhook.Update)
				r.Delete("/webhooks/{id}", h.Webhook.Delete)
			})
		})
	})

	return r
}
