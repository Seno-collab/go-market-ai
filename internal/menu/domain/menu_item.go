package domain

import (
	"errors"
	"go-ai/pkg/utils"
	"time"
)

type MenuItem struct {
	ID           int64
	RestaurantID int32
	TopicID      *TopicID
	Name         string
	Description  string
	BasePrice    utils.Money
	Type         MenuItemType
	ImageUrl     utils.Url
	Sku          string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Variants     []Variant
	OptionGroups []OptionGroup
	ComboGroups  []ComboGroup
}

type SearchMenuItemsParam struct {
	IsActive     *bool
	Limit        int32
	Offset       int32
	RestaurantID int32
	Filter       string
	Category     string
}
type NewMenuItemParams struct {
	Name         string
	Price        utils.Money
	Type         MenuItemType
	ImageURL     utils.Url
	Description  string
	SKU          string
	RestaurantID int32
	TopicID      TopicID
}

func NewMenuItem(params NewMenuItemParams) (*MenuItem, error) {
	if params.Name == "" {
		return nil, ErrNameRequired
	}
	if !params.Type.Valid() {
		return nil, errors.New("Menu: invalid type")
	}
	return &MenuItem{
		Name:         params.Name,
		BasePrice:    params.Price,
		Type:         params.Type,
		ImageUrl:     params.ImageURL,
		Sku:          params.SKU,
		Description:  params.Description,
		RestaurantID: params.RestaurantID,
		TopicID:      &params.TopicID,
		IsActive:     true,
	}, nil
}

func (m *MenuItem) Validate() error {
	if m.Name == "" {
		return ErrNameRequired
	}
	if !m.Type.Valid() {
		return errors.New("Menu: invalid type")
	}
	return nil
}
