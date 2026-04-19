package backend

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/timlinux/baboon/auth"
	"github.com/timlinux/baboon/badge"
	"github.com/timlinux/baboon/database"
	"github.com/timlinux/baboon/filter"
	"github.com/timlinux/baboon/settings"
	"github.com/timlinux/baboon/stats"
)

// Session represents a single game session with its own engine.
type Session struct {
	ID        string
	Engine    *Engine
	CreatedAt time.Time
	LastUsed  time.Time
}

// Server provides a RESTful API for the game engine.
// It supports multiple concurrent sessions, each with their own game state.
type Server struct {
	config       Config
	sessions     map[string]*Session
	mu           sync.RWMutex
	addr         string
	staticDir    string
	db           *database.DB
	authService  *auth.Service
	authHandlers *auth.Handlers
	authMW       *auth.Middleware
}

// NewServer creates a new REST API server with the given configuration.
func NewServer(config Config, addr string) (*Server, error) {
	s := &Server{
		config:   config,
		sessions: make(map[string]*Session),
		addr:     addr,
	}

	// Initialize database if configured
	if config.DatabaseDSN != "" {
		db, err := database.New(database.Config{DSN: config.DatabaseDSN})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", err)
		}
		s.db = db

		// Initialize auth service
		authCfg := auth.Config{
			JWTSecret:             config.JWTSecret,
			BaseURL:               config.BaseURL,
			AccessTokenExpiry:     auth.DefaultConfig().AccessTokenExpiry,
			RefreshTokenExpiry:    auth.DefaultConfig().RefreshTokenExpiry,
			GoogleClientID:        config.GoogleClientID,
			GoogleClientSecret:    config.GoogleClientSecret,
			GitHubClientID:        config.GitHubClientID,
			GitHubClientSecret:    config.GitHubClientSecret,
			AppleClientID:         config.AppleClientID,
			AppleClientSecret:     config.AppleClientSecret,
			MicrosoftClientID:     config.MicrosoftClientID,
			MicrosoftClientSecret: config.MicrosoftClientSecret,
		}

		authService, err := auth.NewService(authCfg, db)
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to initialize auth service: %w", err)
		}
		s.authService = authService
		s.authHandlers = auth.NewHandlers(authService)
		s.authMW = auth.NewMiddleware(authService)
	}

	return s, nil
}

// GetAddr returns the server address.
func (s *Server) GetAddr() string {
	return s.addr
}

// SetStaticDir sets the directory for serving static web files.
func (s *Server) SetStaticDir(dir string) {
	s.staticDir = dir
}

// generateSessionID creates a unique session identifier.
func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// getSession retrieves a session by ID, updating last used time.
func (s *Server) getSession(id string) (*Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[id]
	if exists {
		session.LastUsed = time.Now()
	}
	return session, exists
}

// Start starts the HTTP server. This is a blocking call.
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Session management
	mux.HandleFunc("POST /api/sessions", s.handleCreateSession)
	mux.HandleFunc("DELETE /api/sessions/{id}", s.handleDeleteSession)
	mux.HandleFunc("GET /api/sessions", s.handleListSessions)

	// Game lifecycle (session-specific)
	mux.HandleFunc("POST /api/sessions/{id}/round", s.handleStartRound)

	// Input handling (session-specific)
	mux.HandleFunc("POST /api/sessions/{id}/keystroke", s.handleKeystroke)
	mux.HandleFunc("POST /api/sessions/{id}/backspace", s.handleBackspace)
	mux.HandleFunc("POST /api/sessions/{id}/clearinput", s.handleClearInput)
	mux.HandleFunc("POST /api/sessions/{id}/space", s.handleSpace)

	// State queries (session-specific)
	mux.HandleFunc("GET /api/sessions/{id}/state", s.handleGetState)
	mux.HandleFunc("GET /api/sessions/{id}/stats/session", s.handleGetSessionStats)
	mux.HandleFunc("GET /api/sessions/{id}/stats/historical", s.handleGetHistoricalStats)

	// Persistence (session-specific)
	mux.HandleFunc("POST /api/sessions/{id}/save", s.handleSaveStats)

	// Timing submission (session-specific)
	mux.HandleFunc("POST /api/sessions/{id}/timing", s.handleSubmitTiming)

	// Health check
	mux.HandleFunc("GET /api/health", s.handleHealth)

	// Version info
	mux.HandleFunc("GET /api/version", s.handleVersion)

	// Configuration endpoint for web frontend
	mux.HandleFunc("GET /api/config", s.handleConfig)

	// Authentication routes (only if auth is enabled)
	if s.authService != nil && s.authService.IsEnabled() {
		// OAuth login and callback
		mux.HandleFunc("GET /api/auth/{provider}/login", s.authHandlers.HandleLogin)
		mux.HandleFunc("GET /api/auth/{provider}/callback", s.authHandlers.HandleCallback)

		// Auth management
		mux.HandleFunc("POST /api/auth/logout", s.authMW.RequireAuthFunc(s.authHandlers.HandleLogout))
		mux.HandleFunc("POST /api/auth/refresh", s.authHandlers.HandleRefresh)
		mux.HandleFunc("GET /api/auth/me", s.authMW.RequireAuthFunc(s.authHandlers.HandleMe))
		mux.HandleFunc("DELETE /api/auth/account", s.authMW.RequireAuthFunc(s.authHandlers.HandleDeleteAccount))

		// User stats (authenticated)
		mux.HandleFunc("GET /api/user/stats", s.authMW.RequireAuthFunc(s.handleGetUserStats))
		mux.HandleFunc("POST /api/user/stats/sync", s.authMW.RequireAuthFunc(s.handleSyncUserStats))
		mux.HandleFunc("DELETE /api/user/stats", s.authMW.RequireAuthFunc(s.handleDeleteUserStats))
		mux.HandleFunc("GET /api/user/export", s.authMW.RequireAuthFunc(s.handleExportUserData))

	}

	// Leaderboard routes (v1 API) - available when database is configured
	if s.db != nil && s.db.IsConfigured() {
		mux.HandleFunc("GET /api/v1/leaderboard", s.handleGetLeaderboard)
		mux.HandleFunc("GET /api/v1/leaderboard/months", s.handleGetLeaderboardMonths)
		mux.HandleFunc("GET /api/v1/leaderboard/check", s.handleCheckLeaderboard)
		mux.HandleFunc("POST /api/v1/leaderboard/submit", s.handleSubmitLeaderboard)
		mux.HandleFunc("GET /api/v1/leaderboard/validate", s.handleValidateName)
		mux.HandleFunc("GET /api/v1/leaderboard/badge/{id}", s.handleGetBadgeSVG)
	}

	// Catch-all for undefined API routes - return JSON 404 instead of HTML
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "API endpoint not found. The database may not be configured.",
		})
	})

	// Serve static files if directory is configured
	if s.staticDir != "" {
		// Serve static files for all non-API routes
		fs := http.FileServer(http.Dir(s.staticDir))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Check if the file exists
			path := s.staticDir + r.URL.Path
			if r.URL.Path == "/" {
				path = s.staticDir + "/index.html"
			}

			// Handle /docs/ paths - serve Hugo static site directly
			if strings.HasPrefix(r.URL.Path, "/docs/") || r.URL.Path == "/docs" {
				// Redirect /docs to /docs/
				if r.URL.Path == "/docs" {
					http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
					return
				}
				// Serve docs files directly (no SPA routing)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					// Try adding index.html for directory paths
					indexPath := path + "/index.html"
					if _, err := os.Stat(indexPath); err == nil {
						http.ServeFile(w, r, indexPath)
						return
					}
					http.NotFound(w, r)
					return
				}
				fs.ServeHTTP(w, r)
				return
			}

			if _, err := os.Stat(path); os.IsNotExist(err) {
				// For SPA routing, serve index.html for non-existent paths
				http.ServeFile(w, r, s.staticDir+"/index.html")
				return
			}
			fs.ServeHTTP(w, r)
		})
	}

	return http.ListenAndServe(s.addr, mux)
}

// StartAsync starts the HTTP server in a goroutine.
func (s *Server) StartAsync() {
	go func() {
		if err := s.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()
}

// API Request/Response types

// CreateSessionRequest is the request body for POST /api/sessions
type CreateSessionRequest struct {
	PunctuationMode bool   `json:"punctuation_mode"`
	PracticeMode    string `json:"practice_mode"` // "words" or "ngrams"
}

// CreateSessionResponse is the response body for POST /api/sessions
type CreateSessionResponse struct {
	SessionID string `json:"session_id"`
}

// SessionInfo provides information about a session
type SessionInfo struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}

// ListSessionsResponse is the response body for GET /api/sessions
type ListSessionsResponse struct {
	Sessions []SessionInfo `json:"sessions"`
}

// KeystrokeRequest is the request body for POST /api/sessions/{id}/keystroke
type KeystrokeRequest struct {
	Char       string `json:"char"`
	SeekTimeMs int64  `json:"seek_time_ms,omitempty"` // Frontend-measured seek time (optional)
}

// SpaceRequest is the request body for POST /api/sessions/{id}/space
type SpaceRequest struct {
	SeekTimeMs int64 `json:"seek_time_ms,omitempty"` // Frontend-measured seek time (optional)
}

// TimingRequest is the request body for POST /api/sessions/{id}/timing
type TimingRequest struct {
	StartTimeUnixMs int64 `json:"start_time_unix_ms"` // Unix milliseconds
	EndTimeUnixMs   int64 `json:"end_time_unix_ms"`   // Unix milliseconds
	DurationMs      int64 `json:"duration_ms"`        // Duration in milliseconds
}

// KeystrokeResponse is the response body for POST /api/sessions/{id}/keystroke
type KeystrokeResponse struct {
	IsCorrect     bool `json:"is_correct"`
	TimerStarted  bool `json:"timer_started"`
	CharIndex     int  `json:"char_index"`
	RoundComplete bool `json:"round_complete"`
}

// BackspaceResponse is the response body for POST /api/sessions/{id}/backspace
type BackspaceResponse struct {
	Removed bool `json:"removed"`
}

// SpaceResponse is the response body for POST /api/sessions/{id}/space
type SpaceResponse struct {
	Advanced       bool `json:"advanced"`
	RoundComplete  bool `json:"round_complete"`
	TreatedAsError bool `json:"treated_as_error"`
}

// GameStateResponse is the response body for GET /api/sessions/{id}/state
type GameStateResponse struct {
	Words           []string `json:"words"`
	CurrentWordIdx  int      `json:"current_word_idx"`
	CurrentInput    string   `json:"current_input"`
	TimerStarted    bool     `json:"timer_started"`
	PunctuationMode bool     `json:"punctuation_mode"`
	WordNumber      int      `json:"word_number"`
	TotalWords      int      `json:"total_words"`
	LiveWPM         float64  `json:"live_wpm"`
	CurrentWord     string   `json:"current_word"`
	PreviousWord    string   `json:"previous_word"`
	NextWord        string   `json:"next_word"`
	NextWords       []string `json:"next_words"`
}

// HealthResponse is the response body for GET /api/health
type HealthResponse struct {
	Status         string `json:"status"`
	ActiveSessions int    `json:"active_sessions"`
}

// ConfigResponse is the response body for GET /api/config
type ConfigResponse struct {
	AdsenseEnabled bool     `json:"adsense_enabled"`
	AdsenseKey     string   `json:"adsense_key,omitempty"`
	AuthEnabled    bool     `json:"auth_enabled"`
	AuthProviders  []string `json:"auth_providers,omitempty"`
}

// VersionResponse is the response body for GET /api/version
type VersionResponse struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
}

// Handler implementations

func (s *Server) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	var req CreateSessionRequest
	// Decode request body (optional - use defaults if not provided)
	json.NewDecoder(r.Body).Decode(&req)

	// Create config for this session
	config := s.config
	if req.PunctuationMode {
		config.PunctuationMode = true
	}
	// Set practice mode based on request
	if req.PracticeMode == "ngrams" {
		config.PracticeMode = settings.ModeNgrams
	} else {
		config.PracticeMode = settings.ModeWords
	}

	// Create new engine for this session
	engine, err := NewEngine(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate session ID and store
	sessionID := generateSessionID()
	session := &Session{
		ID:        sessionID,
		Engine:    engine,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}

	s.mu.Lock()
	s.sessions[sessionID] = session
	s.mu.Unlock()

	resp := CreateSessionResponse{SessionID: sessionID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleDeleteSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")

	s.mu.Lock()
	_, exists := s.sessions[sessionID]
	if exists {
		delete(s.sessions, sessionID)
	}
	s.mu.Unlock()

	if !exists {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleListSessions(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	sessions := make([]SessionInfo, 0, len(s.sessions))
	for _, session := range s.sessions {
		sessions = append(sessions, SessionInfo{
			ID:        session.ID,
			CreatedAt: session.CreatedAt,
			LastUsed:  session.LastUsed,
		})
	}
	s.mu.RUnlock()

	resp := ListSessionsResponse{Sessions: sessions}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleStartRound(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	session, exists := s.getSession(sessionID)
	if !exists {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	s.mu.Lock()
	session.Engine.StartRound()
	s.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleKeystroke(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	session, exists := s.getSession(sessionID)
	if !exists {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	var req KeystrokeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	// Use timing-aware method with frontend-provided seek time
	result := session.Engine.ProcessKeystrokeWithTiming(req.Char, req.SeekTimeMs)
	s.mu.Unlock()

	resp := KeystrokeResponse{
		IsCorrect:     result.IsCorrect,
		TimerStarted:  result.TimerStarted,
		CharIndex:     result.CharIndex,
		RoundComplete: result.RoundComplete,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleBackspace(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	session, exists := s.getSession(sessionID)
	if !exists {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	s.mu.Lock()
	removed := session.Engine.ProcessBackspace()
	s.mu.Unlock()

	resp := BackspaceResponse{Removed: removed}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleClearInput(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	session, exists := s.getSession(sessionID)
	if !exists {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	s.mu.Lock()
	cleared := session.Engine.ClearInput()
	s.mu.Unlock()

	resp := BackspaceResponse{Removed: cleared}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleSpace(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	session, exists := s.getSession(sessionID)
	if !exists {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	// Parse optional seek time from request body
	var req SpaceRequest
	json.NewDecoder(r.Body).Decode(&req) // Ignore errors, use zero values

	s.mu.Lock()
	// Use timing-aware method with frontend-provided seek time
	result := session.Engine.ProcessSpaceWithTiming(req.SeekTimeMs)
	s.mu.Unlock()

	resp := SpaceResponse{
		Advanced:       result.Advanced,
		RoundComplete:  result.RoundComplete,
		TreatedAsError: result.TreatedAsError,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleGetState(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	session, exists := s.getSession(sessionID)
	if !exists {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	s.mu.RLock()
	state := session.Engine.GetGameState()
	s.mu.RUnlock()

	resp := GameStateResponse{
		Words:           state.Words,
		CurrentWordIdx:  state.CurrentWordIdx,
		CurrentInput:    state.CurrentInput,
		TimerStarted:    state.TimerStarted,
		PunctuationMode: state.PunctuationMode,
		WordNumber:      state.WordNumber,
		TotalWords:      state.TotalWords,
		LiveWPM:         state.LiveWPM,
		CurrentWord:     state.CurrentWord,
		PreviousWord:    state.PreviousWord,
		NextWord:        state.NextWord,
		NextWords:       state.NextWords,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleGetSessionStats(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	session, exists := s.getSession(sessionID)
	if !exists {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	s.mu.RLock()
	sessionStats := session.Engine.GetSessionStats()
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessionStats)
}

func (s *Server) handleGetHistoricalStats(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	session, exists := s.getSession(sessionID)
	if !exists {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	s.mu.RLock()
	historicalStats := session.Engine.GetHistoricalStats()
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(historicalStats)
}

func (s *Server) handleSaveStats(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	session, exists := s.getSession(sessionID)
	if !exists {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	s.mu.Lock()
	err := session.Engine.SaveStats()
	s.mu.Unlock()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleSubmitTiming(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("id")
	session, exists := s.getSession(sessionID)
	if !exists {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	var req TimingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert Unix milliseconds to time.Time
	startTime := time.UnixMilli(req.StartTimeUnixMs)
	endTime := time.UnixMilli(req.EndTimeUnixMs)

	s.mu.Lock()
	session.Engine.SubmitTiming(startTime, endTime, req.DurationMs)
	s.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	activeCount := len(s.sessions)
	s.mu.RUnlock()

	resp := HealthResponse{
		Status:         "healthy",
		ActiveSessions: activeCount,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	resp := VersionResponse{
		Version:   s.config.Version,
		GitCommit: s.config.GitCommit,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	resp := ConfigResponse{
		AdsenseEnabled: s.config.AdsenseEnabled,
		AdsenseKey:     s.config.AdsenseKey,
	}

	// Add auth configuration if auth is enabled
	if s.authService != nil {
		resp.AuthEnabled = s.authService.IsEnabled()
		resp.AuthProviders = s.authService.GetEnabledProviders()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// extractSessionID extracts the session ID from the URL path.
// Expected format: /api/sessions/{id}/...
func extractSessionID(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 4 && parts[1] == "api" && parts[2] == "sessions" {
		return parts[3]
	}
	return ""
}

// SessionStatsJSON is a JSON-serializable version of stats.Stats
// This is needed because some fields in stats.Stats use json:"-"
type SessionStatsJSON struct {
	WordsCompleted   int                      `json:"words_completed"`
	TotalCharacters  int                      `json:"total_characters"`
	CorrectChars     int                      `json:"correct_chars"`
	IncorrectChars   int                      `json:"incorrect_chars"`
	DurationSeconds  float64                  `json:"duration_seconds"`
	WPM              float64                  `json:"wpm"`
	Accuracy         float64                  `json:"accuracy"`
	SFBCount         int                      `json:"sfb_count"`
	SFBTotalTime     int64                    `json:"sfb_total_time"`
	HandAlternations int                      `json:"hand_alternations"`
	SameHandRuns     int                      `json:"same_hand_runs"`
	SeekTimes        []int64                  `json:"seek_times"`
	FingerStats      map[int]stats.FingerStat `json:"finger_stats"`
	HandStats        map[int]stats.HandStat   `json:"hand_stats"`
	RowStats         map[int]stats.RowStat    `json:"row_stats"`
}

// User stats handlers (authenticated)

// SyncStatsRequest is the request body for POST /api/user/stats/sync
type SyncStatsRequest struct {
	LocalStats *stats.HistoricalStats `json:"local_stats"`
}

// handleGetUserStats returns the authenticated user's stats.
// GET /api/user/stats
func (s *Server) handleGetUserStats(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	userStats, err := s.db.Stats().Get(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userStats)
}

// handleSyncUserStats merges local stats with server stats.
// POST /api/user/stats/sync
func (s *Server) handleSyncUserStats(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	var req SyncStatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get server stats
	serverStats, err := s.db.Stats().Get(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Merge local stats with server stats
	mergedStats := database.MergeStats(serverStats, req.LocalStats)

	// Save merged stats
	if err := s.db.Stats().Save(userID, mergedStats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "synced",
		"stats":  mergedStats,
	})
}

// handleDeleteUserStats deletes all stats for the authenticated user.
// DELETE /api/user/stats
func (s *Server) handleDeleteUserStats(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	if err := s.db.Stats().Delete(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// handleExportUserData exports all user data as JSON (GDPR compliance).
// GET /api/user/export
func (s *Server) handleExportUserData(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	if userID == "" {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	// Get user info
	user, err := s.db.Users().FindByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user stats
	userStats, err := s.db.Stats().Get(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	export := map[string]interface{}{
		"user":        user,
		"stats":       userStats,
		"exported_at": time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=baboon-data-export.json")
	json.NewEncoder(w).Encode(export)
}

// Leaderboard handlers

// LeaderboardResponse is the response for GET /api/v1/leaderboard
type LeaderboardResponse struct {
	Entries   []database.LeaderboardEntry `json:"entries"`
	MonthYear string                      `json:"month_year"`
}

// LeaderboardCheckResponse is the response for GET /api/v1/leaderboard/check
type LeaderboardCheckResponse struct {
	Qualifies bool `json:"qualifies"`
	Rank      int  `json:"rank,omitempty"`
}

// LeaderboardSubmitRequest is the request for POST /api/v1/leaderboard/submit
type LeaderboardSubmitRequest struct {
	DisplayName     string  `json:"display_name"`
	WPM             float64 `json:"wpm"`
	Accuracy        float64 `json:"accuracy"`
	DurationSeconds float64 `json:"duration_seconds"`
}

// handleGetLeaderboard returns the top 10 for a given month.
// GET /api/v1/leaderboard?month=2026-04
func (s *Server) handleGetLeaderboard(w http.ResponseWriter, r *http.Request) {
	monthYear := r.URL.Query().Get("month")
	if monthYear == "" {
		// Default to current month
		monthYear = time.Now().Format("2006-01")
	}

	entries, err := s.db.Leaderboard().GetTop10(monthYear)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := LeaderboardResponse{
		Entries:   entries,
		MonthYear: monthYear,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleGetLeaderboardMonths returns available months with leaderboard entries.
// GET /api/v1/leaderboard/months
func (s *Server) handleGetLeaderboardMonths(w http.ResponseWriter, r *http.Request) {
	months, err := s.db.Leaderboard().GetAvailableMonths()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Always include current month even if no entries
	currentMonth := time.Now().Format("2006-01")
	found := false
	for _, m := range months {
		if m == currentMonth {
			found = true
			break
		}
	}
	if !found {
		months = append([]string{currentMonth}, months...)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"months": months})
}

// handleCheckLeaderboard checks if a WPM and accuracy would qualify for top 10.
// GET /api/v1/leaderboard/check?wpm=75.5&accuracy=98.5
// Score is calculated as WPM * (Accuracy/100).
func (s *Server) handleCheckLeaderboard(w http.ResponseWriter, r *http.Request) {
	wpmStr := r.URL.Query().Get("wpm")
	if wpmStr == "" {
		http.Error(w, "wpm parameter required", http.StatusBadRequest)
		return
	}

	var wpm float64
	if _, err := fmt.Sscanf(wpmStr, "%f", &wpm); err != nil {
		http.Error(w, "invalid wpm value", http.StatusBadRequest)
		return
	}

	// Get accuracy parameter (default to 100 for backwards compatibility)
	accuracy := 100.0
	if accStr := r.URL.Query().Get("accuracy"); accStr != "" {
		if _, err := fmt.Sscanf(accStr, "%f", &accuracy); err != nil {
			http.Error(w, "invalid accuracy value", http.StatusBadRequest)
			return
		}
	}

	monthYear := time.Now().Format("2006-01")
	qualifies, rank, err := s.db.Leaderboard().CheckQualifies(monthYear, wpm, accuracy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := LeaderboardCheckResponse{
		Qualifies: qualifies,
		Rank:      rank,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleSubmitLeaderboard creates or updates a leaderboard entry.
// POST /api/v1/leaderboard/submit
// Supports both authenticated and anonymous submissions.
// Username is set from auth provider or "ANON" for anonymous users.
// DisplayName is the user's custom message.
func (s *Server) handleSubmitLeaderboard(w http.ResponseWriter, r *http.Request) {
	var req LeaderboardSubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get user ID and username if authenticated, otherwise generate anonymous ID
	userID := auth.GetUserID(r)
	username := "ANON"
	displayName := req.DisplayName

	if userID == "" {
		// Anonymous user - generate a session-based ID
		// Use a combination of timestamp and random to create uniqueness
		userID = fmt.Sprintf("anon-%d", time.Now().UnixNano())
	} else {
		// Authenticated user - get their platform display name as username
		if s.db != nil {
			if user, err := s.db.Users().FindByID(userID); err == nil && user.DisplayName != "" {
				username = user.DisplayName
			}
		}
		// Still "ANON"? Use their email prefix
		if username == "ANON" {
			if user, err := s.db.Users().FindByID(userID); err == nil && user.Email != "" {
				parts := strings.Split(user.Email, "@")
				if len(parts) > 0 && parts[0] != "" {
					username = strings.ToUpper(parts[0])
					if len(username) > 10 {
						username = username[:10]
					}
				}
			}
		}
	}

	// If no custom message provided, use a default
	if displayName == "" {
		displayName = "HELLO WORLD"
	}

	// Validate display name (the message)
	validation := filter.Validate(displayName)
	if !validation.Valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":         validation.Reason,
			"invalid_chars": validation.InvalidChars,
		})
		return
	}

	entry := &database.LeaderboardEntry{
		UserID:          userID,
		Username:        username,
		DisplayName:     displayName,
		WPM:             req.WPM,
		Accuracy:        req.Accuracy,
		DurationSeconds: req.DurationSeconds,
		MonthYear:       time.Now().Format("2006-01"),
	}

	submitted, err := s.db.Leaderboard().Submit(entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate badge URL
	baseURL := s.config.BaseURL
	if baseURL == "" {
		baseURL = "https://baboon.kartoza.com"
	}
	submitted.BadgeImageURL = fmt.Sprintf("%s/api/v1/leaderboard/badge/%s", baseURL, submitted.ID)

	// Update badge URL in database
	s.db.Leaderboard().UpdateBadgeURL(submitted.ID, submitted.BadgeImageURL)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(submitted)
}

// handleValidateName validates a display name for profanity.
// GET /api/v1/leaderboard/validate?name=PLAYERNAME
func (s *Server) handleValidateName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name parameter required", http.StatusBadRequest)
		return
	}

	validation := filter.Validate(name)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(validation)
}

// handleGetBadgeSVG generates and returns a badge SVG for a leaderboard entry.
// GET /api/v1/leaderboard/badge/{id}
func (s *Server) handleGetBadgeSVG(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// Remove .svg extension if present
	id = strings.TrimSuffix(id, ".svg")

	entry, err := s.db.Leaderboard().GetByID(id)
	if err != nil {
		if err == database.ErrUserNotFound {
			http.Error(w, "entry not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	baseURL := s.config.BaseURL
	if baseURL == "" {
		baseURL = "https://baboon.kartoza.com"
	}

	badgeEntry := badge.Entry{
		DisplayName: entry.DisplayName,
		WPM:         entry.WPM,
		Accuracy:    entry.Accuracy,
		Rank:        entry.Rank,
		MonthYear:   entry.MonthYear,
		SiteURL:     strings.TrimPrefix(strings.TrimPrefix(baseURL, "https://"), "http://"),
	}

	svg := badge.GenerateSVG(badgeEntry)

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write([]byte(svg))
}
