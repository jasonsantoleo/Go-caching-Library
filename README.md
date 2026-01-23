# Go Multi-Backend Caching Library

This is a progressive implementation of a robust caching library in Go. The project is structured into 10 levels of increasing complexity, from a simple in-memory cache to a production-ready system with multiple backends (Redis, Memcached) and advanced features.

## 🚀 Progress: Levels 1-3 Completed

We have successfully implemented the core in-memory caching logic with LRU eviction.

### ✅ Level 1: Basic Structure
- Defined the fundamental `Cache` struct.
- Created initialization logic `New()`.

### ✅ Level 2: Basic Operations
Implemented standard CRUD operations:
- **Set**: Add items to the cache.
- **Get**: Retrieve items (returns error if missing).
- **Delete**: Remove specific items.
- **Clear**: Flush the entire cache.

### ✅ Level 3: LRU Eviction Policy
Implemented Least Recently Used (LRU) eviction to manage memory usage efficiently.
- Uses **Doubly Linked List** (`container/list`) + **HashMap** for O(1) performance.
- **`SetMaxSize(size int)`**: Configurable capacity.
- Automatically removes the least recently used item when the cache exceeds its limit.
- Accessing an item (`Get` or `Set`) promotes it to the "most recently used" position.

---

## 🛠 Usage Example

```go
package main

import (
	"fmt"
	"Go-library/cache" 
)

func main() {
	// Initialize cache
	c := cache.New()
	
	// Set max size (for LRU eviction)
	c.SetMaxSize(2)

	// Basic Operations
	c.Set("user:1", "Alice")
	c.Set("user:2", "Bob")

	// Get value
	val, _ := c.Get("user:1")
	fmt.Println(val) // Output: Alice

	// LRU Eviction Demo
	// Cache is full [user:1, user:2] (Alice is MRU because we just accessed it)
	
	c.Set("user:3", "Charlie") 
	// "user:2" (Bob) is evicted because it was the Least Recently Used.
	
	_, err := c.Get("user:2")
	if err != nil {
		fmt.Println("Bob was evicted!") 
	}
}
```

## 🧪 Testing

To run the test suite, verify all levels are working:

```bash
go test -v .
```

---

## 🔜 Next Steps: Level 4 (TTL & Expiration)
The next phase will introduce Time-To-Live (TTL) support, allowing entries to expire automatically after a set duration.
