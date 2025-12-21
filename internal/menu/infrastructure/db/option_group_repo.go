package db

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/internal/menu/infrastructure/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OptionGroupRepo struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewOptionGroupRepo(pool *pgxpool.Pool) *OptionGroupRepo {
	return &OptionGroupRepo{
		queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *OptionGroupRepo) GetOptionGroups(ctx context.Context, menuItemID int64, restaurantID int32) ([]domain.OptionGroup, error) {
	rows, err := r.queries.GetOptionGroupsByItem(ctx, sqlc.GetOptionGroupsByItemParams{
		MenuItemID:   menuItemID,
		RestaurantID: restaurantID,
	})
	if err != nil {
		return nil, err
	}
	groups := make([]domain.OptionGroup, 0, len(rows))
	for _, row := range rows {
		groups = append(groups, convertOptionGroupModel(row))
	}
	return groups, nil
}

func (r *OptionGroupRepo) GetOptionGroup(ctx context.Context, id int64, restaurantID int32) (domain.OptionGroup, error) {
	row, err := r.queries.GetOptionGroup(ctx, sqlc.GetOptionGroupParams{
		ID:           id,
		RestaurantID: restaurantID,
	})
	if err != nil {
		return domain.OptionGroup{}, err
	}
	return convertOptionGroupModel(row), nil
}

func (r *OptionGroupRepo) CreateOptionGroup(ctx context.Context, g *domain.OptionGroup) (int64, error) {
	return r.queries.CreateOptionGroup(ctx, sqlc.CreateOptionGroupParams{
		RestaurantID: g.RestaurantID,
		Name:         g.Name,
		MinSelect:    g.MinSelect,
		MaxSelect:    ptrToNullableInt(g.MaxSelect),
		IsRequired:   g.IsRequired,
		SortOrder:    g.SortOrder,
	})
}

func (r *OptionGroupRepo) AttachOptionGroupToItem(ctx context.Context, menuItemID, optionGroupID int64, sortOrder int32) error {
	return r.queries.AttachOptionGroupToItem(ctx, sqlc.AttachOptionGroupToItemParams{
		MenuItemID:    menuItemID,
		OptionGroupID: optionGroupID,
		SortOrder:     sortOrder,
	})
}

func (r *OptionGroupRepo) UpdateOptionGroup(ctx context.Context, g *domain.OptionGroup) error {
	return r.queries.UpdateOptionGroup(ctx, sqlc.UpdateOptionGroupParams{
		ID:           g.ID,
		RestaurantID: g.RestaurantID,
		Name:         g.Name,
		MinSelect:    g.MinSelect,
		MaxSelect:    ptrToNullableInt(g.MaxSelect),
		IsRequired:   g.IsRequired,
		SortOrder:    g.SortOrder,
	})
}

func (r *OptionGroupRepo) DeleteOptionGroup(ctx context.Context, id int64, restaurantID int32) error {
	return r.queries.DeleteOptionGroup(ctx, sqlc.DeleteOptionGroupParams{
		ID:           id,
		RestaurantID: restaurantID,
	})
}
