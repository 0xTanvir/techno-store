package algo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a populated cache for testing
func populateCache(capacity int, itemCount int) *LRUCache {
	cache := NewLRUCache(capacity)
	for i := 0; i < itemCount; i++ {
		letter := 'A' + rune(i)
		cache.Put(string(letter), string(letter))
	}
	return cache
}

func TestPutNewItems(t *testing.T) {
	assert := assert.New(t)
	cache := NewLRUCache(2)
	cache.Put("song1", "Song One")
	cache.Put("song2", "Song Two")

	_, exists := cache.Get("song1")
	assert.True(exists, "song1 should be found in the cache")

	_, exists = cache.Get("song2")
	assert.True(exists, "song2 should be found in the cache")
}

func TestGetUpdatesOrder(t *testing.T) {
	assert := assert.New(t)
	cache := populateCache(2, 2) // Cache with 2 items, 'A' and 'B'
	cache.Get("A")               // Access 'A'

	// Add 'C', expect 'B' to be evicted
	cache.Put("C", "C")
	_, exists := cache.Get("B")
	assert.False(exists, "B should have been evicted")
}

func TestEvictionPolicy(t *testing.T) {
	assert := assert.New(t)
	cache := populateCache(2, 2) // Cache with 2 items, 'A' and 'B'

	// Add 'C', expect 'A' to be evicted
	cache.Put("C", "C")
	_, exists := cache.Get("A")
	assert.False(exists, "A should have been evicted")
}

func TestUpdateExistingItem(t *testing.T) {
	assert := assert.New(t)
	cache := populateCache(2, 1) // Cache with 1 item, 'A'

	// Update 'A' with new value
	cache.Put("A", "Updated A")
	value, exists := cache.Get("A")
	assert.True(exists, "A should exist")
	assert.Equal("Updated A", value, "Value of A should be updated")
}
