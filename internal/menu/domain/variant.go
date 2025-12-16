package domain

import "go-ai/pkg/utils"

type Variant struct {
	ID         int64
	MenuItemID int64
	Name       string
	PriceDelta utils.Money
	IsDefault  bool
}

func NewVariant(name string, delta utils.Money, isDefault bool) (*Variant, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if delta < 0 {
		return nil, ErrInvalidVariantDelta
	}
	return &Variant{Name: name, PriceDelta: delta, IsDefault: isDefault}, nil
}

func ValidateVariants(vars []Variant) error {
	countDefault := 0
	for _, v := range vars {
		if v.IsDefault {
			countDefault++
		}
	}
	if countDefault > 1 {
		return ErrVariantDefaultMultiple
	}
	return nil
}
