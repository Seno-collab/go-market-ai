package menuitemapp

import (
	"context"
	"go-ai/internal/menu/domain"
)

const maxOptionItemsPerGroup int32 = 200

type DetailUseCase struct {
	MenuRepo  domain.MenuItemRepository
	GroupRepo domain.OptionRepository
	ItemRepo  domain.OptionItemRepository
}

func NewDetailUseCase(menuRepo domain.MenuItemRepository, groupRepo domain.OptionRepository, itemRepo domain.OptionItemRepository) *DetailUseCase {
	return &DetailUseCase{
		MenuRepo:  menuRepo,
		GroupRepo: groupRepo,
		ItemRepo:  itemRepo,
	}
}

func (uc *DetailUseCase) Execute(ctx context.Context, menuItemID int64, restaurantID int32) (*MenuItemDetailResponse, error) {
	item, err := uc.MenuRepo.GetMenuItemByID(ctx, menuItemID, restaurantID)
	if err != nil {
		return nil, err
	}

	groups, err := uc.GroupRepo.GetOptionGroups(ctx, menuItemID, restaurantID)
	if err != nil {
		return nil, err
	}

	optionGroups := make([]OptionGroupDetail, 0, len(groups))
	for _, g := range groups {
		items, _, err := uc.ItemRepo.GetOptionItems(ctx, g.ID, restaurantID, maxOptionItemsPerGroup, 0)
		if err != nil {
			return nil, err
		}
		groupItems := make([]OptionItemDetail, 0, len(items))
		for _, it := range items {
			groupItems = append(groupItems, OptionItemDetail{
				ID:             it.ID,
				OptionGroupID:  it.OptionGroupID,
				Name:           it.Name,
				LinkedMenuItem: it.LinkedMenuItem,
				PriceDelta:     int64(it.PriceDelta),
				QuantityMin:    it.QuantityMin,
				QuantityMax:    it.QuantityMax,
				SortOrder:      it.SortOrder,
				IsActive:       it.IsActive,
			})
		}
		optionGroups = append(optionGroups, OptionGroupDetail{
			ID:         g.ID,
			Name:       g.Name,
			MinSelect:  g.MinSelect,
			MaxSelect:  g.MaxSelect,
			IsRequired: g.IsRequired,
			SortOrder:  g.SortOrder,
			Items:      groupItems,
		})
	}

	return &MenuItemDetailResponse{
		GetMenuItemResponse: GetMenuItemResponse{
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
		},
		OptionGroups: optionGroups,
	}, nil
}
