package container

import (
	"go-ai/internal/identity/domain/rbac"
	middlewares "go-ai/internal/identity/transport/middlewares"
	restaurantapp "go-ai/internal/restaurant/application/restaurant"
	restaurantrepo "go-ai/internal/restaurant/infrastructure/db"
	restauranthttp "go-ai/internal/restaurant/transport/http"
	middlewaresRestaurant "go-ai/internal/restaurant/transport/middlewares"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type RestaurantModule struct {
	Handler               *restauranthttp.RestaurantHandler
	Auth                  *middlewares.IdentityMiddleware
	RBAC                  rbac.Service
	MiddlewaresRestaurant *middlewaresRestaurant.RestaurantRequired
}

func InitRestaurantModule(pool *pgxpool.Pool, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service, log zerolog.Logger) *RestaurantModule {
	repo := restaurantrepo.NewRestaurantRepo(pool)

	middlewaresRestaurant := middlewaresRestaurant.New(repo)
	createRestaurantUseCase := restaurantapp.NewCreateRestaurantUseCase(repo)
	updateRestaurantUseCase := restaurantapp.NewUpdateRestaurantUseCase(repo)
	deleteRestaurantUseCase := restaurantapp.NewDeleteUseCase(repo)
	getByIDUseCase := restaurantapp.NewGetByIDUseCase(repo)
	getComboboxRestaurant := restaurantapp.NewGetRestaurantItemComboboxUseCase(repo)
	handler := restauranthttp.NewRestaurantHandler(
		createRestaurantUseCase,
		getByIDUseCase,
		updateRestaurantUseCase,
		deleteRestaurantUseCase,
		getComboboxRestaurant,
		log,
	)
	return &RestaurantModule{
		Handler:               handler,
		Auth:                  auth,
		RBAC:                  rbacSvc,
		MiddlewaresRestaurant: middlewaresRestaurant,
	}
}
