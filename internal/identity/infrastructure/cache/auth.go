package cache

import (
	"context"
	"encoding/json"
	"go-ai/pkg/metrics"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	UserID   uuid.UUID
	Email    string
	FullName string
	Role     string
	IsActive bool
	ImageUrl string
}

type AuthCache struct {
	Redis *redis.Client
}

func NewAuthCache(redis *redis.Client) *AuthCache {
	return &AuthCache{
		Redis: redis,
	}
}

func (authCache *AuthCache) GetAuthCache(ctx context.Context, key string) (*UserCache, error) {
	start := time.Now()
	val, err := authCache.Redis.Get(ctx, key).Result()
	metrics.RecordCacheOperation("get", "session", time.Since(start).Seconds())
	if err == redis.Nil {
		metrics.RecordCacheMiss("session")
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	metrics.RecordCacheHit("session")
	authData := &UserCache{}
	if err := json.Unmarshal([]byte(val), authData); err != nil {
		return nil, err
	}
	return authData, nil
}

func (authCache *AuthCache) SetAuthCache(ctx context.Context, key string, value *UserCache, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	start := time.Now()
	err = authCache.Redis.Set(ctx, key, string(b), ttl).Err()
	metrics.RecordCacheOperation("set", "session", time.Since(start).Seconds())
	return err
}

func (authCache *AuthCache) SetRefreshTokenCache(ctx context.Context, key string, refresh string, ttl time.Duration) error {
	start := time.Now()
	err := authCache.Redis.Set(ctx, key, refresh, ttl).Err()
	metrics.RecordCacheOperation("set", "refresh_token", time.Since(start).Seconds())
	return err
}

func (authCache *AuthCache) GetRefreshTokenCache(ctx context.Context, key string) (string, error) {
	start := time.Now()
	value, err := authCache.Redis.Get(ctx, key).Result()
	metrics.RecordCacheOperation("get", "refresh_token", time.Since(start).Seconds())
	if err == redis.Nil {
		metrics.RecordCacheMiss("refresh_token")
		return "", nil
	}
	if err != nil {
		return "", err
	}
	metrics.RecordCacheHit("refresh_token")
	return value, nil
}

func (authCache *AuthCache) DeleteAuthCache(ctx context.Context, key string) error {
	start := time.Now()
	err := authCache.Redis.Del(ctx, key).Err()
	metrics.RecordCacheOperation("del", "session", time.Since(start).Seconds())
	return err
}

func (authCache *AuthCache) DeleteRefreshTokenCache(ctx context.Context, key string) error {
	start := time.Now()
	err := authCache.Redis.Del(ctx, key).Err()
	metrics.RecordCacheOperation("del", "refresh_token", time.Since(start).Seconds())
	return err
}
