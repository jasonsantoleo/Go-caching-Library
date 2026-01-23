package cache

import (
	"container/list"
	"errors"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrEmptyKey    = errors.New("key is empty")
)

// in-memory cache with LRU eviction.
type Cache struct {
	maxSize int
	ll      *list.List
	data    map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
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
	c.maxSize = size
	if c.maxSize > 0 && c.ll.Len() > c.maxSize {
		c.evict()
	}
}

// Set adds or updates a value in the cache.
func (c *Cache) Set(key string, value interface{}) error {
	if key == "" {
		return ErrEmptyKey
	}

	if elem, ok := c.data[key]; ok {
		c.ll.MoveToFront(elem)
		elem.Value.(*entry).value = value
		return nil
	}

	elem := c.ll.PushFront(&entry{key, value})
	c.data[key] = elem

	if c.maxSize > 0 && c.ll.Len() > c.maxSize {
		c.evict()
	}
	return nil
}

// Get retrieves a value from the cache.
func (c *Cache) Get(key string) (interface{}, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}
	if elem, ok := c.data[key]; ok {
		c.ll.MoveToFront(elem)
		return elem.Value.(*entry).value, nil
	}
	return nil, ErrKeyNotFound
}

// Delete removes a key from the cache.
func (c *Cache) Delete(key string) error {
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

// Clear removes all keys from the cache.
func (c *Cache) Clear() {
	c.ll.Init()
	c.data = make(map[string]*list.Element)
}

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
