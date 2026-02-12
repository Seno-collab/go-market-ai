package authapp

import (
	"context"
	"fmt"
	"go-ai/internal/identity/domain/auth"
	"go-ai/internal/identity/infrastructure/cache"
	"go-ai/internal/platform/config"
	domainerr "go-ai/pkg/domain_err"
	"go-ai/pkg/helpers"
	"go-ai/pkg/metrics"
	"time"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshTokenUseCase struct {
	Repo   auth.Repository
	Cache  *cache.AuthCache
	Config *config.Config
}

func NewRefreshTokenUseCase(repo auth.Repository, cache *cache.AuthCache, config *config.Config) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		Repo:   repo,
		Cache:  cache,
		Config: config,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, request RefreshTokenRequest) (*RefreshTokenResponse, error) {
	if request.RefreshToken == "" {
		return nil, auth.ErrTokenInvalid
	}
	claims, err := helpers.VerifyToken(request.RefreshToken, uc.Config.JwtRefreshSecret)
	if err != nil {
		return nil, auth.ErrTokenNotActive
	}

	oldSid := claims.Sid
	blacklistKey := fmt.Sprintf("blacklist_%s", oldSid)

	// Check if the refresh token has already been used (blacklisted)
	blacklisted, err := uc.Cache.IsTokenBlacklisted(ctx, blacklistKey)
	if err != nil {
		return nil, domainerr.ErrInternalServerError
	}
	if blacklisted {
		return nil, auth.ErrTokenInvalid
	}

	oldRefreshKey := fmt.Sprintf("refresh_token_%s", oldSid)
	cachedRefreshToken, err := uc.Cache.GetRefreshTokenCache(ctx, oldRefreshKey)
	if err != nil {
		return nil, err
	}
	if cachedRefreshToken != request.RefreshToken {
		return nil, auth.ErrTokenMalformed
	}

	oldSessionKey := fmt.Sprintf("session_%s", oldSid)
	authData, err := uc.Cache.GetAuthCache(ctx, oldSessionKey)
	if err != nil {
		return nil, err
	}
	record, err := uc.Repo.GetByEmail(ctx, authData.Email)
	if err != nil {
		return nil, err
	}
	if !record.IsActive {
		return nil, auth.ErrUserInactive
	}

	// Generate new session
	newSid := helpers.GenerateKey()
	accessToken, err := helpers.GenerateToken(newSid, uc.Config.JwtAccessSecret, uc.Config.JwtExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	refreshToken, err := helpers.GenerateToken(newSid, uc.Config.JwtRefreshSecret, uc.Config.JwtRefreshExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}

	// Blacklist the old refresh token
	blacklistTTL := time.Duration(uc.Config.JwtRefreshExpiresIn) * time.Second
	if err := uc.Cache.BlacklistToken(ctx, blacklistKey, blacklistTTL); err != nil {
		return nil, domainerr.ErrInternalServerError
	}

	// Delete old session and refresh token caches
	_ = uc.Cache.DeleteAuthCache(ctx, oldSessionKey)
	_ = uc.Cache.DeleteRefreshTokenCache(ctx, oldRefreshKey)

	// Store new session and refresh token with new sid
	newSessionKey := fmt.Sprintf("session_%s", newSid)
	newRefreshKey := fmt.Sprintf("refresh_token_%s", newSid)
	dataCache := &cache.UserCache{
		UserID:   record.ID,
		Role:     record.Role,
		Email:    record.Email.String(),
		IsActive: record.IsActive,
		FullName: record.FullName,
		ImageUrl: record.ImageUrl,
	}
	if err := uc.Cache.SetLoginCaches(
		ctx,
		newSessionKey, dataCache, time.Duration(uc.Config.JwtExpiresIn*int(time.Second)),
		newRefreshKey, refreshToken, time.Duration(uc.Config.JwtRefreshExpiresIn*int(time.Second)),
	); err != nil {
		return nil, domainerr.ErrInternalServerError
	}

	metrics.AuthTokenRefreshes.Inc()
	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    uc.Config.JwtExpiresIn,
	}, nil
}
