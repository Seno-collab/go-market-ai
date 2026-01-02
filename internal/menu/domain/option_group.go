package domain

import "errors"

type OptionGroup struct {
	ID           int64
	RestaurantID int32
	Name         string
	MinSelect    int32
	MaxSelect    int32
	IsRequired   bool
	SortOrder    int32
	Options      []OptionItem
}

func NewOptionGroup(restaurantID int32, name string, min int32, max int32, required bool, sortOrder int32) (*OptionGroup, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if max > 0 && max < min {
		return nil, ErrOptionGroupInvalidMinMax
	}
	if required && min < 1 {
		return nil, errors.New("Option_group: required must have min >= 1")
	}
	return &OptionGroup{
		RestaurantID: restaurantID,
		Name:         name,
		MinSelect:    min,
		MaxSelect:    max,
		IsRequired:   required,
		SortOrder:    sortOrder,
	}, nil
}
