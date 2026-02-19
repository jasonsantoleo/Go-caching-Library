# Multi-Backend Caching Library in Go

A robust, thread-safe caching library for Go that supports multiple backends (Memory, Redis, Memcached) with a unified API.

## Features

- **Unified API**: Switch backends easily using a factory pattern.
- **In-Memory Cache**: Built-in LRU (Least Recently Used) eviction policy.
- **Redis Support**: Seamless integration with Redis (v9).
- **Memcached Support**: Full support for Memcached servers.
- **TTL Support**: Time-To-Live expiration for all backends.
- **Thread Safe**: Safe for concurrent use in high-performance applications.

## Installation

```bash
go get github.com/jasonsantoleo/Go-caching-Library
```

## Quick Start

### 1. Import the Factory
```go
import "Go-library/cache/factory"
```

### 2. Initialize a Cache
Use `factory.New` with your desired backend type and configuration.

#### In-Memory (No dependencies)
```go
config := factory.Config{
    MemoryMaxSize: 100, // Max items before eviction (0 = unlimited)
}
cache, err := factory.New(factory.Memory, config)
```

#### Redis
```go
config := factory.Config{
    RedisAddr: "localhost:6379",
    RedisPassword: "", // set if your server requires auth
    RedisDB: 0,
}
cache, err := factory.New(factory.Redis, config)
```

#### Memcached
```go
config := factory.Config{
    MemcachedServers: []string{"localhost:11211"},
}
cache, err := factory.New(factory.Memcached, config)
```

### 3. Usage Example

```go
// Set a value with TTL
err := cache.SetWithTTL("session:1", "active", 5*time.Minute)

// Get a value
val, err := cache.Get("session:1")
if err != nil {
    // Check if it's a "key not found" error.
    fmt.Println("Key not found or expired")
} else {
    fmt.Println("Value:", val)
}

// Delete
err := cache.Delete("session:1")

// Clear all
err := cache.Clear()
```

## API Reference

### Configuration (`factory.Config`)
```go
type Config struct {
    MemoryMaxSize    int      // Max items for Memory cache
    RedisAddr        string   // Redis address "host:port"
    RedisPassword    string   // Redis password
    RedisDB          int      // Redis DB index
    MemcachedServers []string // List of Memcached servers
}
```

### Cache Interface
All backends implement the `cache.Cache` interface:
- `Set(key string, value interface{}) error`
- `SetWithTTL(key string, value interface{}, ttl time.Duration) error`
- `Get(key string) (interface{}, error)`
- `Delete(key string) error`
- `Clear() error`

### Time Complexity (In-Memory)
| Method | Complexity | Notes |
| :--- | :--- | :--- |
| `Set` | **O(1)** | Constant time map/list operations. |
| `SetWithTTL` | **O(1)** | Constant time (same as Set). |
| `Get` | **O(1)** | Map lookup + LRU update (O(1)). |
| `Delete` | **O(1)** | Map delete + list remove. |
| `Clear` | **O(1)** | Constant time re-initialization. |
| `SetMaxSize` | **O(N)** | Linear if resizing requires eviction (N = items to evict). |

## Tests & Verification

### Running Tests
To run the full suite, including integration tests, you need Redis and Memcached running.

**1. Start Services (Docker)**
Use these commands to start clean containers on standard test ports:
```bash
# Redis on port 6380 (to avoid conflicts with local Redis)
docker run -d -p 6380:6379 --name my-test-redis redis:alpine

# Memcached on port 11211
docker run -d -p 11211:11211 --name my-test-memcached memcached:alpine
```

**2. Run Commands**
```bash
# Run Unit & Integration Tests
go test -v ./cache/tests/...

# Run Example Application
go run examples/main.go
```

### Troubleshooting
- **"connection refused"**: Ensure Docker containers are running (`docker ps`). If stopped, restart them or re-run the `docker run` commands above.
- **"NOAUTH Authentication required"**: This means you are connecting to a Redis instance that requires a password (e.g., a preexisting local service). Our tests use `localhost:6380` to target the password-less container we created.