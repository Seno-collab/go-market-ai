package authapp

import (
	"context"
	"fmt"
	"go-ai/internal/identity/domain/auth"
	"go-ai/internal/identity/infrastructure/cache"
	"go-ai/pkg/metrics"

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
	sessionKey := fmt.Sprintf("session_%s", userID)
	refreshKey := fmt.Sprintf("refresh_token_%s", userID)
	err := uc.Cache.DeleteAuthCache(ctx, sessionKey)
	errRefresh := uc.Cache.DeleteRefreshTokenCache(ctx, refreshKey)
	if err != nil {
		return err
	}
	if errRefresh != nil {
		return errRefresh
	}
	metrics.ActiveSessions.Dec()
	return nil
}
