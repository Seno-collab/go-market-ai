package restaurantapp

import (
	"context"
	"go-ai/internal/restaurant/domain/restaurant"
	"go-ai/internal/transport/response"

	"github.com/google/uuid"
)

type GetRestaurantItemComboboxUseCase struct {
	Repo restaurant.Repository
}

func NewGetRestaurantItemComboboxUseCase(repo restaurant.Repository) *GetRestaurantItemComboboxUseCase {
	return &GetRestaurantItemComboboxUseCase{
		Repo: repo,
	}
}

func (useCase *GetRestaurantItemComboboxUseCase) Execute(ctx context.Context, userID uuid.UUID) (*[]response.Combobox, error) {
	return useCase.Repo.GetRestaurantItemCombobox(ctx, userID)
}
