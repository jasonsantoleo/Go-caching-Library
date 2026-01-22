package cache

import (
	"testing"
)

// TestNew tests the basic cache initialization (Level 1)
func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("New() returned nil")
	}

