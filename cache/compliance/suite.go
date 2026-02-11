package compliance

import (
	"Go-library/cache"
	"testing"
	"time"
)

// resuable test for the (set,get,delete,clear )
func RunTest(t *testing.T, setup func(t *testing.T) (cache.Cache, func(time.Duration))) {
	t.Run("SetGet", func(t *testing.T) {
		// setup is a function that returns a fresh instance of the Cache and a time advancer.
		c, _ := setup(t)
		testSetGet(t, c)
	})
	t.Run("Delete", func(t *testing.T) {
		c, _ := setup(t)
		testDelete(t, c)
	})
	t.Run("TTL", func(t *testing.T) {
		c, advanceTime := setup(t)
		if advanceTime == nil {
			advanceTime = time.Sleep
		}
		testTTL(t, c, advanceTime)
	})
	t.Run("TTLOverwrite", func(t *testing.T) {
		c, advanceTime := setup(t)
		if advanceTime == nil {
			// if the advance time is nil the time.sleep is used a defined
			advanceTime = time.Sleep
		}
		testTTLOverwrite(t, c, advanceTime)
	})
	t.Run("Clear", func(t *testing.T) {
		c, _ := setup(t)
		testClear(t, c)
	})
}

func testSetGet(t *testing.T, c cache.Cache) {
	//non-existent key
	_, err := c.Get("non-existent")
	if err != cache.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}

	// doing set
	err = c.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	//should get a value1 for key1
	val, err := c.Get("key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "value1" {
		t.Errorf("Expected 'value1', got %v", val)
	}

	// overwriting the key1 withh value2
	err = c.Set("key1", "value2")
	if err != nil {
		t.Fatalf("Set (overwrite) failed: %v", err)
	}

	val, err = c.Get("key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "value2" {
		t.Errorf("Expected 'value2', got %v", val)
	}

	// empty-key check
	err = c.Set("", "val")
	if err != cache.ErrEmptyKey {
		t.Errorf("Expected ErrEmptyKey for Set, got %v", err)
	}
	_, err = c.Get("")
	if err != cache.ErrEmptyKey {
		t.Errorf("Expected ErrEmptyKey for Get, got %v", err)
	}
}

func testDelete(t *testing.T, c cache.Cache) {
	//  non-existent key
	err := c.Delete("non-existent")
	if err != cache.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}

	// setting the value to check  Delete
	c.Set("key-delete", "val")
	err = c.Delete("key-delete")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// verify deletion
	_, err = c.Get("key-delete")
	if err != cache.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound after delete, got %v", err)
	}

	// empty key
	err = c.Delete("")
	if err != cache.ErrEmptyKey {
		t.Errorf("Expected ErrEmptyKey for Delete, got %v", err)
	}
}

func testTTL(t *testing.T, c cache.Cache, advanceTime func(time.Duration)) {
	// set with ttl
	ttl := 1 * time.Second
	err := c.SetWithTTL("key-ttl", "val", ttl)
	// verify the set
	if err != nil {
		t.Fatalf("SetWithTTL failed: %v", err)
	}

	// verify if value exist
	val, err := c.Get("key-ttl")
	if err != nil {
		t.Fatalf("Get failed immediate: %v", err)
	}
	if val != "val" {
		t.Errorf("Expected 'val', got %v", val)
	}

	// do expiration
	advanceTime(2 * ttl)

	// good if its gone
	_, err = c.Get("key-ttl")
	if err != cache.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound after expiration, got %v", err)
	}
}

func testClear(t *testing.T, c cache.Cache) {
	//set two element for clear verification
	c.Set("k1", "v1")
	c.Set("k2", "v2")
	err := c.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}
	// shoudl not get k1
	_, err = c.Get("k1")
	if err != cache.ErrKeyNotFound {
		t.Errorf("Expected k1 to be gone after Clear")
	}
	// should not get k2
	_, err = c.Get("k2")
	if err != cache.ErrKeyNotFound {
		t.Errorf("Expected k2 to be gone after Clear")
	}
}

func testTTLOverwrite(t *testing.T, c cache.Cache, advanceTime func(time.Duration)) {
	// apply a small ttl
	err := c.SetWithTTL("key", "val1", 1*time.Second)
	if err != nil {
		t.Fatalf("SetWithTTL failed: %v", err)
	}

	// apply a large ttl
	err = c.SetWithTTL("key", "val2", 3*time.Second)
	if err != nil {
		t.Fatalf("SetWithTTL overwite failed: %v", err)
	}

	// skip througth some time after 1st ttl abd before 2nd ttl
	advanceTime(1500 * time.Millisecond)

	// make sure the key is not deleted
	val, err := c.Get("key")
	if err != nil {
		t.Fatalf("Expected key to exist after overwrite extension, got error: %v", err)
	}
	if val != "val2" {
		t.Errorf("Expected 'val2', got %v", val)
	}

	// advance more thgan 2nd ttl
	advanceTime(2000 * time.Millisecond)

	// now the ttl should be gone
	_, err = c.Get("key")
	if err != cache.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound after extended expiration, got %v", err)
	}
}
