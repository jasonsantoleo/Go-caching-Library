package memory

import (
	"Go-library/cache"
	"container/list"
	"sync"
	"time"
)

// in-memory cache with LRU eviction.
type Memorycache struct {
	maxSize int
	ll      *list.List
	data    map[string]*list.Element
	mu      sync.Mutex
}

type entry struct {
	key       string
	value     interface{}
	expiresAt time.Time
}

// New creates a new instance of Cache.
func NewMemorycache() *Memorycache {
	return &Memorycache{
		maxSize: 0, // 0 means no limit
		ll:      list.New(),
		data:    make(map[string]*list.Element),
	}
}

var _ cache.Cache = (*Memorycache)(nil)

// change 1 addeds mutex
// If valid is less than current size, eviction will happen.
func (c *Memorycache) SetMaxSize(size int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.maxSize = size
	if c.maxSize > 0 && c.ll.Len() > c.maxSize {
		c.evict()
	}
}

// change 1 addeds mutex
// Set adds or updates a value in the cache. o(1)
func (c *Memorycache) Set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if key == "" {
		return cache.ErrEmptyKey
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

// change 1 addeds mutex
// have a seperate set with TTL,follows go-idiometic pattern
func (c *Memorycache) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if key == "" {
		return cache.ErrEmptyKey
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

// change 1 addeds mutex
// Get retrieves a value from the cache. o(1)
func (c *Memorycache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if key == "" {
		return nil, cache.ErrEmptyKey
	}
	if elem, ok := c.data[key]; ok {
		//check if the key is exppired
		if !elem.Value.(*entry).expiresAt.IsZero() && time.Now().After(elem.Value.(*entry).expiresAt) {
			// c.Delete(key) //self dead loack happend because of nesterd calling
			c.ll.Remove(elem)
			delete(c.data, key)
			// return nil,ErrKeyExpired //(for debugging key expired is not something to be exposed )
			return nil, cache.ErrKeyNotFound
		}
		c.ll.MoveToFront(elem)
		return elem.Value.(*entry).value, nil
	}
	return nil, cache.ErrKeyNotFound
}

// change 1 addeds mutex
// Delete removes a key from the cache.
func (c *Memorycache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if key == "" {
		return cache.ErrEmptyKey
	}
	if elem, ok := c.data[key]; ok {
		c.ll.Remove(elem)
		delete(c.data, key)
		return nil
	}
	return cache.ErrKeyNotFound
}

// change 1 addeds mutex
// Clear removes all keys from the cache.o(1)
func (c *Memorycache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ll.Init()
	c.data = make(map[string]*list.Element)
	return nil
}

// helper function (no change required for TTL implementation)
func (c *Memorycache) evict() {
	for c.ll.Len() > c.maxSize {
		elem := c.ll.Back()
		if elem != nil {
			c.ll.Remove(elem)
			kv := elem.Value.(*entry)
			delete(c.data, kv.key)
		}
	}
}
