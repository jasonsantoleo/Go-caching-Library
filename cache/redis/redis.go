package redis

import (
	"Go-library/cache"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// redis sturcture for cache.cache
type RedisCache struct {
	client *redis.Client
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// constructor for redisCache
func NewRedisCache(cfc RedisConfig) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfc.Addr,
		Password: cfc.Password,
		DB:       cfc.DB,
	})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: rdb,
	}, nil
}

// chekf id redis cache can create interface with Cache
var _ cache.Cache = (*RedisCache)(nil)

// redis client setup

// adds or updates a value in the cache.
func (c *RedisCache) Set(key string, value interface{}) error {
	if key == "" {
		return cache.ErrEmptyKey
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx := context.Background()

	// implementation
	c.client.Set(ctx, key, data, 0).Err()
	return nil
}

// adds or updates a value in the cache with a TTL.
func (c *RedisCache) SetWithTTL(key string, value interface{}, ttl time.Duration) error {

	if key == "" {
		return cache.ErrEmptyKey
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx := context.Background()
	return c.client.Set(ctx, key, data, ttl).Err()
}

// retrieves a value from the cache.
func (c *RedisCache) Get(key string) (interface{}, error) {
	if key == "" {
		return nil, cache.ErrEmptyKey
	}
	ctx := context.Background()
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, cache.ErrKeyNotFound
	}
	if err != nil {
		return nil, err
	}
	var out interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// removes a key from the cache.
func (c *RedisCache) Delete(key string) error {
	if key == "" {
		return cache.ErrEmptyKey
	}
	ctx := context.Background()
	isDeleted, err := c.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	if isDeleted == 0 {
		return cache.ErrKeyNotFound
	}
	return nil
}

// removes all keys from the cache.
func (c *RedisCache) Clear() error {
	ctx := context.Background()
	return c.client.FlushDB(ctx).Err()
}
