package proxy

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ProxyServer struct {
	HttpClient *http.Client
	URL        string
	Logger     RequestLogger
}

type ProxyHandler struct {
	HttpClient *http.Client
	Logger     io.Writer
	URL        string
	Port       int
}

type RequestLogger interface {
	Log(request *http.Request)
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	// Setup Upstream request
	newRequest := createUpStreamRequest(*request, p.URL)

	// Make Network Request
	response, err := p.HttpClient.Do(newRequest)
	if err != nil {
		log.Fatalf("An error occurred: %+v", err)
	}

	writeResponse(w, *response)

	if p.Logger != nil {
		p.Logger.Log(request)
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
		w.Header()[key] = values
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
	}
	return server
}
