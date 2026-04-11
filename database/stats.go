package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/timlinux/baboon/stats"
)

// StatsRepository provides access to user stats data.
type StatsRepository struct {
	db *DB
}

// Stats returns the stats repository.
func (db *DB) Stats() *StatsRepository {
	return &StatsRepository{db: db}
}

// Get retrieves the historical stats for a user.
func (r *StatsRepository) Get(userID string) (*stats.HistoricalStats, error) {
	if r.db == nil || r.db.conn == nil {
		return nil, ErrNotConfigured
	}

	query := `
		SELECT best_wpm, best_accuracy, best_time, total_wpm, total_accuracy, total_time,
			   total_sessions, last_session_date, letter_accuracy, letter_seek_time,
			   bigram_seek_time, finger_stats, hand_stats, row_stats, error_substitution,
			   sfb_stats, hand_alternations, same_hand_runs, rhythm_stats
		FROM user_stats WHERE user_id = $1
	`
	if r.db.driver == "sqlite3" {
		query = `
			SELECT best_wpm, best_accuracy, best_time, total_wpm, total_accuracy, total_time,
				   total_sessions, last_session_date, letter_accuracy, letter_seek_time,
				   bigram_seek_time, finger_stats, hand_stats, row_stats, error_substitution,
				   sfb_stats, hand_alternations, same_hand_runs, rhythm_stats
			FROM user_stats WHERE user_id = ?
		`
	}

	var hs stats.HistoricalStats
	var lastSessionDate sql.NullTime
	var letterAccuracy, letterSeekTime, bigramSeekTime sql.NullString
	var fingerStats, handStats, rowStats, errorSubstitution sql.NullString
	var sfbStats, rhythmStats sql.NullString

	err := r.db.conn.QueryRow(query, userID).Scan(
		&hs.BestWPM, &hs.BestAccuracy, &hs.BestTime,
		&hs.TotalWPM, &hs.TotalAccuracy, &hs.TotalTime,
		&hs.TotalSessions, &lastSessionDate,
		&letterAccuracy, &letterSeekTime, &bigramSeekTime,
		&fingerStats, &handStats, &rowStats, &errorSubstitution,
		&sfbStats, &hs.HandAlternations, &hs.SameHandRuns, &rhythmStats,
	)

	if err == sql.ErrNoRows {
		// Return empty stats for new users
		return &stats.HistoricalStats{
			LetterAccuracy:    make(map[string]stats.LetterStats),
			LetterSeekTime:    make(map[string]stats.LetterSeekStats),
			BigramSeekTime:    make(map[string]stats.BigramSeekStats),
			FingerStats:       make(map[int]stats.FingerStat),
			HandStats:         make(map[int]stats.HandStat),
			RowStats:          make(map[int]stats.RowStat),
			ErrorSubstitution: make(map[string]map[string]int),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	if lastSessionDate.Valid {
		hs.LastSessionDate = lastSessionDate.Time
	}

	// Unmarshal JSON fields
	if letterAccuracy.Valid && letterAccuracy.String != "" {
		if err := json.Unmarshal([]byte(letterAccuracy.String), &hs.LetterAccuracy); err != nil {
			hs.LetterAccuracy = make(map[string]stats.LetterStats)
		}
	} else {
		hs.LetterAccuracy = make(map[string]stats.LetterStats)
	}

	if letterSeekTime.Valid && letterSeekTime.String != "" {
		if err := json.Unmarshal([]byte(letterSeekTime.String), &hs.LetterSeekTime); err != nil {
			hs.LetterSeekTime = make(map[string]stats.LetterSeekStats)
		}
	} else {
		hs.LetterSeekTime = make(map[string]stats.LetterSeekStats)
	}

	if bigramSeekTime.Valid && bigramSeekTime.String != "" {
		if err := json.Unmarshal([]byte(bigramSeekTime.String), &hs.BigramSeekTime); err != nil {
			hs.BigramSeekTime = make(map[string]stats.BigramSeekStats)
		}
	} else {
		hs.BigramSeekTime = make(map[string]stats.BigramSeekStats)
	}

	if fingerStats.Valid && fingerStats.String != "" {
		if err := json.Unmarshal([]byte(fingerStats.String), &hs.FingerStats); err != nil {
			hs.FingerStats = make(map[int]stats.FingerStat)
		}
	} else {
		hs.FingerStats = make(map[int]stats.FingerStat)
	}

	if handStats.Valid && handStats.String != "" {
		if err := json.Unmarshal([]byte(handStats.String), &hs.HandStats); err != nil {
			hs.HandStats = make(map[int]stats.HandStat)
		}
	} else {
		hs.HandStats = make(map[int]stats.HandStat)
	}

	if rowStats.Valid && rowStats.String != "" {
		if err := json.Unmarshal([]byte(rowStats.String), &hs.RowStats); err != nil {
			hs.RowStats = make(map[int]stats.RowStat)
		}
	} else {
		hs.RowStats = make(map[int]stats.RowStat)
	}

	if errorSubstitution.Valid && errorSubstitution.String != "" {
		if err := json.Unmarshal([]byte(errorSubstitution.String), &hs.ErrorSubstitution); err != nil {
			hs.ErrorSubstitution = make(map[string]map[string]int)
		}
	} else {
		hs.ErrorSubstitution = make(map[string]map[string]int)
	}

	if sfbStats.Valid && sfbStats.String != "" {
		json.Unmarshal([]byte(sfbStats.String), &hs.SFBStats)
	}

	if rhythmStats.Valid && rhythmStats.String != "" {
		json.Unmarshal([]byte(rhythmStats.String), &hs.RhythmStats)
	}

	return &hs, nil
}

// Save saves the historical stats for a user (upsert).
func (r *StatsRepository) Save(userID string, hs *stats.HistoricalStats) error {
	if r.db == nil || r.db.conn == nil {
		return ErrNotConfigured
	}

	// Marshal JSON fields
	letterAccuracy, _ := json.Marshal(hs.LetterAccuracy)
	letterSeekTime, _ := json.Marshal(hs.LetterSeekTime)
	bigramSeekTime, _ := json.Marshal(hs.BigramSeekTime)
	fingerStats, _ := json.Marshal(hs.FingerStats)
	handStats, _ := json.Marshal(hs.HandStats)
	rowStats, _ := json.Marshal(hs.RowStats)
	errorSubstitution, _ := json.Marshal(hs.ErrorSubstitution)
	sfbStats, _ := json.Marshal(hs.SFBStats)
	rhythmStats, _ := json.Marshal(hs.RhythmStats)

	now := time.Now()

	if r.db.driver == "postgres" {
		query := `
			INSERT INTO user_stats (
				user_id, best_wpm, best_accuracy, best_time, total_wpm, total_accuracy, total_time,
				total_sessions, last_session_date, letter_accuracy, letter_seek_time, bigram_seek_time,
				finger_stats, hand_stats, row_stats, error_substitution, sfb_stats, hand_alternations,
				same_hand_runs, rhythm_stats, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
			ON CONFLICT (user_id) DO UPDATE SET
				best_wpm = EXCLUDED.best_wpm,
				best_accuracy = EXCLUDED.best_accuracy,
				best_time = EXCLUDED.best_time,
				total_wpm = EXCLUDED.total_wpm,
				total_accuracy = EXCLUDED.total_accuracy,
				total_time = EXCLUDED.total_time,
				total_sessions = EXCLUDED.total_sessions,
				last_session_date = EXCLUDED.last_session_date,
				letter_accuracy = EXCLUDED.letter_accuracy,
				letter_seek_time = EXCLUDED.letter_seek_time,
				bigram_seek_time = EXCLUDED.bigram_seek_time,
				finger_stats = EXCLUDED.finger_stats,
				hand_stats = EXCLUDED.hand_stats,
				row_stats = EXCLUDED.row_stats,
				error_substitution = EXCLUDED.error_substitution,
				sfb_stats = EXCLUDED.sfb_stats,
				hand_alternations = EXCLUDED.hand_alternations,
				same_hand_runs = EXCLUDED.same_hand_runs,
				rhythm_stats = EXCLUDED.rhythm_stats,
				updated_at = EXCLUDED.updated_at
		`
		_, err := r.db.conn.Exec(query,
			userID, hs.BestWPM, hs.BestAccuracy, hs.BestTime,
			hs.TotalWPM, hs.TotalAccuracy, hs.TotalTime,
			hs.TotalSessions, hs.LastSessionDate,
			string(letterAccuracy), string(letterSeekTime), string(bigramSeekTime),
			string(fingerStats), string(handStats), string(rowStats), string(errorSubstitution),
			string(sfbStats), hs.HandAlternations, hs.SameHandRuns, string(rhythmStats), now,
		)
		return err
	}

	// SQLite - use INSERT OR REPLACE
	query := `
		INSERT OR REPLACE INTO user_stats (
			user_id, best_wpm, best_accuracy, best_time, total_wpm, total_accuracy, total_time,
			total_sessions, last_session_date, letter_accuracy, letter_seek_time, bigram_seek_time,
			finger_stats, hand_stats, row_stats, error_substitution, sfb_stats, hand_alternations,
			same_hand_runs, rhythm_stats, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.conn.Exec(query,
		userID, hs.BestWPM, hs.BestAccuracy, hs.BestTime,
		hs.TotalWPM, hs.TotalAccuracy, hs.TotalTime,
		hs.TotalSessions, hs.LastSessionDate,
		string(letterAccuracy), string(letterSeekTime), string(bigramSeekTime),
		string(fingerStats), string(handStats), string(rowStats), string(errorSubstitution),
		string(sfbStats), hs.HandAlternations, hs.SameHandRuns, string(rhythmStats), now,
	)
	return err
}

// Delete removes all stats for a user.
func (r *StatsRepository) Delete(userID string) error {
	if r.db == nil || r.db.conn == nil {
		return ErrNotConfigured
	}

	query := `DELETE FROM user_stats WHERE user_id = $1`
	if r.db.driver == "sqlite3" {
		query = `DELETE FROM user_stats WHERE user_id = ?`
	}

	_, err := r.db.conn.Exec(query, userID)
	return err
}

// MergeStats merges local stats with server stats using the merge algorithm:
// - Cumulative values (total sessions, total time): ADD
// - Best values (best WPM, accuracy): MAX
// - Maps (letter stats, finger stats): Merge additively
func MergeStats(server, local *stats.HistoricalStats) *stats.HistoricalStats {
	if server == nil {
		return local
	}
	if local == nil {
		return server
	}

	merged := &stats.HistoricalStats{
		// Take MAX of best values
		BestWPM:      max(server.BestWPM, local.BestWPM),
		BestAccuracy: max(server.BestAccuracy, local.BestAccuracy),
		BestTime:     minPositive(server.BestTime, local.BestTime),

		// ADD cumulative values
		TotalWPM:      server.TotalWPM + local.TotalWPM,
		TotalAccuracy: server.TotalAccuracy + local.TotalAccuracy,
		TotalTime:     server.TotalTime + local.TotalTime,
		TotalSessions: server.TotalSessions + local.TotalSessions,

		// Take most recent session date
		LastSessionDate: maxTime(server.LastSessionDate, local.LastSessionDate),

		// Merge maps
		LetterAccuracy:    mergeLetterAccuracy(server.LetterAccuracy, local.LetterAccuracy),
		LetterSeekTime:    mergeLetterSeekTime(server.LetterSeekTime, local.LetterSeekTime),
		BigramSeekTime:    mergeBigramSeekTime(server.BigramSeekTime, local.BigramSeekTime),
		FingerStats:       mergeFingerStats(server.FingerStats, local.FingerStats),
		HandStats:         mergeHandStats(server.HandStats, local.HandStats),
		RowStats:          mergeRowStats(server.RowStats, local.RowStats),
		ErrorSubstitution: mergeErrorSubstitution(server.ErrorSubstitution, local.ErrorSubstitution),

		// Merge SFB stats
		SFBStats: stats.SFBStats{
			Count:       server.SFBStats.Count + local.SFBStats.Count,
			TotalTimeMs: server.SFBStats.TotalTimeMs + local.SFBStats.TotalTimeMs,
		},

		// ADD hand/rhythm stats
		HandAlternations: server.HandAlternations + local.HandAlternations,
		SameHandRuns:     server.SameHandRuns + local.SameHandRuns,
		RhythmStats: stats.RhythmStats{
			TotalSeekTimeMs: server.RhythmStats.TotalSeekTimeMs + local.RhythmStats.TotalSeekTimeMs,
			TotalSeekTimeSq: server.RhythmStats.TotalSeekTimeSq + local.RhythmStats.TotalSeekTimeSq,
			Count:           server.RhythmStats.Count + local.RhythmStats.Count,
		},
	}

	return merged
}

// Helper functions for merging

func minPositive(a, b float64) float64 {
	if a <= 0 {
		return b
	}
	if b <= 0 {
		return a
	}
	if a < b {
		return a
	}
	return b
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func mergeLetterAccuracy(a, b map[string]stats.LetterStats) map[string]stats.LetterStats {
	result := make(map[string]stats.LetterStats)

	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		existing := result[k]
		existing.Presented += v.Presented
		existing.Correct += v.Correct
		result[k] = existing
	}

	return result
}

func mergeLetterSeekTime(a, b map[string]stats.LetterSeekStats) map[string]stats.LetterSeekStats {
	result := make(map[string]stats.LetterSeekStats)

	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		existing := result[k]
		existing.TotalTimeMs += v.TotalTimeMs
		existing.Count += v.Count
		result[k] = existing
	}

	return result
}

func mergeBigramSeekTime(a, b map[string]stats.BigramSeekStats) map[string]stats.BigramSeekStats {
	result := make(map[string]stats.BigramSeekStats)

	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		existing := result[k]
		existing.TotalTimeMs += v.TotalTimeMs
		existing.Count += v.Count
		result[k] = existing
	}

	return result
}

func mergeFingerStats(a, b map[int]stats.FingerStat) map[int]stats.FingerStat {
	result := make(map[int]stats.FingerStat)

	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		existing := result[k]
		existing.Presented += v.Presented
		existing.Correct += v.Correct
		existing.TotalTimeMs += v.TotalTimeMs
		existing.Count += v.Count
		result[k] = existing
	}

	return result
}

func mergeHandStats(a, b map[int]stats.HandStat) map[int]stats.HandStat {
	result := make(map[int]stats.HandStat)

	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		existing := result[k]
		existing.Presented += v.Presented
		existing.Correct += v.Correct
		existing.TotalTimeMs += v.TotalTimeMs
		existing.Count += v.Count
		result[k] = existing
	}

	return result
}

func mergeRowStats(a, b map[int]stats.RowStat) map[int]stats.RowStat {
	result := make(map[int]stats.RowStat)

	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		existing := result[k]
		existing.Presented += v.Presented
		existing.Correct += v.Correct
		existing.TotalTimeMs += v.TotalTimeMs
		existing.Count += v.Count
		result[k] = existing
	}

	return result
}

func mergeErrorSubstitution(a, b map[string]map[string]int) map[string]map[string]int {
	result := make(map[string]map[string]int)

	for expected, typedMap := range a {
		result[expected] = make(map[string]int)
		for typed, count := range typedMap {
			result[expected][typed] = count
		}
	}
	for expected, typedMap := range b {
		if result[expected] == nil {
			result[expected] = make(map[string]int)
		}
		for typed, count := range typedMap {
			result[expected][typed] += count
		}
	}

	return result
}
