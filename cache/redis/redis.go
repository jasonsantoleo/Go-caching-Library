package redis

import (
	"Go-library/cache"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

type RedisConfig struct {
	addr     string
	password string
	db       int
}

// chekf id redis cache can create interface with Cache
var _ cache.Cache = (*RedisCache)(nil)

// redis client setup
func NewRedisCache(cfc RedisConfig) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfc.addr,
		Password: cfc.password,
		DB:       cfc.db,
	})

	return &RedisCache{
		client: rdb,
	}
}

// Set adds or updates a value in the cache.
func (c *RedisCache) Set(key string, value interface{}) error {
	//yet to implement
	return cache.ErrEmptyKey
}

// SetWithTTL adds or updates a value in the cache with a TTL.
func (c *RedisCache) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	//yet to implement
	return cache.ErrEmptyKey
}

// Get retrieves a value from the cache.
func (c *RedisCache) Get(key string) (interface{}, error) {
	//yet to implement
	return nil, nil
}

// Delete removes a key from the cache.
func (c *RedisCache) Delete(key string) error {
	//yet to implement
	return cache.ErrEmptyKey
}

// Clear removes all keys from the cache.
func (c *RedisCache) Clear() error {
	//yet to implement
	return cache.ErrEmptyKey
}
