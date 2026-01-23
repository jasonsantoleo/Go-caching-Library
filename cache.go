package cache

import "errors"

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrEmptyKey    = errors.New("key is empty")
)

// Cache represents a simple in-memory cache.
type Cache struct {
	data map[string]interface{}
}

// New creates a new instance of Cache.
func New() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

// Set adds or updates a value in the cache.
func (c *Cache) Set(key string, value interface{}) error {
	if key == "" {
		return ErrEmptyKey
	}
	c.data[key] = value
	return nil
}

// Get retrieves a value from the cache.
func (c *Cache) Get(key string) (interface{}, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}
	value, ok := c.data[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return value, nil
}

// Delete removes a key from the cache.
func (c *Cache) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	if _, ok := c.data[key]; !ok {
		return ErrKeyNotFound
	}
	delete(c.data, key)
	return nil
}

// Clear removes all keys from the cache.
func (c *Cache) Clear() {
	c.data = make(map[string]interface{})
}
