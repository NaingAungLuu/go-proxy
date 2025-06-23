package main

import (
	"errors"
	"flag"
	"fmt"
	"go-proxy/proxy"
	"go-proxy/ui"
	"log"
	"net/http"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// parseArgs parses and validates command line arguments
// Returns the URL, port, and any error encountered during parsing or validation
func parseArgs() (string, int, error) {
	// Define Program Flags
	var (
		baseUrl = flag.String("url", "", "Target Base URL of the proxy server")
		port    = flag.Int("port", -1, "Port number of the proxy server")
	)

	flag.Parse()

	if !flag.Parsed() {
		return "", 0, errors.New("invalid flags/params")
	}

	if *baseUrl == "" {
		return "", 0, fmt.Errorf("URL cannot be empty")
	}

	finalPort := *port
	if finalPort == -1 {
		fmt.Println("Port number not specified; Using Port 3000 by default")
		finalPort = 3000
	}

	// Additional validation could be added here
	if finalPort < 1 || finalPort > 65535 {
		return "", 0, fmt.Errorf("port must be between 1 and 65535")
	}

	return *baseUrl, finalPort, nil
}

type UILogger struct {
	program tea.Program
}

func (l *UILogger) Log(request *http.Request) {
	l.program.Send(ui.LogEvent{Request: request})
}

func startServer(url string, port int, p *tea.Program) {
	go func() {
		server := proxy.NewServer(url)
		server.AttachLogger(&UILogger{program: *p})
		http.ListenAndServe(":"+strconv.Itoa(port), http.HandlerFunc(server.ServeHTTP))
	}()
}
func main() {
	url, port, err := parseArgs()

	if err != nil {
		fmt.Println(err)
		return
	}

	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	startServer(url, port, p)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
