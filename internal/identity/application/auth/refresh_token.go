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
	keyRefreshToken := fmt.Sprintf("refresh_token_%s", claims.Sid)
	cachedRefreshToken, err := uc.Cache.GetRefreshTokenCache(ctx, keyRefreshToken)
	if err != nil {
		return nil, err
	}
	if cachedRefreshToken != request.RefreshToken {
		return nil, auth.ErrTokenMalformed
	}
	keyAuthCache := fmt.Sprintf("profile_%s", claims.Sid)
	authData, err := uc.Cache.GetAuthCache(ctx, keyAuthCache)
	record, err := uc.Repo.GetByEmail(ctx, authData.Email)
	if err != nil {
		return nil, err
	}
	if !record.IsActive {
		return nil, auth.ErrUserInactive
	}
	sid := security.GenerateKey()
	accessToken, err := security.GenerateToken(sid, uc.Config.JwtAccessSecret, uc.Config.JwtExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	refreshToken, err := security.GenerateToken(sid, uc.Config.JwtRefreshSecret, uc.Config.JwtRefreshExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	dataCache := &cache.AuthData{
		UserID:   record.ID,
		Role:     record.Role,
		Email:    record.Email.String(),
		IsActive: record.IsActive,
		FullName: record.FullName,
		ImageUrl: record.ImageUrl,
	}
	if err := uc.Cache.SetAuthCache(ctx, keyAuthCache, dataCache, time.Duration(uc.Config.JwtExpiresIn*int(time.Second))); err != nil {
		return nil, domainerr.ErrInternalServerError
	}
	if err := uc.Cache.SetRefreshTokenCache(ctx, keyRefreshToken, refreshToken, time.Duration(uc.Config.JwtRefreshExpiresIn*int(time.Second))); err != nil {
		return nil, domainerr.ErrInternalServerError
	}
	return &RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    uc.Config.JwtExpiresIn,
	}, nil

}
