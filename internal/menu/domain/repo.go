package domain

import (
	"context"
	"go-ai/internal/transport/response"
)

type MenuItemRepository interface {
	GetMenuItems(ctx context.Context, param SearchMenuItemsParam) ([]MenuItem, int64, error)
	GetMenuItemByID(ctx context.Context, id int64, restaurantID int32) (MenuItem, error)
	CreateMenuItem(ctx context.Context, item *MenuItem) (int64, error)
	UpdateMenuItem(ctx context.Context, item *MenuItem) error
	DeleteMenuItem(ctx context.Context, id int64, restaurantID int32) error
	UpdateStatusMenuItem(ctx context.Context, restaurantID int32, id int64, isActive bool) error
}

type TopicRepository interface {
	GetTopics(ctx context.Context, name string, restaurantID, limit, offset int32) ([]Topic, int64, error)
	CreateTopic(ctx context.Context, t *Topic) (TopicID, error)
	GetTopic(ctx context.Context, id TopicID, restaurantID int32) (Topic, error)
	UpdateTopic(ctx context.Context, t *Topic) error
	DeleteTopic(ctx context.Context, id TopicID, restaurantID int32) error
	GetTopicsByRestaurantCombobox(ctx context.Context, restaurantID int32) (*[]response.Combobox, error)
}

type ComboRepository interface {
	GetComboGroups(ctx context.Context, itemID int64) ([]ComboGroup, error)
	GetComboGroupItems(ctx context.Context, groupID int64) ([]ComboGroupItem, error)
	CreateComboGroup(ctx context.Context, g *ComboGroup) (int64, error)
	CreateComboGroupItem(ctx context.Context, i *ComboGroupItem) (int64, error)
	DeleteComboGroup(ctx context.Context, id int64) error
	DeleteComboGroupItem(ctx context.Context, id int64) error
}

type VariantsRepository interface {
	GetVariants(ctx context.Context, itemID int64) ([]Variant, error)
	CreateVariant(ctx context.Context, v *Variant) (int64, error)
	UpdateVariant(ctx context.Context, v *Variant) error
	DeleteVariant(ctx context.Context, id int64) error
}

type OptionItemRepository interface {
	GetOptionItems(ctx context.Context, groupID int64, restaurantID, limit, offset int32) ([]OptionItem, int64, error)
	GetOptionItem(ctx context.Context, id int64, restaurantID int32) (OptionItem, error)
	CreateOptionItem(ctx context.Context, o *OptionItem) (int64, error)
	UpdateOptionItem(ctx context.Context, o *OptionItem, restaurantID int32) error
	DeleteOptionItem(ctx context.Context, id int64, restaurantID int32) error
}

type OptionRepository interface {
	GetOptionGroups(ctx context.Context, itemID int64, restaurantID int32) ([]OptionGroup, error)
	GetOptionGroup(ctx context.Context, id int64, restaurantID int32) (OptionGroup, error)
	CreateOptionGroup(ctx context.Context, g *OptionGroup) (int64, error)
	AttachOptionGroupToItem(ctx context.Context, menuItemID, optionGroupID int64, sortOrder int32) error
	UpdateOptionGroup(ctx context.Context, g *OptionGroup) error
	DeleteOptionGroup(ctx context.Context, id int64, restaurantID int32) error
}
