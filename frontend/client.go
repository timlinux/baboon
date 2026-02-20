package frontend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/timlinux/baboon/backend"
	"github.com/timlinux/baboon/stats"
)

// Client is a REST API client that implements the backend.GameAPI interface.
// It communicates with the backend server via HTTP using a session ID.
type Client struct {
	baseURL         string
	sessionID       string
	punctuationMode bool
	httpClient      *http.Client

	// Cached state to reduce HTTP calls during rendering
	cachedState      *backend.GameState
	cachedSession    *stats.Stats
	cachedHistorical *stats.HistoricalStats
}

// NewClient creates a new REST API client.
func NewClient(baseURL string, punctuationMode bool) *Client {
	return &Client{
		baseURL:         baseURL,
		punctuationMode: punctuationMode,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetSessionID returns the current session ID.
func (c *Client) GetSessionID() string {
	return c.sessionID
}

// WaitForServer waits until the server is ready, with a timeout.
func (c *Client) WaitForServer(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := c.httpClient.Get(c.baseURL + "/api/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(10 * time.Millisecond)
	}
	return fmt.Errorf("server not ready after %v", timeout)
}

// CreateSession creates a new session on the server.
func (c *Client) CreateSession() error {
	body, _ := json.Marshal(map[string]bool{"punctuation_mode": c.punctuationMode})
	req, _ := http.NewRequest("POST", c.baseURL+"/api/sessions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create session: status %d", resp.StatusCode)
	}

	var result struct {
		SessionID string `json:"session_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode session response: %w", err)
	}

	c.sessionID = result.SessionID
	return nil
}

// DeleteSession deletes the current session from the server.
func (c *Client) DeleteSession() error {
	if c.sessionID == "" {
		return nil
	}

	req, _ := http.NewRequest("DELETE", c.baseURL+"/api/sessions/"+c.sessionID, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	c.sessionID = ""
	return nil
}

// sessionURL returns the base URL for session-specific endpoints.
func (c *Client) sessionURL() string {
	return c.baseURL + "/api/sessions/" + c.sessionID
}

// StartRound starts a new round via the REST API.
func (c *Client) StartRound() {
	if c.sessionID == "" {
		return
	}

	req, _ := http.NewRequest("POST", c.sessionURL()+"/round", nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()

	// Invalidate cache
	c.cachedState = nil
	c.cachedSession = nil
}

// ProcessKeystroke sends a keystroke to the server (legacy, no timing).
func (c *Client) ProcessKeystroke(char string) backend.KeystrokeResult {
	return c.ProcessKeystrokeWithTiming(char, 0)
}

// ProcessKeystrokeWithTiming sends a keystroke with frontend-measured seek time.
func (c *Client) ProcessKeystrokeWithTiming(char string, seekTimeMs int64) backend.KeystrokeResult {
	if c.sessionID == "" {
		return backend.KeystrokeResult{}
	}

	body, _ := json.Marshal(map[string]interface{}{
		"char":         char,
		"seek_time_ms": seekTimeMs,
	})
	req, _ := http.NewRequest("POST", c.sessionURL()+"/keystroke", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return backend.KeystrokeResult{}
	}
	defer resp.Body.Close()

	var result struct {
		IsCorrect    bool `json:"is_correct"`
		TimerStarted bool `json:"timer_started"`
		CharIndex    int  `json:"char_index"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	// Invalidate cache
	c.cachedState = nil
	c.cachedSession = nil

	return backend.KeystrokeResult{
		IsCorrect:    result.IsCorrect,
		TimerStarted: result.TimerStarted,
		CharIndex:    result.CharIndex,
	}
}

// ProcessBackspace sends a backspace to the server.
func (c *Client) ProcessBackspace() bool {
	if c.sessionID == "" {
		return false
	}

	req, _ := http.NewRequest("POST", c.sessionURL()+"/backspace", nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Removed bool `json:"removed"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	// Invalidate cache
	c.cachedState = nil

	return result.Removed
}

// ProcessSpace sends a space to the server (legacy, no timing).
func (c *Client) ProcessSpace() backend.SpaceResult {
	return c.ProcessSpaceWithTiming(0)
}

// ProcessSpaceWithTiming sends a space with frontend-measured seek time.
func (c *Client) ProcessSpaceWithTiming(seekTimeMs int64) backend.SpaceResult {
	if c.sessionID == "" {
		return backend.SpaceResult{}
	}

	body, _ := json.Marshal(map[string]int64{"seek_time_ms": seekTimeMs})
	req, _ := http.NewRequest("POST", c.sessionURL()+"/space", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return backend.SpaceResult{}
	}
	defer resp.Body.Close()

	var result struct {
		Advanced       bool `json:"advanced"`
		RoundComplete  bool `json:"round_complete"`
		TreatedAsError bool `json:"treated_as_error"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	// Invalidate cache
	c.cachedState = nil
	c.cachedSession = nil
	c.cachedHistorical = nil

	return backend.SpaceResult{
		Advanced:       result.Advanced,
		RoundComplete:  result.RoundComplete,
		TreatedAsError: result.TreatedAsError,
	}
}

// SubmitTiming sends the final timing data from the frontend when a round completes.
func (c *Client) SubmitTiming(startTime, endTime time.Time, durationMs int64) {
	if c.sessionID == "" {
		return
	}

	body, _ := json.Marshal(map[string]int64{
		"start_time_unix_ms": startTime.UnixMilli(),
		"end_time_unix_ms":   endTime.UnixMilli(),
		"duration_ms":        durationMs,
	})
	req, _ := http.NewRequest("POST", c.sessionURL()+"/timing", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()

	// Invalidate cache
	c.cachedSession = nil
	c.cachedHistorical = nil
}

// GetGameState fetches the current game state from the server.
func (c *Client) GetGameState() backend.GameState {
	if c.sessionID == "" {
		return backend.GameState{}
	}

	resp, err := c.httpClient.Get(c.sessionURL() + "/state")
	if err != nil {
		if c.cachedState != nil {
			return *c.cachedState
		}
		return backend.GameState{}
	}
	defer resp.Body.Close()

	var state struct {
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
	json.NewDecoder(resp.Body).Decode(&state)

	result := backend.GameState{
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

	c.cachedState = &result
	return result
}

// GetSessionStats fetches the session statistics from the server.
func (c *Client) GetSessionStats() *stats.Stats {
	if c.sessionID == "" {
		return &stats.Stats{}
	}

	resp, err := c.httpClient.Get(c.sessionURL() + "/stats/session")
	if err != nil {
		if c.cachedSession != nil {
			return c.cachedSession
		}
		return &stats.Stats{}
	}
	defer resp.Body.Close()

	var sessionStats stats.Stats
	json.NewDecoder(resp.Body).Decode(&sessionStats)

	c.cachedSession = &sessionStats
	return &sessionStats
}

// GetHistoricalStats fetches the historical statistics from the server.
func (c *Client) GetHistoricalStats() *stats.HistoricalStats {
	if c.sessionID == "" {
		return &stats.HistoricalStats{}
	}

	resp, err := c.httpClient.Get(c.sessionURL() + "/stats/historical")
	if err != nil {
		if c.cachedHistorical != nil {
			return c.cachedHistorical
		}
		return &stats.HistoricalStats{}
	}
	defer resp.Body.Close()

	var historicalStats stats.HistoricalStats
	json.NewDecoder(resp.Body).Decode(&historicalStats)

	// Ensure maps are initialised
	if historicalStats.LetterAccuracy == nil {
		historicalStats.LetterAccuracy = make(map[string]stats.LetterStats)
	}
	if historicalStats.LetterSeekTime == nil {
		historicalStats.LetterSeekTime = make(map[string]stats.LetterSeekStats)
	}
	if historicalStats.BigramSeekTime == nil {
		historicalStats.BigramSeekTime = make(map[string]stats.BigramSeekStats)
	}
	if historicalStats.FingerStats == nil {
		historicalStats.FingerStats = make(map[int]stats.FingerStat)
	}
	if historicalStats.HandStats == nil {
		historicalStats.HandStats = make(map[int]stats.HandStat)
	}
	if historicalStats.RowStats == nil {
		historicalStats.RowStats = make(map[int]stats.RowStat)
	}
	if historicalStats.ErrorSubstitution == nil {
		historicalStats.ErrorSubstitution = make(map[string]map[string]int)
	}

	c.cachedHistorical = &historicalStats
	return &historicalStats
}

// SaveStats saves the statistics via the server.
func (c *Client) SaveStats() error {
	if c.sessionID == "" {
		return fmt.Errorf("no session")
	}

	req, _ := http.NewRequest("POST", c.sessionURL()+"/save", nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("save failed with status %d", resp.StatusCode)
	}

	return nil
}

// Ensure Client implements backend.GameAPI
var _ backend.GameAPI = (*Client)(nil)
