package container

import (
	"context"
	"fmt"
	"go-ai/internal/platform/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis DB numbers — managed in code, not via environment variables,
// so each module can use a dedicated database.
const (
	RedisDBAuth    = 0 // sessions, refresh tokens, blacklist
	RedisDBDefault = 0
)

// RedisPoolConfig holds Redis connection pool settings
type RedisPoolConfig struct {
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DefaultRedisPoolConfig returns default Redis pool configuration
func DefaultRedisPoolConfig() RedisPoolConfig {
	return RedisPoolConfig{
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// RedisPoolConfigFromConfig creates RedisPoolConfig from application config
func RedisPoolConfigFromConfig(cfg *config.Config) RedisPoolConfig {
	return RedisPoolConfig{
		PoolSize:     cfg.RedisPoolSize,
		MinIdleConns: cfg.RedisMinIdleConns,
		MaxRetries:   cfg.RedisMaxRetries,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

func ConnectRedis(dsn string) (*redis.Client, error) {
	return ConnectRedisWithConfig(dsn, DefaultRedisPoolConfig())
}

func ConnectRedisWithConfig(dsn string, poolCfg RedisPoolConfig) (*redis.Client, error) {
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis dsn: %w", err)
	}

	// Apply connection pool settings
	opt.PoolSize = poolCfg.PoolSize
	opt.MinIdleConns = poolCfg.MinIdleConns
	opt.MaxRetries = poolCfg.MaxRetries
	opt.DialTimeout = poolCfg.DialTimeout
	opt.ReadTimeout = poolCfg.ReadTimeout
	opt.WriteTimeout = poolCfg.WriteTimeout

	redisClient := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), poolCfg.DialTimeout)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return redisClient, nil
}
