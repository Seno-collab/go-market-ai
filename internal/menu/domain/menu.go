package domain

import (
	"context"
	"go-ai/pkg/utils"
)

// MenuType is an alias of MenuItemType to reuse the same allowed values.
type MenuType = MenuItemType

type Menu struct {
	ID        int64
	Name      string
	Type      MenuType
	ImageURL  utils.Url
	BasePrice utils.Money
	Topics    []string
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
