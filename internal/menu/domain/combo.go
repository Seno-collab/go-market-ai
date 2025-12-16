package domain

import "go-ai/pkg/utils"

type ComboGroup struct {
	ID          int64
	ComboItemID int64
	Name        string
	MinSelect   int
	MaxSelect   int
	Items       []ComboGroupItem
}

type ComboGroupItem struct {
	ID              int64
	ComboGroupID    int64
	MenuItemID      int64
	PriceDelta      utils.Money
	QuantityDefault int
	QuantityMin     int
	QuantityMax     int
}

func NewComboGroup(name string, min, max int) (*ComboGroup, error) {
	if max < min {
		return nil, ErrOptionGroupInvalidMinMax
	}
	return &ComboGroup{Name: name, MinSelect: min, MaxSelect: max}, nil
}
