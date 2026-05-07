// Baboon - Terminal-based typing practice application
//
// This is the entry point for the typing practice application.
//
// Architecture:
//   - TUI mode (default): Embedded engine with Bubble Tea frontend (single binary)
//   - Server mode: REST API server for web frontend or remote clients
//   - Client mode: TUI connecting to a remote server
//
// Usage:
//
//	baboon              # TUI mode with embedded engine (default)
//	baboon -p           # Punctuation mode (words separated by punctuation)
//	baboon -server      # Run REST API server only (for web frontend)
//	baboon -client      # Run TUI connecting to existing server
//	baboon -port 8080   # Use custom port (for -server or -client modes)
package main

import (
	"flag"
	"fmt"
	"net/http"
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

// Version information set by ldflags during build
var (
	// Version is the semantic version (e.g., "1.11.0")
	Version = "dev"
	// GitCommit is the short git commit hash (e.g., "abc1234")
	GitCommit = "unknown"
)

// runEmbedded runs the TUI with the engine embedded directly (default mode).
// This is a single-binary mode with no client-server architecture.
func runEmbedded(punctuationMode bool) {
	config := backend.DefaultConfig()
	config.PunctuationMode = punctuationMode

	engine, err := backend.NewEngine(config)
	if err != nil {
		fmt.Printf("Error creating engine: %v\n", err)
		os.Exit(1)
	}

	// Create and run TUI with embedded engine
	model := frontend.NewModel(engine, Version, GitCommit)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	// Check for subcommands first
	if len(os.Args) > 1 && os.Args[1] == "web" {
		runWebCommand(os.Args[2:])
		return
	}

	// Parse command line flags
	punctuationMode := flag.Bool("p", false, "Enable punctuation mode (words separated by punctuation + space)")
	port := flag.Int("port", 8787, "Port for the REST API server")
	serverOnly := flag.Bool("server", false, "Run REST API server only (for web frontend)")
	clientOnly := flag.Bool("client", false, "Run TUI connecting to existing server")
	flag.Parse()

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	baseURL := fmt.Sprintf("http://%s", addr)

	// Validate flags
	if *serverOnly && *clientOnly {
		fmt.Println("Error: cannot use both -server and -client flags")
		os.Exit(1)
	}

	// Server-only mode: run REST API server and block
	if *serverOnly {
		runServerOnly(addr, *punctuationMode)
		return
	}

	// Client-only mode: connect to existing backend
	if *clientOnly {
		runClientOnly(baseURL, *punctuationMode)
		return
	}

	// Default mode: TUI with embedded engine (no server)
	runEmbedded(*punctuationMode)
}

// runServerOnly starts the backend server and blocks until interrupted.
func runServerOnly(addr string, punctuationMode bool) {
	config := backend.DefaultConfig()
	config.PunctuationMode = punctuationMode

	// Load database and auth configuration from environment
	loadConfigFromEnv(&config)

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
	model := frontend.NewModel(client, Version, GitCommit)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

// runWebCommand handles the 'web' subcommand for serving the web frontend with optional AdSense.
func runWebCommand(args []string) {
	webFlags := flag.NewFlagSet("web", flag.ExitOnError)
	port := webFlags.Int("port", 8787, "Port for the web server")
	adsenseKey := webFlags.String("adsense", "", "Google AdSense publisher ID (e.g., ca-pub-1234567890) or 'preview' to show ad placeholder")
	webDir := webFlags.String("dir", "web/dist", "Directory containing built web frontend")
	punctuationMode := webFlags.Bool("p", false, "Enable punctuation mode by default")

	webFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: baboon web [options]\n\n")
		fmt.Fprintf(os.Stderr, "Serve the web frontend with the backend API.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		webFlags.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  baboon web -adsense ca-pub-1234567890   # Show real ads\n")
		fmt.Fprintf(os.Stderr, "  baboon web -adsense preview             # Show ad preview placeholder\n")
	}

	if err := webFlags.Parse(args); err != nil {
		os.Exit(1)
	}

	addr := fmt.Sprintf("0.0.0.0:%d", *port)

	config := backend.DefaultConfig()
	config.PunctuationMode = *punctuationMode

	// Load database and auth configuration from environment
	loadConfigFromEnv(&config)

	// Handle adsense flag: "preview" shows placeholder, actual key shows real ads
	if *adsenseKey == "preview" {
		config.AdsenseEnabled = true
		config.AdsenseKey = ""
	} else if *adsenseKey != "" {
		config.AdsenseEnabled = true
		config.AdsenseKey = *adsenseKey
	}

	server, err := backend.NewServer(config, addr)
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		os.Exit(1)
	}

	// Check if web directory exists
	if _, err := os.Stat(*webDir); os.IsNotExist(err) {
		fmt.Printf("Warning: Web directory '%s' not found.\n", *webDir)
		fmt.Println("Run 'make web-build' first to build the web frontend.")
		fmt.Println("Starting API server only...")
	} else {
		// Set the static file directory
		server.SetStaticDir(*webDir)
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

	fmt.Printf("🐒 Baboon web server starting on http://%s\n", addr)
	if *adsenseKey != "" {
		fmt.Printf("   AdSense enabled: %s\n", *adsenseKey)
	}
	fmt.Println("   Press Ctrl+C to stop")

	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Server error: %v\n", err)
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

// getDefaultDatabasePath returns the default database path in the config directory.
func getDefaultDatabasePath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			// Fallback to current directory
			return "./baboon.db"
		}
		configDir = filepath.Join(home, ".config")
	}
	baboonDir := filepath.Join(configDir, "baboon")

	// Ensure the directory exists
	if err := os.MkdirAll(baboonDir, 0755); err != nil {
		fmt.Printf("Warning: could not create config directory: %v\n", err)
		return "./baboon.db"
	}

	return filepath.Join(baboonDir, "baboon.db")
}

// loadConfigFromEnv loads configuration from environment variables.
func loadConfigFromEnv(config *backend.Config) {
	// Set version information from build-time ldflags
	config.Version = Version
	config.GitCommit = GitCommit

	// Database configuration - use default path if not specified
	if dsn := os.Getenv("BABOON_DATABASE_DSN"); dsn != "" {
		config.DatabaseDSN = dsn
	} else {
		// Set default database path for automatic creation
		config.DatabaseDSN = getDefaultDatabasePath()
	}

	// JWT secret for authentication
	if secret := os.Getenv("BABOON_JWT_SECRET"); secret != "" {
		config.JWTSecret = secret
	}

	// Base URL for OAuth callbacks
	if baseURL := os.Getenv("BABOON_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	// OAuth providers
	if id := os.Getenv("BABOON_GOOGLE_CLIENT_ID"); id != "" {
		config.GoogleClientID = id
	}
	if secret := os.Getenv("BABOON_GOOGLE_CLIENT_SECRET"); secret != "" {
		config.GoogleClientSecret = secret
	}
	if id := os.Getenv("BABOON_GITHUB_CLIENT_ID"); id != "" {
		config.GitHubClientID = id
	}
	if secret := os.Getenv("BABOON_GITHUB_CLIENT_SECRET"); secret != "" {
		config.GitHubClientSecret = secret
	}
	if id := os.Getenv("BABOON_APPLE_CLIENT_ID"); id != "" {
		config.AppleClientID = id
	}
	if secret := os.Getenv("BABOON_APPLE_CLIENT_SECRET"); secret != "" {
		config.AppleClientSecret = secret
	}
	if id := os.Getenv("BABOON_MICROSOFT_CLIENT_ID"); id != "" {
		config.MicrosoftClientID = id
	}
	if secret := os.Getenv("BABOON_MICROSOFT_CLIENT_SECRET"); secret != "" {
		config.MicrosoftClientSecret = secret
	}
}
