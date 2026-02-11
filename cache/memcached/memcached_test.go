package memcached

import (
	"Go-library/cache"
	"Go-library/cache/cache/compliance"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func TestCompliance(t *testing.T) {
	// local memcache instance created in docker
	client := memcache.New("localhost:11211")
	err := client.Ping()
	if err != nil {
		// If it fails, skip the test
		t.Skip("Memcached is not running on localhost:11211, skipping compliance tests")
	}

	compliance.RunTest(t, func(t *testing.T) (cache.Cache, func(time.Duration)) {
		c := New(client)
		// clear for clean state and avoiding extra space
		_ = c.Clear()
		return c, func(d time.Duration) {
			time.Sleep(d)
		}
	})
}
