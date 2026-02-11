package memory

import (
	"Go-library/cache"
	"Go-library/cache/cache/compliance"
	"testing"
	"time"
)

// TestCompliance runs the shared test suite
func TestCompliance(t *testing.T) {
	compliance.RunTest(t, func(t *testing.T) (cache.Cache, func(time.Duration)) {
		return NewMemorycache(), nil
	})
}

//some extra test which is not part of the standard function (set,get,delete,cleat)
//since these are automatically applied in the redis and memcache this test specifically for the in-memory cache logic

// initialization
func TestNew(t *testing.T) {
	c := NewMemorycache()
	if c == nil {
		t.Fatalf("New() returned nil")
	}

	if c.data == nil {
		t.Fatalf("Cache data map is nil")
	}

	if len(c.data) != 0 {
		t.Fatalf("Expected empty cache, got %d entries", len(c.data))
	}
}

// checks eviction  with set
func TestLRUEviction(t *testing.T) {
	c := NewMemorycache()
	c.SetMaxSize(2)

	c.Set("a", 1)
	// Order: [a]
	c.Set("b", 2)
	// Order: [b, a]

	c.Set("c", 3)
	// Should evict 'a': [c, b]

	_, err := c.Get("a")
	if err != cache.ErrKeyNotFound {
		t.Fatalf("Expected 'a' to be evicted, but got error: %v", err)
	}

	if _, err := c.Get("b"); err != nil {
		t.Fatalf("Expected 'b' to exist, got error: %v", err)
	}
	if _, err := c.Get("c"); err != nil {
		t.Fatalf("Expected 'c' to exist, got error: %v", err)
	}
}

// checks eviction logic with set and get
func TestLRUAccessOrder(t *testing.T) {
	c := NewMemorycache()
	c.SetMaxSize(2)

	c.Set("a", 1)
	c.Set("b", 2)
	// Order: [b, a]

	// Access 'a' to move it to front
	c.Get("a")
	// Order: [a, b]

	c.Set("c", 3)
	// Should evict 'b' (LRU): [c, a]

	_, err := c.Get("b")
	if err != cache.ErrKeyNotFound {
		t.Fatalf("Expected 'b' to be evicted, but got error: %v", err)
	}

	if _, err := c.Get("a"); err != nil {
		t.Fatalf("Expected 'a' to exist, got error: %v", err)
	}
}

// get eviction logic withh update
func TestLRUUpdateOrder(t *testing.T) {
	c := NewMemorycache()
	c.SetMaxSize(2)

	c.Set("a", 1)
	c.Set("b", 2)
	// Order: [b, a]

	// Update 'a'
	c.Set("a", 10)
	// Order: [a, b]

	c.Set("c", 3)
	// Should evict 'b': [c, a]

	_, err := c.Get("b")
	if err != cache.ErrKeyNotFound {
		t.Fatalf("Expected 'b' to be evicted, but got error: %v", err)
	}

	val, err := c.Get("a")
	if err != nil {
		t.Fatalf("Expected 'a' to exist")
	}
	if val != 10 {
		t.Fatalf("Expected 'a' value to be 10, got %v", val)
	}
}

// TestTTLWithLRU - testing interaction between TTL and LRU (implementation specific)
func TestTTLWithLRU(t *testing.T) {
	c := NewMemorycache()
	c.SetMaxSize(2)
	c.SetWithTTL("a", 1, 50*time.Millisecond)
	c.SetWithTTL("b", 2, 150*time.Millisecond)
	//order [b,a]
	time.Sleep(100 * time.Millisecond)
	c.Set("c", 3)
	//should evict a since it is expired anyway
	_, err := c.Get("a")
	if err != cache.ErrKeyNotFound {
		t.Fatalf("expected a to be expired")
	}
	if _, err2 := c.Get("b"); err2 != nil {
		t.Fatalf("expected b to be exist,got %v", err2)
	}
	if _, err3 := c.Get("c"); err3 != nil {
		t.Fatalf("expected c to be exist,got %v", err3)
	}
}
