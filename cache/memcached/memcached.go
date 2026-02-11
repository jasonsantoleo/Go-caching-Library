package memcached

import (
	"Go-library/cache"
	"errors"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcachedCache (cache.Cache interface)
type MemcachedCache struct {
	client *memcache.Client
}

// Ensure MemcachedCache implements cache.Cache
var _ cache.Cache = (*MemcachedCache)(nil)

// constructor for memcache
func New(client *memcache.Client) *MemcachedCache {
	return &MemcachedCache{
		client: client,
	}
}

// sets add new value or update the old value
func (c *MemcachedCache) Set(key string, value interface{}) error {
	return c.SetWithTTL(key, value, 0)
}

// SetWithTTL adds or update a new valuye withh a expiration time
func (c *MemcachedCache) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	if key == "" {
		return cache.ErrEmptyKey
	}

	valStr, ok := value.(string)
	if !ok {
		return errors.New("memcached client only supports string values for now")
	}

	sec := int32(ttl.Seconds())
	if ttl > 0 && sec == 0 {
		sec = 1
	}

	item := &memcache.Item{
		Key:        key,
		Value:      []byte(valStr),
		Expiration: sec,
	}

	return c.client.Set(item)
}

// Get retrieves a value from the cache.
func (c *MemcachedCache) Get(key string) (interface{}, error) {
	if key == "" {
		return nil, cache.ErrEmptyKey
	}

	item, err := c.client.Get(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return nil, cache.ErrKeyNotFound
		}
		return nil, err
	}

	return string(item.Value), nil
}

// Delete removes a key from the cache.
func (c *MemcachedCache) Delete(key string) error {
	if key == "" {
		return cache.ErrEmptyKey
	}

	err := c.client.Delete(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return cache.ErrKeyNotFound
		}
		return err
	}
	return nil
}

// Clear removes all keys from the cache.
func (c *MemcachedCache) Clear() error {
	return c.client.DeleteAll()
}
