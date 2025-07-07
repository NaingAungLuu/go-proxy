package cache

import (
	"testing"
)

func TestCacheStore(t *testing.T) {
	t.Run("CacheStore returns correct key", func(t *testing.T) {
		cache := NewCacheStore()
		key := "test"
		expected := "Hello, World"
		cache.Put(key, expected)

		actual := cache.Get(key)

		if actual == "" {
			t.Errorf("Not expecting an empty string for value")
		}

		if actual != expected {
			t.Errorf("Expected: %v, Actual Value: %v", expected, actual)
		}
	})
}

