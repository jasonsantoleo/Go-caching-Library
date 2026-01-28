# Go Multi-Backend Caching Library

implementation of a robust caching library in Go.

##  Progress

We have successfully implemented the core in-memory caching logic with LRU eviction.

###  Basic Structure
- Defined the fundamental `Cache` struct.
- Created initialization logic `New()`.

### Basic Operations
Implemented standard CRUD operations:
- **Set**: Add items to the cache.
- **Get**: Retrieve items (returns error if missing).
- **Delete**: Remove specific items.
- **Clear**: Flush the entire cache.

###  LRU Eviction Policy
Implemented Least Recently Used (LRU) eviction to manage memory usage efficiently.
- Uses **Doubly Linked List** (`container/list`) + **HashMap** for O(1) performance.
- **`SetMaxSize(size int)`**: Configurable capacity.
- Automatically removes the least recently used item when the cache exceeds its limit.
- Accessing an item (`Get` or `Set`) promotes it to the "most recently used" position.

### TTL and Expiration
Implemented Time-to-Live (TTL) for cache entries.
- **`SetWithTTL(key, value, duration)`**: Set a value that automatically expires after the specified duration.
- **Lazy Expiration**: Checks for expiration on retrieval (`Get`).
- **Integration with LRU**: Expired items are treated as non-existent and can be evicted normally or overwritten.

---

## Example

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
	c.Set("user:1", "Jason")
	c.Set("user:2", "Dhanalakshmi")

	// Get value
	val, _ := c.Get("user:1")
	fmt.Println(val) // Output: Jason   

	// LRU Eviction Demo
	// Cache is full [user:1, user:2] (Jason is MRU because we just accessed it)
	
	c.Set("user:3", "Rahul") 
	// "user:2" (Dhanalakshmi) is evicted because it was the Least Recently Used.
	
	_, err := c.Get("user:2")
	if err != nil {
		fmt.Println("Dhanalakshmi was evicted!") 
	}

	// TTL Demo
	c.SetWithTTL("session:1", "active", 50 * time.MilliSeconds)
	
	val, _ = c.Get("session:1")
	fmt.Println("Session:", val) // Output: active

	//  wait
	// time.Sleep(60 * time.MilliSeconds) 
	// _, err = c.Get("session:1") // Returns "key not found"
}
```

##  Testing

To run the test suite, verify all levels are working:

```bash
go test -v .
```

---

##  Next Steps
- ImplementingThread Safety (Mutex) 