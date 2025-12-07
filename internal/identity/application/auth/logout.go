package authapp

import (
	"context"
	"fmt"
	"go-ai/internal/identity/domain/auth"
	"go-ai/internal/identity/infra/cache"

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
	key := fmt.Sprintf("profile_%s", userID)
	return uc.Cache.DeleteAuthCache(key)
}
