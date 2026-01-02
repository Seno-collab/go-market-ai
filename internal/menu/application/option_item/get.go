package optionitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/internal/transport/response"
)

type GetUseCase struct {
	Repo domain.OptionItemRepository
}

func NewGetUseCase(repo domain.OptionItemRepository) *GetUseCase {
	return &GetUseCase{Repo: repo}
}

func (uc *GetUseCase) Execute(ctx context.Context, id int64, restaurantID int32) (*GetOptionItemResponse, error) {
	item, err := uc.Repo.GetOptionItem(ctx, id, restaurantID)
	if err != nil {
		return nil, err
	}
	resp := mapOptionItem(item)
	return &resp, nil
}

type GetByGroupUseCase struct {
	Repo      domain.OptionItemRepository
	GroupRepo domain.OptionRepository
}

func NewGetByGroupUseCase(repo domain.OptionItemRepository, groupRepo domain.OptionRepository) *GetByGroupUseCase {
	return &GetByGroupUseCase{
		Repo:      repo,
		GroupRepo: groupRepo,
	}
}

func (uc *GetByGroupUseCase) Execute(ctx context.Context, groupID int64, restaurantID int32, req GetOptionItemsRequest) (*GetOptionItemsResponse, error) {
	if _, err := uc.GroupRepo.GetOptionGroup(ctx, groupID, restaurantID); err != nil {
		return nil, err
	}
	page, limit, offset := response.ApplyDefaultPaginated(req.Page, req.Limit)
	items, total, err := uc.Repo.GetOptionItems(ctx, groupID, restaurantID, limit, offset)
	if err != nil {
		return nil, err
	}
	result := make([]GetOptionItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, mapOptionItem(item))
	}

	return &GetOptionItemsResponse{
		PaginatedResponse: response.PaginatedResponse[[]GetOptionItemResponse]{
			Page:       page,
			Limit:      limit,
			TotalPages: total,
			TotalItems: response.CalculateTotalPages(total, int64(limit)),
			Items:      result,
		},
	}, nil
}

func mapOptionItem(item domain.OptionItem) GetOptionItemResponse {
	return GetOptionItemResponse{
		ID:             item.ID,
		OptionGroupID:  item.OptionGroupID,
		Name:           item.Name,
		LinkedMenuItem: item.LinkedMenuItem,
		PriceDelta:     int64(item.PriceDelta),
		QuantityMin:    item.QuantityMin,
		QuantityMax:    item.QuantityMax,
		SortOrder:      item.SortOrder,
		IsActive:       item.IsActive,
	}
}
