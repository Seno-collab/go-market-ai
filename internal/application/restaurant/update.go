package restaurantapp

import (
	"context"
	"go-ai/internal/domain/restaurant"
	"go-ai/internal/transport/http/status"
	"strings"

	"github.com/google/uuid"
)

type UpdateRestaurantUseCase struct {
	repo restaurant.Repository
}

func NewUpdateRestaurantUseCase(repo restaurant.Repository) *UpdateRestaurantUseCase {
	return &UpdateRestaurantUseCase{
		repo: repo,
	}
}

func (uc *UpdateRestaurantUseCase) Execute(ctx context.Context, request CreateRestaurantRequest, userID uuid.UUID, id int32) error {
	if !strings.Contains(request.Email, "@") {
		return status.ErrInvalidEmail
	}
	if request.Name == "" {
		return status.ErrInvalidName
	}
	if request.Address == "" {
		return restaurant.ErrInvalidAddress
	}
	if request.BannerUrl == "" {
		return restaurant.ErrInvalidBanner
	}
	if request.LogoUrl == "" {
		return restaurant.ErrInvalidLogo
	}
	if request.BannerUrl == "" {
		return restaurant.ErrInvalidBanner
	}
	if request.PhoneNumber == "" {
		return restaurant.ErrInvalidPhoneNumber
	}
	record, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	if record == nil {
		return restaurant.ErrRestaurantNoExitis
	}
	hours := make([]restaurant.Hours, 0, len(request.Hours))
	for _, hour := range request.Hours {
		hours = append(hours, restaurant.Hours{
			Day:       hour.Day,
			OpenTime:  hour.OpenTime,
			CloseTime: hour.CloseTime,
		})
	}
	err = uc.repo.Update(ctx, &restaurant.Entity{
		Email:       request.Email,
		Name:        request.Name,
		Description: request.Description,
		Category:    request.Category,
		WebsiteUrl:  request.WebsiteUrl,
		LogoUrl:     request.LogoUrl,
		BannerUrl:   request.BannerUrl,
		PhoneNumber: request.PhoneNumber,
		Address:     request.PhoneNumber,
		City:        request.City,
		District:    request.District,
		UserID:      userID,
		Hours:       hours,
	}, id)
	if err != nil {
		return err
	}
	return nil
}
