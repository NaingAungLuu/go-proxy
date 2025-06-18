package proxy

import (
	"bytes"
	"net/http"
	"testing"
)

type CustomLogger struct {
	Buffer *bytes.Buffer
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
	server := NewServer(proxyUrl, proxyServerPort)

	// Setup custom log buffer catcher
	logBuffer := &bytes.Buffer{}
	customLogger := CustomLogger{Buffer: logBuffer}
	// Attach custom logger buffer to proxy server
	server.AttachLogger(&customLogger)

	// Make an api request to the proxy proxy server
	_, err := http.Get("http://localhost:3000")

	if err != nil {
		t.Errorf("Unexpected error occurred: %+v", err)
	}

	if logBuffer.Len() <= 0 {
		t.Error("Log buffer not received")
	}
}
