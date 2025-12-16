package domain

type MenuItemType string

const (
	ItemDish  MenuItemType = "dish"
	ItemDrink MenuItemType = "beverage"
	ItemExtra MenuItemType = "extra"
	ItemCombo MenuItemType = "combo"
)

func (t MenuItemType) Valid() bool {
	switch t {
	case ItemDish, ItemDrink, ItemExtra, ItemCombo:
		return true
	}
	return false
}
