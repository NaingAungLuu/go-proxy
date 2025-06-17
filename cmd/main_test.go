package main

import (
	"flag"
	"os"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	// Save original arguments and restore them after the test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// We need to import flag directly in the test
	originalFlagCommandLine := flag.CommandLine
	defer func() { flag.CommandLine = originalFlagCommandLine }()

	testCases := []struct {
		name          string
		args          []string
		expectedURL   string
		expectedPort  int
		shouldFail    bool
		errorContains string
	}{
		{
			name:          "Error when no url provided",
			args:          []string{"cmd"},
			expectedURL:   "",
			expectedPort:  -1,
			shouldFail:    true,
			errorContains: "URL cannot be empty",
		},
		{
			name:         "URL provided, default port used",
			args:         []string{"cmd", "-url", "http://example.com"},
			expectedURL:  "http://example.com",
			expectedPort: 3000, // Default port in parseArgs
			shouldFail:   false,
		},
		{
			name:         "Both URL and port provided",
			args:         []string{"cmd", "-url", "http://example.com", "-port", "8080"},
			expectedURL:  "http://example.com",
			expectedPort: 8080,
			shouldFail:   false,
		},
		{
			name:          "Invalid port provided (out of range)",
			args:          []string{"cmd", "-url", "http://example.com", "-port", "70000"},
			expectedURL:   "",
			expectedPort:  0,
			shouldFail:    true,
			errorContains: "must be between",
		},
		{
			name:          "Negative port provided",
			args:          []string{"cmd", "-url", "http://example.com", "-port", "-5"},
			expectedURL:   "",
			expectedPort:  0,
			shouldFail:    true,
			errorContains: "must be between",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset flag.CommandLine for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			// Mock command line arguments
			os.Args = tc.args

			// Call the parseArgs function directly
			url, port, err := parseArgs()

			// Check error conditions
			if tc.shouldFail {
				if err == nil {
					t.Errorf("Expected an error but got nil")
				} else if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Expected error containing %q but got %q", tc.errorContains, err.Error())
				}
				return
			}

			// If we shouldn't fail, check the results
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if url != tc.expectedURL {
				t.Errorf("Expected URL %q but got %q", tc.expectedURL, url)
			}

			if port != tc.expectedPort {
				t.Errorf("Expected port %d but got %d", tc.expectedPort, port)
			}
		})
	}
}
