package restaurantapp

import (
	"context"
	"errors"
	"go-ai/internal/restaurant/domain/restaurant"
	"go-ai/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CreateRestaurantUseCase struct {
	Repo restaurant.Repository
}

func NewCreateRestaurantUseCase(repo restaurant.Repository) *CreateRestaurantUseCase {
	return &CreateRestaurantUseCase{
		Repo: repo,
	}
}

func (uc *CreateRestaurantUseCase) Execute(ctx context.Context, req CreateRestaurantRequest, userID uuid.UUID) (int32, error) {
	email, err := utils.NewEmail(req.Email)
	if err != nil {
		return 0, err
	}
	phone, err := restaurant.NewPhone(req.PhoneNumber)
	if err != nil {
		return 0, err
	}

	logoUrl, err := utils.NewUrl(req.LogoUrl)
	if err != nil {
		return 0, err
	}

	websiteUrl, err := utils.NewUrl(req.WebsiteUrl)
	if err != nil {
		return 0, err
	}
	bannerUrl, err := utils.NewUrl(req.BannerUrl)
	if err != nil {
		return 0, err
	}
	hours := make([]restaurant.Hours, 0, len(req.Hours))
	for _, hour := range req.Hours {
		hours = append(hours, restaurant.Hours{
			Day:       hour.Day,
			OpenTime:  hour.OpenTime,
			CloseTime: hour.CloseTime,
		})
	}
	record, err := uc.Repo.GetByName(ctx, req.Name)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return 0, err
		}
	}
	if record != nil {
		return 0, restaurant.ErrRestaurantExists
	}

	entity, err := restaurant.NewEntity(
		req.Name,
		req.Description,
		req.Address,
		req.Category,
		req.City,
		req.District,
		logoUrl,
		bannerUrl, // banner optional example
		phone,
		websiteUrl,
		email,
		userID,
		hours,
	)
	if err != nil {
		return 0, err
	}
	id, err := uc.Repo.Create(ctx, entity, userID)
	if err != nil {
		return 0, err
	}
	return id, nil
}
