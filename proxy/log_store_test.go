package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogStoreEntry(t *testing.T) {
	t.Run("Inserts the logs correctly", func(t *testing.T) {
		var logStore LogStore = NewLogStore()
		request := httptest.NewRequest("GET", "/", nil)
		response := &http.Response{}
		logStore.Save(request, response)
	})
}

func TestLogStoreQuery(t *testing.T) {
	t.Run("Returns the correct log", func(t *testing.T) {

	})
}
