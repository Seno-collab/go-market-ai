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

func NewMenuItem(name string, price utils.Money, t MenuItemType, url utils.Url, description, sku string, restaurantID int32, topicID TopicID) (*MenuItem, error) {
	if name == "" {
		return nil, ErrNameRequired
	}
	if !t.Valid() {
		return nil, errors.New("menu: invalid type")
	}
	return &MenuItem{
		Name:         name,
		BasePrice:    price,
		Type:         t,
		ImageUrl:     url,
		Sku:          sku,
		Description:  description,
		RestaurantID: restaurantID,
		TopicID:      &topicID,
		IsActive:     true,
	}, nil
}

func (m *MenuItem) Validate() error {
	if m.Name == "" {
		return ErrNameRequired
	}
	if !m.Type.Valid() {
		return errors.New("menu: invalid type")
	}
	return nil
}
