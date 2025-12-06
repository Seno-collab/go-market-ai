package container

import (
	"go-ai/internal/identity/domain/rbac"
	middlewares "go-ai/internal/identity/transport/middelwares"
	restaurantapp "go-ai/internal/restaurant/application/restaurant"
	restaurantrepo "go-ai/internal/restaurant/infra/db"
	restauranthttp "go-ai/internal/restaurant/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type RestaurantModule struct {
	Handler *restauranthttp.RestaurantHandler
	Auth    *middlewares.IdentityMiddleware
	RBAC    rbac.Service
}

func InitRestaurantModule(pool *pgxpool.Pool, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service, log zerolog.Logger) *RestaurantModule {
	repo := restaurantrepo.NewRestaurantRepo(pool)
	createRestaurantUC := restaurantapp.NewCreateRestaurantUseCase(repo)
	updateRestaurantUC := restaurantapp.NewUpdateRestaurantUseCase(repo)
	deleteRestaurantUC := restaurantapp.NewDeleteUseCase(repo)
	getByIdUC := restaurantapp.NewGetByIDUseCase(repo)
	handler := restauranthttp.NewRestaurantHandler(
		createRestaurantUC,
		getByIdUC,
		updateRestaurantUC,
		deleteRestaurantUC,
		log,
	)
	return &RestaurantModule{
		Handler: handler,
		Auth:    auth,
		RBAC:    rbacSvc,
	}
}
