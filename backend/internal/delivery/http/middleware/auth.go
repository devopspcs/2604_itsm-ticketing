package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
	jwtpkg "github.com/org/itsm/pkg/jwt"
	"github.com/org/itsm/pkg/apperror"
)

type contextKey string

const ClaimsKey contextKey = "user_claims"

func JWTAuth(jwtManager *jwtpkg.Manager, userRepo repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				apperror.WriteError(w, apperror.ErrTokenInvalid)
				return
			}
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtManager.ValidateAccessToken(tokenStr)
			if err != nil {
				apperror.WriteError(w, apperror.ErrTokenInvalid)
				return
			}
			userClaims := domainUC.UserClaims{
				Role: claims.Role,
			}
			if id, err := parseUUID(claims.UserID); err == nil {
				userClaims.UserID = id
			}

			// Check if user is still active (immediate revoke on deactivation)
			if userRepo != nil {
				user, err := userRepo.FindByID(r.Context(), userClaims.UserID)
				if err != nil || !user.IsActive {
					apperror.WriteError(w, apperror.ErrForbidden)
					return
				}
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, userClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...entity.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsKey).(domainUC.UserClaims)
			if !ok {
				apperror.WriteError(w, apperror.ErrForbidden)
				return
			}
			for _, role := range roles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			apperror.WriteError(w, apperror.ErrForbidden)
		})
	}
}

func GetClaims(r *http.Request) (domainUC.UserClaims, bool) {
	claims, ok := r.Context().Value(ClaimsKey).(domainUC.UserClaims)
	return claims, ok
}
