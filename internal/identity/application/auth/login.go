package authapp

import (
	"context"
	"fmt"
	"go-ai/internal/identity/domain/auth"
	"go-ai/internal/identity/infrastructure/cache"
	"go-ai/internal/platform/config"
	"go-ai/internal/platform/security"
	"go-ai/pkg/utils"

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

func (s *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	email, err := utils.NewEmail(req.Email)
	storedUser, err := s.Repo.GetByEmail(ctx, email.String())
	if err != nil {
		return nil, auth.ErrInvalidCredentials
	}
	if storedUser.IsActive == false {
		return nil, auth.ErrUserInactive
	}
	if !security.CheckPasswordHash(req.Password, storedUser.Password.String()) {
		return nil, auth.ErrInvalidCredentials
	}
	accessToken, err := security.GenerateToken(storedUser.ID, storedUser.Email.String(), s.Config.JwtAccessSecret, s.Config.JwtExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	refreshToken, err := security.GenerateToken(storedUser.ID, storedUser.Email.String(), s.Config.JwtRefreshSecret, s.Config.JwtRefreshExpiresIn)
	if err != nil {
		return nil, auth.ErrTokenGenerateFail
	}
	dataCache := &cache.AuthData{
		UserID:   storedUser.ID,
		Role:     storedUser.Role,
		Email:    storedUser.Email.String(),
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
