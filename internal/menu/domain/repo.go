package domain

import "context"

type MenuRepository interface {
	// Menu Items
	// GetMenuItems(ctx context.Context, restaurantID int64) ([]MenuItem, error)
	// GetMenuItemByID(ctx context.Context, id int64) (*MenuItem, error)
	CreateMenuItem(ctx context.Context, item *MenuItem) (int64, error)
	// UpdateMenuItem(ctx context.Context, item *MenuItem) error
	// DeleteMenuItem(ctx context.Context, id int64) error
}

type TopicRepository interface {
	GetTopics(ctx context.Context, restaurantID int64) ([]Topic, error)
	CreateTopic(ctx context.Context, t *Topic) (TopicID, error)
	UpdateTopic(ctx context.Context, t *Topic) error
	DeleteTopic(ctx context.Context, id TopicID) error
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
	// Option Items
	GetOptionItems(ctx context.Context, groCreateMenuItempID int64) ([]OptionItem, error)
	CreateOptionItem(ctx context.Context, o *OptionItem) (int64, error)
	UpdateOptionItem(ctx context.Context, o *OptionItem) error
	DeleteOptionItem(ctx context.Context, id int64) error
}

type OptionRepository interface {
	GetOptionGroups(ctx context.Context, itemID int64) ([]OptionGroup, error)
	CreateOptionGroup(ctx context.Context, g *OptionGroup) (int64, error)
	UpdateOptionGroup(ctx context.Context, g *OptionGroup) error
	DeleteOptionGroup(ctx context.Context, id int64) error
}
