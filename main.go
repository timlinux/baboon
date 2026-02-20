// Baboon - Terminal-based typing practice application
//
// This is the entry point that starts the REST API server (backend)
// and connects the terminal UI (frontend) to it via HTTP.
//
// Architecture:
//   - Backend: REST API server handling game logic and statistics
//   - Frontend: Bubble Tea TUI communicating via REST client
//
// Usage:
//
//	baboon              # Normal mode (starts backend + frontend)
//	baboon -p           # Punctuation mode (words separated by punctuation)
//	baboon -port 8080   # Use custom port for REST API
//	baboon -server      # Run backend server only (blocking)
//	baboon -client      # Run frontend only (connect to existing backend)
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/timlinux/baboon/backend"
	"github.com/timlinux/baboon/frontend"
)

func main() {
	// Parse command line flags
	punctuationMode := flag.Bool("p", false, "Enable punctuation mode (words separated by punctuation + space)")
	port := flag.Int("port", 8787, "Port for the REST API server")
	serverOnly := flag.Bool("server", false, "Run backend server only (no TUI)")
	clientOnly := flag.Bool("client", false, "Run frontend only (connect to existing backend)")
	flag.Parse()

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	baseURL := fmt.Sprintf("http://%s", addr)

	// Validate flags
	if *serverOnly && *clientOnly {
		fmt.Println("Error: cannot use both -server and -client flags")
		os.Exit(1)
	}

	// Server-only mode: run backend and block
	if *serverOnly {
		runServerOnly(addr, *punctuationMode)
		return
	}

	// Client-only mode: connect to existing backend
	if *clientOnly {
		runClientOnly(baseURL, *punctuationMode)
		return
	}

	// Default mode: start backend and frontend together
	runCombined(addr, baseURL, *punctuationMode)
}

// runServerOnly starts the backend server and blocks until interrupted.
func runServerOnly(addr string, punctuationMode bool) {
	config := backend.DefaultConfig()
	config.PunctuationMode = punctuationMode

	server, err := backend.NewServer(config, addr)
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		os.Exit(1)
	}

	// Write PID file for management scripts
	pidFile := getPIDFilePath()
	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644); err != nil {
		fmt.Printf("Warning: could not write PID file: %v\n", err)
	}
	defer os.Remove(pidFile)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down server...")
		os.Remove(pidFile)
		os.Exit(0)
	}()

	fmt.Printf("Baboon backend server starting on %s\n", addr)
	fmt.Printf("PID: %d (written to %s)\n", os.Getpid(), pidFile)
	fmt.Println("Press Ctrl+C to stop")

	if err := server.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

// runClientOnly connects to an existing backend server.
func runClientOnly(baseURL string, punctuationMode bool) {
	client := frontend.NewClient(baseURL, punctuationMode)

	// Wait for server to be ready
	fmt.Printf("Connecting to backend at %s...\n", baseURL)
	if err := client.WaitForServer(5 * time.Second); err != nil {
		fmt.Printf("Error: Could not connect to backend: %v\n", err)
		fmt.Println("Make sure the backend is running with: baboon -server")
		os.Exit(1)
	}

	// Create a session on the server
	if err := client.CreateSession(); err != nil {
		fmt.Printf("Error creating session: %v\n", err)
		os.Exit(1)
	}
	defer client.DeleteSession()

	// Create and run TUI
	model := frontend.NewModel(client)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

// runCombined starts both backend and frontend together (default mode).
func runCombined(addr, baseURL string, punctuationMode bool) {
	config := backend.DefaultConfig()
	config.PunctuationMode = punctuationMode

	server, err := backend.NewServer(config, addr)
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		os.Exit(1)
	}

	// Start server in background
	server.StartAsync()

	client := frontend.NewClient(baseURL, punctuationMode)

	// Wait for server to be ready
	if err := client.WaitForServer(2 * time.Second); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create a session on the server
	if err := client.CreateSession(); err != nil {
		fmt.Printf("Error creating session: %v\n", err)
		os.Exit(1)
	}
	defer client.DeleteSession()

	// Create and run TUI
	model := frontend.NewModel(client)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

// getPIDFilePath returns the path to the PID file.
func getPIDFilePath() string {
	// Use XDG runtime dir if available, otherwise /tmp
	runDir := os.Getenv("XDG_RUNTIME_DIR")
	if runDir == "" {
		runDir = "/tmp"
	}
	return filepath.Join(runDir, "baboon.pid")
}
