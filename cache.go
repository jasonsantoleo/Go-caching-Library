package cache

import (
	"time"
)

type Cache interface {
	// input : key,value output : error
	Set(key string, value interface{}) error
	// input : key,value,time to live output:error
	SetWithTTL(key string, value interface{}, ttl time.Duration) error
	//input : key,value. output:error
	Get(key string) (interface{}, error)
	//input : key output: (interface{},error)
	Delete(key string) error
	//input : output:error
	Clear() error
}
