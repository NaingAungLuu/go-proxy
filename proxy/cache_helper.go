package proxy

import (
	"bufio"
	"bytes"
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

func SerializeResponse(response http.Response) string {
	// body, err := io.ReadAll(response.Body)
	// defer response.Body.Close()

	// if err != nil {
	// log.Fatalf("Error while reading response body: %v", err)
	// }

	buffer := bytes.Buffer{}
	response.Write(&buffer)

	return buffer.String()
}

func ResponseFromString(input string) http.Response {
	reader := bytes.NewReader([]byte(input))
	response, _ := http.ReadResponse(bufio.NewReader(reader), nil)
	return *response
}
