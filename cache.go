package cache

import "errors"

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrEmptyKey    = errors.New("key is empty")
)

type cache struct {
	data map[string]interface{}
}

func New() *cache {
	return &cache{
		data: make(map[string]interface{}),
	}
}

func (c *cache) Set(key string, value interface{}) error {
	if key == "" {
		return ErrKeyNotFound
	}
	c.data[key] = value
	return nil
}

func (c *cache) Get(key string) (interface{}, error) {
	if key == "" {
		return nil, ErrEmptyKey
	}
	value, ok := c.data[key]
	if !ok {
		return nil, ErrEmptyKey
	}
	return value, nil
}
func (c *cache) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	_, empty := c.data[key]
	if !empty {
		return ErrKeyNotFound
	}
	delete(c.data, key)
	return nil
}
func (c *cache) Clear() {
	c.data = make(map[string]interface{})
}
