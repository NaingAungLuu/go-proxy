package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type ProxyServer struct {
	HttpClient *http.Client
	URL        string
	Port       int
	Logger     io.Writer
}

type ProxyHandler struct {
	HttpClient *http.Client
	Logger     io.Writer
	URL        string
	Port       int
}

func (p *ProxyServer) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	finalUrl := p.URL + request.URL.Path
	newRequest, _ := http.NewRequest(request.Method, finalUrl, request.Body)

	// Copy headers from the original request to the new request
	for key, values := range request.Header {
		for _, value := range values {
			newRequest.Header.Add(key, value)
		}
	}

	response, err := p.HttpClient.Do(newRequest)

	if err != nil {
		log.Fatalf("An error occurred: %+v", err)
	}

	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		log.Fatalf("An error occurred: %+v", err)
	}

	writeHeaders(w, *response)
	w.WriteHeader(response.StatusCode)
	w.Write(body)

	if p.Logger != nil {
		initialMessage := fmt.Sprintf("New request received\nURL:%+v\nMethod:%+v\nHeaders:%+v\nBody:%+v", request.URL, request.Method, request.Header, request.Body)
		p.Logger.Write([]byte(initialMessage))

		responseMessage := fmt.Sprintf("RESPONSE:\nHeaders\n%+v\nBody:\n%+v", response.Header, body)

		p.Logger.Write([]byte(responseMessage))
	}
}

func writeHeaders(w http.ResponseWriter, response http.Response) {
	// Copy headers from the response to the ResponseWriter
	for key, values := range response.Header {
		w.Header()[key] = values
	}
}

func (p *ProxyServer) AttachLogger(logger io.Writer) {
	p.Logger = logger
}

func NewServer(destinationUrl string) *ProxyServer {
	server := &ProxyServer{
		URL:        destinationUrl,
		Logger:     nil,
		Port:       3000,
		HttpClient: &http.Client{},
	}
	return server
}
