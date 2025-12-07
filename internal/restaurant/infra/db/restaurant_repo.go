package restaurantrepo

import (
	"context"
	"errors"
	"go-ai/internal/restaurant/domain/restaurant"
	sqlc "go-ai/internal/restaurant/infra/sqlc/restaurant"
	"go-ai/pkg/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RestaurantRepo struct {
	Pool *pgxpool.Pool
	Sqlc *sqlc.Queries
}

func NewRestaurantRepo(pool *pgxpool.Pool) *RestaurantRepo {
	return &RestaurantRepo{
		Sqlc: sqlc.New(pool),
		Pool: pool,
	}
}

func (rr *RestaurantRepo) Create(ctx context.Context, r *restaurant.Entity, userID uuid.UUID) (int32, error) {

	tx, err := rr.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	qtx := rr.Sqlc.WithTx(tx)
	_, err = qtx.GetByName(ctx, r.Name)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return 0, err
		}
	}
	banner := r.BannerUrl.String()
	phone := r.PhoneNumber.String()
	website := r.WebsiteUrl.String()
	email := r.Email.String()
	logo := r.LogoUrl.String()
	id, err := qtx.CreateRestaurant(ctx, sqlc.CreateRestaurantParams{
		Name:        r.Name,
		Description: &r.Description,
		Address:     &r.Address,
		Category:    &r.Category,
		City:        &r.City,
		District:    &r.District,
		LogoUrl:     &logo,
		BannerUrl:   &banner,
		PhoneNumber: &phone,
		WebsiteUrl:  &website,
		Email:       &email,
		CreatedBy:   r.CreatedBy,
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
	err = qtx.UpsertRestaurantUser(ctx, sqlc.UpsertRestaurantUserParams{
		RestaurantID: id,
		UserID:       userID,
		Role:         restaurant.Owner,
	})
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return id, nil
}

func (rr *RestaurantRepo) GetById(ctx context.Context, id int32) (*restaurant.Entity, error) {
	records, err := rr.Sqlc.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if records == nil {
		return nil, restaurant.ErrRestaurantNotExists
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
	logoUrl, err := restaurant.NewUrl(*first.LogoUrl)
	if err != nil {
		return nil, err
	}

	bannerUrl, err := restaurant.NewUrl(*first.BannerUrl)
	if err != nil {
		return nil, err
	}

	phone, err := restaurant.NewPhone(*first.PhoneNumber)
	if err != nil {
		return nil, err
	}

	websiteUrl, err := restaurant.NewUrl(*first.WebsiteUrl)
	if err != nil {
		return nil, err
	}

	email, err := utils.NewEmail(*first.Email)
	if err != nil {
		return nil, err
	}
	entity := &restaurant.Entity{
		Name:        first.Name,
		Description: *first.Description,
		Address:     *first.Address,
		Category:    *first.Category,
		City:        *first.City,
		District:    *first.District,
		LogoUrl:     logoUrl,
		BannerUrl:   bannerUrl,
		PhoneNumber: phone,
		WebsiteUrl:  websiteUrl,
		Email:       email,
		CreatedBy:   first.CreatedBy,
		Hours:       hours,
	}
	return entity, nil
}

func (rr *RestaurantRepo) GetByName(ctx context.Context, name string) (*restaurant.Entity, error) {
	records, err := rr.Sqlc.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if records == nil {
		return nil, restaurant.ErrRestaurantNotExists
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
	logoUrl, err := restaurant.NewUrl(*first.LogoUrl)
	if err != nil {
		return nil, err
	}

	bannerUrl, err := restaurant.NewUrl(*first.BannerUrl)
	if err != nil {
		return nil, err
	}

	phone, err := restaurant.NewPhone(*first.PhoneNumber)
	if err != nil {
		return nil, err
	}

	websiteUrl, err := restaurant.NewUrl(*first.WebsiteUrl)
	if err != nil {
		return nil, err
	}

	email, err := utils.NewEmail(*first.Email)
	if err != nil {
		return nil, err
	}
	entity := &restaurant.Entity{
		Name:        first.Name,
		Description: *first.Description,
		Address:     *first.Address,
		Category:    *first.Category,
		City:        *first.City,
		District:    *first.District,
		LogoUrl:     logoUrl,
		BannerUrl:   bannerUrl,
		PhoneNumber: phone,
		WebsiteUrl:  websiteUrl,
		Email:       email,
		CreatedBy:   first.CreatedBy,
		Hours:       hours,
	}
	return entity, nil
}

func (rr *RestaurantRepo) Update(ctx context.Context, r *restaurant.Entity, id int32) error {
	tx, err := rr.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := rr.Sqlc.WithTx(tx)
	banner := r.BannerUrl.String()
	phone := r.PhoneNumber.String()
	website := r.WebsiteUrl.String()
	email := r.Email.String()
	logo := r.LogoUrl.String()
	err = rr.Sqlc.UpdateRestaurant(ctx, sqlc.UpdateRestaurantParams{
		ID:          id,
		Name:        r.Name,
		Description: &r.Description,
		Address:     &r.Address,
		Category:    &r.Category,
		City:        &r.City,
		District:    &r.District,
		LogoUrl:     &logo,
		BannerUrl:   &banner,
		PhoneNumber: &phone,
		WebsiteUrl:  &website,
		Email:       &email,
		UpdatedBy:   r.UpdateBy,
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
	tx, err := rr.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := rr.Sqlc.WithTx(tx)
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
