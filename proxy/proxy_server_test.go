package proxy_test

import (
	"go-proxy/proxy"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestProxyTunnel(t *testing.T) {
	t.Run("Connects to the right destination", func(t *testing.T) {
		const (
			successfulStatusCode = 200
			proxyServerPort      = 3000
		)
		mockedSourceServer := setupMockedServer(t)
		defer mockedSourceServer.Close()

		proxyServer := proxy.NewServer(mockedSourceServer.URL)

		proxyRequest, err := http.NewRequest("GET", "/test", nil)

		if err != nil {
			t.Errorf("An unexpected error occurred: %+v", err)
		}

		rr := httptest.NewRecorder()
		proxyServer.ServeHTTP(rr, proxyRequest)
		shadowResponse, err := http.Get(mockedSourceServer.URL + "/test")

		if err != nil {
			t.Errorf("An unexpected error occurred: %+v", err)
		}

		if rr.Code != successfulStatusCode {
			t.Errorf("Expected status code value of %v, but got %v", successfulStatusCode, rr.Code)
		}

		if !reflect.DeepEqual(rr.Header(), shadowResponse.Header) {
			t.Errorf("Proxied Response has different header values than the shadowed request:\nExpected: %+v\nActual: %+v", shadowResponse.Header, rr.Header())
		}

	})

	t.Run("Strips the proxy headers", func(t *testing.T) {
		mockedServer := setupMockedServer(t)
		defer mockedServer.Close()

		proxyServer := proxy.NewServer(mockedServer.URL)
		proxyRequest, err := http.NewRequest("GET", "/test", nil)

		if err != nil {
			t.Fatalf("An unexpected error occurred %+v", err)
		}

		// Attach proxy headers to the request
		proxyHeaderList := map[string]string{
			"Proxy-Connection":    "Keep-Alive",
			"Proxy-Authorization": "Basic",
		}
		for key, value := range proxyHeaderList {
			proxyRequest.Header.Add(key, value)
		}
		// Make Request
		rr := httptest.NewRecorder()
		proxyServer.ServeHTTP(rr, proxyRequest)

		// Check for headers starting with Proxy
		for key, _ := range rr.Result().Header {
			lowerCaseKey := strings.ToLower(key)
			if strings.HasPrefix(lowerCaseKey, "proxy-") {
				t.Errorf("Key: %v shouldn't be present in the response", key)
			}
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
