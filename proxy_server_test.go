package proxy_test

import (
	"go-proxy"
	"net/http/httptest"
	"testing"
)

func TestProxyServerCreation(t *testing.T) {
	const (
		proxyUrl        = "http://dummyjson.com"
		proxyServerPort = 3000
	)
	t.Run("NewServer creates proper proxy", func(t *testing.T) {
		server := proxy.NewServer(proxyUrl, proxyServerPort)

		expectedAddress := ":3000"
		// expectedUrl := proxyUrl

		if server.Server.Addr != expectedAddress {
			t.Errorf("Expected Server port address to be %q, but got %q", expectedAddress, server.Server.Addr)
		}

	})
}

func TestProxyTunnel(t *testing.T) {

	const (
		successfulStatusCode = 200
		proxyUrl             = "http://dummyjson.com"
		proxyServerPort      = 3000
	)

	// Setup proxy server & response recorder
	server := proxy.NewServer(proxyUrl, proxyServerPort)
	rr := httptest.NewRecorder()

	t.Run("Returns OK for /test", func(t *testing.T) {
		rr.Flush()
		server.Server.Handler.ServeHTTP(rr, httptest.NewRequest("GET", "http://localhost:3000/test", nil))

		if rr.Code != successfulStatusCode {
			t.Errorf("Expected status code value of %v, but got %v", successfulStatusCode, rr.Code)
		}
	})

	t.Run("Returns todo list for /todos/1", func(t *testing.T) {
		const (
			requiredContentType = "application/json; charset=utf-8"
		)

		rr.Flush()
		server.Server.Handler.ServeHTTP(rr, httptest.NewRequest("GET", "http://localhost:3000/todos/random", nil))

		// Test for successful status code
		if rr.Code != successfulStatusCode {
			t.Errorf("Expected status code value of %v, but got %v", successfulStatusCode, rr.Code)
		}

		responseContentType := rr.Header().Get("Content-Type")

		if responseContentType != requiredContentType {
			t.Errorf("Expected content type : %q, but got %q", requiredContentType, responseContentType)
		}
	})
}
