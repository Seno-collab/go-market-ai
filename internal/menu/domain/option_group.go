package domain

import "errors"

type OptionGroup struct {
	ID         int64
	MenuItemID int64
	Name       string
	MinSelect  int
	MaxSelect  *int
	IsRequired bool
	Options    []OptionItem
}

func NewOptionGroup(name string, min, max *int, required bool) (*OptionGroup, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if max != nil && *max < minValue(min) {
		return nil, ErrOptionGroupInvalidMinMax
	}
	if required && minValue(min) < 1 {
		return nil, errors.New("option_group: required must have min >= 1")
	}
	return &OptionGroup{
		Name:       name,
		MinSelect:  minValue(min),
		MaxSelect:  max,
		IsRequired: required,
	}, nil
}

func minValue(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}
