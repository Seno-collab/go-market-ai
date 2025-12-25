package db

import (
	"context"
	"go-ai/internal/menu/domain"
	"go-ai/internal/menu/infrastructure/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VariantRepo struct {
	Sqlc *sqlc.Queries
	Pool *pgxpool.Pool
}

func NewVariantRepo(pool *pgxpool.Pool) *VariantRepo {
	return &VariantRepo{
		Sqlc: sqlc.New(pool),
		Pool: pool,
	}
}

func (vr *VariantRepo) Create(ctx context.Context, v domain.Variant) error {
	err := vr.Sqlc.CreateVariant(ctx, sqlc.CreateVariantParams{
		Name:       v.Name,
		PriceDelta: v.PriceDelta.Numeric(),
		IsDefault:  v.IsDefault,
		SortOrder:  0,
	})
	if err != nil {
		return err
	}
	return nil
}

func (vr *VariantRepo) Update(ctx context.Context, v domain.Variant, id int64) error {
	err := vr.Sqlc.UpdateVariant(ctx, sqlc.UpdateVariantParams{
		ID:         id,
		Name:       v.Name,
		PriceDelta: v.PriceDelta.Numeric(),
		IsDefault:  v.IsDefault,
		SortOrder:  0,
	})
	if err != nil {
		return err
	}
	return nil
}

func (vr *VariantRepo) Delete(ctx context.Context, id int64) error {
	return vr.Delete(ctx, id)
}
