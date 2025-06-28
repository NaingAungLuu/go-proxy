package proxy_test

import (
	"go-proxy/proxy"
	"io"
	"net/http"
	"net/http/httptest"
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
		// _, err := http.Get(mockedSourceServer.URL + "/test")

		if err != nil {
			t.Errorf("An unexpected error occurred: %+v", err)
		}

		if rr.Code != successfulStatusCode {
			t.Errorf("Expected status code value of %v, but got %v", successfulStatusCode, rr.Code)
		}

	})
}

func TestProxyHandler(t *testing.T) {
	mockedServer := setupMockedServer(t)
	defer mockedServer.Close()

	t.Run("Strips the proxy headers", func(t *testing.T) {
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
		assertHeaderPrefixNotExist(t, rr.Result().Header, "proxy-")

	})

	t.Run("Attaches the correct proxy headers", func(t *testing.T) {
		proxyServer := proxy.NewServer(mockedServer.URL)
		proxyRequest, err := http.NewRequest("GET", "http://localhost:3001/test", nil)

		if err != nil {
			t.Errorf("An unepxected error occurred : %+v", err)
		}

		proxyRequest.Header.Add("X-Forwarded-Host", "dummyjson.com")

		t.Logf("HeaderList:\n")
		for key, value := range proxyRequest.Header {
			t.Logf(key + ":" + strings.Join(value, ","))
		}

		rr := httptest.NewRecorder()
		proxyServer.ServeHTTP(rr, proxyRequest)
	})

	t.Run("Returns the correct content headers", func(t *testing.T) {
		jsonMockServer := setupMockedContentServer(t)
		proxyServer := proxy.NewServer(jsonMockServer.URL)
		proxyRequest, err := http.NewRequest("GET", "http://localhost:3001/test", nil)

		if err != nil {
			t.Errorf("An error occurred while creating the test request: %+v", err)
		}

		proxyRequest.Header.Add("Accept", "application/json")
		rr := httptest.NewRecorder()

		proxyServer.ServeHTTP(rr, proxyRequest)
		contentTypeHeader := rr.Header().Get("Content-Type")

		if contentTypeHeader != "application/json" {
			t.Errorf("Expecting 'application/json' for 'Content-Type' Header, got '%+v'", contentTypeHeader)
		}

		t.Logf("Content-Type Header Value: %v", contentTypeHeader)

	})
}

func assertHeaderPrefixNotExist(t *testing.T, headerList http.Header, prefix string) {
	for key := range headerList {
		lowerCaseKey := strings.ToLower(key)
		if strings.HasPrefix(lowerCaseKey, prefix) {
			t.Errorf("Key: %v shouldn't be present in the response", key)
		}
	}
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

func setupMockedContentServer(t *testing.T) *httptest.Server {
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

		w.Header()["Content-Type"] = []string{"application/json"}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(body)
	}

	return httptest.NewServer(http.HandlerFunc(mockHttpHandler))
}
