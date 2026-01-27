package app

import (
	"context"
	"fmt"
	"go-ai/internal/container"
	"go-ai/internal/platform/config"
	"go-ai/internal/transport/middlewares"
	"go-ai/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type AppServer struct {
	Echo   *echo.Echo
	DB     *pgxpool.Pool
	Redis  *redis.Client
	Config *config.Config
	Logger zerolog.Logger
}

func NewServer() *echo.Echo {
	e := echo.New()
	logger := logger.NewLogger()

	e.Use(echoprometheus.NewMiddleware("echo"))
	e.GET("/metrics", echoprometheus.NewHandler())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://157.66.218.138:3000"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middlewares.RequestIDMiddleware(logger))
	e.Use(middleware.ContextTimeout(30 * time.Second))
	e.Use(middleware.Recover())
	e.Use(middlewares.RequestLoggerMiddleware(logger))
	e.Use(middlewares.ResponseLoggerMiddleware())

	return e
}

func Run(e *echo.Echo) error {

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config failed: %w", err)
	}

	log := logger.NewLogger()
	log.Info().
		Str("environment", cfg.Environment).
		Str("server_port", cfg.ServerPort).
		Str("db_host", cfg.DBHost).
		Msg("Configuration loaded")

	dsnPg := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	pool, err := container.ConnectPostgres(dsnPg)
	if err != nil {
		return fmt.Errorf("connect postgres failed: %w", err)
	}

	dsnRedis := fmt.Sprintf("redis://%s:%s@%s:%d/%d",
		"", cfg.RedisPassword, cfg.RedisHost, cfg.RedisPort, cfg.RedisDB,
	)
	redisClient, err := container.ConnectRedis(dsnRedis)
	if err != nil {
		return fmt.Errorf("connect redis failed: %w", err)
	}

	BuildApp(e, pool, redisClient, cfg, log)
	chServer := make(chan error, 1)
	serverAddr := ":" + cfg.ServerPort
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: e,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server startup failed")
			chServer <- err
		}
		close(chServer)
	}()
	ctxSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	select {
	case <-ctxSignal.Done():
		log.Info().Str("signal", ctxSignal.Err().Error()).Msg("Server shutting down gracefully")
	case err := <-chServer:
		log.Error().Err(err).Msg("Server encountered an error")
		return err
	}
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Forced shutdown")
	}

	log.Info().Msg("Server exited gracefully")
	return nil
}
