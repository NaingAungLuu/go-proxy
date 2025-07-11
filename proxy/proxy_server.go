package proxy

import (
	"go-proxy/cache"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ProxyServer struct {
	HttpClient *http.Client
	URL        string
	Logger     RequestLogger
	Cache      cache.CacheStore
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	// Setup Upstream request
	newRequest := createUpStreamRequest(*request, p.URL)

	// Mark Start Time
	startTime := time.Now()

	// Make Network Request
	response, err := p.HttpClient.Do(newRequest)

	// Mark End Time
	timeTaken := time.Since(startTime)

	if err != nil {
		log.Fatalf("An error occurred: %+v", err)
	}

	writeResponse(w, *response)

	if p.Logger != nil {
		log := NewRequestLog(*request, timeTaken)
		p.Logger.Log(log)
	}
}

func createUpStreamRequest(request http.Request, destinationUrl string) *http.Request {
	newURL, _ := url.Parse(destinationUrl)
	newRequest, _ := http.NewRequest(request.Method, newURL.String(), request.Body)

	// Copy headers from the original request to the new request
	for key, values := range request.Header {
		for _, value := range values {
			newRequest.Header.Add(key, value)
		}
	}

	// Attach proxy related headers
	newRequest.Header.Add("X-Forwarded-Host", request.Header.Get("Host"))
	newRequest.URL.Host = newURL.Host
	newRequest.URL.Scheme = newURL.Scheme
	newRequest.Host = newURL.Host

	// Preprocess Request
	stripProxyHeaders(newRequest)

	return newRequest
}

func stripProxyHeaders(request *http.Request) {
	for key := range request.Header {
		if strings.HasPrefix(key, "Proxy") {
			request.Header.Del(key)
		}
	}
}

func writeResponse(w http.ResponseWriter, response http.Response) {
	body := readResponseBody(response)

	// Write Headers
	writeHeaders(w, response)
	w.WriteHeader(response.StatusCode)

	// Write Response & Body
	w.Write(body)

}

func readResponseBody(response http.Response) []byte {
	// Catch Response Body
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		log.Fatalf("Error while reading response body: %+v", err)
	}

	return body
}

func writeHeaders(w http.ResponseWriter, response http.Response) {
	// Copy headers from the response to the ResponseWriter
	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Set(key, value)
		}
	}
}

func (p *ProxyServer) AttachLogger(logger RequestLogger) {
	p.Logger = logger
}

func NewServer(destinationUrl string) *ProxyServer {
	server := &ProxyServer{
		URL:        destinationUrl,
		Logger:     nil,
		HttpClient: &http.Client{},
		Cache:      cache.NewCacheStore(),
	}
	return server
}
