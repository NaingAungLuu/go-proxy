package proxy

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
)

type ProxyServer struct {
	Server *http.Server
}

type ProxyHandler struct {
	URL        string
	HttpClient *http.Client
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
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
}

func writeHeaders(w http.ResponseWriter, response http.Response) {
	// Copy headers from the response to the ResponseWriter
	for key, values := range response.Header {
		w.Header()[key] = values
	}
}

func NewServer(destinationUrl string, port int) *ProxyServer {
	serverAddress := ":" + strconv.Itoa(port)
	proxyServer := &ProxyServer{
		Server: &http.Server{
			Addr: serverAddress,
			Handler: &ProxyHandler{
				URL:        destinationUrl,
				HttpClient: &http.Client{},
			},
		},
	}

	return proxyServer
}

func (p *ProxyServer) Start() {
	err := p.Server.ListenAndServe()
	if err != nil {
		log.Printf("An error occurred: %+v", err)
	}
}

func (p *ProxyServer) Stop() {
	// log.Println("Shutting Down Server")
	// defer log.Println("Server gracefully shutdown")
	p.Server.Shutdown(context.Background())
}
