package metrics

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// PoolMetricsCollector collects metrics from database and cache pools
type PoolMetricsCollector struct {
	pgPool      *pgxpool.Pool
	redisClient *redis.Client
	logger      zerolog.Logger
	stopChan    chan struct{}
}

// NewPoolMetricsCollector creates a new pool metrics collector
func NewPoolMetricsCollector(pgPool *pgxpool.Pool, redisClient *redis.Client, logger zerolog.Logger) *PoolMetricsCollector {
	return &PoolMetricsCollector{
		pgPool:      pgPool,
		redisClient: redisClient,
		logger:      logger.With().Str("component", "metrics_collector").Logger(),
		stopChan:    make(chan struct{}),
	}
}

// Start begins collecting metrics at regular intervals
func (c *PoolMetricsCollector) Start(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Collect immediately on start
		c.collect()

		for {
			select {
			case <-ticker.C:
				c.collect()
			case <-c.stopChan:
				c.logger.Info().Msg("Metrics collector stopped")
				return
			}
		}
	}()

	c.logger.Info().Dur("interval", interval).Msg("Metrics collector started")
}

// Stop stops the metrics collector
func (c *PoolMetricsCollector) Stop() {
	close(c.stopChan)
}

func (c *PoolMetricsCollector) collect() {
	// Collect PostgreSQL pool stats
	if c.pgPool != nil {
		stats := c.pgPool.Stat()
		DBConnectionsActive.Set(float64(stats.AcquiredConns()))
		DBConnectionsIdle.Set(float64(stats.IdleConns()))
	}

	// Collect Redis pool stats
	if c.redisClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		stats := c.redisClient.PoolStats()

		// Update gauge metrics for Redis
		// Note: Redis pool stats are more limited
		_ = stats // Stats available: Hits, Misses, Timeouts, TotalConns, IdleConns, StaleConns

		// Ping Redis to check health
		if err := c.redisClient.Ping(ctx).Err(); err != nil {
			c.logger.Warn().Err(err).Msg("Redis health check failed")
		}
	}
}
