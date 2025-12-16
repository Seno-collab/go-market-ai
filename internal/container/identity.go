package container

import (
	authapp "go-ai/internal/identity/application/auth"
	"go-ai/internal/identity/domain/rbac"
	"go-ai/internal/identity/infrastructure/cache"
	"go-ai/internal/identity/infrastructure/db"
	identityhttp "go-ai/internal/identity/transport/http"
	middlewares "go-ai/internal/identity/transport/middlewares"
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

	registerUseCase := authapp.NewRegisterUseCase(authRepo, authCache)
	loginUseCase := authapp.NewLoginUseCase(authRepo, authCache, config)
	refreshUseCase := authapp.NewRefreshTokenUseCase(authRepo, authCache, config)
	profileUseCase := authapp.NewGetProfileUseCase(authRepo, authCache)
	changePasswordUseCase := authapp.NewChangePasswordUseCase(authRepo)
	logoutUseCase := authapp.NewLogoutUseCase(authRepo, authCache)
	handler := identityhttp.NewAuthHandler(
		registerUseCase,
		loginUseCase,
		refreshUseCase,
		profileUseCase,
		changePasswordUseCase,
		logoutUseCase,
		log,
	)

	middleware := middlewares.NewIdentityMiddleware(authCache, config)

	return &IdentityModule{
		Handler:     handler,
		Middleware:  middleware,
		RbacService: rbacService,
	}
}
