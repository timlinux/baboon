package stats

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Stats represents typing statistics for a session
type Stats struct {
	WordsCompleted   int           `json:"words_completed"`
	TotalCharacters  int           `json:"total_characters"`
	CorrectChars     int           `json:"correct_chars"`
	IncorrectChars   int           `json:"incorrect_chars"`
	StartTime        time.Time     `json:"start_time"`
	EndTime          time.Time     `json:"end_time"`
	Duration         time.Duration `json:"duration"`
	WPM              float64       `json:"wpm"`
	Accuracy         float64       `json:"accuracy"`
}

// HistoricalStats stores best performance data
type HistoricalStats struct {
	BestWPM         float64   `json:"best_wpm"`
	BestAccuracy    float64   `json:"best_accuracy"`
	TotalSessions   int       `json:"total_sessions"`
	LastSessionDate time.Time `json:"last_session_date"`
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
			return &HistoricalStats{}, nil
		}
		return &HistoricalStats{}, err
	}

	var stats HistoricalStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return &HistoricalStats{}, err
	}

	return &stats, nil
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

	if session.WPM > h.BestWPM {
		h.BestWPM = session.WPM
	}
	if session.Accuracy > h.BestAccuracy {
		h.BestAccuracy = session.Accuracy
	}
}
