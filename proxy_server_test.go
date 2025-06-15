package proxy_test

import (
	"go-proxy"
	"net/http/httptest"
	"testing"
)

func TestProxyStartServer(t *testing.T) {
	const (
		expectedStatusCode = 200
	)

	server := proxy.NewServer("http://dummyjson.com", 3000)
	rr := httptest.NewRecorder()

	server.Server.Handler.ServeHTTP(rr, httptest.NewRequest("GET", "http://localhost:3000/test", nil))

	if rr.Code != expectedStatusCode {
		t.Errorf("Expected status code value of %v, but got %v", expectedStatusCode, rr.Code)
	}
}
