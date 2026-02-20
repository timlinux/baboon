package backend

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

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
	config   Config
	sessions map[string]*Session
	mu       sync.RWMutex
	addr     string
}

// NewServer creates a new REST API server with the given configuration.
func NewServer(config Config, addr string) (*Server, error) {
	return &Server{
		config:   config,
		sessions: make(map[string]*Session),
		addr:     addr,
	}, nil
}

// GetAddr returns the server address.
func (s *Server) GetAddr() string {
	return s.addr
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
	PunctuationMode bool `json:"punctuation_mode"`
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
	IsCorrect    bool `json:"is_correct"`
	TimerStarted bool `json:"timer_started"`
	CharIndex    int  `json:"char_index"`
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
}

// HealthResponse is the response body for GET /api/health
type HealthResponse struct {
	Status        string `json:"status"`
	ActiveSessions int   `json:"active_sessions"`
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
		IsCorrect:    result.IsCorrect,
		TimerStarted: result.TimerStarted,
		CharIndex:    result.CharIndex,
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
		Status:        "healthy",
		ActiveSessions: activeCount,
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
