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

type LoginUseCase struct {
	Repo   auth.Repository
	Cache  *cache.AuthCache
	Config *config.Config
}

func NewLoginUseCase(repo auth.Repository, cache *cache.AuthCache, config *config.Config) *LoginUseCase {
	return &LoginUseCase{
		Repo:   repo,
		Cache:  cache,
		Config: config,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpresIn     int    `json:"expires_in"`
}

func (s *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	email, err := helpers.NewEmail(req.Email)
	if err != nil || email.String() == "" {
		return nil, auth.ErrInvalidCredentials
	}
	storedUser, err := s.Repo.GetByEmail(ctx, email.String())
	if err != nil {
		metrics.RecordAuthAttempt(false)
		return nil, auth.ErrInvalidCredentials
	}
	if storedUser.IsActive == false {
		metrics.RecordAuthAttempt(false)
		return nil, auth.ErrUserInactive
	}
	if !helpers.CheckPasswordHash(req.Password, storedUser.Password.String()) {
		metrics.RecordAuthAttempt(false)
		return nil, auth.ErrInvalidCredentials
	}
	sid := helpers.GenerateKey()
	accessToken, err := helpers.GenerateToken(sid, s.Config.JwtAccessSecret, s.Config.JwtExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	refreshToken, err := helpers.GenerateToken(sid, s.Config.JwtRefreshSecret, s.Config.JwtRefreshExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	dataCache := &cache.UserCache{
		UserID:   storedUser.ID,
		Role:     storedUser.Role,
		Email:    storedUser.Email.String(),
		IsActive: storedUser.IsActive,
		FullName: storedUser.FullName,
		ImageUrl: storedUser.ImageUrl,
	}
	keyAuthCache := fmt.Sprintf("session_%s", sid)
	if err := s.Cache.SetAuthCache(ctx, keyAuthCache, dataCache, time.Duration(s.Config.JwtExpiresIn*int(time.Minute))); err != nil {
		return nil, domainerr.ErrInternalServerError
	}
	keyRefreshToken := fmt.Sprintf("refresh_token_%s", sid)
	if err := s.Cache.SetRefreshTokenCache(ctx, keyRefreshToken, refreshToken, time.Duration(s.Config.JwtRefreshExpiresIn*int(time.Second))); err != nil {
		return nil, domainerr.ErrInternalServerError
	}
	metrics.RecordAuthAttempt(true)
	metrics.ActiveSessions.Inc()
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpresIn:     s.Config.JwtExpiresIn,
	}, nil
}
