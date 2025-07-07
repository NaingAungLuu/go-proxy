package cache

import (
	"testing"
)

func TestCacheStore(t *testing.T) {
	t.Run("CacheStore returns correct key", func(t *testing.T) {
		cache := NewCacheStore()
		testCases := []struct {
			key   string
			value string
		}{
			{
				key:   "test",
				value: "Hello, World",
			},
			{
				key:   "second test",
				value: "Another Test",
			},
			{
				key:   "123",
				value: "123",
			},
		}

		for _, testCase := range testCases {
			cache.Put(testCase.key, testCase.value)
			expected := testCase.value
			actual := cache.Get(testCase.key)

			if actual == "" {
				t.Errorf("Not expecting an empty string for value")
			}

			if actual != expected {
				t.Errorf("Expected: %v, Actual Value: %v", expected, actual)
			}
		}

	})
}

