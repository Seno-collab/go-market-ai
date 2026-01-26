package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthData struct {
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

func (authCache *AuthCache) GetAuthCache(ctx context.Context, key string) (*AuthData, error) {
	val, err := authCache.Redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	authData := &AuthData{}
	if err := json.Unmarshal([]byte(val), authData); err != nil {
		return nil, err
	}
	return authData, nil
}

func (authCache *AuthCache) SetAuthCache(ctx context.Context, key string, value *AuthData, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return authCache.Redis.Set(ctx, key, string(b), ttl).Err()
}

func (authCache *AuthCache) SetRefreshTokenCache(ctx context.Context, key string, refresh string, ttl time.Duration) error {
	return authCache.Redis.Set(ctx, key, refresh, ttl).Err()
}

func (authCache *AuthCache) GetRefreshTokenCache(ctx context.Context, key string) (string, error) {
	value, err := authCache.Redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

func (authCache *AuthCache) DeleteAuthCache(ctx context.Context, key string) error {
	return authCache.Redis.Del(ctx, key).Err()
}

func (authCache *AuthCache) DeleteRefreshTokenCache(ctx context.Context, key string) error {
	return authCache.Redis.Del(ctx, key).Err()
}
