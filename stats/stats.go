package stats

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Stats represents typing statistics for a session
type Stats struct {
	WordsCompleted  int                        `json:"words_completed"`
	TotalCharacters int                        `json:"total_characters"`
	CorrectChars    int                        `json:"correct_chars"`
	IncorrectChars  int                        `json:"incorrect_chars"`
	StartTime       time.Time                  `json:"start_time"`
	EndTime         time.Time                  `json:"end_time"`
	Duration        time.Duration              `json:"duration"`
	WPM             float64                    `json:"wpm"`
	Accuracy        float64                    `json:"accuracy"`
	LetterAccuracy  map[string]LetterStats     `json:"-"` // Per-letter accuracy for this session
	LetterSeekTime  map[string]LetterSeekStats `json:"-"` // Per-letter seek time for this session
	BigramSeekTime  map[string]BigramSeekStats `json:"-"` // Per-bigram seek time for this session
	LastKeyTime     time.Time                  `json:"-"` // Time of last keystroke for seek time calc
	LastLetter      string                     `json:"-"` // Last letter typed (for bigram tracking)
}

// LetterStats tracks per-letter accuracy
type LetterStats struct {
	Presented int `json:"presented"` // Number of times this letter was presented
	Correct   int `json:"correct"`   // Number of times typed correctly
}

// LetterSeekStats tracks per-letter seek time (time to find and press the key)
type LetterSeekStats struct {
	TotalTimeMs int64 `json:"total_time_ms"` // Total time in milliseconds
	Count       int   `json:"count"`         // Number of measurements
}

// AverageMs returns the average seek time in milliseconds
func (s LetterSeekStats) AverageMs() float64 {
	if s.Count == 0 {
		return 0
	}
	return float64(s.TotalTimeMs) / float64(s.Count)
}

// BigramSeekStats tracks seek time for letter pairs (e.g., "th", "he", "in")
type BigramSeekStats struct {
	TotalTimeMs int64 `json:"total_time_ms"` // Total time in milliseconds
	Count       int   `json:"count"`         // Number of measurements
}

// AverageMs returns the average seek time in milliseconds for the bigram
func (s BigramSeekStats) AverageMs() float64 {
	if s.Count == 0 {
		return 0
	}
	return float64(s.TotalTimeMs) / float64(s.Count)
}

// HistoricalStats stores best performance data
type HistoricalStats struct {
	BestWPM         float64                    `json:"best_wpm"`
	BestAccuracy    float64                    `json:"best_accuracy"`
	BestTime        float64                    `json:"best_time"` // Best (fastest) time in seconds
	TotalWPM        float64                    `json:"total_wpm"`
	TotalAccuracy   float64                    `json:"total_accuracy"`
	TotalTime       float64                    `json:"total_time"` // Total time across all sessions
	TotalSessions   int                        `json:"total_sessions"`
	LastSessionDate time.Time                  `json:"last_session_date"`
	LetterAccuracy  map[string]LetterStats     `json:"letter_accuracy"`  // Per-letter accuracy tracking
	LetterSeekTime  map[string]LetterSeekStats `json:"letter_seek_time"` // Per-letter seek time tracking
	BigramSeekTime  map[string]BigramSeekStats `json:"bigram_seek_time"` // Per-bigram seek time tracking
}

// RecordLetterPresented records that a letter was presented to the user
func (s *Stats) RecordLetterPresented(letter string) {
	if s.LetterAccuracy == nil {
		s.LetterAccuracy = make(map[string]LetterStats)
	}
	stats := s.LetterAccuracy[letter]
	stats.Presented++
	s.LetterAccuracy[letter] = stats
}

// RecordLetterCorrect records that a letter was typed correctly
func (s *Stats) RecordLetterCorrect(letter string) {
	if s.LetterAccuracy == nil {
		s.LetterAccuracy = make(map[string]LetterStats)
	}
	stats := s.LetterAccuracy[letter]
	stats.Correct++
	s.LetterAccuracy[letter] = stats
}

// RecordLetterSeekTime records the time taken to type a letter
func (s *Stats) RecordLetterSeekTime(letter string, durationMs int64) {
	if s.LetterSeekTime == nil {
		s.LetterSeekTime = make(map[string]LetterSeekStats)
	}
	stats := s.LetterSeekTime[letter]
	stats.TotalTimeMs += durationMs
	stats.Count++
	s.LetterSeekTime[letter] = stats
}

// RecordBigramSeekTime records the time taken to type a letter pair (bigram)
func (s *Stats) RecordBigramSeekTime(bigram string, durationMs int64) {
	if s.BigramSeekTime == nil {
		s.BigramSeekTime = make(map[string]BigramSeekStats)
	}
	stats := s.BigramSeekTime[bigram]
	stats.TotalTimeMs += durationMs
	stats.Count++
	s.BigramSeekTime[bigram] = stats
}

// Calculate computes WPM and accuracy from raw stats
func (s *Stats) Calculate() {
	s.EndTime = time.Now()
	s.Duration = s.EndTime.Sub(s.StartTime)

	// WPM calculation: (characters / 5) / minutes
	// Standard word length is 5 characters
	minutes := s.Duration.Minutes()
	if minutes > 0 {
		s.WPM = (float64(s.CorrectChars) / 5.0) / minutes
	}

	// Accuracy calculation
	if s.TotalCharacters > 0 {
		s.Accuracy = (float64(s.CorrectChars) / float64(s.TotalCharacters)) * 100
	}
}

// GetStatsPath returns the path to the stats file
func GetStatsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	statsDir := filepath.Join(homeDir, ".config", "baboon")
	if err := os.MkdirAll(statsDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(statsDir, "stats.json"), nil
}

// LoadHistoricalStats loads historical stats from disk
func LoadHistoricalStats() (*HistoricalStats, error) {
	path, err := GetStatsPath()
	if err != nil {
		return &HistoricalStats{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &HistoricalStats{
				LetterAccuracy: make(map[string]LetterStats),
				LetterSeekTime: make(map[string]LetterSeekStats),
				BigramSeekTime: make(map[string]BigramSeekStats),
			}, nil
		}
		return &HistoricalStats{}, err
	}

	var stats HistoricalStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return &HistoricalStats{}, err
	}

	// Initialize maps if nil
	if stats.LetterAccuracy == nil {
		stats.LetterAccuracy = make(map[string]LetterStats)
	}
	if stats.LetterSeekTime == nil {
		stats.LetterSeekTime = make(map[string]LetterSeekStats)
	}
	if stats.BigramSeekTime == nil {
		stats.BigramSeekTime = make(map[string]BigramSeekStats)
	}

	// Validate and fix corrupted averages from older versions
	stats.validateAndFix()

	return &stats, nil
}

// validateAndFix checks for inconsistent data and resets if needed
func (h *HistoricalStats) validateAndFix() {
	if h.TotalSessions == 0 {
		return
	}

	// Check if averages are inconsistent (e.g., average time < best time is impossible
	// if we have more than 1 session, or totals are 0 when they shouldn't be)
	needsReset := false

	// If we have sessions but no totals tracked, data is from old version
	if h.TotalSessions > 0 && h.TotalWPM == 0 && h.BestWPM > 0 {
		needsReset = true
	}
	if h.TotalSessions > 0 && h.TotalTime == 0 && h.BestTime > 0 {
		needsReset = true
	}
	if h.TotalSessions > 0 && h.TotalAccuracy == 0 && h.BestAccuracy > 0 {
		needsReset = true
	}

	// If average is less than best (impossible), data is corrupted
	if h.TotalSessions > 1 {
		avgWPM := h.TotalWPM / float64(h.TotalSessions)
		if avgWPM > 0 && avgWPM < h.BestWPM*0.5 {
			// Average WPM less than half of best is suspicious
			needsReset = true
		}
	}

	if needsReset {
		// Reset totals based on best values as estimates
		// This gives a reasonable starting point
		h.TotalWPM = h.BestWPM * float64(h.TotalSessions)
		h.TotalAccuracy = h.BestAccuracy * float64(h.TotalSessions)
		h.TotalTime = h.BestTime * float64(h.TotalSessions)
	}
}

// SaveHistoricalStats saves historical stats to disk
func SaveHistoricalStats(stats *HistoricalStats) error {
	path, err := GetStatsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// UpdateHistorical updates historical stats with new session data
func (h *HistoricalStats) UpdateHistorical(session *Stats) {
	h.TotalSessions++
	h.LastSessionDate = time.Now()

	// Update totals for averages
	h.TotalWPM += session.WPM
	h.TotalAccuracy += session.Accuracy
	h.TotalTime += session.Duration.Seconds()

	// Update bests
	if session.WPM > h.BestWPM {
		h.BestWPM = session.WPM
	}
	if session.Accuracy > h.BestAccuracy {
		h.BestAccuracy = session.Accuracy
	}
	// Best time is the fastest (lowest) time
	if h.BestTime == 0 || session.Duration.Seconds() < h.BestTime {
		h.BestTime = session.Duration.Seconds()
	}

	// Merge session letter accuracy into historical
	if h.LetterAccuracy == nil {
		h.LetterAccuracy = make(map[string]LetterStats)
	}
	for letter, sessionStats := range session.LetterAccuracy {
		histStats := h.LetterAccuracy[letter]
		histStats.Presented += sessionStats.Presented
		histStats.Correct += sessionStats.Correct
		h.LetterAccuracy[letter] = histStats
	}

	// Merge session letter seek time into historical
	if h.LetterSeekTime == nil {
		h.LetterSeekTime = make(map[string]LetterSeekStats)
	}
	for letter, sessionStats := range session.LetterSeekTime {
		histStats := h.LetterSeekTime[letter]
		histStats.TotalTimeMs += sessionStats.TotalTimeMs
		histStats.Count += sessionStats.Count
		h.LetterSeekTime[letter] = histStats
	}

	// Merge session bigram seek time into historical
	if h.BigramSeekTime == nil {
		h.BigramSeekTime = make(map[string]BigramSeekStats)
	}
	for bigram, sessionStats := range session.BigramSeekTime {
		histStats := h.BigramSeekTime[bigram]
		histStats.TotalTimeMs += sessionStats.TotalTimeMs
		histStats.Count += sessionStats.Count
		h.BigramSeekTime[bigram] = histStats
	}
}

// AverageWPM returns the average WPM across all sessions
func (h *HistoricalStats) AverageWPM() float64 {
	if h.TotalSessions == 0 {
		return 0
	}
	return h.TotalWPM / float64(h.TotalSessions)
}

// AverageAccuracy returns the average accuracy across all sessions
func (h *HistoricalStats) AverageAccuracy() float64 {
	if h.TotalSessions == 0 {
		return 0
	}
	return h.TotalAccuracy / float64(h.TotalSessions)
}

// AverageTime returns the average time across all sessions in seconds
func (h *HistoricalStats) AverageTime() float64 {
	if h.TotalSessions == 0 {
		return 0
	}
	return h.TotalTime / float64(h.TotalSessions)
}
