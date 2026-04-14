package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	"github.com/org/itsm/pkg/apperror"
	"github.com/org/itsm/pkg/config"
	jwtpkg "github.com/org/itsm/pkg/jwt"
	"github.com/org/itsm/pkg/password"
)

type SSOHandler struct {
	cfg        *config.Config
	userRepo   repository.UserRepository
	jwtManager *jwtpkg.Manager
}

func NewSSOHandler(cfg *config.Config, userRepo repository.UserRepository, jwtManager *jwtpkg.Manager) *SSOHandler {
	return &SSOHandler{cfg: cfg, userRepo: userRepo, jwtManager: jwtManager}
}

// GetLoginURL returns the Keycloak authorization URL for frontend redirect
func (h *SSOHandler) GetLoginURL(w http.ResponseWriter, r *http.Request) {
	redirectURI := h.cfg.BaseURL + "/sso/callback"
	authURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email+profile",
		h.cfg.KeycloakURL, h.cfg.KeycloakRealm,
		url.QueryEscape(h.cfg.KeycloakClientID),
		url.QueryEscape(redirectURI))

	apperror.WriteJSON(w, http.StatusOK, map[string]string{
		"login_url": authURL,
	})
}

// Callback handles the OAuth2 callback from Keycloak
func (h *SSOHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		apperror.WriteError(w, apperror.ErrValidation.WithDetails(map[string]interface{}{"error": "missing authorization code"}))
		return
	}

	// Exchange code for token
	tokenResp, err := h.exchangeCode(code)
	if err != nil {
		apperror.WriteError(w, apperror.ErrInternal.WithDetails(map[string]interface{}{"error": "failed to exchange code: " + err.Error()}))
		return
	}

	// Get user info from Keycloak
	userInfo, err := h.getUserInfo(tokenResp.AccessToken)
	if err != nil {
		apperror.WriteError(w, apperror.ErrInternal.WithDetails(map[string]interface{}{"error": "failed to get user info"}))
		return
	}

	// Find or create user in our database
	user, err := h.findOrCreateUser(r.Context(), userInfo)
	if err != nil {
		apperror.WriteError(w, apperror.ErrInternal.WithDetails(map[string]interface{}{"error": "failed to create user"}))
		return
	}

	// Generate our own JWT tokens
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		apperror.WriteError(w, apperror.ErrInternal)
		return
	}
	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		apperror.WriteError(w, apperror.ErrInternal)
		return
	}

	// Redirect to frontend with tokens in URL fragment
	redirectURL := fmt.Sprintf("%s/sso/callback?access_token=%s&refresh_token=%s&expires_in=%d",
		h.cfg.BaseURL, url.QueryEscape(accessToken), url.QueryEscape(refreshToken), int64(jwtpkg.AccessTokenTTL.Seconds()))

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

type keycloakTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type keycloakUserInfo struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	PreferredUsername  string `json:"preferred_username"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	EmailVerified     bool   `json:"email_verified"`
}

func (h *SSOHandler) exchangeCode(code string) (*keycloakTokenResponse, error) {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", h.cfg.KeycloakURL, h.cfg.KeycloakRealm)
	redirectURI := h.cfg.BaseURL + "/sso/callback"

	data := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {h.cfg.KeycloakClientID},
		"client_secret": {h.cfg.KeycloakSecret},
		"code":          {code},
		"redirect_uri":  {redirectURI},
	}

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp keycloakTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (h *SSOHandler) getUserInfo(accessToken string) (*keycloakUserInfo, error) {
	userInfoURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", h.cfg.KeycloakURL, h.cfg.KeycloakRealm)

	req, _ := http.NewRequest("GET", userInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("userinfo failed: %s", string(body))
	}

	var info keycloakUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (h *SSOHandler) findOrCreateUser(ctx context.Context, info *keycloakUserInfo) (*entity.User, error) {
	email := info.Email
	if email == "" {
		email = info.PreferredUsername + "@sso.local"
	}

	// Try to find existing user
	user, err := h.userRepo.FindByEmail(ctx, email)
	if err == nil && user != nil {
		return user, nil
	}

	// Create new user from Keycloak info
	fullName := info.Name
	if fullName == "" {
		fullName = info.GivenName + " " + info.FamilyName
	}
	if strings.TrimSpace(fullName) == "" {
		fullName = info.PreferredUsername
	}

	// Generate random password (user won't use it, they login via SSO)
	randomPass, _ := password.Hash(uuid.New().String())

	now := time.Now().UTC()
	newUser := &entity.User{
		ID:           uuid.New(),
		FullName:     fullName,
		Email:        email,
		PasswordHash: randomPass,
		Role:         entity.RoleUser, // default role for SSO users
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
