package main

import (
	"Go-library/cache"
	"fmt"
)

func main() {
	// Level 1: Create a new cache instance
	c := cache.New()

	// Level 1: Just structure, no operations yet
	fmt.Printf("Cache created successfully!\n")
	fmt.Printf("Cache structure: %+v\n", c)

	// The basic structure.
}
