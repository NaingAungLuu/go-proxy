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

func startServer(url string, port int, p *tea.Program) {
	// Starting the function as a "Goroutine" to allow both server and terminal UI
	// to run without blocking each other
	go func() {
		logger := ui.NewUiLogger(*p)
		server := proxy.NewServer(url)
		server.AttachLogger(&logger)
		http.ListenAndServe(":"+strconv.Itoa(port), http.HandlerFunc(server.ServeHTTP))
	}()
}

func main() {
	url, port, err := parseArgs()

	if err != nil {
		fmt.Println(err)
		return
	}

	p := tea.NewProgram(ui.NewTestModel(url, port), tea.WithAltScreen(), tea.WithMouseCellMotion())
	startServer(url, port, p)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
