package redis

import (
	"Go-library/cache"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

type RedisConfig struct {
	addr     string
	password string
	DB       int
}

func NewRedisCache(cfc RedisConfig) (*RedisCache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfc.addr,
		Password: cfc.password,
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

// Set adds or updates a value in the cache.
func (c *RedisCache) Set(key string, value interface{}) error {
	if key == "" {
		return cache.ErrEmptyKey
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	ctx := context.Background()

	//yet to implement
	return c.client.Set(ctx, key, data, 0).Err()
}

// SetWithTTL adds or updates a value in the cache with a TTL.
func (c *RedisCache) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	//yet to implement
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

// Get retrieves a value from the cache.
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
	//yet to implement
	return out, nil
}

// Delete removes a key from the cache.
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
	//yet to implement
	return nil
}

// Clear removes all keys from the cache.
func (c *RedisCache) Clear() error {
	//yet to implement
	ctx := context.Background()
	return c.client.FlushDB(ctx).Err()
}
