package proxy

import (
	"net/http/httptest"
	"testing"
)

func TestCacheKeyHelper(t *testing.T) {
	// We'll use the format `method;uri` for key
	testCases := []struct {
		method string
		url    string
		key    string
	}{
		{
			method: "GET",
			url:    "https://google.com",
			key:    "get;https://google.com",
		},
		{
			method: "POST",
			url:    "http://localhost:3000/login",
			key:    "post;http://localhost:3000/login",
		},
		{
			method: "PUT",
			url:    "https://dummyjson.com/1",
			key:    "put;https://dummyjson.com/1",
		},
		{
			method: "delete",
			url:    "https://test.com/resources/1",
			key:    "delete;https://test.com/resources/1",
		},
	}

	t.Run("Returns the correct key format", func(t *testing.T) {

		for _, testCase := range testCases {
			mockedRequest := httptest.NewRequest(testCase.method, testCase.url, nil)
			expected := testCase.key
			actual := UniqueKeyForRequest(*mockedRequest)

			if actual != expected {
				t.Errorf("Unique Key is not correct, expected: %q, actual: %q", expected, actual)
			}

		}
	})
}

