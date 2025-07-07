package proxy

import (
	"net/http"
	"strings"
)

func UniqueKeyForRequest(request http.Request) string {
	key := ""
	key += strings.ToLower(request.Method)
	key += ";"
	key += request.URL.String()
	return key
}
