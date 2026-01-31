package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	restaurantCachePrefix = "restaurant:"
	restaurantCacheTTL    = 5 * time.Minute
)

// CachedRestaurant represents restaurant data stored in cache
type CachedRestaurant struct {
	ID          int32             `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Address     string            `json:"address"`
	Category    string            `json:"category"`
	City        string            `json:"city"`
	District    string            `json:"district"`
	LogoUrl     string            `json:"logo_url"`
	BannerUrl   string            `json:"banner_url"`
	PhoneNumber string            `json:"phone_number"`
	WebsiteUrl  string            `json:"website_url"`
	Email       string            `json:"email"`
	Hours       []CachedHours     `json:"hours"`
	CachedAt    time.Time         `json:"cached_at"`
}

// CachedHours represents operating hours stored in cache
type CachedHours struct {
	Day       string `json:"day"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
}

// RestaurantCache provides caching for restaurant data
type RestaurantCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRestaurantCache creates a new RestaurantCache instance
func NewRestaurantCache(client *redis.Client) *RestaurantCache {
	return &RestaurantCache{
		client: client,
		ttl:    restaurantCacheTTL,
	}
}

// NewRestaurantCacheWithTTL creates a new RestaurantCache with custom TTL
func NewRestaurantCacheWithTTL(client *redis.Client, ttl time.Duration) *RestaurantCache {
	return &RestaurantCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *RestaurantCache) key(id int32) string {
	return fmt.Sprintf("%s%d", restaurantCachePrefix, id)
}

// Get retrieves restaurant from cache
func (c *RestaurantCache) Get(ctx context.Context, id int32) (*CachedRestaurant, error) {
	data, err := c.client.Get(ctx, c.key(id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("cache get error: %w", err)
	}

	var cached CachedRestaurant
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, fmt.Errorf("cache unmarshal error: %w", err)
	}

	return &cached, nil
}

// Set stores restaurant in cache
func (c *RestaurantCache) Set(ctx context.Context, restaurant *CachedRestaurant) error {
	restaurant.CachedAt = time.Now()
	data, err := json.Marshal(restaurant)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}

	if err := c.client.Set(ctx, c.key(restaurant.ID), data, c.ttl).Err(); err != nil {
		return fmt.Errorf("cache set error: %w", err)
	}

	return nil
}

// Delete removes restaurant from cache
func (c *RestaurantCache) Delete(ctx context.Context, id int32) error {
	if err := c.client.Del(ctx, c.key(id)).Err(); err != nil {
		return fmt.Errorf("cache delete error: %w", err)
	}
	return nil
}

// DeletePattern removes all restaurants matching a pattern
func (c *RestaurantCache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, restaurantCachePrefix+pattern, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("cache scan error: %w", err)
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("cache delete pattern error: %w", err)
		}
	}

	return nil
}

// Invalidate clears all restaurant cache
func (c *RestaurantCache) Invalidate(ctx context.Context) error {
	return c.DeletePattern(ctx, "*")
}
