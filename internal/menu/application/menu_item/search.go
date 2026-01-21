package menuitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/internal/transport/response"
	"math"
)

type GetMenuItemsUseCase struct {
	Repo domain.MenuItemRepository
}

func NewGetMenuItemsUseCase(repo domain.MenuItemRepository) *GetMenuItemsUseCase {
	return &GetMenuItemsUseCase{
		Repo: repo,
	}
}

func (uc *GetMenuItemsUseCase) Execute(ctx context.Context, restaurantID int32, req GetMenuItemsRequest) (*GetMenuItemsResponse, error) {
	page, limit, offset := response.ApplyDefaultPaginated(req.Page, req.Limit)

	filter := ""
	if req.Filter != nil {
		filter = *req.Filter
	}

	category := ""
	if req.Category != nil {
		category = *req.Category
	}
	param := domain.SearchMenuItemsParam{
		Limit:        limit,
		Offset:       offset,
		Filter:       filter,
		Category:     category,
		IsActive:     req.IsActive,
		RestaurantID: restaurantID,
	}
	items, total, err := uc.Repo.GetMenuItems(ctx, param)
	if err != nil {
		return nil, err
	}
	resp := make([]GetMenuItemResponse, 0, len(items))
	for _, item := range items {
		resp = append(resp, GetMenuItemResponse{
			ID:          item.ID,
			Name:        item.Name,
			Sku:         item.Sku,
			Description: item.Description,
			Type:        item.Type,
			ImageUrl:    item.ImageUrl.String(),
			Price:       int64(item.BasePrice),
			IsActive:    item.IsActive,
			UpdatedAt:   item.UpdatedAt,
			CreateAt:    item.CreatedAt,
		})
	}
	totalPages := int64(math.Ceil(float64(total) / float64(limit)))

	return &GetMenuItemsResponse{
		PaginatedResponse: response.PaginatedResponse[[]GetMenuItemResponse]{
			Items:      resp,
			Page:       page,
			Limit:      limit,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil

}
