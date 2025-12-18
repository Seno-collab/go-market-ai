package restaurantapp

import (
	"context"
	"go-ai/internal/restaurant/domain/restaurant"
	domainerr "go-ai/pkg/domain_err"
)

type GetByIDUseCase struct {
	Repo restaurant.Repository
}

func NewGetByIDUseCase(repo restaurant.Repository) *GetByIDUseCase {
	return &GetByIDUseCase{
		Repo: repo,
	}
}

func (uc *GetByIDUseCase) Execute(ctx context.Context, id int32) (*GetRestaurantByIDResponse, error) {
	if id <= 0 {
		return nil, domainerr.ErrInvalidField
	}
	record, err := uc.Repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	// key := fmt.Sprintf("profile_%s", record.UserID.String())
	// profile, err := uc.Cache.GetAuthCache(key)
	// if err != nil {
	// 	return nil, err
	// }
	hours := make([]RestaurantHoursBase, 0, len(record.Hours))
	for _, hour := range record.Hours {
		hours = append(hours, RestaurantHoursBase{
			Day:       hour.Day,
			OpenTime:  hour.OpenTime,
			CloseTime: hour.CloseTime,
		})
	}
	resp := &GetRestaurantByIDResponse{
		RestaurantBase: RestaurantBase{
			Name:        record.Name,
			Description: record.Description,
			Address:     record.Address,
			Category:    record.Category,
			City:        record.City,
			District:    record.District,
			LogoUrl:     record.LogoUrl.String(),
			BannerUrl:   record.BannerUrl.String(),
			PhoneNumber: record.PhoneNumber.String(),
			WebsiteUrl:  record.WebsiteUrl.String(),
			Email:       record.Email.String(),
		},
		Hours:    hours,
		IsActive: len(hours) > 0,
	}

	return resp, nil
}
