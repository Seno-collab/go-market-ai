package restaurantrepo

import (
	"context"
	"errors"
	"go-ai/internal/domain/restaurant"
	sqlc "go-ai/internal/infra/sqlc/restaurant"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RestaurantRepo struct {
	q *sqlc.Queries
}

func NewRestaurantRepo(pool *pgxpool.Pool) *RestaurantRepo {
	return &RestaurantRepo{
		q: sqlc.New(pool),
	}
}

func (rr *RestaurantRepo) Create(ctx context.Context, r *restaurant.Entity) (int32, error) {
	_, err := rr.q.GetByName(ctx, r.Name)
	if err != nil || !errors.Is(err, pgx.ErrNoRows) {
		return 0, err
	}
	id, err := rr.q.CreateRestaurant(context.Background(), sqlc.CreateRestaurantParams{
		Name:        r.Name,
		Description: r.Description,
		Address:     r.Address,
		Category:    r.Category,
		City:        r.City,
		District:    r.District,
		LogoUrl:     r.LogoUrl,
		BannerUrl:   r.BannerUrl,
		PhoneNumber: r.PhoneNumber,
		WebsiteUrl:  r.WebsiteUrl,
		Email:       r.Email,
		UserID:      r.UserID,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (rr *RestaurantRepo) GetById(ctx context.Context, id int32) (*restaurant.Entity, error) {
	record, err := rr.q.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return MappingRestaurant(&record), nil
}

func (rr *RestaurantRepo) GetByName(ctx context.Context, name string) (*restaurant.Entity, error) {
	record, err := rr.q.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return MappingRestaurant(&record), nil

}

func MappingRestaurant(record *sqlc.Restaurant) *restaurant.Entity {
	return &restaurant.Entity{
		Name:        record.Name,
		Description: record.Description,
		Address:     record.Address,
		Category:    record.Category,
		City:        record.City,
		District:    record.District,
		LogoUrl:     record.LogoUrl,
		BannerUrl:   record.BannerUrl,
		PhoneNumber: record.PhoneNumber,
		WebsiteUrl:  record.WebsiteUrl,
		Email:       record.Email,
		UserID:      record.UserID,
	}
}
