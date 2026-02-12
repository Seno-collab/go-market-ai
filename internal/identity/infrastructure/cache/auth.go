package cache

import (
	"context"
	"encoding/json"
	pkgcache "go-ai/pkg/cache"
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
	client       *redis.Client
	session      *pkgcache.Cache[UserCache]
	refreshToken *pkgcache.Cache[string]
	blacklist    *pkgcache.Cache[string]
}

func NewAuthCache(client *redis.Client) *AuthCache {
	return &AuthCache{
		client: client,
		session: pkgcache.New[UserCache](client, pkgcache.Options{
			CacheType: "session",
		}),
		refreshToken: pkgcache.New[string](client, pkgcache.Options{
			CacheType: "refresh_token",
			RawString: true,
		}),
		blacklist: pkgcache.New[string](client, pkgcache.Options{
			CacheType: "blacklist",
			RawString: true,
		}),
	}
}

func (a *AuthCache) GetAuthCache(ctx context.Context, key string) (*UserCache, error) {
	return a.session.Get(ctx, key)
}

func (a *AuthCache) SetAuthCache(ctx context.Context, key string, value *UserCache, ttl time.Duration) error {
	return a.session.Set(ctx, key, value, ttl)
}

// SetLoginCaches writes session and refresh token in a single round-trip.
func (a *AuthCache) SetLoginCaches(
	ctx context.Context,
	sessionKey string,
	sessionValue *UserCache,
	sessionTTL time.Duration,
	refreshKey string,
	refreshToken string,
	refreshTTL time.Duration,
) error {
	data, err := json.Marshal(sessionValue)
	if err != nil {
		return err
	}

	start := time.Now()
	pipe := a.client.Pipeline()
	pipe.Set(ctx, sessionKey, data, sessionTTL)
	pipe.Set(ctx, refreshKey, refreshToken, refreshTTL)
	_, err = pipe.Exec(ctx)
	metrics.RecordCacheOperation("set_batch", "auth", time.Since(start).Seconds())
	return err
}

func (a *AuthCache) DeleteAuthCache(ctx context.Context, key string) error {
	return a.session.Delete(ctx, key)
}

func (a *AuthCache) GetRefreshTokenCache(ctx context.Context, key string) (string, error) {
	val, err := a.refreshToken.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if val == nil {
		return "", nil
	}
	return *val, nil
}

func (a *AuthCache) SetRefreshTokenCache(ctx context.Context, key string, refresh string, ttl time.Duration) error {
	return a.refreshToken.Set(ctx, key, &refresh, ttl)
}

func (a *AuthCache) DeleteRefreshTokenCache(ctx context.Context, key string) error {
	return a.refreshToken.Delete(ctx, key)
}

func (a *AuthCache) BlacklistToken(ctx context.Context, key string, ttl time.Duration) error {
	val := "1"
	return a.blacklist.Set(ctx, key, &val, ttl)
}

func (a *AuthCache) IsTokenBlacklisted(ctx context.Context, key string) (bool, error) {
	val, err := a.blacklist.Get(ctx, key)
	if err != nil {
		return false, err
	}
	return val != nil, nil
}
