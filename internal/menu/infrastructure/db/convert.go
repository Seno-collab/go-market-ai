package db

import (
	"go-ai/internal/menu/domain"
	"go-ai/internal/menu/infrastructure/sqlc"
	"go-ai/pkg/utils"
)

func convertOptionGroupModel(row sqlc.OptionGroup) domain.OptionGroup {
	return domain.OptionGroup{
		ID:           row.ID,
		RestaurantID: row.RestaurantID,
		Name:         row.Name,
		MinSelect:    row.MinSelect,
		MaxSelect:    nullableIntToPtr(row.MaxSelect),
		IsRequired:   row.IsRequired,
		SortOrder:    row.SortOrder,
	}
}

func convertOptionItemModel(row sqlc.OptionItem) (domain.OptionItem, error) {
	price, err := utils.NumericToMoney(row.PriceDelta)
	if err != nil {
		return domain.OptionItem{}, err
	}
	name := ""
	if row.Name != nil {
		name = *row.Name
	}
	return domain.OptionItem{
		ID:             row.ID,
		OptionGroupID:  row.OptionGroupID,
		Name:           name,
		LinkedMenuItem: row.LinkedMenuItem,
		PriceDelta:     price,
		QuantityMin:    row.QuantityMin,
		QuantityMax:    nullableIntToPtr(row.QuantityMax),
		SortOrder:      row.SortOrder,
		IsActive:       row.IsActive,
	}, nil
}

func nullableIntToPtr(v int) *int32 {
	if v == 0 {
		return nil
	}
	value := int32(v)
	return &value
}

func ptrToNullableInt(v *int32) int {
	if v == nil {
		return 0
	}
	return int(*v)
}

func stringPtr(v string) *string {
	if v == "" {
		return nil
	}
	value := v
	return &value
}
