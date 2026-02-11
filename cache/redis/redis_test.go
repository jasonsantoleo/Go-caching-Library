package redis

import (
	"Go-library/cache"
	"Go-library/cache/cache/compliance"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

// testing with miniredis
func newTestCacheStructure(t *testing.T) (*RedisCache, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	t.Cleanup(mr.Close)

	c, err := NewRedisCache(RedisConfig{
		Addr: mr.Addr(),
		DB:   10,
	})
	if err != nil {
		t.Fatalf("failed to create redis cache: %v", err)
	}

	return c, mr
}

// run the shared test for (set,get,post,delete)
func TestCompliance(t *testing.T) {
	compliance.RunTest(t, func(t *testing.T) (cache.Cache, func(time.Duration)) {
		c, mr := newTestCacheStructure(t)
		return c, func(d time.Duration) {
			mr.FastForward(d)
		}
	})
}

// some extra test for redis specific
// TestSetWithTTLSpecific tests Redis-specific TTL behavior if needed,
func TestSetWithTTLSpecific(t *testing.T) {
	c, mr := newTestCacheStructure(t)
	//set a with a 50 milllisecond as ttl
	err := c.SetWithTTL("a", 1, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("setwithTTL failed,%v", err)
	}
	//skip 100 millisecond
	mr.FastForward(100 * time.Millisecond)
	//shoduld get keynotfound
	get, err := c.Get("a")
	if err != cache.ErrKeyNotFound {
		t.Fatalf("Expected key not found, %v", get)
	}
}
