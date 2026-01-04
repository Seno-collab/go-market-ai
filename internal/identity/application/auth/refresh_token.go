package authapp

import (
	"context"
	"fmt"
	"go-ai/internal/identity/domain/auth"
	"go-ai/internal/identity/infrastructure/cache"
	"go-ai/internal/platform/config"
	"go-ai/internal/platform/security"
	domainerr "go-ai/pkg/domain_err"
	"time"

	"github.com/google/uuid"
)

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
	claims, err := security.VerifyToken(request.RefreshToken, uc.Config.JwtRefreshSecret)
	if err != nil {
		return nil, auth.ErrTokenNotActive
	}
	userID := claims.UserID
	if userID == uuid.Nil {
		return nil, auth.ErrTokenMissing
	}
	email := claims.Email
	if email == "" {
		return nil, auth.ErrTokenMissing
	}
	keyRefreshToken := fmt.Sprintf("refresh_token_%s", userID)
	cachedRefreshToken, err := uc.Cache.GetRefreshTokenCache(keyRefreshToken)
	if err != nil {
		return nil, err
	}
	if cachedRefreshToken != request.RefreshToken {
		return nil, auth.ErrTokenMalformed
	}
	record, err := uc.Repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if !record.IsActive {
		return nil, auth.ErrUserInactive
	}
	accessToken, err := security.GenerateToken(userID, email, uc.Config.JwtAccessSecret, uc.Config.JwtExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	refreshToken, err := security.GenerateToken(userID, email, uc.Config.JwtRefreshSecret, uc.Config.JwtRefreshExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	dataCache := &cache.AuthData{
		UserID:   record.ID,
		Role:     record.Role,
		Email:    record.Email.String(),
		IsActive: record.IsActive,
		FullName: record.FullName,
	}
	keyAuthCache := fmt.Sprintf("profile_%s", record.ID.String())
	if err := uc.Cache.SetAuthCache(keyAuthCache, dataCache, time.Duration(uc.Config.JwtExpiresIn*int(time.Second))); err != nil {
		return nil, domainerr.ErrInternalServerError
	}
	if err := uc.Cache.SetRefreshTokenCache(keyRefreshToken, refreshToken, time.Duration(uc.Config.JwtRefreshExpiresIn*int(time.Second))); err != nil {
		return nil, domainerr.ErrInternalServerError
	}
	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    uc.Config.JwtExpiresIn,
	}, nil

}
