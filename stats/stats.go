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

	// Advanced typing theory stats
	FingerStats       map[int]FingerStat         `json:"-"` // Per-finger accuracy and speed
	HandStats         map[int]HandStat           `json:"-"` // Per-hand statistics
	RowStats          map[int]RowStat            `json:"-"` // Per-row statistics
	ErrorSubstitution map[string]map[string]int  `json:"-"` // Expected -> Typed -> Count
	SFBCount          int                        `json:"-"` // Same-finger bigram count
	SFBTotalTime      int64                      `json:"-"` // Total time for SFBs
	HandAlternations  int                        `json:"-"` // Count of hand alternations
	SameHandRuns      int                        `json:"-"` // Count of same-hand consecutive pairs
	SeekTimes         []int64                    `json:"-"` // All seek times for variance calculation
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

// FingerStat tracks per-finger statistics
type FingerStat struct {
	Presented   int   `json:"presented"`     // Times this finger was needed
	Correct     int   `json:"correct"`       // Correct keypresses
	TotalTimeMs int64 `json:"total_time_ms"` // Total seek time in ms
	Count       int   `json:"count"`         // Number of timed keypresses
}

// Accuracy returns accuracy as a percentage
func (f FingerStat) Accuracy() float64 {
	if f.Presented == 0 {
		return 0
	}
	return (float64(f.Correct) / float64(f.Presented)) * 100
}

// AverageMs returns average seek time in milliseconds
func (f FingerStat) AverageMs() float64 {
	if f.Count == 0 {
		return 0
	}
	return float64(f.TotalTimeMs) / float64(f.Count)
}

// HandStat tracks per-hand statistics
type HandStat struct {
	Presented   int   `json:"presented"`     // Times this hand was needed
	Correct     int   `json:"correct"`       // Correct keypresses
	TotalTimeMs int64 `json:"total_time_ms"` // Total seek time in ms
	Count       int   `json:"count"`         // Number of timed keypresses
}

// Accuracy returns accuracy as a percentage
func (h HandStat) Accuracy() float64 {
	if h.Presented == 0 {
		return 0
	}
	return (float64(h.Correct) / float64(h.Presented)) * 100
}

// AverageMs returns average seek time in milliseconds
func (h HandStat) AverageMs() float64 {
	if h.Count == 0 {
		return 0
	}
	return float64(h.TotalTimeMs) / float64(h.Count)
}

// RowStat tracks per-row statistics
type RowStat struct {
	Presented   int   `json:"presented"`     // Times this row was needed
	Correct     int   `json:"correct"`       // Correct keypresses
	TotalTimeMs int64 `json:"total_time_ms"` // Total seek time in ms
	Count       int   `json:"count"`         // Number of timed keypresses
}

// Accuracy returns accuracy as a percentage
func (r RowStat) Accuracy() float64 {
	if r.Presented == 0 {
		return 0
	}
	return (float64(r.Correct) / float64(r.Presented)) * 100
}

// AverageMs returns average seek time in milliseconds
func (r RowStat) AverageMs() float64 {
	if r.Count == 0 {
		return 0
	}
	return float64(r.TotalTimeMs) / float64(r.Count)
}

// SFBStats tracks same-finger bigram statistics
type SFBStats struct {
	Count       int   `json:"count"`         // Number of SFBs encountered
	TotalTimeMs int64 `json:"total_time_ms"` // Total time for SFBs
}

// AverageMs returns average SFB seek time
func (s SFBStats) AverageMs() float64 {
	if s.Count == 0 {
		return 0
	}
	return float64(s.TotalTimeMs) / float64(s.Count)
}

// RhythmStats tracks typing rhythm consistency
type RhythmStats struct {
	TotalSeekTimeMs  int64   `json:"total_seek_time_ms"`
	TotalSeekTimeSq  float64 `json:"total_seek_time_sq"` // Sum of squared seek times
	Count            int     `json:"count"`
	LastVariance     float64 `json:"last_variance"` // Last calculated variance
}

// Variance returns the variance of seek times
func (r RhythmStats) Variance() float64 {
	if r.Count < 2 {
		return 0
	}
	mean := float64(r.TotalSeekTimeMs) / float64(r.Count)
	return (r.TotalSeekTimeSq / float64(r.Count)) - (mean * mean)
}

// StdDev returns the standard deviation of seek times
func (r RhythmStats) StdDev() float64 {
	variance := r.Variance()
	if variance <= 0 {
		return 0
	}
	return sqrt(variance)
}

// sqrt helper using Newton's method (avoid math import in this file)
func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
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

	// Advanced typing theory stats
	FingerStats       map[int]FingerStat        `json:"finger_stats"`       // Per-finger accuracy and speed
	HandStats         map[int]HandStat          `json:"hand_stats"`         // Per-hand statistics
	RowStats          map[int]RowStat           `json:"row_stats"`          // Per-row statistics
	ErrorSubstitution map[string]map[string]int `json:"error_substitution"` // Expected -> Typed -> Count
	SFBStats          SFBStats                  `json:"sfb_stats"`          // Same-finger bigram stats
	HandAlternations  int                       `json:"hand_alternations"`  // Total hand alternations
	SameHandRuns      int                       `json:"same_hand_runs"`     // Total same-hand consecutive pairs
	RhythmStats       RhythmStats               `json:"rhythm_stats"`       // Rhythm consistency tracking
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

// RecordFingerPresented records that a key was presented for a specific finger
func (s *Stats) RecordFingerPresented(finger int) {
	if s.FingerStats == nil {
		s.FingerStats = make(map[int]FingerStat)
	}
	stat := s.FingerStats[finger]
	stat.Presented++
	s.FingerStats[finger] = stat
}

// RecordFingerCorrect records a correct keypress for a specific finger with optional seek time
func (s *Stats) RecordFingerCorrect(finger int, seekTimeMs int64) {
	if s.FingerStats == nil {
		s.FingerStats = make(map[int]FingerStat)
	}
	stat := s.FingerStats[finger]
	stat.Correct++
	if seekTimeMs > 0 {
		stat.TotalTimeMs += seekTimeMs
		stat.Count++
	}
	s.FingerStats[finger] = stat
}

// RecordHandPresented records that a key was presented for a specific hand
func (s *Stats) RecordHandPresented(hand int) {
	if s.HandStats == nil {
		s.HandStats = make(map[int]HandStat)
	}
	stat := s.HandStats[hand]
	stat.Presented++
	s.HandStats[hand] = stat
}

// RecordHandCorrect records a correct keypress for a specific hand with optional seek time
func (s *Stats) RecordHandCorrect(hand int, seekTimeMs int64) {
	if s.HandStats == nil {
		s.HandStats = make(map[int]HandStat)
	}
	stat := s.HandStats[hand]
	stat.Correct++
	if seekTimeMs > 0 {
		stat.TotalTimeMs += seekTimeMs
		stat.Count++
	}
	s.HandStats[hand] = stat
}

// RecordRowPresented records that a key was presented for a specific row
func (s *Stats) RecordRowPresented(row int) {
	if s.RowStats == nil {
		s.RowStats = make(map[int]RowStat)
	}
	stat := s.RowStats[row]
	stat.Presented++
	s.RowStats[row] = stat
}

// RecordRowCorrect records a correct keypress for a specific row with optional seek time
func (s *Stats) RecordRowCorrect(row int, seekTimeMs int64) {
	if s.RowStats == nil {
		s.RowStats = make(map[int]RowStat)
	}
	stat := s.RowStats[row]
	stat.Correct++
	if seekTimeMs > 0 {
		stat.TotalTimeMs += seekTimeMs
		stat.Count++
	}
	s.RowStats[row] = stat
}

// RecordErrorSubstitution records when an expected letter was mistyped as another letter
func (s *Stats) RecordErrorSubstitution(expected, typed string) {
	if s.ErrorSubstitution == nil {
		s.ErrorSubstitution = make(map[string]map[string]int)
	}
	if s.ErrorSubstitution[expected] == nil {
		s.ErrorSubstitution[expected] = make(map[string]int)
	}
	s.ErrorSubstitution[expected][typed]++
}

// RecordSFB records a same-finger bigram with its seek time
func (s *Stats) RecordSFB(seekTimeMs int64) {
	s.SFBCount++
	s.SFBTotalTime += seekTimeMs
}

// RecordHandTransition records whether the hand alternated or stayed the same
func (s *Stats) RecordHandTransition(alternated bool) {
	if alternated {
		s.HandAlternations++
	} else {
		s.SameHandRuns++
	}
}

// RecordSeekTime records a seek time for rhythm variance calculation
func (s *Stats) RecordSeekTime(seekTimeMs int64) {
	s.SeekTimes = append(s.SeekTimes, seekTimeMs)
}

// CalculateRhythmVariance calculates the rhythm variance from collected seek times
func (s *Stats) CalculateRhythmVariance() float64 {
	if len(s.SeekTimes) < 2 {
		return 0
	}
	var sum, sumSq float64
	for _, t := range s.SeekTimes {
		sum += float64(t)
		sumSq += float64(t) * float64(t)
	}
	mean := sum / float64(len(s.SeekTimes))
	variance := (sumSq / float64(len(s.SeekTimes))) - (mean * mean)
	if variance < 0 {
		return 0
	}
	return variance
}

// CalculateRhythmStdDev calculates the rhythm standard deviation from collected seek times
func (s *Stats) CalculateRhythmStdDev() float64 {
	variance := s.CalculateRhythmVariance()
	return sqrt(variance)
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
				LetterAccuracy:    make(map[string]LetterStats),
				LetterSeekTime:    make(map[string]LetterSeekStats),
				BigramSeekTime:    make(map[string]BigramSeekStats),
				FingerStats:       make(map[int]FingerStat),
				HandStats:         make(map[int]HandStat),
				RowStats:          make(map[int]RowStat),
				ErrorSubstitution: make(map[string]map[string]int),
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
	if stats.FingerStats == nil {
		stats.FingerStats = make(map[int]FingerStat)
	}
	if stats.HandStats == nil {
		stats.HandStats = make(map[int]HandStat)
	}
	if stats.RowStats == nil {
		stats.RowStats = make(map[int]RowStat)
	}
	if stats.ErrorSubstitution == nil {
		stats.ErrorSubstitution = make(map[string]map[string]int)
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

	// Merge session finger stats into historical
	if h.FingerStats == nil {
		h.FingerStats = make(map[int]FingerStat)
	}
	for finger, sessionStats := range session.FingerStats {
		histStats := h.FingerStats[finger]
		histStats.Presented += sessionStats.Presented
		histStats.Correct += sessionStats.Correct
		histStats.TotalTimeMs += sessionStats.TotalTimeMs
		histStats.Count += sessionStats.Count
		h.FingerStats[finger] = histStats
	}

	// Merge session hand stats into historical
	if h.HandStats == nil {
		h.HandStats = make(map[int]HandStat)
	}
	for hand, sessionStats := range session.HandStats {
		histStats := h.HandStats[hand]
		histStats.Presented += sessionStats.Presented
		histStats.Correct += sessionStats.Correct
		histStats.TotalTimeMs += sessionStats.TotalTimeMs
		histStats.Count += sessionStats.Count
		h.HandStats[hand] = histStats
	}

	// Merge session row stats into historical
	if h.RowStats == nil {
		h.RowStats = make(map[int]RowStat)
	}
	for row, sessionStats := range session.RowStats {
		histStats := h.RowStats[row]
		histStats.Presented += sessionStats.Presented
		histStats.Correct += sessionStats.Correct
		histStats.TotalTimeMs += sessionStats.TotalTimeMs
		histStats.Count += sessionStats.Count
		h.RowStats[row] = histStats
	}

	// Merge session error substitution into historical
	if h.ErrorSubstitution == nil {
		h.ErrorSubstitution = make(map[string]map[string]int)
	}
	for expected, typedMap := range session.ErrorSubstitution {
		if h.ErrorSubstitution[expected] == nil {
			h.ErrorSubstitution[expected] = make(map[string]int)
		}
		for typed, count := range typedMap {
			h.ErrorSubstitution[expected][typed] += count
		}
	}

	// Merge SFB stats
	h.SFBStats.Count += session.SFBCount
	h.SFBStats.TotalTimeMs += session.SFBTotalTime

	// Merge hand alternation stats
	h.HandAlternations += session.HandAlternations
	h.SameHandRuns += session.SameHandRuns

	// Merge rhythm stats
	for _, seekTime := range session.SeekTimes {
		h.RhythmStats.TotalSeekTimeMs += seekTime
		h.RhythmStats.TotalSeekTimeSq += float64(seekTime) * float64(seekTime)
		h.RhythmStats.Count++
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
