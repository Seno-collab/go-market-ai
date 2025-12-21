package optionitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
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

func (uc *GetByGroupUseCase) Execute(ctx context.Context, groupID int64, restaurantID int32) (*GetOptionItemsResponse, error) {
	if _, err := uc.GroupRepo.GetOptionGroup(ctx, groupID, restaurantID); err != nil {
		return nil, err
	}
	items, err := uc.Repo.GetOptionItems(ctx, groupID, restaurantID)
	if err != nil {
		return nil, err
	}
	resp := GetOptionItemsResponse{
		Items: make([]GetOptionItemResponse, 0, len(items)),
	}
	for _, item := range items {
		resp.Items = append(resp.Items, mapOptionItem(item))
	}
	return &resp, nil
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
