package container

import (
	"context"
	"fmt"
	"go-ai/internal/platform/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	postgresMaxRetries     = 8
	postgresInitialBackoff = time.Second
	postgresDialTimeout    = 15 * time.Second
)

// PostgresPoolConfig holds connection pool settings
type PostgresPoolConfig struct {
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

// DefaultPostgresPoolConfig returns default pool configuration
func DefaultPostgresPoolConfig() PostgresPoolConfig {
	return PostgresPoolConfig{
		MaxConns:          25,
		MinConns:          5,
		MaxConnLifetime:   time.Hour,
		MaxConnIdleTime:   5 * time.Minute,
		HealthCheckPeriod: 30 * time.Second,
	}
}

// PoolConfigFromConfig creates PostgresPoolConfig from application config
func PoolConfigFromConfig(cfg *config.Config) PostgresPoolConfig {
	return PostgresPoolConfig{
		MaxConns:          int32(cfg.DBMaxConns),
		MinConns:          int32(cfg.DBMinConns),
		MaxConnLifetime:   time.Duration(cfg.DBMaxConnLifetime) * time.Second,
		MaxConnIdleTime:   time.Duration(cfg.DBMaxConnIdleTime) * time.Second,
		HealthCheckPeriod: time.Duration(cfg.DBHealthCheckPeriod) * time.Second,
	}
}

func ConnectPostgres(dsn string) (*pgxpool.Pool, error) {
	return ConnectPostgresWithConfig(dsn, DefaultPostgresPoolConfig())
}

func ConnectPostgresWithConfig(dsn string, poolCfg PostgresPoolConfig) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}

	// Apply connection pool settings
	pgxConfig.MaxConns = poolCfg.MaxConns
	pgxConfig.MinConns = poolCfg.MinConns
	pgxConfig.MaxConnLifetime = poolCfg.MaxConnLifetime
	pgxConfig.MaxConnIdleTime = poolCfg.MaxConnIdleTime
	pgxConfig.HealthCheckPeriod = poolCfg.HealthCheckPeriod

	backoff := postgresInitialBackoff
	for attempt := 1; attempt <= postgresMaxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), postgresDialTimeout)

		pool, err = pgxpool.NewWithConfig(ctx, pgxConfig)
		if err == nil {
			if pingErr := pool.Ping(ctx); pingErr == nil {
				cancel()
				return pool, nil
			} else {
				err = pingErr
			}
			pool.Close()
		}
		cancel()

		if attempt == postgresMaxRetries {
			break
		}

		time.Sleep(backoff)
		backoff *= 2
	}

	return nil, fmt.Errorf("cannot connect to postgres after %d attempts: %w", postgresMaxRetries, err)
}
