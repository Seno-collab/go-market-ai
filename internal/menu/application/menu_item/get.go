package menuitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type GetUseCase struct {
	Repo domain.MenuItemRepository
}

func NewGetUseCase(repo domain.MenuItemRepository) *GetUseCase {
	return &GetUseCase{
		Repo: repo,
	}
}

func (useCase *GetUseCase) Execute(ctx context.Context, id int64, restaurantID int32) (*GetMenuItemResponse, error) {
	item, err := useCase.Repo.GetMenuItemByID(ctx, id, restaurantID)
	if err != nil {
		return nil, err
	}
	return &GetMenuItemResponse{
		Name:        item.Name,
		Sku:         item.Sku,
		Description: item.Description,
		Type:        item.Type,
		ImageUrl:    item.ImageUrl.String(),
		Price:       int64(item.BasePrice),
	}, nil
}
