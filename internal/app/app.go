package app

import (
	"go-ai/internal/container"
	identityhttp "go-ai/internal/identity/transport/http"
	uploadhttp "go-ai/internal/media/transport/http"
	"go-ai/internal/platform/config"
	restauranthttp "go-ai/internal/restaurant/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func BuildApp(e *echo.Echo, pool *pgxpool.Pool, redis *redis.Client, cfg *config.Config, log zerolog.Logger) {
	api := e.Group("/api")

	initIdentityModule := container.InitIdentityModule(pool, redis, cfg, log)
	identityhttp.RegisterIdentityRoutes(api, initIdentityModule.Handler, initIdentityModule.Middleware)

	restaurantModule := container.InitRestaurantModule(pool, initIdentityModule.Middleware, initIdentityModule.RbacService, log)
	restauranthttp.RegisterRestaurantRoutes(api, restaurantModule.Handler, initIdentityModule.Middleware, initIdentityModule.RbacService)

	mediaModule := container.InitMediaModule(initIdentityModule.Middleware, log)
	uploadhttp.RegisterMediaRoutes(api, mediaModule.Handler, mediaModule.Auth)
}
