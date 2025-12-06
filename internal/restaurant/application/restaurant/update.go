package restaurantapp

import (
	"context"
	"go-ai/internal/restaurant/domain/restaurant"
	"go-ai/internal/transport/response"
	"strings"

	"github.com/google/uuid"
)

type UpdateRestaurantUseCase struct {
	Repo restaurant.Repository
}

func NewUpdateRestaurantUseCase(repo restaurant.Repository) *UpdateRestaurantUseCase {
	return &UpdateRestaurantUseCase{
		Repo: repo,
	}
}

func (uc *UpdateRestaurantUseCase) Execute(ctx context.Context, request CreateRestaurantRequest, userID uuid.UUID, id int32) error {
	if !strings.Contains(request.Email, "@") {
		return response.ErrInvalidEmail
	}
	if request.Name == "" {
		return response.ErrInvalidName
	}
	if request.Address == "" {
		return response.ErrInvalidAddress
	}
	if request.BannerUrl == "" {
		return response.ErrInvalidBanner
	}
	if request.LogoUrl == "" {
		return response.ErrInvalidLogo
	}
	if request.BannerUrl == "" {
		return response.ErrInvalidBanner
	}
	if request.PhoneNumber == "" {
		return response.ErrInvalidPhoneNumber
	}
	record, err := uc.Repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	if record == nil {
		return response.ErrRestaurantNoExitis
	}
	hours := make([]restaurant.Hours, 0, len(request.Hours))
	for _, hour := range request.Hours {
		hours = append(hours, restaurant.Hours{
			Day:       hour.Day,
			OpenTime:  hour.OpenTime,
			CloseTime: hour.CloseTime,
		})
	}
	err = uc.Repo.Update(ctx, &restaurant.Entity{
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
