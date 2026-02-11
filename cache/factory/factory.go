package factory

import (
	"Go-library/cache"
	"Go-library/cache/cache/memcached"
	"Go-library/cache/cache/memory"
	"Go-library/cache/cache/redis"
	"errors"

	gormemcache "github.com/bradfitz/gomemcache/memcache"
)

// New creates a new Cache instance based on the provided type and configuration.
func New(t BackendType, cfg Config) (cache.Cache, error) {
	switch t {
	case Memory:
		c := memory.NewMemorycache()
		if cfg.MemoryMaxSize > 0 {
			c.SetMaxSize(cfg.MemoryMaxSize)
		}
		return c, nil

	case Redis:
		if cfg.RedisAddr == "" {
			return nil, errors.New("redis address is required")
		}
		// Redis package expects its own RedisConfig struct
		rConfig := redis.RedisConfig{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		}
		return redis.NewRedisCache(rConfig)

	case Memcached:
		if len(cfg.MemcachedServers) == 0 {
			return nil, errors.New("at least one memcached server is required")
		}
		client := gormemcache.New(cfg.MemcachedServers...)
		return memcached.New(client), nil

	default:
		return nil, errors.New("unsupported backend type")
	}
}
