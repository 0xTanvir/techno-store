package main

import (
	"fmt"

	"techno-store/internal/domain/algo"
)

func main() {
	fmt.Println("Creating LRU Cache with 2 items: song1, song2")
	cache := algo.NewLRUCache(2)
	cache.Put("song1", "Song One")
	cache.Put("song2", "Song Two")

	fmt.Println("Retrieving song1 from cache")
	s1, exists := cache.Get("song1")
	if exists {
		fmt.Println("song1:", s1)
	} else {
		fmt.Println("song1 not found")
	}

	fmt.Println("For more detail, see internal/domain/algo/lru_test.go")
}
