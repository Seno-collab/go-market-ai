package authapp

import (
	"context"
	"fmt"
	"go-ai/internal/identity/domain/auth"
	"go-ai/internal/identity/infrastructure/cache"

	"github.com/google/uuid"
)

type LogoutUseCase struct {
	Repo  auth.Repository
	Cache *cache.AuthCache
}

func NewLogoutUseCase(repo auth.Repository, cache *cache.AuthCache) *LogoutUseCase {
	return &LogoutUseCase{
		Repo:  repo,
		Cache: cache,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, userID uuid.UUID) error {
	profileKey := fmt.Sprintf("profile_%s", userID)
	refreshKey := fmt.Sprintf("refresh_token_%s", userID)
	err := uc.Cache.DeleteAuthCache(profileKey)
	errRefresh := uc.Cache.DeleteRefreshTokenCache(refreshKey)
	if err != nil {
		return err
	}
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
