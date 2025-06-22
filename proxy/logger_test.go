package proxy_test

import (
	"go-proxy/proxy"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CustomLoggerSpy struct {
	TimesCalled int
}

func (l *CustomLoggerSpy) Log(request *http.Request) {
	l.TimesCalled++
}

func TestLogger(t *testing.T) {
	const (
		proxyUrl        = "http://dummyjson.com"
		proxyServerPort = 3000
	)
	mockedServer := setupMockedServer(t)
	server := proxy.NewServer(mockedServer.URL)

	// Setup custom log buffer catcher
	customLogger := CustomLoggerSpy{}
	// Attach custom logger buffer to proxy server
	server.AttachLogger(&customLogger)

	// Make an api request to the proxy proxy server
	request, err := http.NewRequest("GET", "/test", nil)

	if err != nil {
		t.Errorf("Unexpected error occurred: %+v", err)
	}

	server.ServeHTTP(httptest.NewRecorder(), request)

	if customLogger.TimesCalled < 1 {
		t.Error("Log buffer not received")
	}
}
