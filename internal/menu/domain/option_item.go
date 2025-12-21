package domain

import "go-ai/pkg/utils"

type OptionItem struct {
	ID             int64
	OptionGroupID  int64
	Name           string
	LinkedMenuItem *int64
	PriceDelta     utils.Money
	QuantityMin    int32
	QuantityMax    *int32
	SortOrder      int32
	IsActive       bool
}

func NewOptionItem(name string, delta utils.Money) (*OptionItem, error) {
	if delta < 0 {
		return nil, ErrOptionItemInvalidPrice
	}
	return &OptionItem{Name: name, PriceDelta: delta}, nil
}
