package factory

import (
	"testing"
)

func TestNewMemory(t *testing.T) {
	c, err := New(Memory, Config{
		MemoryMaxSize: 100,
	})
	if err != nil {
		t.Fatalf("Failed to create memory cache: %v", err)
	}
	if c == nil {
		t.Fatal("Returned cache is nil")
	}

	// Basic check
	err = c.Set("foo", "bar")
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}
}

func TestNewRedis(t *testing.T) {
	c, err := New(Redis, Config{
		RedisAddr: "localhost:6380",
	})
	if err != nil {
		t.Skipf("Skipping redis test (connection failed): %v", err)
	}
	if c == nil {
		t.Fatal("Returned cache is nil")
	}
}

func TestNewMemcached(t *testing.T) {
	c, err := New(Memcached, Config{
		MemcachedServers: []string{"localhost:11211"},
	})
	if err != nil {
		t.Skipf("Skipping memcached test (connection failed): %v", err)
	}
	if c == nil {
		t.Fatal("Returned cache is nil")
	}
}

func TestUnknownBackend(t *testing.T) {
	_, err := New("unknown", Config{})
	if err == nil {
		t.Fatal("Expected error for unknown backend, got nil")
	}
}
