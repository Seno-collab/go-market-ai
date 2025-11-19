package main

import (
	"context"
	"fmt"
	_ "go-ai/docs"
	"go-ai/internal/config"
	"go-ai/internal/infra/cache"
	"go-ai/internal/infra/db"
	httpHandler "go-ai/internal/transport/http"
	"go-ai/pkg/common"
	"go-ai/pkg/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// go ai
// mission using golang build ai. Development application AI

func main() {
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	common.NewLogger()
	cfg, err := config.LoadConfig()
	if err != nil {
		common.Logger.Fatal().Fields(map[string]interface{}{
			"Error": err,
		}).Msg("Load config error")
		return
	}
	dsnPg := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	pool, err := db.ConnDbPgPool(dsnPg)
	if err != nil {
		common.Logger.Fatal().Fields(map[string]interface{}{
			"Error": err,
		}).Msg("Connect db fail")
		return
	}
	httpHandler.Router(pool, e)
	dsnRedis := fmt.Sprintf("redis://%s:%s@%s:%d/%d", "", cfg.RedisPassword, cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	_, err = cache.ConnectRedis(dsnRedis)
	if err != nil {
		common.Logger.Fatal().Fields(map[string]interface{}{
			"Error": err,
		}).Msg("Connect redis fail")
	}

	e.Use(middleware.LoggingMiddleware)
	port := fmt.Sprintf(":%s", cfg.ServerPort)
	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			common.Logger.Fatal().Fields(map[string]interface{}{
				"Error": err,
			}).Msg("Shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	common.Logger.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		common.Logger.Fatal().Fields(map[string]interface{}{
			"Error": err,
		}).Msg("Server forced to shutdown")
	}
	common.Logger.Info().Msg("Server exited gracefully")
}
