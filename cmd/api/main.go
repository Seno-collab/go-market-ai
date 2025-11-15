package main

import (
	"context"
	"go-ai/pkg/common"
	"go-ai/pkg/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

// go ai
// mission using golang build ai. Development application AI

func main() {
	e := echo.New()

	common.NewLogger()
	e.Use(middleware.LoggingMiddleware)
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
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
