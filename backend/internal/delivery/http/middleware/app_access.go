package middleware

import (
	"net/http"

	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	"github.com/org/itsm/pkg/apperror"
)

// RequireAppAccess creates a middleware that checks if the authenticated user
// has access to the specified application (by app code).
func RequireAppAccess(appRepo repository.ApplicationRepository, accessRepo repository.UserAppAccessRepository, appCode string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r)
			if !ok {
				apperror.WriteError(w, apperror.ErrForbidden)
				return
			}

			// Admin always has access
			if claims.Role == entity.RoleAdmin {
				next.ServeHTTP(w, r)
				return
			}

			// Find the application by code
			app, err := appRepo.FindByCode(r.Context(), appCode)
			if err != nil {
				// If app not found, allow access (app not configured yet)
				next.ServeHTTP(w, r)
				return
			}

			// Check if user has access
			hasAccess, err := accessRepo.HasAccess(r.Context(), claims.UserID, app.ID)
			if err != nil || !hasAccess {
				apperror.WriteError(w, apperror.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
