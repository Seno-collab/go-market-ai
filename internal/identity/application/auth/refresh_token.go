package authapp

import (
	"context"
	"fmt"
	"go-ai/internal/identity/domain/auth"
	"go-ai/internal/identity/infra/cache"
	"go-ai/internal/platform/config"
	"go-ai/internal/platform/security"
	"go-ai/internal/transport/response"
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
		return nil, response.ErrTokenInvalid
	}
	claims, err := security.VerifyToken(request.RefreshToken, uc.Config.JwtRefreshSecret)
	if err != nil {
		return nil, response.ErrTokenNotActive
	}
	userID := claims.UserID
	if userID == uuid.Nil {
		return nil, response.ErrTokenMissing
	}
	email := claims.Email
	if email == "" {
		return nil, response.ErrTokenMissing
	}
	keyRefreshToken := fmt.Sprintf("refresh_token_%s", userID)
	cachedRefreshToken, err := uc.Cache.GetRefreshTokenCache(keyRefreshToken)
	if err != nil {
		return nil, err
	}
	if cachedRefreshToken != request.RefreshToken {
		return nil, response.ErrTokenMalformed
	}
	accessToken, err := security.GenerateToken(userID, email, uc.Config.JwtAccessSecret, uc.Config.JwtExpiresIn)
	if err != nil {
		return nil, response.ErrTokenGenerateFail
	}
	refreshToken, err := security.GenerateToken(userID, email, uc.Config.JwtRefreshSecret, uc.Config.JwtRefreshExpiresIn)
	if err != nil {
		return nil, response.ErrTokenGenerateFail
	}
	record, err := uc.Repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	dataCache := &cache.AuthData{
		UserID:   record.ID,
		Role:     record.Role,
		Email:    record.Email,
		IsActive: record.IsActive,
		FullName: record.FullName,
	}
	keyAuthCache := fmt.Sprintf("profile_%s", record.ID.String())
	uc.Cache.SetAuthCache(keyAuthCache, dataCache, time.Duration(uc.Config.JwtExpiresIn*int(time.Second)))
	uc.Cache.SetRefreshTokenCache(keyRefreshToken, refreshToken, time.Duration(uc.Config.JwtRefreshExpiresIn*int(time.Second)))
	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    uc.Config.JwtExpiresIn,
	}, nil

}
