package db

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/internal/menu/infrastructure/sqlc"

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

func (r *OptionItemRepo) GetOptionItems(ctx context.Context, groupID int64, restaurantID, limit, offset int32) ([]domain.OptionItem, int64, error) {
	rows, err := r.queries.GetOptionItemsByGroup(ctx, sqlc.GetOptionItemsByGroupParams{
		OptionGroupID: groupID,
		RestaurantID:  restaurantID,
		Limit:         limit,
		Offset:        offset,
	})
	if err != nil {
		return nil, 0, err
	}
	items := make([]domain.OptionItem, 0, len(rows))
	for _, row := range rows {
		item, err := convertOptionItemModel(row)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	total, err := r.queries.CountOptionItems(ctx, sqlc.CountOptionItemsParams{
		OptionGroupID: groupID,
		RestaurantID:  restaurantID,
	})
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
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
		PriceDelta:     item.PriceDelta.Numeric(),
		QuantityMin:    item.QuantityMin,
		QuantityMax:    item.QuantityMax,
		SortOrder:      item.SortOrder,
	})
}

func (r *OptionItemRepo) UpdateOptionItem(ctx context.Context, item *domain.OptionItem, restaurantID int32) error {
	name := stringPtr(item.Name)
	return r.queries.UpdateOptionItem(ctx, sqlc.UpdateOptionItemParams{
		ID:             item.ID,
		Name:           name,
		LinkedMenuItem: item.LinkedMenuItem,
		PriceDelta:     item.PriceDelta.Numeric(),
		QuantityMin:    item.QuantityMin,
		QuantityMax:    item.QuantityMax,
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
