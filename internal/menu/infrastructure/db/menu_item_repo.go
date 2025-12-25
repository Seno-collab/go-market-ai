package db

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/internal/menu/infrastructure/sqlc"
	"go-ai/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MenuItemRepo struct {
	Sqlc *sqlc.Queries
	Pool *pgxpool.Pool
}

func NewMenuRepo(pool *pgxpool.Pool) *MenuItemRepo {
	return &MenuItemRepo{
		Sqlc: sqlc.New(pool),
		Pool: pool,
	}
}

func (r *MenuItemRepo) GetMenuItemByID(ctx context.Context, id int64, restaurantID int32) (domain.MenuItem, error) {
	row, err := r.Sqlc.GetMenuItemByID(ctx, sqlc.GetMenuItemByIDParams{
		ID:           id,
		RestaurantID: restaurantID,
	})
	if err != nil {
		return domain.MenuItem{}, err
	}
	price, _ := utils.NewMoney(row.BasePrice.Int.Int64())
	return domain.MenuItem{
		ID:           row.ID,
		RestaurantID: row.RestaurantID,
		Name:         row.Name,
		Description:  *row.Description,
		BasePrice:    price,
		Type:         domain.MenuItemType(row.Type.(string)),
		IsActive:     row.IsActive,
	}, nil
}

func (r *MenuItemRepo) CreateMenuItem(ctx context.Context, item *domain.MenuItem) (int64, error) {
	price := item.BasePrice.Numeric()
	url := item.ImageUrl.String()
	id := int64(*item.TopicID)
	id, err := r.Sqlc.CreateMenuItem(ctx, sqlc.CreateMenuItemParams{
		RestaurantID: item.RestaurantID,
		TopicID:      &id,
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

func (r *MenuItemRepo) UpdateMenuItem(ctx context.Context, item *domain.MenuItem) error {
	price := item.BasePrice.Numeric()
	url := item.ImageUrl.String()
	err := r.Sqlc.UpdateMenuItem(ctx, sqlc.UpdateMenuItemParams{
		ID:           item.ID,
		Description:  &item.Description,
		Type:         item.Type,
		Name:         item.Name,
		ImageUrl:     &url,
		Sku:          &item.Sku,
		BasePrice:    price,
		IsActive:     item.IsActive,
		RestaurantID: item.RestaurantID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *MenuItemRepo) DeleteMenuItem(ctx context.Context, id int64, restaurantID int32) error {
	return r.Sqlc.DeleteMenuItem(ctx, sqlc.DeleteMenuItemParams{
		ID:           id,
		RestaurantID: restaurantID,
	})
}

func (r *MenuItemRepo) GetMenuItems(ctx context.Context, restaurantID int32) ([]domain.MenuItem, error) {
	rows, err := r.Sqlc.GetMenuItemsByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	items := make([]domain.MenuItem, 0, len(rows))
	for _, row := range rows {
		price, err := utils.NumericToMoney(row.BasePrice)
		if err != nil {
			return nil, err
		}
		url, err := utils.NewUrl(*row.ImageUrl)
		if err != nil {
			return nil, err
		}
		items = append(items, domain.MenuItem{
			ID:           row.ID,
			Description:  *row.Description,
			Type:         domain.MenuItemType(row.Type.(string)),
			Name:         row.Name,
			ImageUrl:     url,
			Sku:          *row.Sku,
			BasePrice:    price,
			IsActive:     row.IsActive,
			RestaurantID: row.RestaurantID,
		})
	}
	return items, nil
}
