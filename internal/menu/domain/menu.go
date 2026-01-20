package domain

import "context"

// MenuType is an alias of MenuItemType to reuse the same allowed values.
type MenuType = MenuItemType

type Menu struct {
	ID           int64
	Name         string
	RestaurantID int32
	Type         MenuType
	Topics       []string
}

type ListMenusParams struct {
	RestaurantID int32
	Type         MenuType
	Topics       []string
	Limit        int32
	Cursor       *int64
}

type MenuRepository interface {
	ListMenus(ctx context.Context, params ListMenusParams) ([]Menu, *int64, error)
}
