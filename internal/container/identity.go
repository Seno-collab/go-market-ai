package container

import (
	authapp "go-ai/internal/identity/application/auth"
	"go-ai/internal/identity/domain/rbac"
	"go-ai/internal/identity/infra/cache"
	"go-ai/internal/identity/infra/db"
	identityhttp "go-ai/internal/identity/transport/http"
	middlewares "go-ai/internal/identity/transport/middelwares"
	"go-ai/internal/platform/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type IdentityModule struct {
	Handler     *identityhttp.AuthHandler
	Middleware  *middlewares.IdentityMiddleware
	RbacService rbac.Service
}

func InitIdentityModule(pool *pgxpool.Pool, redis *redis.Client, config *config.Config, log zerolog.Logger) *IdentityModule {

	rbacRepo := db.NewRbacRepo(pool)
	rbacService := rbac.Service{Repo: rbacRepo}

	authRepo := db.NewAuthRepo(pool)
	authCache := cache.NewAuthCache(redis)

	registerUC := authapp.NewRegisterUseCase(authRepo, authCache)
	loginUC := authapp.NewLoginUseCase(authRepo, authCache, config)
	refreshUC := authapp.NewRefreshTokenUseCase(authRepo, authCache, config)
	profileUC := authapp.NewGetProfileUseCase(authRepo, authCache)

	handler := identityhttp.NewAuthHandler(
		registerUC,
		loginUC,
		refreshUC,
		profileUC,
		log,
	)

	middleware := middlewares.NewIdentityMiddleware(authCache, config)

	return &IdentityModule{
		Handler:     handler,
		Middleware:  middleware,
		RbacService: rbacService,
	}
}
