package proxy_test

import (
	"go-proxy"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
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
		proxyServerPort      = 3000
	)
	t.Run("Returns OK for /test", func(t *testing.T) {
		mockedSourceServer := setupMockedServer(t)
		defer mockedSourceServer.Close()

		server := proxy.NewServer(mockedSourceServer.URL, proxyServerPort)
		rr := httptest.NewRecorder()

		proxyRequest := httptest.NewRequest("GET", "http://localhost:3000/test", nil)

		// Execute Requests
		server.Server.Handler.ServeHTTP(rr, proxyRequest)
		shadowResponse, err := http.Get(mockedSourceServer.URL + "/test")

		if err != nil {
			t.Fatalf("An unexpected error occurred while executing shadow request %+v", err)
		}

		if rr.Code != successfulStatusCode {
			t.Errorf("Expected status code value of %v, but got %v", successfulStatusCode, rr.Code)
		}

		if !reflect.DeepEqual(rr.Header(), shadowResponse.Header) {
			t.Errorf("Proxied Response has different header values than the shadowed request:\nExpected: %+v\nActual: %+v", shadowResponse.Header, rr.Header())
		}
	})

	t.Run("Returns todo list for /todos/1", func(t *testing.T) {
		mockedSourceServer := setupMockedServer(t)
		defer mockedSourceServer.Close()

		server := proxy.NewServer(mockedSourceServer.URL, proxyServerPort)
		rr := httptest.NewRecorder()

		proxyRequest := httptest.NewRequest("GET", "http://localhost:3000/todos/1", nil)

		server.Server.Handler.ServeHTTP(rr, proxyRequest)
		shadowResponse, err := http.Get(mockedSourceServer.URL + "/todos/1")

		if err != nil {
			t.Fatalf("An unexpected error occurred while executing shadow request %+v", err)
		}

		// Test for successful status code
		if rr.Code != successfulStatusCode {
			t.Errorf("Expected status code value of %v, but got %v", successfulStatusCode, rr.Code)
		}

		if !reflect.DeepEqual(rr.Header(), shadowResponse.Header) {
			t.Errorf("Proxied Response has different header values than the shadowed request:\nExpected: %+v\nActual: %+v", shadowResponse.Header, rr.Header())
		}
	})
}

func setupMockedServer(t *testing.T) *httptest.Server {
	mockHttpHandler := func(w http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		body, err := io.ReadAll(request.Body)

		if err != nil {
			w.WriteHeader(500)
			t.Errorf("An error occurred in mocked http handler : %+v", err)
		}

		for key, value := range request.Header {
			t.Logf("%+v -> %+v\n", key, value)
			w.Header()[key] = value
		}
		w.WriteHeader(200)
		w.Write(body)
	}

	return httptest.NewServer(http.HandlerFunc(mockHttpHandler))
}
