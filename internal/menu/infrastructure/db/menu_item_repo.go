package db

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/internal/menu/infrastructure/sqlc"
	"go-ai/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MenuRepo struct {
	Sqlc *sqlc.Queries
	Pool *pgxpool.Pool
}

func NewMenuRepo(pool *pgxpool.Pool) *MenuRepo {
	return &MenuRepo{
		Sqlc: sqlc.New(pool),
		Pool: pool,
	}
}

func (r *MenuRepo) GetMenuItemByID(ctx context.Context, id int64) (*domain.MenuItem, error) {
	row, err := r.Sqlc.GetMenuItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	price, _ := utils.NewMoney(row.BasePrice.Int.Int64())
	return &domain.MenuItem{
		ID:           row.ID,
		RestaurantID: row.RestaurantID,
		Name:         row.Name,
		Description:  *row.Description,
		BasePrice:    price,
		Type:         domain.MenuItemType(row.Type.(string)),
		IsActive:     row.IsActive,
	}, nil
}

func (r *MenuRepo) CreateMenuItem(ctx context.Context, item *domain.MenuItem) (int64, error) {
	price := utils.NumericFromMoney(item.BasePrice)
	url := item.ImageUrl.String()
	id, err := r.Sqlc.CreateMenuItem(ctx, sqlc.CreateMenuItemParams{
		RestaurantID: item.RestaurantID,
		TopicID:      item.TopicID,
		Type:         item.Type,
		Name:         item.Name,
		Description:  &item.Description,
		ImageUrl:     &url,
		Sku:          &item.Sku,
		BasePrice:    price,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}
