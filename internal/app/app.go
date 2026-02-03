package app

import (
	"go-ai/internal/container"
	healthhttp "go-ai/internal/health/transport/http"
	identityhttp "go-ai/internal/identity/transport/http"
	uploadhttp "go-ai/internal/media/transport/http"
	"go-ai/internal/platform/config"

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

	mediaModule, err := container.InitMediaModule(initIdentityModule.Middleware, cfg, log)
	if err != nil {
		log.Error().Err(err).Msg("failed to init media module")
	} else {
		uploadhttp.RegisterMediaRoutes(api, mediaModule.Handler, mediaModule.Auth)
	}
}
