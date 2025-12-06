package restaurantrepo

import (
	"context"
	"go-ai/internal/restaurant/domain/restaurant"
	sqlc "go-ai/internal/restaurant/infra/sqlc/restaurant"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RestaurantRepo struct {
	pool *pgxpool.Pool
	q    *sqlc.Queries
}

func NewRestaurantRepo(pool *pgxpool.Pool) *RestaurantRepo {
	return &RestaurantRepo{
		q:    sqlc.New(pool),
		pool: pool,
	}
}

func (rr *RestaurantRepo) Create(ctx context.Context, r *restaurant.Entity) (int32, error) {
	tx, err := rr.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	qtx := rr.q.WithTx(tx)
	id, err := qtx.CreateRestaurant(ctx, sqlc.CreateRestaurantParams{
		Name:        r.Name,
		Description: &r.Description,
		Address:     &r.Address,
		Category:    &r.Category,
		City:        &r.City,
		District:    &r.District,
		LogoUrl:     &r.LogoUrl,
		BannerUrl:   &r.BannerUrl,
		PhoneNumber: &r.PhoneNumber,
		WebsiteUrl:  &r.WebsiteUrl,
		Email:       &r.Email,
		CreatedBy:   r.UserID,
	})
	if err != nil {
		return 0, err
	}
	for _, h := range r.Hours {
		err := qtx.CreateRestaurantHours(ctx, sqlc.CreateRestaurantHoursParams{
			RestaurantID: id,
			DayOfWeek:    int32(h.Day),
			OpenTime:     h.OpenTime,
			CloseTime:    h.CloseTime,
		})
		if err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return id, nil
}

func (rr *RestaurantRepo) GetById(ctx context.Context, id int32) (*restaurant.Entity, error) {
	records, err := rr.q.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	hours := []restaurant.Hours{}
	for _, r := range records {
		dayOfWeek, err := restaurant.ParseDayOfWeek(r.DayOfWeek)
		if err != nil {
			return nil, err
		}
		hours = append(hours, restaurant.Hours{
			Day:       dayOfWeek,
			OpenTime:  r.OpenTime,
			CloseTime: r.CloseTime,
		})
	}
	first := records[0]
	entity := &restaurant.Entity{
		Name:        first.Name,
		Description: *first.Description,
		Address:     *first.Address,
		Category:    *first.Category,
		City:        *first.City,
		District:    *first.District,
		LogoUrl:     *first.LogoUrl,
		BannerUrl:   *first.BannerUrl,
		PhoneNumber: *first.PhoneNumber,
		WebsiteUrl:  *first.WebsiteUrl,
		Email:       *first.Email,
		UserID:      first.CreatedBy,
		Hours:       hours,
	}
	return entity, nil
}

func (rr *RestaurantRepo) GetByName(ctx context.Context, name string) (*restaurant.Entity, error) {
	records, err := rr.q.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	hours := []restaurant.Hours{}
	for _, r := range records {
		dayOfWeek, err := restaurant.ParseDayOfWeek(r.DayOfWeek)
		if err != nil {
			return nil, err
		}
		hours = append(hours, restaurant.Hours{
			Day:       dayOfWeek,
			OpenTime:  r.OpenTime,
			CloseTime: r.CloseTime,
		})
	}
	first := records[0]
	entity := &restaurant.Entity{
		Name:        first.Name,
		Description: *first.Description,
		Address:     *first.Address,
		Category:    *first.Category,
		City:        *first.City,
		District:    *first.District,
		LogoUrl:     *first.LogoUrl,
		BannerUrl:   *first.BannerUrl,
		PhoneNumber: *first.PhoneNumber,
		WebsiteUrl:  *first.WebsiteUrl,
		Email:       *first.Email,
		UserID:      first.CreatedBy,
		Hours:       hours,
	}
	return entity, nil
}

func (rr *RestaurantRepo) Update(ctx context.Context, r *restaurant.Entity, id int32) error {
	tx, err := rr.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := rr.q.WithTx(tx)
	err = rr.q.UpdateRestaurant(ctx, sqlc.UpdateRestaurantParams{
		ID:          id,
		Name:        r.Name,
		Description: &r.Description,
		Address:     &r.Address,
		Category:    &r.Category,
		City:        &r.City,
		District:    &r.District,
		LogoUrl:     &r.LogoUrl,
		BannerUrl:   &r.BannerUrl,
		PhoneNumber: &r.PhoneNumber,
		WebsiteUrl:  &r.WebsiteUrl,
		Email:       &r.Email,
		UpdatedBy:   r.UserID,
	})
	if err != nil {
		return err
	}
	if err := qtx.SoftDeleteRestaurantHours(ctx, id); err != nil {
		return err
	}
	for _, h := range r.Hours {
		err = qtx.CreateRestaurantHours(ctx, sqlc.CreateRestaurantHoursParams{
			RestaurantID: id,
			OpenTime:     h.OpenTime,
			DayOfWeek:    int32(h.Day),
			CloseTime:    h.CloseTime,
		})
		if err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (rr *RestaurantRepo) SoftDelete(ctx context.Context, id int32, userID uuid.UUID) error {
	tx, err := rr.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := rr.q.WithTx(tx)
	if err := qtx.SoftDeleteRestaurant(ctx, sqlc.SoftDeleteRestaurantParams{
		ID:        id,
		UpdatedBy: userID,
	}); err != nil {
		return err
	}
	if err := qtx.SoftDeleteRestaurantHours(ctx, id); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
