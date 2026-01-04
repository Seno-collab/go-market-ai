package db

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/internal/menu/infrastructure/sqlc"
	"go-ai/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MenuItemRepo struct {
	queries *sqlc.Queries
	pool    *pgxpool.Pool
}

func NewMenuRepo(pool *pgxpool.Pool) *MenuItemRepo {
	return &MenuItemRepo{
		queries: sqlc.New(pool),
		pool:    pool,
	}
}

func (r *MenuItemRepo) GetMenuItemByID(ctx context.Context, id int64, restaurantID int32) (domain.MenuItem, error) {
	row, err := r.queries.GetMenuItemByID(ctx, sqlc.GetMenuItemByIDParams{
		ID:           id,
		RestaurantID: restaurantID,
	})
	if err != nil {
		return domain.MenuItem{}, err
	}
	price, err := utils.NumericToMoney(row.BasePrice)
	if err != nil {
		return domain.MenuItem{}, err
	}
	return domain.MenuItem{
		ID:           row.ID,
		RestaurantID: row.RestaurantID,
		Name:         row.Name,
		Description:  *row.Description,
		BasePrice:    price,
		Type:         domain.MenuItemType(row.Type),
		IsActive:     row.IsActive,
	}, nil
}

func (r *MenuItemRepo) CreateMenuItem(ctx context.Context, item *domain.MenuItem) (int64, error) {
	price := item.BasePrice.Numeric()
	url := item.ImageUrl.String()
	id := int64(*item.TopicID)
	id, err := r.queries.CreateMenuItem(ctx, sqlc.CreateMenuItemParams{
		RestaurantID: item.RestaurantID,
		TopicID:      &id,
		Type:         string(item.Type),
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
	err := r.queries.UpdateMenuItem(ctx, sqlc.UpdateMenuItemParams{
		ID:           item.ID,
		Description:  &item.Description,
		Type:         string(item.Type),
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
	return r.queries.DeleteMenuItem(ctx, sqlc.DeleteMenuItemParams{
		ID:           id,
		RestaurantID: restaurantID,
	})
}

func (r *MenuItemRepo) GetMenuItems(ctx context.Context, param domain.SearchMenuItemsParam) ([]domain.MenuItem, int64, error) {
	var rows []sqlc.MenuItem
	var err error
	if param.IsActive == nil {
		rows, err = r.queries.GetMenuItemsByRestaurant(ctx, sqlc.GetMenuItemsByRestaurantParams{
			RestaurantID: param.RestaurantID,
			Limit:        param.Limit,
			Offset:       param.Offset,
			Column2:      param.Filter,
			Column3:      param.Category,
		})
	} else {
		isActive := *param.IsActive
		rows, err = r.queries.GetMenuItemsByRestaurantAndActive(ctx, sqlc.GetMenuItemsByRestaurantAndActiveParams{
			RestaurantID: param.RestaurantID,
			Limit:        param.Limit,
			Offset:       param.Offset,
			IsActive:     isActive,
			Column3:      param.Filter,
			Column4:      param.Category,
		})
	}

	if err != nil {
		return nil, 0, err
	}
	items := make([]domain.MenuItem, 0, len(rows))
	for _, row := range rows {
		price, err := utils.NumericToMoney(row.BasePrice)
		if err != nil {
			return nil, 0, err
		}
		url, err := utils.NewUrl(*row.ImageUrl)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, domain.MenuItem{
			ID:           row.ID,
			Description:  *row.Description,
			Type:         domain.MenuItemType(row.Type),
			Name:         row.Name,
			ImageUrl:     url,
			Sku:          *row.Sku,
			BasePrice:    price,
			IsActive:     row.IsActive,
			RestaurantID: row.RestaurantID,
			UpdatedAt:    row.UpdatedAt,
			CreatedAt:    row.CreatedAt,
		})
	}
	total, err := r.queries.CountMenuItems(ctx, param.RestaurantID)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *MenuItemRepo) UpdateStatusMenuItem(ctx context.Context, restaurantID int32, id int64, isActive bool) error {
	return r.queries.UpdateStatusMenuItem(ctx, sqlc.UpdateStatusMenuItemParams{
		RestaurantID: restaurantID,
		IsActive:     isActive,
		ID:           id,
	})
}
