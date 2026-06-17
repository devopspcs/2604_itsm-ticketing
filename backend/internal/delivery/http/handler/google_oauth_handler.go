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

const (
	googleAuthURL  = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL = "https://oauth2.googleapis.com/token"
	googleUserURL  = "https://www.googleapis.com/oauth2/v2/userinfo"
	allowedDomain  = "pcsindonesia.co.id"
)

type GoogleOAuthHandler struct {
	cfg        *config.Config
	userRepo   repository.UserRepository
	jwtManager *jwtpkg.Manager
}

func NewGoogleOAuthHandler(cfg *config.Config, userRepo repository.UserRepository, jwtManager *jwtpkg.Manager) *GoogleOAuthHandler {
	return &GoogleOAuthHandler{cfg: cfg, userRepo: userRepo, jwtManager: jwtManager}
}

// GetLoginURL returns the Google OAuth authorization URL
func (h *GoogleOAuthHandler) GetLoginURL(w http.ResponseWriter, r *http.Request) {
	redirectURI := h.cfg.BaseURL + "/api/v1/auth/google/callback"
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email+profile&hd=%s&access_type=offline&prompt=consent",
		googleAuthURL,
		url.QueryEscape(h.cfg.GoogleClientID),
		url.QueryEscape(redirectURI),
		allowedDomain,
	)

	apperror.WriteJSON(w, http.StatusOK, map[string]string{
		"login_url": authURL,
	})
}

// Redirect does a server-side 302 redirect to Google login
func (h *GoogleOAuthHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	redirectURI := h.cfg.BaseURL + "/api/v1/auth/google/callback"
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email+profile&hd=%s&access_type=offline&prompt=consent",
		googleAuthURL,
		url.QueryEscape(h.cfg.GoogleClientID),
		url.QueryEscape(redirectURI),
		allowedDomain,
	)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// Callback handles the Google OAuth callback
func (h *GoogleOAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Redirect(w, r, h.cfg.BaseURL+"/login?error=no_code", http.StatusFound)
		return
	}

	// Exchange code for tokens
	tokenResp, err := h.exchangeCode(code)
	if err != nil {
		http.Redirect(w, r, h.cfg.BaseURL+"/login?error=token_exchange_failed", http.StatusFound)
		return
	}

	// Get user info from Google
	userInfo, err := h.getUserInfo(tokenResp.AccessToken)
	if err != nil {
		http.Redirect(w, r, h.cfg.BaseURL+"/login?error=userinfo_failed", http.StatusFound)
		return
	}

	// Validate domain
	if !strings.HasSuffix(userInfo.Email, "@"+allowedDomain) {
		http.Redirect(w, r, h.cfg.BaseURL+"/login?error=domain_not_allowed", http.StatusFound)
		return
	}

	// Find or create user
	user, err := h.findOrCreateUser(r.Context(), userInfo)
	if err != nil {
		http.Redirect(w, r, h.cfg.BaseURL+"/login?error=user_creation_failed", http.StatusFound)
		return
	}

	// Check if user is active
	if !user.IsActive {
		http.Redirect(w, r, h.cfg.BaseURL+"/login?error=account_disabled", http.StatusFound)
		return
	}

	// Generate JWT tokens
	accessToken, err := h.jwtManager.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		http.Redirect(w, r, h.cfg.BaseURL+"/login?error=token_generation_failed", http.StatusFound)
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Redirect(w, r, h.cfg.BaseURL+"/login?error=token_generation_failed", http.StatusFound)
		return
	}

	// Redirect to frontend with tokens
	redirectURL := fmt.Sprintf("%s/sso/callback?access_token=%s&refresh_token=%s&expires_in=%d",
		h.cfg.BaseURL, url.QueryEscape(accessToken), url.QueryEscape(refreshToken), int64(jwtpkg.AccessTokenTTL.Seconds()))

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

type googleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func (h *GoogleOAuthHandler) exchangeCode(code string) (*googleTokenResponse, error) {
	redirectURI := h.cfg.BaseURL + "/api/v1/auth/google/callback"

	data := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {h.cfg.GoogleClientID},
		"client_secret": {h.cfg.GoogleClientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURI},
	}

	resp, err := http.Post(googleTokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("google token exchange failed: %s", string(body))
	}

	var tokenResp googleTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

func (h *GoogleOAuthHandler) getUserInfo(accessToken string) (*googleUserInfo, error) {
	req, _ := http.NewRequest("GET", googleUserURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("google userinfo failed: %s", string(body))
	}

	var info googleUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (h *GoogleOAuthHandler) findOrCreateUser(ctx context.Context, info *googleUserInfo) (*entity.User, error) {
	// Try to find existing user
	user, err := h.userRepo.FindByEmail(ctx, info.Email)
	if err == nil && user != nil {
		// Update avatar if changed
		if info.Picture != "" && user.AvatarURL != info.Picture {
			user.AvatarURL = info.Picture
			h.userRepo.Update(ctx, user)
		}
		return user, nil
	}

	// Create new user
	fullName := info.Name
	if fullName == "" {
		fullName = info.GivenName + " " + info.FamilyName
	}

	randomPass, _ := password.Hash(uuid.New().String())

	now := time.Now().UTC()
	newUser := &entity.User{
		ID:           uuid.New(),
		FullName:     fullName,
		Email:        info.Email,
		PasswordHash: randomPass,
		Role:         entity.RoleUser,
		IsActive:     true,
		AvatarURL:    info.Picture,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
