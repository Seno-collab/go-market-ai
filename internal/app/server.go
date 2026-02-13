package app

import (
	"context"
	"fmt"
	"go-ai/internal/container"
	"go-ai/internal/platform/config"
	"go-ai/internal/transport/middlewares"
	"go-ai/internal/transport/swagger"
	"go-ai/pkg/logger"
	"go-ai/pkg/metrics"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "go-ai/docs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

// ServerConfig holds configurable server settings
type ServerConfig struct {
	RateLimitRequests int
	RateLimitBurst    int
	RequestTimeout    time.Duration
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		RateLimitRequests: 20,
		RateLimitBurst:    50,
		RequestTimeout:    30 * time.Second,
	}
}

// ServerConfigFromConfig creates ServerConfig from application config
func ServerConfigFromConfig(cfg *config.Config) ServerConfig {
	return ServerConfig{
		RateLimitRequests: cfg.RateLimitRequests,
		RateLimitBurst:    cfg.RateLimitBurst,
		RequestTimeout:    time.Duration(cfg.RequestTimeout) * time.Second,
	}
}

func NewServer() *echo.Echo {
	return NewServerWithConfig(DefaultServerConfig())
}

func NewServerWithConfig(srvCfg ServerConfig) *echo.Echo {
	e := echo.New()
	log := logger.NewLogger()

	httpInFlight := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: metrics.Namespace,
		Name:      "http_requests_in_flight",
		Help:      "Number of HTTP requests currently being processed",
	})

	promMiddleware := echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Namespace: metrics.Namespace,
		Subsystem: "http",
		Skipper: func(c *echo.Context) bool {
			path := c.Request().URL.Path
			return path == "/metrics" || path == "/api/auth/login" || strings.HasPrefix(path, "/swagger")
		},
		DoNotUseRequestPathFor404: true,
		BeforeNext: func(c *echo.Context) {
			httpInFlight.Inc()
		},
		AfterNext: func(c *echo.Context, err error) {
			httpInFlight.Dec()
		},
		HistogramOptsFunc: func(opts prometheus.HistogramOpts) prometheus.HistogramOpts {
			opts.Buckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5}
			return opts
		},
	})
	e.Use(promMiddleware)
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

	// Configurable rate limiter with burst support
	rateLimiterStore := middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate:      float64(srvCfg.RateLimitRequests),
			Burst:     srvCfg.RateLimitBurst,
			ExpiresIn: time.Minute,
		},
	)
	e.Use(middleware.RateLimiter(rateLimiterStore))

	e.Use(middlewares.RequestIDMiddleware(log))
	e.Use(middleware.ContextTimeout(srvCfg.RequestTimeout))
	e.Use(middleware.Recover())
	e.Use(middlewares.RequestLoggerMiddleware(log))
	e.Use(middlewares.ResponseLoggerMiddleware())

	// Swagger UI
	e.GET("/swagger/*", swagger.Handler(nil))

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
		Int("db_max_conns", cfg.DBMaxConns).
		Int("redis_pool_size", cfg.RedisPoolSize).
		Int("rate_limit", cfg.RateLimitRequests).
		Msg("Configuration loaded")

	dsnPg := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	// Connect to PostgreSQL with optimized pool settings
	pgPoolCfg := container.PoolConfigFromConfig(cfg)
	pool, err := container.ConnectPostgresWithConfig(dsnPg, pgPoolCfg)
	if err != nil {
		return fmt.Errorf("connect postgres failed: %w", err)
	}

	dsnRedis := fmt.Sprintf("redis://%s:%s@%s:%d/%d",
		"", cfg.RedisPassword, cfg.RedisHost, cfg.RedisPort, container.RedisDBAuth,
	)

	// Connect to Redis with optimized pool settings
	redisPoolCfg := container.RedisPoolConfigFromConfig(cfg)
	redisClient, err := container.ConnectRedisWithConfig(dsnRedis, redisPoolCfg)
	if err != nil {
		return fmt.Errorf("connect redis failed: %w", err)
	}

	BuildApp(e, pool, redisClient, cfg, log)

	// Register pool collectors for on-demand Prometheus scraping
	metrics.RegisterPoolCollectors(pool, redisClient)

	chServer := make(chan error, 1)
	serverAddr := ":" + cfg.ServerPort
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      e,
		ReadTimeout:  time.Duration(cfg.RequestTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.RequestTimeout) * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Str("addr", serverAddr).Msg("Starting HTTP server")
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

	// Graceful shutdown with configurable timeout
	shutdownTimeout := time.Duration(cfg.ShutdownTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown HTTP server first
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}

	// Clean up database connections
	log.Info().Msg("Closing database connections...")
	pool.Close()

	// Clean up Redis connections
	log.Info().Msg("Closing Redis connections...")
	if err := redisClient.Close(); err != nil {
		log.Error().Err(err).Msg("Redis close error")
	}

	log.Info().Msg("Server exited gracefully")
	return nil
}
