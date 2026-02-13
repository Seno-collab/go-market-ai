package cache

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
	"time"
)

var errRawStringType = errors.New("cache: RawString mode requires Cache[string]")

// Options configures a Cache instance.
type Options struct {
	CacheType  string        // Label (e.g. "session", "refresh_token")
	KeyPrefix  string        // Prefix prepended to every key
	DefaultTTL time.Duration // Default TTL for Set when 0 is passed
	RawString  bool          // Skip JSON marshal/unmarshal (for Cache[string])
}

// Cache is a generic Redis cache.
type Cache[T any] struct {
	client *redis.Client
	opts   Options
}

// New creates a new Cache with the given Redis client and options.
func New[T any](client *redis.Client, opts Options) *Cache[T] {
	return &Cache[T]{client: client, opts: opts}
}

// Get retrieves a value by key. Returns (nil, nil) on cache miss.
func (c *Cache[T]) Get(ctx context.Context, key string) (*T, error) {
	fullKey := c.opts.KeyPrefix + key

	val, err := c.client.Get(ctx, fullKey).Result()

	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var result T
	if c.opts.RawString {
		v, ok := any(&result).(*string)
		if !ok {
			return nil, errRawStringType
		}
		*v = val
	} else {
		if err := json.Unmarshal([]byte(val), &result); err != nil {
			return nil, err
		}
	}
	return &result, nil
}

// Set stores a value with the given TTL. If ttl is 0, DefaultTTL is used.
func (c *Cache[T]) Set(ctx context.Context, key string, value *T, ttl time.Duration) error {
	fullKey := c.opts.KeyPrefix + key

	if ttl == 0 {
		ttl = c.opts.DefaultTTL
	}

	var data string
	if c.opts.RawString {
		v, ok := any(value).(*string)
		if !ok {
			return errRawStringType
		}
		data = *v
	} else {
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}
		data = string(b)
	}

	return c.client.Set(ctx, fullKey, data, ttl).Err()
}

// Delete removes a key from the cache.
func (c *Cache[T]) Delete(ctx context.Context, key string) error {
	fullKey := c.opts.KeyPrefix + key
	return c.client.Del(ctx, fullKey).Err()
}
