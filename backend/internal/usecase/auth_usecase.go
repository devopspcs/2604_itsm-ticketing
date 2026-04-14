package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	"github.com/org/itsm/internal/domain/usecase"
	"github.com/org/itsm/pkg/apperror"
	"github.com/org/itsm/pkg/jwt"
	"github.com/org/itsm/pkg/password"
)

type authUseCase struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtManager       *jwt.Manager
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtManager *jwt.Manager,
) usecase.AuthUseCase {
	return &authUseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
	}
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", sum)
}

// Login validates credentials, checks is_active, generates a TokenPair, and saves the refresh token hash.
func (a *authUseCase) Login(ctx context.Context, req usecase.LoginRequest) (*usecase.TokenPair, error) {
	user, err := a.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		// Return generic error to avoid revealing whether email exists
		return nil, apperror.ErrInvalidCredentials
	}

	if !password.Verify(req.Password, user.PasswordHash) {
		return nil, apperror.ErrInvalidCredentials
	}

	// Check is_active — return same generic error to avoid leaking info (Req 1.2, 3.4)
	if !user.IsActive {
		return nil, apperror.ErrInvalidCredentials
	}

	accessToken, err := a.jwtManager.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	refreshTokenStr, err := a.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	tokenHash := hashToken(refreshTokenStr)

	rt := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(jwt.RefreshTokenTTL),
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	if err := a.refreshTokenRepo.Save(ctx, rt); err != nil {
		return nil, apperror.ErrInternal
	}

	return &usecase.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    int64(jwt.AccessTokenTTL.Seconds()),
	}, nil
}

// RefreshToken validates the refresh token hash, checks revoked/expiry, and returns a new access token.
func (a *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (*usecase.TokenPair, error) {
	// Validate JWT signature first
	userID, err := a.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, apperror.ErrTokenInvalid
	}

	tokenHash := hashToken(refreshToken)

	stored, err := a.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, apperror.ErrTokenInvalid
	}

	if stored.Revoked {
		return nil, apperror.ErrTokenInvalid
	}

	if time.Now().After(stored.ExpiresAt) {
		return nil, apperror.ErrTokenExpired
	}

	user, err := a.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.ErrTokenInvalid
	}

	if !user.IsActive {
		return nil, apperror.ErrInvalidCredentials
	}

	newAccessToken, err := a.jwtManager.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	return &usecase.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(jwt.AccessTokenTTL.Seconds()),
	}, nil
}

// Logout revokes the refresh token.
func (a *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)

	stored, err := a.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		// Token not found — treat as already logged out (idempotent)
		return nil
	}

	return a.refreshTokenRepo.Revoke(ctx, stored.ID)
}
