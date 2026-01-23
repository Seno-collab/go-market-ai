package container

import (
	healthapp "go-ai/internal/health/application/health"
	healthhttp "go-ai/internal/health/transport/http"
	"go-ai/internal/platform/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type HealthModule struct {
	Handler *healthhttp.HealthHandler
}

func InitHealthModule(pool *pgxpool.Pool, redis *redis.Client, cfg *config.Config, log zerolog.Logger) *HealthModule {
	checker := healthapp.NewCheckHealthUseCase(pool, redis, cfg.Environment)
	handler := healthhttp.NewHealthHandler(checker, log)
	return &HealthModule{
		Handler: handler,
	}
}
