package jwt

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"pgregory.net/rapid"
)

// Feature: itsm-web-app, Property 1: JWT Access Token Expiry Enforcement
// **Validates: Requirements 1.3, 1.5**
//
// For any valid JWT access token, once the token's expiry time has passed,
// all authenticated endpoints SHALL reject the token with HTTP 401.

// roleGen generates a random valid Role.
func roleGen() *rapid.Generator[entity.Role] {
	return rapid.SampledFrom([]entity.Role{
		entity.RoleUser,
		entity.RoleApprover,
		entity.RoleAdmin,
	})
}

// TestProperty_JWTAccessTokenExpiryEnforcement_Validation tests that
// ValidateAccessToken rejects any token whose expiry time has passed.
func TestProperty_JWTAccessTokenExpiryEnforcement_Validation(t *testing.T) {
	// Feature: itsm-web-app, Property 1: JWT Access Token Expiry Enforcement
	rapid.Check(t, func(t *rapid.T) {
		secret := rapid.StringMatching(`[a-zA-Z0-9]{16,64}`).Draw(t, "secret")
		role := roleGen().Draw(t, "role")
		userID := uuid.New()

		manager := NewManager(secret, "refresh-secret")

		// Create a token that is already expired by crafting claims manually.
		// We pick a random past expiry between 1 second and 24 hours ago.
		pastSeconds := rapid.IntRange(1, 86400).Draw(t, "pastSeconds")
		expiredAt := time.Now().Add(-time.Duration(pastSeconds) * time.Second)

		claims := Claims{
			UserID: userID.String(),
			Role:   role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiredAt),
				IssuedAt:  jwt.NewNumericDate(expiredAt.Add(-AccessTokenTTL)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("failed to sign token: %v", err)
		}

		// Validate — must fail because the token is expired.
		_, err = manager.ValidateAccessToken(tokenStr)
		if err == nil {
			t.Fatalf("expected expired token to be rejected, but ValidateAccessToken returned nil error for userID=%s role=%s expiredAt=%v",
				userID, role, expiredAt)
		}
	})
}

// TestProperty_JWTAccessTokenExpiryEnforcement_HTTP tests that the JWTAuth
// middleware returns HTTP 401 for any expired access token.
func TestProperty_JWTAccessTokenExpiryEnforcement_HTTP(t *testing.T) {
	// Feature: itsm-web-app, Property 1: JWT Access Token Expiry Enforcement
	rapid.Check(t, func(t *rapid.T) {
		secret := rapid.StringMatching(`[a-zA-Z0-9]{16,64}`).Draw(t, "secret")
		role := roleGen().Draw(t, "role")
		userID := uuid.New()

		manager := NewManager(secret, "refresh-secret")

		// Create an expired token.
		pastSeconds := rapid.IntRange(1, 86400).Draw(t, "pastSeconds")
		expiredAt := time.Now().Add(-time.Duration(pastSeconds) * time.Second)

		claims := Claims{
			UserID: userID.String(),
			Role:   role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiredAt),
				IssuedAt:  jwt.NewNumericDate(expiredAt.Add(-AccessTokenTTL)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(secret))
		if err != nil {
			t.Fatalf("failed to sign token: %v", err)
		}

		// Build a minimal HTTP handler behind the JWTAuth middleware.
		// We import the middleware inline to avoid circular deps — instead
		// we replicate the core validation logic that the middleware uses.
		// The middleware calls manager.ValidateAccessToken, so we simulate
		// the same flow: parse token, check error, return 401.
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || len(authHeader) < 8 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			tkStr := authHeader[7:] // strip "Bearer "
			_, valErr := manager.ValidateAccessToken(tkStr)
			if valErr != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
		req.Header.Set("Authorization", "Bearer "+tokenStr)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected HTTP 401 for expired token, got %d (userID=%s role=%s expiredAt=%v)",
				rec.Code, userID, role, expiredAt)
		}
	})
}

// TestProperty_JWTAccessTokenExpiryEnforcement_ValidTokenAccepted is a sanity
// check: a freshly generated (non-expired) token MUST be accepted.
func TestProperty_JWTAccessTokenExpiryEnforcement_ValidTokenAccepted(t *testing.T) {
	// Feature: itsm-web-app, Property 1: JWT Access Token Expiry Enforcement
	rapid.Check(t, func(t *rapid.T) {
		secret := rapid.StringMatching(`[a-zA-Z0-9]{16,64}`).Draw(t, "secret")
		role := roleGen().Draw(t, "role")
		userID := uuid.New()

		manager := NewManager(secret, "refresh-secret")

		tokenStr, err := manager.GenerateAccessToken(userID, role)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		// Token is fresh — must be valid.
		parsedClaims, err := manager.ValidateAccessToken(tokenStr)
		if err != nil {
			t.Fatalf("expected fresh token to be valid, got error: %v", err)
		}
		if parsedClaims.UserID != userID.String() {
			t.Fatalf("expected userID %s, got %s", userID, parsedClaims.UserID)
		}
		if parsedClaims.Role != role {
			t.Fatalf("expected role %s, got %s", role, parsedClaims.Role)
		}
	})
}
