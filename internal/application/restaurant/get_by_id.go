package restaurantapp

import (
	"context"
	"fmt"
	"go-ai/internal/domain/restaurant"
	"go-ai/internal/infra/cache"
	"go-ai/internal/transport/http/status"
)

type GetByIDUseCase struct {
	repo  restaurant.Repository
	cache *cache.AuthCache
}

func NewGetByIDUseCase(repo restaurant.Repository, cache *cache.AuthCache) *GetByIDUseCase {
	return &GetByIDUseCase{
		repo:  repo,
		cache: cache,
	}
}

func (uc *GetByIDUseCase) Execute(ctx context.Context, id int32) (*GetRestaurantByIDResponse, error) {
	if id == 0 {
		return nil, status.ErrInvalidField
	}
	record, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("profile_%s", record.UserID.String())
	profile, err := uc.cache.GetAuthCache(key)
	if err != nil {
		return nil, err
	}
	return &GetRestaurantByIDResponse{
		RestaurantBase: RestaurantBase{
			Name:        record.Name,
			Description: *record.Description,
			Address:     *record.Address,
			Category:    *record.Category,
			City:        *record.City,
			District:    *record.District,
			LogoUrl:     *record.LogoUrl,
			BannerUrl:   *record.BannerUrl,
			PhoneNumber: *record.PhoneNumber,
			WebsiteUrl:  *record.WebsiteUrl,
			Email:       *record.Email,
		},
		UserName: profile.FullName,
	}, nil
}
