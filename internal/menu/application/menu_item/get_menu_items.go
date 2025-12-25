package menuitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

type GetMenuItemsUseCase struct {
	Repo domain.MenuItemRepository
}

func NewGetMenuItemsUseCase(repo domain.MenuItemRepository) *GetMenuItemsUseCase {
	return &GetMenuItemsUseCase{
		Repo: repo,
	}
}

func (uc *GetMenuItemsUseCase) Execute(ctx context.Context, restaurantID int32) (*GetMenuItemsResponse, error) {
	items, err := uc.Repo.GetMenuItems(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	resp := make([]GetMenuItemResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, GetMenuItemResponse{
			Name:        item.Name,
			Sku:         item.Sku,
			Description: item.Description,
			Type:        item.Type,
			ImageUrl:    item.ImageUrl.String(),
			Price:       int64(item.BasePrice),
		})
	}
	return &GetMenuItemsResponse{
		Items: resp,
	}, nil
}
