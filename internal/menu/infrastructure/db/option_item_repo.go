package db

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/internal/menu/infrastructure/sqlc"
	"go-ai/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OptionItemRepo struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewOptionItemRepo(pool *pgxpool.Pool) *OptionItemRepo {
	return &OptionItemRepo{
		queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *OptionItemRepo) GetOptionItems(ctx context.Context, groupID int64, restaurantID int32) ([]domain.OptionItem, error) {
	rows, err := r.queries.GetOptionItemsByGroup(ctx, sqlc.GetOptionItemsByGroupParams{
		OptionGroupID: groupID,
		RestaurantID:  restaurantID,
	})
	if err != nil {
		return nil, err
	}
	items := make([]domain.OptionItem, 0, len(rows))
	for _, row := range rows {
		item, err := convertOptionItemModel(row)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *OptionItemRepo) GetOptionItem(ctx context.Context, id int64, restaurantID int32) (domain.OptionItem, error) {
	row, err := r.queries.GetOptionItem(ctx, sqlc.GetOptionItemParams{
		ID:           id,
		RestaurantID: restaurantID,
	})
	if err != nil {
		return domain.OptionItem{}, err
	}
	return convertOptionItemModel(row)
}

func (r *OptionItemRepo) CreateOptionItem(ctx context.Context, item *domain.OptionItem) (int64, error) {
	name := stringPtr(item.Name)
	return r.queries.CreateOptionItem(ctx, sqlc.CreateOptionItemParams{
		OptionGroupID:  item.OptionGroupID,
		Name:           name,
		LinkedMenuItem: item.LinkedMenuItem,
		PriceDelta:     utils.NumericFromMoney(item.PriceDelta),
		QuantityMin:    item.QuantityMin,
		QuantityMax:    ptrToNullableInt(item.QuantityMax),
		SortOrder:      item.SortOrder,
	})
}

func (r *OptionItemRepo) UpdateOptionItem(ctx context.Context, item *domain.OptionItem, restaurantID int32) error {
	name := stringPtr(item.Name)
	return r.queries.UpdateOptionItem(ctx, sqlc.UpdateOptionItemParams{
		ID:             item.ID,
		Name:           name,
		LinkedMenuItem: item.LinkedMenuItem,
		PriceDelta:     utils.NumericFromMoney(item.PriceDelta),
		QuantityMin:    item.QuantityMin,
		QuantityMax:    ptrToNullableInt(item.QuantityMax),
		SortOrder:      item.SortOrder,
		RestaurantID:   restaurantID,
	})
}

func (r *OptionItemRepo) DeleteOptionItem(ctx context.Context, id int64, restaurantID int32) error {
	return r.queries.DeleteOptionItem(ctx, sqlc.DeleteOptionItemParams{
		ID:           id,
		RestaurantID: restaurantID,
	})
}
