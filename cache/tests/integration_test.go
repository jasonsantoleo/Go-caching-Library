package tests

import (
	"Go-library/cache"
	"Go-library/cache/cache/compliance"
	"Go-library/cache/cache/factory"
	"testing"
	"time"
)

// Integration testing that uses all the api (redis,in-memory,memcache) at will
func TestUnifiedAPI(t *testing.T) {
	// test in-memroy backend tests
	t.Run("Memory", func(t *testing.T) {
		compliance.RunTest(t, func(t *testing.T) (cache.Cache, func(time.Duration)) {
			c, err := factory.New(factory.Memory, factory.Config{
				MemoryMaxSize: 100,
			})
			if err != nil {
				t.Fatalf("Failed to create memory cache: %v", err)
			}
			return c, nil
		})
	})

	// test redis-backend
	t.Run("Redis", func(t *testing.T) {
		// creating a connection
		_, err := factory.New(factory.Redis, factory.Config{
			RedisAddr: "localhost:6380",
		})
		// check connectivity
		if err != nil {
			t.Skipf("Skipping Redis tests: %v", err)
		}
		//run all the test for redis
		compliance.RunTest(t, func(t *testing.T) (cache.Cache, func(time.Duration)) {
			c, err := factory.New(factory.Redis, factory.Config{
				RedisAddr: "localhost:6380",
			})
			if err != nil {
				t.Fatalf("Failed to create redis cache: %v", err)
			}
			// Clear before each test
			_ = c.Clear()
			return c, func(d time.Duration) {
				time.Sleep(d)
			}
		})
	})

	// test memcached backend
	t.Run("Memcached", func(t *testing.T) {
		// creatinif a connection for memcache
		c, err := factory.New(factory.Memcached, factory.Config{
			MemcachedServers: []string{"localhost:11211"},
		})
		if err != nil {
			t.Skipf("Skipping Memcached tests: %v", err)
		}

		// check connectivity
		err = c.Set("ping", "pong")
		if err != nil {
			t.Skipf("Skipping Memcached tests (connection check failed): %v", err)
		}
		//run the standard test for memcache
		compliance.RunTest(t, func(t *testing.T) (cache.Cache, func(time.Duration)) {
			c, _ := factory.New(factory.Memcached, factory.Config{
				MemcachedServers: []string{"localhost:11211"},
			})
			_ = c.Clear()
			return c, func(d time.Duration) {
				time.Sleep(d)
			}
		})
	})
}
