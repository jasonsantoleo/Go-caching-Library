package memory

import (
	"Go-library/cache"
	"testing"
	"time"
)

// TestNew tests the basic cache initialization (Level 1)
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

// TestCacheStructure tests that the cache structure is correctly initialized
func TestCacheStructure(t *testing.T) {
	c := NewMemorycache()

	// Level 1: Just verify the structure exists
	if c.data == nil {
		t.Fatalf("Cache data map should not be nil")
	}
}

// TestSet tests the Set operation
func TestSet(t *testing.T) {
	c := NewMemorycache()
	// success case
	err := c.Set("a", 1)
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
	if len(c.data) != 1 {
		t.Fatalf("Expected len 1, got %d", len(c.data))
	}

	// empty key
	err = c.Set("", 1)
	if err != cache.ErrEmptyKey {
		t.Fatalf("Expected ErrEmptyKey, got %v", err)
	}

	// overwrite
	err = c.Set("a", 2)
	if err != nil {
		t.Fatalf("Expected nil error on overwrite, got %v", err)
	}
	if val, ok := c.data["a"]; !ok || val.Value.(*entry).value != 2 {
		t.Fatalf("Expected value 2, got %v", val.Value.(*entry).value)
	}
}

// TestGet tests the Get operation
func TestGet(t *testing.T) {
	c := NewMemorycache()

	// empty key
	_, err := c.Get("")
	if err != cache.ErrEmptyKey {
		t.Fatalf("Expected ErrEmptyKey, got %v", err)
	}

	// non-existent key
	_, err = c.Get("non-exist")
	if err != cache.ErrKeyNotFound {
		t.Fatalf("Expected ErrKeyNotFound, got %v", err)
	}

	// success case
	c.Set("a", 1)
	val, err := c.Get("a")
	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}
	if val != 1 {
		t.Fatalf("Expected value 1, got %v", val)
	}
}

// TestDelete tests the Delete operation
func TestDelete(t *testing.T) {
	c := NewMemorycache()

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
	if len(c.data) != 0 {
		t.Fatalf("Expected empty cache, got len %d", len(c.data))
	}
}

// TestClear tests the Clear operation
func TestClear(t *testing.T) {
	c := NewMemorycache()
	c.Set("a", 1)
	c.Set("b", 1)

	if len(c.data) != 2 {
		t.Fatalf("Expected len 2, got %d", len(c.data))
	}

	c.Clear()
	if len(c.data) != 0 {
		t.Fatalf("Expected len 0 after clear, got %d", len(c.data))
	}
}

// TestIntegration performs a sequence of operations
func TestIntegration(t *testing.T) {
	c := NewMemorycache()
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)

	if len(c.data) != 3 {
		t.Fatalf("Expected len 3, got %d", len(c.data))
	}

	val, err := c.Get("a")
	if err != nil || val != 1 {
		t.Fatalf("Expected 1, got %v, err: %v", val, err)
	}

	c.Delete("a")
	if len(c.data) != 2 {
		t.Fatalf("Expected len 2, got %d", len(c.data))
	}

	_, err = c.Get("a")
	if err != cache.ErrKeyNotFound {
		t.Fatalf("Expected ErrKeyNotFound for deleted key, got %v", err)
	}
}

func TestLRUEviction(t *testing.T) {
	c := NewMemorycache()
	c.SetMaxSize(2)

	c.Set("a", 1)
	c.Set("b", 2)

	// Cache is full: [b, a]

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

// ttl tests
// test if expired
func TestTTLExpired(t *testing.T) {
	c := NewMemorycache()
	c.SetWithTTL("a", 1, 50*time.Millisecond)
	time.Sleep(100 * time.Millisecond)
	_, err := c.Get("a")
	if err != cache.ErrKeyNotFound {
		t.Fatalf("Expected ErrKeyNotFound for expired key,got %v", err)
	}
}

// test if not expired
func TestTTLNotExpired(t *testing.T) {
	c := NewMemorycache()
	c.SetWithTTL("a", 1, 150*time.Millisecond)
	time.Sleep(50 * time.Millisecond)
	val, err := c.Get("a")
	if err != nil {
		t.Fatalf("Expected nil error,got %v", err)
	}
	if val != 1 {
		t.Fatalf("Expected value 1 ,got %v", val)
	}
}

// test if ttl works with lru
func TestTTLWithLRU(t *testing.T) {
	c := NewMemorycache()
	c.SetMaxSize(2)
	c.SetWithTTL("a", 1, 50*time.Millisecond)
	c.SetWithTTL("b", 2, 150*time.Millisecond)
	//order [b,a]
	time.Sleep(100 * time.Millisecond)
	c.Set("c", 3)
	//should evict a
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

// test if ttl overwrites
func TestTTLOverwrite(t *testing.T) {
	c := NewMemorycache()
	c.SetWithTTL("a", 1, 50*time.Millisecond)
	c.SetWithTTL("a", 2, 150*time.Millisecond)
	time.Sleep(100 * time.Millisecond)
	val, err := c.Get("a")
	if err != nil {
		t.Fatalf("expected nil error,got %v", err)
	}
	if val != 2 {
		t.Fatalf("expected value 2 ,got %v", val)
	}
}
