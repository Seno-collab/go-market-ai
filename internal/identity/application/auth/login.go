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

func (s *LoginUseCase) Execute(ctx context.Context, request LoginRequest) (*LoginResponse, error) {

	if request.Email == "" {
		return nil, response.ErrInvalidEmail
	}

	if request.Password == "" {
		return nil, response.ErrInvalidPassword
	}

	storedUser, err := s.Repo.GetByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if storedUser.IsActive == false {
		return nil, response.ErrUserInactive
	}
	if !security.CheckPasswordHash(request.Password, storedUser.Password) {
		return nil, response.ErrPasswordVerifyFail
	}
	accessToken, err := security.GenerateToken(storedUser.ID, storedUser.Email, s.Config.JwtAccessSecret, s.Config.JwtExpiresIn)
	if err != nil {
		return nil, response.ErrTokenGenerateFail
	}
	refreshToken, err := security.GenerateToken(storedUser.ID, storedUser.Email, s.Config.JwtRefreshSecret, s.Config.JwtRefreshExpiresIn)
	if err != nil {
		return nil, response.ErrTokenGenerateFail
	}
	dataCache := &cache.AuthData{
		UserID:   storedUser.ID,
		Role:     storedUser.Role,
		Email:    storedUser.Email,
		IsActive: storedUser.IsActive,
		FullName: storedUser.FullName,
	}
	keyAuthCache := fmt.Sprintf("profile_%s", storedUser.ID.String())
	s.Cache.SetAuthCache(keyAuthCache, dataCache, time.Duration(s.Config.JwtExpiresIn*int(time.Second)))
	keyRefreshToken := fmt.Sprintf("refresh_token_%s", storedUser.ID.String())
	s.Cache.SetRefreshTokenCache(keyRefreshToken, refreshToken, time.Duration(s.Config.JwtRefreshExpiresIn*int(time.Second)))
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpresIn:     s.Config.JwtExpiresIn,
	}, nil
}
