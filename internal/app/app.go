package app

import (
	"go-ai/internal/container"
	healthhttp "go-ai/internal/health/transport/http"
	identityhttp "go-ai/internal/identity/transport/http"
	uploadhttp "go-ai/internal/media/transport/http"
	menuhttp "go-ai/internal/menu/transport/http"
	"go-ai/internal/platform/config"
	restauranthttp "go-ai/internal/restaurant/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func BuildApp(e *echo.Echo, pool *pgxpool.Pool, redis *redis.Client, cfg *config.Config, log zerolog.Logger) {
	api := e.Group("/api")

	healthModule := container.InitHealthModule(pool, redis, cfg, log)
	healthhttp.RegisterHealthRoutes(api, healthModule.Handler)

	initIdentityModule := container.InitIdentityModule(pool, redis, cfg, log)
	identityhttp.RegisterIdentityRoutes(api, initIdentityModule.Handler, initIdentityModule.Middleware)

	restaurantModule := container.InitRestaurantModule(pool, initIdentityModule.Middleware, initIdentityModule.RbacService, log)
	restauranthttp.RegisterRestaurantRoutes(api, restaurantModule.Handler, initIdentityModule.Middleware, initIdentityModule.RbacService)

	mediaModule, err := container.InitMediaModule(initIdentityModule.Middleware, cfg, log)
	if err != nil {
		log.Error().Err(err).Msg("failed to init media module")
	} else {
		uploadhttp.RegisterMediaRoutes(api, mediaModule.Handler, mediaModule.Auth)
	}

	menuModule := container.InitMenuItemModule(pool, log)
	menuhttp.RegisterMenuItemRoutes(api, menuModule.MenuItemHandler, initIdentityModule.Middleware, initIdentityModule.RbacService)

	menusModule := container.InitMenuModule(pool, log)
	menuhttp.RegisterMenusRoutes(api, menusModule.Handler)

	topicModule := container.InitTopicModule(pool, log)
	menuhttp.RegisterTopicRoutes(api, topicModule.TopicHandler, initIdentityModule.Middleware, initIdentityModule.RbacService)

	optionModule := container.InitOptionModule(pool, log)
	menuhttp.RegisterOptionGroupRoutes(api, optionModule.OptionGroupHandler, initIdentityModule.Middleware, initIdentityModule.RbacService)
	menuhttp.RegisterOptionItemRoutes(api, optionModule.OptionItemHandler, initIdentityModule.Middleware, initIdentityModule.RbacService)
}
