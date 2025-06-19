package proxy_test

import (
	"bytes"
	"go-proxy/proxy"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CustomLogger struct {
	Buffer *bytes.Buffer
	Test   string
}

func (l *CustomLogger) Write(b []byte) (n int, err error) {
	l.Buffer.Write(b)
	return l.Buffer.Len(), nil
}

func TestLogger(t *testing.T) {
	const (
		proxyUrl        = "http://dummyjson.com"
		proxyServerPort = 3000
	)
	mockedServer := setupMockedServer(t)
	server := proxy.NewServer(mockedServer.URL)

	// Setup custom log buffer catcher
	logBuffer := &bytes.Buffer{}
	customLogger := CustomLogger{Buffer: logBuffer}
	// Attach custom logger buffer to proxy server
	server.AttachLogger(&customLogger)

	// Make an api request to the proxy proxy server
	request, err := http.NewRequest("GET", "/test", nil)

	if err != nil {
		t.Errorf("Unexpected error occurred: %+v", err)
	}

	server.ServeHTTP(httptest.NewRecorder(), request)

	if logBuffer.Len() <= 0 {
		t.Error("Log buffer not received")
	}
}
