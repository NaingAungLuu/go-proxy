package proxy

import (
	"net/http"
	"strings"
)

const (
	separator = ";"
)

func UniqueKeyForRequest(request http.Request) string {
	keySegments := []string{
		strings.ToLower(request.Method),
		request.URL.String(),
	}
	return strings.Join(keySegments, separator)
}
