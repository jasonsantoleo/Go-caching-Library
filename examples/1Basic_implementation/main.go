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

	// Note: At Level 1, we only have the basic structure.
	// Operations will be added in Level 2.
}
