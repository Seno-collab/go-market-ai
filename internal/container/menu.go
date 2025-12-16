package container

import (
	"go-ai/internal/identity/domain/rbac"
	middlewares "go-ai/internal/identity/transport/middlewares"
	menuitemapp "go-ai/internal/menu/application/menu_item"
	"go-ai/internal/menu/infrastructure/db"
	menuitemhttp "go-ai/internal/menu/transport/http"
	middlewaresRestaurant "go-ai/internal/restaurant/transport/middlewares"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type MenuModule struct {
	MenuHandler        *menuitemhttp.MenuItemHandler
	Auth               *middlewares.IdentityMiddleware
	RBAC               rbac.Service
	RestaurantRequired *middlewaresRestaurant.RestaurantRequired
}

func InitMenuModule(pool *pgxpool.Pool, auth *middlewares.IdentityMiddleware, rbacSvc rbac.Service, restaurantRequired *middlewaresRestaurant.RestaurantRequired, log zerolog.Logger) *MenuModule {
	repo := db.NewMenuRepo(pool)
	creatMenuHandler := menuitemapp.NewCreateUseCase(repo)
	menuHandler := menuitemhttp.NewMenuItemHandler(creatMenuHandler, log)
	return &MenuModule{
		MenuHandler:        menuHandler,
		Auth:               auth,
		RBAC:               rbacSvc,
		RestaurantRequired: restaurantRequired,
	}
}
