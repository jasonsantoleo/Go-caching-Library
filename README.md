# Go Multi-Backend Caching Library

A robust, thread-safe caching library in Go that supports multiple backends (In-Memory and Redis), unified under a single interface.

## Features

- **Unified Interface**: All backends implement the `cache.Cache` interface, allowing you to switch storage mechanisms easily.
- **In-Memory Cache**:
    -   **LRU Eviction**: Automatically removes least recently used items when capacity is reached.
    -   **TTL Support**: Time-To-Live expiration for keys.
    -   **Thread-Safe**: Safe for concurrent use with `sync.Mutex`.
- **Redis Cache**:
    -   **Connector**: Standard Redis implementation using `go-redis`.
    -   **Persistence**: Leverages Redis for data durability and shared cache across instances.

## Installation

```bash
go get github.com/jasonsantoleo/Go-caching-Library
```

*(Note: Ensure your `go.mod` matches the import paths)*

## Usage

### 1. The Interface

Your application should rely on the `cache.Cache` interface to remain backend-agnostic:

```go
import "Go-library/cache"

var myCache cache.Cache
```

### 2. Initialize a Backend

#### In-Memory Cache

```go
import "Go-library/cache/cache/memory"

// Create a new in-memory cache
memCache := memory.NewMemorycache()

// Optional: Set max size (e.g., 100 items) for LRU eviction
memCache.SetMaxSize(100)

myCache = memCache
```

#### Redis Cache

```go
import "Go-library/cache/cache/redis"

// Create a new Redis cache
// You can pass configuration for your Redis server
redisCache := redis.NewRedisCache(redis.RedisConfig{
    Addr:     "localhost:6379",
    Password:("", // no password set
    DB:       0,  // use default DB
})

myCache = redisCache
```

### 3. Operations

The API is consistent across all backends:

```go
import (
    "fmt"
    "time"
    "Go-library/cache" 
)

func main() {
    // 1. Set a value
    err := myCache.Set("user:123", "Jason")
    if err != nil {
        panic(err)
    }

    // 2. Set with TTL (expires in 5 minutes)
    err = myCache.SetWithTTL("session:abc", "active", 5*time.Minute)

    // 3. Get a value
    val, err := myCache.Get("user:123")
    if err != nil {
        if err == cache.ErrKeyNotFound {
            fmt.Println("Key not found!")
        } else {
            fmt.Printf("Error: %v\n", err)
        }
    } else {
        fmt.Printf("Got value: %v\n", val)
    }

    // 4. Delete a key
    err = myCache.Delete("user:123")

    // 5. Clear the entire cache
    err = myCache.Clear()
}
```

## Project Structure

- `cache.go`: Defines the central `Cache` interface and common errors.
- `cache/memory/`: In-memory implementation (LRU, TTL).
- `cache/redis/`: Redis implementation.

## Testing

To run the test suite for all packages:

```bash
go test ./...
```