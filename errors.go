package cache

import "errors"

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrEmptyKey    = errors.New("key is empty")
	ErrKeyExpired  = errors.New("key has expired")
)
