package domain

import "go-ai/pkg/utils"

type Selection struct {
	VariantID       *int64
	OptionItemIDs   []int64
	ComboSelections map[int64]int64
}

func CalculatePrice(item MenuItem, sel Selection) (utils.Money, error) {
	price := item.BasePrice

	// Variant
	for _, v := range item.Variants {
		if sel.VariantID != nil && v.ID == *sel.VariantID {
			price = price.Add(v.PriceDelta)
		}
	}

	// Options
	for _, grp := range item.OptionGroups {
		for _, opt := range grp.Options {
			for _, id := range sel.OptionItemIDs {
				if opt.ID == id {
					price = price.Add(opt.PriceDelta)
				}
			}
		}
	}

	// Combo
	for _, grp := range item.ComboGroups {
		chosen, ok := sel.ComboSelections[grp.ID]
		if !ok {
			continue
		}
		for _, ci := range grp.Items {
			if ci.ID == chosen {
				price = price.Add(ci.PriceDelta)
			}
		}
	}

	return price, nil
}
