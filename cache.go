package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrEmptyKey    = errors.New("key is empty")
	ErrKeyExpired  = errors.New("key has expired")
)

// in-memory cache with LRU eviction.
type Cache struct {
	maxSize int
	ll      *list.List
	data    map[string]*list.Element
	mu      sync.Mutex
}

// type entry struct {
// 	key   string
// 	value interface{}
// }

// changed 1
type entry struct {
	key       string
	value     interface{}
	expiresAt time.Time
}

// New creates a new instance of Cache.
func New() *Cache {
	return &Cache{
		maxSize: 0, // 0 means no limit
		ll:      list.New(),
		data:    make(map[string]*list.Element),
	}
}

// If valid is less than current size, eviction will happen.
func (c *Cache) SetMaxSize(size int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.maxSize = size
	if c.maxSize > 0 && c.ll.Len() > c.maxSize {
		c.evict()
	}
}

// changed 2(added default value for expiresAt)
// Set adds or updates a value in the cache. o(1)
func (c *Cache) Set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if key == "" {
		return ErrEmptyKey
	}

	if elem, ok := c.data[key]; ok {
		c.ll.MoveToFront(elem)
		elem.Value.(*entry).value = value
		return nil
	}

	elem := c.ll.PushFront(&entry{key, value, time.Time{}})
	c.data[key] = elem

	if c.maxSize > 0 && c.ll.Len() > c.maxSize {
		c.evict()
	}
	return nil
}

// change3
// have a seperate set with TTL,follows go-idiometic pattern
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if key == "" {
		return ErrEmptyKey
	}
	if elem, ok := c.data[key]; ok {
		c.ll.MoveToFront(elem)
		elem.Value.(*entry).value = value
		elem.Value.(*entry).expiresAt = time.Now().Add(ttl)
		return nil
	}
	elem := c.ll.PushFront(&entry{key, value, time.Now().Add(ttl)})
	c.data[key] = elem

	if c.maxSize > 0 && c.ll.Len() > c.maxSize {
		c.evict()
	}
	return nil
}

// change 4 (get revokes if expired)
// Get retrieves a value from the cache. o(1)
func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if key == "" {
		return nil, ErrEmptyKey
	}
	if elem, ok := c.data[key]; ok {
		//check if the key is exppired
		if !elem.Value.(*entry).expiresAt.IsZero() && time.Now().After(elem.Value.(*entry).expiresAt) {
			// c.Delete(key) //self dead loack happend because of nesterd calling
			c.ll.Remove(elem)
			delete(c.data, key)
			// return nil,ErrKeyExpired //(for debugging key expired is not something to be exposed )
			return nil, ErrKeyNotFound
		}
		c.ll.MoveToFront(elem)
		return elem.Value.(*entry).value, nil
	}
	return nil, ErrKeyNotFound
}

// no change required for TTL implementation
// Delete removes a key from the cache.
func (c *Cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if key == "" {
		return ErrEmptyKey
	}
	if elem, ok := c.data[key]; ok {
		c.ll.Remove(elem)
		delete(c.data, key)
		return nil
	}
	return ErrKeyNotFound
}

// no change required for TTL implementation
// Clear removes all keys from the cache.o(1)
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ll.Init()
	c.data = make(map[string]*list.Element)
}

// no change required for TTL implementation
func (c *Cache) evict() {
	for c.ll.Len() > c.maxSize {
		elem := c.ll.Back()
		if elem != nil {
			c.ll.Remove(elem)
			kv := elem.Value.(*entry)
			delete(c.data, kv.key)
		}
	}
}
