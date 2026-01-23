package cache

import (
	"testing"
)

// TestNew tests the basic cache initialization (Level 1)
func TestNew(t *testing.T) {
	c := New()
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
	c := New()

	// Level 1: Just verify the structure exists
	if c.data == nil {
		t.Fatalf("Cache data map should not be nil")
	}
}

// TestSet tests the Set operation
func TestSet(t *testing.T) {
	c := New()
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
	if err != ErrEmptyKey {
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
	c := New()

	// empty key
	_, err := c.Get("")
	if err != ErrEmptyKey {
		t.Fatalf("Expected ErrEmptyKey, got %v", err)
	}

	// non-existent key
	_, err = c.Get("non-exist")
	if err != ErrKeyNotFound {
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
	c := New()

	// empty key
	err := c.Delete("")
	if err != ErrEmptyKey {
		t.Fatalf("Expected ErrEmptyKey, got %v", err)
	}

	// delete non-existent
	err = c.Delete("a")
	if err != ErrKeyNotFound {
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
	c := New()
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
	c := New()
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
	if err != ErrKeyNotFound {
		t.Fatalf("Expected ErrKeyNotFound for deleted key, got %v", err)
	}
}

func TestLRUEviction(t *testing.T) {
	c := New()
	c.SetMaxSize(2)

	c.Set("a", 1)
	c.Set("b", 2)

	// Cache is full: [b, a]

	c.Set("c", 3)
	// Should evict 'a': [c, b]

	_, err := c.Get("a")
	if err != ErrKeyNotFound {
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
	c := New()
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
	if err != ErrKeyNotFound {
		t.Fatalf("Expected 'b' to be evicted, but got error: %v", err)
	}

	if _, err := c.Get("a"); err != nil {
		t.Fatalf("Expected 'a' to exist, got error: %v", err)
	}
}

func TestLRUUpdateOrder(t *testing.T) {
	c := New()
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
	if err != ErrKeyNotFound {
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
