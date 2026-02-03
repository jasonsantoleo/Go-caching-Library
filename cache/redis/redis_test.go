package redis

import (
	"Go-library/cache"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

// since real redis server is not yet set up we eun on miniredis to duplicate the behavior and check if test passess
// TestCacheStructure tests that the cache structure is correctly initialized
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

// TestSetGet tests the Set operation
func TestSetGet(t *testing.T) {

	c, _ := newTestCacheStructure(t)
	//empty key
	_, err := c.Get("")
	if err != cache.ErrEmptyKey {
		t.Fatalf("Expected ErrEmptyKey, got %v", err)
	}
	//key not found
	_, err1 := c.Get("non exist")
	if err1 != cache.ErrKeyNotFound {
		t.Fatalf("Expected ErrKeyNotFOund, got %v", err1)
	}
	// success case
	c.Set("a", 1)
	get2, err2 := c.Get("a")
	if get2 == nil {
		t.Fatalf("Expected val nil, got %v", err)
	}
	if err2 != nil {
		t.Fatalf("Error fetching the data %v", err)
	}
	// overwrite
	c.Set("a", 2)
	get3, err3 := c.Get("a")
	if err3 != nil {
		t.Fatalf("Error fetching the data %v", err)
	}
	// redis uses json -> number is float64
	if get3.(float64) != 2 {
		t.Fatalf("Expected 2,but got %v", get3)
	}
}
func TestSetWithTTL(t *testing.T) {
	c, mr := newTestCacheStructure(t)
	err := c.SetWithTTL("a", 1, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("setwithTTL failed,%v", err)
	}
	mr.FastForward(100 * time.Millisecond)
	get, err := c.Get("a")
	if err != cache.ErrKeyNotFound {
		t.Fatalf("Expected key not found, %v", get)
	}
}

// // TestDelete tests the Delete operation
func TestDelete(t *testing.T) {
	c, _ := newTestCacheStructure(t)

	// empty key
	err := c.Delete("")
	if err != cache.ErrEmptyKey {
		t.Fatalf("Expected ErrEmptyKey, got %v", err)
	}

	// delete non-existent
	err = c.Delete("a")
	if err != cache.ErrKeyNotFound {
		t.Fatalf("Expected ErrKeyNotFound, got %v", err)
	}
	// success case
	c.Set("a", 1)
	err = c.Delete("a")
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
}

// TestClear tests the Clear operation
func TestClear(t *testing.T) {
	c, _ := newTestCacheStructure(t)
	c.Set("a", 1)
	c.Set("b", 1)

	err := c.Clear()
	if err != nil {
		t.Fatalf("expected nil error, %v", err)
	}
	_, err = c.Get("a")
	if err != cache.ErrKeyNotFound {
		t.Fatalf("expected ErrKeyNotFound after Clear, got %v", err)
	}
}
