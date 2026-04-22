// Package database provides database connectivity and operations for the leaderboard.
package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// LeaderboardEntry represents a single entry on the monthly leaderboard.
type LeaderboardEntry struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Username        string    `json:"username"`     // Platform username or "ANON"
	DisplayName     string    `json:"display_name"` // User's custom message
	WPM             float64   `json:"wpm"`
	Accuracy        float64   `json:"accuracy"`
	Score           float64   `json:"score"` // Combined score: WPM * (Accuracy/100)
	DurationSeconds float64   `json:"duration_seconds"`
	MonthYear       string    `json:"month_year"` // Format: "2026-04"
	CreatedAt       time.Time `json:"created_at"`
	BadgeImageURL   string    `json:"badge_image_url,omitempty"`
	Rank            int       `json:"rank,omitempty"` // Populated when fetching top 10
}

// LeaderboardRepository handles database operations for leaderboard entries.
type LeaderboardRepository struct {
	db     *sql.DB
	driver string
}

// Leaderboard returns a LeaderboardRepository for the database.
func (db *DB) Leaderboard() *LeaderboardRepository {
	if db == nil {
		return nil
	}
	return &LeaderboardRepository{
		db:     db.conn,
		driver: db.driver,
	}
}

// GetTop10 returns the top 10 entries for a given month.
// monthYear should be in format "2026-04".
// Entries are ranked by combined score (WPM * Accuracy/100).
func (r *LeaderboardRepository) GetTop10(monthYear string) ([]LeaderboardEntry, error) {
	if r == nil || r.db == nil {
		return nil, ErrNotConfigured
	}

	query := `
		SELECT id, user_id, COALESCE(username, 'ANON'), display_name, wpm, accuracy, COALESCE(score, wpm * accuracy / 100), duration_seconds, month_year, created_at, badge_image_url
		FROM leaderboard_entries
		WHERE month_year = $1
		ORDER BY COALESCE(score, wpm * accuracy / 100) DESC
		LIMIT 10
	`
	if r.driver == "sqlite3" {
		query = `
			SELECT id, user_id, COALESCE(username, 'ANON'), display_name, wpm, accuracy, COALESCE(score, wpm * accuracy / 100), duration_seconds, month_year, created_at, COALESCE(badge_image_url, '')
			FROM leaderboard_entries
			WHERE month_year = ?
			ORDER BY COALESCE(score, wpm * accuracy / 100) DESC
			LIMIT 10
		`
	}

	rows, err := r.db.Query(query, monthYear)
	if err != nil {
		return nil, fmt.Errorf("failed to query leaderboard: %w", err)
	}
	defer rows.Close()

	entries := make([]LeaderboardEntry, 0, 10)
	rank := 1
	for rows.Next() {
		var entry LeaderboardEntry
		var badgeURL sql.NullString
		if r.driver == "sqlite3" {
			var badgeURLStr string
			if err := rows.Scan(&entry.ID, &entry.UserID, &entry.Username, &entry.DisplayName, &entry.WPM, &entry.Accuracy, &entry.Score, &entry.DurationSeconds, &entry.MonthYear, &entry.CreatedAt, &badgeURLStr); err != nil {
				return nil, fmt.Errorf("failed to scan leaderboard entry: %w", err)
			}
			entry.BadgeImageURL = badgeURLStr
		} else {
			if err := rows.Scan(&entry.ID, &entry.UserID, &entry.Username, &entry.DisplayName, &entry.WPM, &entry.Accuracy, &entry.Score, &entry.DurationSeconds, &entry.MonthYear, &entry.CreatedAt, &badgeURL); err != nil {
				return nil, fmt.Errorf("failed to scan leaderboard entry: %w", err)
			}
			if badgeURL.Valid {
				entry.BadgeImageURL = badgeURL.String
			}
		}
		entry.Rank = rank
		rank++
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

// CheckQualifies checks if a given WPM and accuracy would qualify for the top 10 in a month.
// Returns whether it qualifies and what rank it would achieve.
// Score is calculated as WPM * (Accuracy/100).
func (r *LeaderboardRepository) CheckQualifies(monthYear string, wpm, accuracy float64) (bool, int, error) {
	if r == nil || r.db == nil {
		return false, 0, ErrNotConfigured
	}

	// Calculate the combined score
	candidateScore := wpm * accuracy / 100

	// Get current top 10 scores for the month
	query := `
		SELECT COALESCE(score, wpm * accuracy / 100) FROM leaderboard_entries
		WHERE month_year = $1
		ORDER BY COALESCE(score, wpm * accuracy / 100) DESC
		LIMIT 10
	`
	if r.driver == "sqlite3" {
		query = `
			SELECT COALESCE(score, wpm * accuracy / 100) FROM leaderboard_entries
			WHERE month_year = ?
			ORDER BY COALESCE(score, wpm * accuracy / 100) DESC
			LIMIT 10
		`
	}

	rows, err := r.db.Query(query, monthYear)
	if err != nil {
		return false, 0, fmt.Errorf("failed to query leaderboard: %w", err)
	}
	defer rows.Close()

	scores := make([]float64, 0, 10)
	for rows.Next() {
		var s float64
		if err := rows.Scan(&s); err != nil {
			return false, 0, fmt.Errorf("failed to scan score: %w", err)
		}
		scores = append(scores, s)
	}

	if err := rows.Err(); err != nil {
		return false, 0, err
	}

	// If less than 10 entries, any score qualifies
	if len(scores) < 10 {
		// Find the rank position
		rank := 1
		for _, s := range scores {
			if candidateScore <= s {
				rank++
			}
		}
		return true, rank, nil
	}

	// Check if this score beats the 10th place
	lowestTop10 := scores[len(scores)-1]
	if candidateScore > lowestTop10 {
		// Find the rank position
		rank := 1
		for _, s := range scores {
			if candidateScore <= s {
				rank++
			}
		}
		return true, rank, nil
	}

	return false, 0, nil
}

// Submit creates or updates a leaderboard entry for a user.
// Only one entry per user per month is allowed - updates if the new score is higher.
// Score is calculated as WPM * (Accuracy/100).
func (r *LeaderboardRepository) Submit(entry *LeaderboardEntry) (*LeaderboardEntry, error) {
	if r == nil || r.db == nil {
		return nil, ErrNotConfigured
	}

	// Calculate the combined score
	entry.Score = entry.WPM * entry.Accuracy / 100

	// Default username to ANON if not set
	if entry.Username == "" {
		entry.Username = "ANON"
	}

	// Check if user already has an entry for this month
	existing, err := r.GetUserEntry(entry.UserID, entry.MonthYear)
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}

	if existing != nil {
		// Only update if the new score is higher
		existingScore := existing.WPM * existing.Accuracy / 100
		if entry.Score <= existingScore {
			return existing, nil // Return existing entry unchanged
		}

		// Update existing entry with better score
		query := `
			UPDATE leaderboard_entries
			SET username = $1, display_name = $2, wpm = $3, accuracy = $4, score = $5, duration_seconds = $6, created_at = $7
			WHERE id = $8
		`
		if r.driver == "sqlite3" {
			query = `
				UPDATE leaderboard_entries
				SET username = ?, display_name = ?, wpm = ?, accuracy = ?, score = ?, duration_seconds = ?, created_at = ?
				WHERE id = ?
			`
		}

		now := time.Now()
		_, err := r.db.Exec(query, entry.Username, entry.DisplayName, entry.WPM, entry.Accuracy, entry.Score, entry.DurationSeconds, now, existing.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update leaderboard entry: %w", err)
		}

		entry.ID = existing.ID
		entry.CreatedAt = now
		return entry, nil
	}

	// Create new entry
	entry.ID = uuid.New().String()
	entry.CreatedAt = time.Now()

	query := `
		INSERT INTO leaderboard_entries (id, user_id, username, display_name, wpm, accuracy, score, duration_seconds, month_year, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	if r.driver == "sqlite3" {
		query = `
			INSERT INTO leaderboard_entries (id, user_id, username, display_name, wpm, accuracy, score, duration_seconds, month_year, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
	}

	_, err = r.db.Exec(query, entry.ID, entry.UserID, entry.Username, entry.DisplayName, entry.WPM, entry.Accuracy, entry.Score, entry.DurationSeconds, entry.MonthYear, entry.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert leaderboard entry: %w", err)
	}

	return entry, nil
}

// GetByID retrieves a leaderboard entry by its ID.
func (r *LeaderboardRepository) GetByID(id string) (*LeaderboardEntry, error) {
	if r == nil || r.db == nil {
		return nil, ErrNotConfigured
	}

	query := `
		SELECT id, user_id, COALESCE(username, 'ANON'), display_name, wpm, accuracy, COALESCE(score, wpm * accuracy / 100), duration_seconds, month_year, created_at, badge_image_url
		FROM leaderboard_entries
		WHERE id = $1
	`
	if r.driver == "sqlite3" {
		query = `
			SELECT id, user_id, COALESCE(username, 'ANON'), display_name, wpm, accuracy, COALESCE(score, wpm * accuracy / 100), duration_seconds, month_year, created_at, COALESCE(badge_image_url, '')
			FROM leaderboard_entries
			WHERE id = ?
		`
	}

	var entry LeaderboardEntry
	var badgeURL sql.NullString

	var err error
	if r.driver == "sqlite3" {
		var badgeURLStr string
		err = r.db.QueryRow(query, id).Scan(&entry.ID, &entry.UserID, &entry.Username, &entry.DisplayName, &entry.WPM, &entry.Accuracy, &entry.Score, &entry.DurationSeconds, &entry.MonthYear, &entry.CreatedAt, &badgeURLStr)
		entry.BadgeImageURL = badgeURLStr
	} else {
		err = r.db.QueryRow(query, id).Scan(&entry.ID, &entry.UserID, &entry.Username, &entry.DisplayName, &entry.WPM, &entry.Accuracy, &entry.Score, &entry.DurationSeconds, &entry.MonthYear, &entry.CreatedAt, &badgeURL)
		if badgeURL.Valid {
			entry.BadgeImageURL = badgeURL.String
		}
	}

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query leaderboard entry: %w", err)
	}

	// Calculate rank for this entry based on score
	rankQuery := `
		SELECT COUNT(*) + 1 FROM leaderboard_entries
		WHERE month_year = $1 AND COALESCE(score, wpm * accuracy / 100) > $2
	`
	if r.driver == "sqlite3" {
		rankQuery = `
			SELECT COUNT(*) + 1 FROM leaderboard_entries
			WHERE month_year = ? AND COALESCE(score, wpm * accuracy / 100) > ?
		`
	}
	var rank int
	if err := r.db.QueryRow(rankQuery, entry.MonthYear, entry.Score).Scan(&rank); err != nil {
		return nil, fmt.Errorf("failed to calculate rank: %w", err)
	}
	entry.Rank = rank

	return &entry, nil
}

// GetUserEntry retrieves a user's entry for a specific month.
func (r *LeaderboardRepository) GetUserEntry(userID, monthYear string) (*LeaderboardEntry, error) {
	if r == nil || r.db == nil {
		return nil, ErrNotConfigured
	}

	query := `
		SELECT id, user_id, COALESCE(username, 'ANON'), display_name, wpm, accuracy, COALESCE(score, wpm * accuracy / 100), duration_seconds, month_year, created_at, badge_image_url
		FROM leaderboard_entries
		WHERE user_id = $1 AND month_year = $2
	`
	if r.driver == "sqlite3" {
		query = `
			SELECT id, user_id, COALESCE(username, 'ANON'), display_name, wpm, accuracy, COALESCE(score, wpm * accuracy / 100), duration_seconds, month_year, created_at, COALESCE(badge_image_url, '')
			FROM leaderboard_entries
			WHERE user_id = ? AND month_year = ?
		`
	}

	var entry LeaderboardEntry
	var badgeURL sql.NullString

	var err error
	if r.driver == "sqlite3" {
		var badgeURLStr string
		err = r.db.QueryRow(query, userID, monthYear).Scan(&entry.ID, &entry.UserID, &entry.Username, &entry.DisplayName, &entry.WPM, &entry.Accuracy, &entry.Score, &entry.DurationSeconds, &entry.MonthYear, &entry.CreatedAt, &badgeURLStr)
		entry.BadgeImageURL = badgeURLStr
	} else {
		err = r.db.QueryRow(query, userID, monthYear).Scan(&entry.ID, &entry.UserID, &entry.Username, &entry.DisplayName, &entry.WPM, &entry.Accuracy, &entry.Score, &entry.DurationSeconds, &entry.MonthYear, &entry.CreatedAt, &badgeURL)
		if badgeURL.Valid {
			entry.BadgeImageURL = badgeURL.String
		}
	}

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query user leaderboard entry: %w", err)
	}

	return &entry, nil
}

// UpdateBadgeURL updates the badge URL for a leaderboard entry.
func (r *LeaderboardRepository) UpdateBadgeURL(id, badgeURL string) error {
	if r == nil || r.db == nil {
		return ErrNotConfigured
	}

	query := `UPDATE leaderboard_entries SET badge_image_url = $1 WHERE id = $2`
	if r.driver == "sqlite3" {
		query = `UPDATE leaderboard_entries SET badge_image_url = ? WHERE id = ?`
	}

	_, err := r.db.Exec(query, badgeURL, id)
	if err != nil {
		return fmt.Errorf("failed to update badge URL: %w", err)
	}

	return nil
}

// GetAvailableMonths returns all months that have leaderboard entries.
func (r *LeaderboardRepository) GetAvailableMonths() ([]string, error) {
	if r == nil || r.db == nil {
		return nil, ErrNotConfigured
	}

	query := `
		SELECT DISTINCT month_year FROM leaderboard_entries
		ORDER BY month_year DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query months: %w", err)
	}
	defer rows.Close()

	months := make([]string, 0)
	for rows.Next() {
		var month string
		if err := rows.Scan(&month); err != nil {
			return nil, fmt.Errorf("failed to scan month: %w", err)
		}
		months = append(months, month)
	}

	return months, rows.Err()
}

// migrateLeaderboard returns the SQL for creating the leaderboard_entries table.
func (db *DB) migrateLeaderboard() string {
	if db.driver == "postgres" {
		return `
		CREATE TABLE IF NOT EXISTS leaderboard_entries (
			id TEXT PRIMARY KEY,
			user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
			username TEXT DEFAULT 'ANON',
			display_name TEXT NOT NULL,
			wpm DOUBLE PRECISION NOT NULL,
			accuracy DOUBLE PRECISION NOT NULL,
			score DOUBLE PRECISION,
			duration_seconds DOUBLE PRECISION NOT NULL,
			month_year TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			badge_image_url TEXT,
			UNIQUE(user_id, month_year)
		);
		CREATE INDEX IF NOT EXISTS idx_leaderboard_month_score ON leaderboard_entries(month_year, score DESC);
		CREATE INDEX IF NOT EXISTS idx_leaderboard_user_month ON leaderboard_entries(user_id, month_year);
		`
	}
	// SQLite
	return `
	CREATE TABLE IF NOT EXISTS leaderboard_entries (
		id TEXT PRIMARY KEY,
		user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
		username TEXT DEFAULT 'ANON',
		display_name TEXT NOT NULL,
		wpm REAL NOT NULL,
		accuracy REAL NOT NULL,
		score REAL,
		duration_seconds REAL NOT NULL,
		month_year TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		badge_image_url TEXT
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_leaderboard_user_month ON leaderboard_entries(user_id, month_year);
	CREATE INDEX IF NOT EXISTS idx_leaderboard_month_score ON leaderboard_entries(month_year, score DESC);
	`
}

// migrateLeaderboardV2 adds username and score columns to existing leaderboard_entries table.
func (db *DB) migrateLeaderboardV2() string {
	if db.driver == "postgres" {
		return `
		ALTER TABLE leaderboard_entries ADD COLUMN IF NOT EXISTS username TEXT DEFAULT 'ANON';
		ALTER TABLE leaderboard_entries ADD COLUMN IF NOT EXISTS score DOUBLE PRECISION;
		UPDATE leaderboard_entries SET score = wpm * accuracy / 100 WHERE score IS NULL;
		UPDATE leaderboard_entries SET username = 'ANON' WHERE username IS NULL;
		`
	}
	// SQLite doesn't support IF NOT EXISTS for ALTER TABLE, so we need to check first
	return `
	-- SQLite migration handled in code
	`
}
