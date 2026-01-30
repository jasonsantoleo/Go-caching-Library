package redis

import (
	"Go-library/cache"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// Compile-time check to ensure RedisCache implements cache.Cache
var _ cache.Cache = (*RedisCache)(nil)

// NewRedisCache creates a new RedisCache instance.
func NewRedisCache(addr string, password string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Set adds or updates a value in the cache.
func (c *RedisCache) Set(key string, value interface{}) error {
	if key == "" {
		return cache.ErrEmptyKey
	}
	return c.client.Set(c.ctx, key, value, 0).Err()
}

// SetWithTTL adds or updates a value in the cache with a TTL.
func (c *RedisCache) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	if key == "" {
		return cache.ErrEmptyKey
	}
	return c.client.Set(c.ctx, key, value, ttl).Err()
}

// Get retrieves a value from the cache.
func (c *RedisCache) Get(key string) (interface{}, error) {
	if key == "" {
		return nil, cache.ErrEmptyKey
	}
	val, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return nil, cache.ErrKeyNotFound
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Delete removes a key from the cache.
func (c *RedisCache) Delete(key string) error {
	if key == "" {
		return cache.ErrEmptyKey
	}
	n, err := c.client.Del(c.ctx, key).Result()
	if err != nil {
		return err
	}
	if n == 0 {
		return cache.ErrKeyNotFound
	}
	return nil
}

// Clear removes all keys from the cache.
func (c *RedisCache) Clear() error {
	return c.client.FlushDB(c.ctx).Err()
}
