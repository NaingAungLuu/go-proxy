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
	mockedServer := setupMockedServer(t)
	server := proxy.NewServer(mockedServer.URL)

	// Setup custom log buffer catcher
	customLogger := CustomLoggerSpy{}
	// Attach custom logger buffer to proxy server
	server.AttachLogger(&customLogger)

	// Make an api request to the proxy proxy server
	request := httptest.NewRequest("GET", "/test", nil)
	server.ServeHTTP(httptest.NewRecorder(), request)

	if customLogger.TimesCalled < 1 {
		t.Error("Log buffer not received")
	}
}
