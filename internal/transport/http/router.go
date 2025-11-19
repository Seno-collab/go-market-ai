package http

import (
	"go-ai/internal/infra/db"
	authservice "go-ai/internal/service/auth"
	"go-ai/internal/transport/http/handler"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func Router(pool *pgxpool.Pool, e *echo.Echo) {
	api := e.Group("/api")
	// ---- API AUTH ----
	authRepo := db.NewAuthRepo(pool)
	authService := authservice.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)
	authGroup := api.Group("/auth")
	authGroup.POST("/register", authHandler.Register)

}
