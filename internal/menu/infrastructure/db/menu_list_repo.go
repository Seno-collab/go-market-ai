package db

import (
	"context"

	"go-ai/internal/menu/domain"
	"go-ai/internal/menu/infrastructure/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MenuListRepo struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewMenuListRepo(pool *pgxpool.Pool) *MenuListRepo {
	return &MenuListRepo{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *MenuListRepo) ListMenus(ctx context.Context, params domain.ListMenusParams) ([]domain.Menu, *int64, error) {
	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	cursor := int64(0)
	if params.Cursor != nil {
		cursor = *params.Cursor
	}

	rows, err := r.queries.ListMenus(ctx, sqlc.ListMenusParams{
		RestaurantID: params.RestaurantID,
		MenuType:     string(params.Type),
		Cursor:       cursor,
		TopicNames:   params.Topics,
		PageSize:     limit + 1,
	})
	if err != nil {
		return nil, nil, err
	}

	menus := make([]domain.Menu, 0, len(rows))
	for _, row := range rows {
		menus = append(menus, domain.Menu{
			ID:           row.ID,
			Name:         row.Name,
			RestaurantID: row.RestaurantID,
			Type:         domain.MenuType(row.Type),
		})
	}

	var nextCursor *int64
	if int32(len(menus)) > limit {
		menus = menus[:limit]
		next := menus[len(menus)-1].ID
		nextCursor = &next
	}
	return menus, nextCursor, nil
}
