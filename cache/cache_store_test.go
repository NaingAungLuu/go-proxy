package cache

import (
	"testing"
)

func TestCacheStore(t *testing.T) {
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

	assertKeyValue := func(t *testing.T, cache CacheStore, key, value string) {
		t.Helper()

		actual := cache.Get(key)
		if actual == "" {
			t.Errorf("Not expecting an empty string for value")
		}

		if actual != value {
			t.Errorf("Expected: %v, Actual Value: %v", value, actual)
		}
	}

	t.Run("CacheStore stores and fetches keys correctly", func(t *testing.T) {
		cache := NewCacheStore()
		for _, testCase := range testCases {
			cache.Put(testCase.key, testCase.value)
			assertKeyValue(t, cache, testCase.key, testCase.value)
		}
	})

}

