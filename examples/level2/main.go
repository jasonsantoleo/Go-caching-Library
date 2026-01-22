package main

import (
	"fmt"

	"Go-library/cache"
)

func main() {
	// Level 2: Create a new cache instance
	c := cache.New()
	fmt.Println("✓ Cache created successfully!")

	// Level 2: Set operations
	fmt.Println("\n--- Setting Values ---")
	err := c.Set("name", "John Doe")
	if err != nil {
		fmt.Printf("Error setting 'name': %v\n", err)
	} else {
		fmt.Println("✓ Set 'name' = 'John Doe'")
	}

	c.Set("age", 30)
	fmt.Println("✓ Set 'age' = 30")

	c.Set("city", "New York")
	fmt.Println("✓ Set 'city' = 'New York'")

	c.Set("active", true)
	fmt.Println("✓ Set 'active' = true")

	// Level 2: Get operations
	fmt.Println("\n--- Getting Values ---")
	name, err := c.Get("name")
	if err != nil {
		fmt.Printf("Error getting 'name': %v\n", err)
	} else {
		fmt.Printf("✓ Get 'name' = %v\n", name)
	}

	age, err := c.Get("age")
	if err != nil {
		fmt.Printf("Error getting 'age': %v\n", err)
	} else {
		fmt.Printf("✓ Get 'age' = %v\n", age)
	}

	// Try to get a non-existent key
	_, err = c.Get("nonexistent")
	if err != nil {
		fmt.Printf("✓ Correctly returned error for non-existent key: %v\n", err)
	}

	// Level 2: Delete operation
	fmt.Println("\n--- Deleting Values ---")
	err = c.Delete("city")
	if err != nil {
		fmt.Printf("Error deleting 'city': %v\n", err)
	} else {
		fmt.Println("✓ Deleted 'city'")
	}

	// Verify deletion
	_, err = c.Get("city")
	if err != nil {
		fmt.Printf("✓ Verified 'city' is deleted: %v\n", err)
	}

	// Level 2: Clear operation
	fmt.Println("\n--- Clearing Cache ---")
	c.Clear()
	fmt.Println("✓ Cache cleared")

	// Verify cache is empty
	_, err = c.Get("name")
	if err != nil {
		fmt.Printf("✓ Verified cache is empty: %v\n", err)
	}

	fmt.Println("\n--- Level 2 Complete! ---")
	fmt.Println("All basic cache operations (Set, Get, Delete, Clear) are working!")
}
