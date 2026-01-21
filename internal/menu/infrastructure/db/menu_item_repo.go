package db

import (
	"context"
	"errors"
	"go-ai/internal/menu/domain"
	"go-ai/internal/menu/infrastructure/sqlc"
	"go-ai/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	var description string
	if row.Description != nil {
		description = *row.Description
	}
	return domain.MenuItem{
		ID:           row.ID,
		RestaurantID: row.RestaurantID,
		Name:         row.Name,
		Description:  description,
		BasePrice:    price,
		Type:         domain.MenuItemType(row.Type),
		IsActive:     row.IsActive,
	}, nil
}

func (r *MenuItemRepo) CreateMenuItem(ctx context.Context, item *domain.MenuItem) (int64, error) {
	price := item.BasePrice.Numeric()
	url := item.ImageUrl.String()
	var topicID *int64
	if item.TopicID != nil && *item.TopicID > 0 {
		id := int64(*item.TopicID)
		topicID = &id
		// Validate topic belongs to the same restaurant to avoid FK violations.
		if _, err := r.queries.GetTopic(ctx, sqlc.GetTopicParams{
			ID:           id,
			RestaurantID: item.RestaurantID,
		}); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return 0, domain.ErrTopicNotFound
			}
			return 0, err
		}
	}
	id, err := r.queries.CreateMenuItem(ctx, sqlc.CreateMenuItemParams{
		RestaurantID: item.RestaurantID,
		TopicID:      topicID,
		Type:         string(item.Type),
		Name:         item.Name,
		Description:  &item.Description,
		ImageUrl:     &url,
		Sku:          &item.Sku,
		BasePrice:    price,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" && pgErr.ConstraintName == "menu_item_topic_id_fkey" {
			return 0, domain.ErrTopicNotFound
		}
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
	rows, err := r.queries.GetMenuItemsByRestaurant(ctx, sqlc.GetMenuItemsByRestaurantParams{
		RestaurantID: param.RestaurantID,
		Name:         param.Filter,
		Type:         param.Category,
		OffsetValue:  param.Offset,
		LimitValue:   param.Limit,
		IsActive:     param.IsActive,
	})

	if err != nil {
		return nil, 0, err
	}
	items := make([]domain.MenuItem, 0, len(rows))
	for _, row := range rows {
		price, err := utils.NumericToMoney(row.BasePrice)
		if err != nil {
			return nil, 0, err
		}
		var imageUrl string
		if row.ImageUrl != nil {
			imageUrl = *row.ImageUrl
		}
		url, err := utils.NewUrl(imageUrl)
		if err != nil {
			return nil, 0, err
		}
		var description string
		if row.Description != nil {
			description = *row.Description
		}
		var sku string
		if row.Sku != nil {
			sku = *row.Sku
		}
		items = append(items, domain.MenuItem{
			ID:           row.ID,
			Description:  description,
			Type:         domain.MenuItemType(row.Type),
			Name:         row.Name,
			ImageUrl:     url,
			Sku:          sku,
			BasePrice:    price,
			IsActive:     row.IsActive,
			RestaurantID: row.RestaurantID,
			UpdatedAt:    row.UpdatedAt,
			CreatedAt:    row.CreatedAt,
		})
	}
	total, err := r.queries.CountMenuItems(ctx, sqlc.CountMenuItemsParams{
		RestaurantID: param.RestaurantID,
		IsActive:     param.IsActive,
		Name:         &param.Filter,
		Type:         &param.Category,
	})
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
