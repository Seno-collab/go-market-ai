package authapp

import (
	"context"
	"fmt"
	"go-ai/internal/identity/domain/auth"
	"go-ai/internal/identity/infra/cache"

	"time"

	"github.com/google/uuid"
)

type GetProfileUseCase struct {
	Repo  auth.Repository
	Cache *cache.AuthCache
}

func NewGetProfileUseCase(repo auth.Repository, cache *cache.AuthCache) *GetProfileUseCase {
	return &GetProfileUseCase{
		Repo:  repo,
		Cache: cache,
	}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, userID uuid.UUID) (*GetProfileResponse, error) {
	keyAuth := fmt.Sprintf("profile_%s", userID.String())
	cacheData, err := uc.Cache.GetAuthCache(keyAuth)
	if err != nil {
		return nil, err
	}
	profile := &GetProfileResponse{}
	if cacheData == nil {
		record, err := uc.Repo.GetById(ctx, userID)
		if err != nil {
			return nil, err
		}
		profile = &GetProfileResponse{
			Email:    record.Email,
			FullName: record.FullName,
			Role:     record.Role,
			IsActive: record.IsActive,
		}
		authData := &cache.AuthData{
			Email:    record.Email,
			FullName: record.FullName,
			Role:     record.Role,
			IsActive: record.IsActive,
		}
		uc.Cache.SetAuthCache(keyAuth, authData, time.Duration(60*int(time.Minute)))
		return profile, nil
	}
	profile = &GetProfileResponse{
		Email:    cacheData.Email,
		FullName: cacheData.FullName,
		Role:     cacheData.Role,
		IsActive: cacheData.IsActive,
	}
	return profile, nil
}
