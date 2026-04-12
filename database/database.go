// Package database provides database connectivity and operations for user authentication
// and stats storage. It supports both SQLite and PostgreSQL, automatically detecting
// the driver based on the DSN format.
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var (
	// ErrNotConfigured is returned when database operations are attempted without configuration.
	ErrNotConfigured = errors.New("database not configured")

	// ErrUserNotFound is returned when a user cannot be found.
	ErrUserNotFound = errors.New("user not found")

	// ErrDuplicateUser is returned when attempting to create a user that already exists.
	ErrDuplicateUser = errors.New("user already exists")
)

// DB represents the database connection and provides repository access.
type DB struct {
	conn   *sql.DB
	driver string
	mu     sync.RWMutex
}

// Config holds database configuration options.
type Config struct {
	// DSN is the database connection string.
	// For SQLite: path to the database file (e.g., "./baboon.db")
	// For PostgreSQL: postgres://user:pass@host/dbname
	DSN string
}

// New creates a new database connection based on the provided configuration.
// Returns nil if DSN is empty (auth features disabled).
func New(cfg Config) (*DB, error) {
	if cfg.DSN == "" {
		return nil, nil
	}

	driver := detectDriver(cfg.DSN)
	dsn := cfg.DSN

	// For SQLite, add query parameters for better concurrency
	if driver == "sqlite3" && !strings.Contains(dsn, "?") {
		dsn = dsn + "?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_cache_size=10000"
	}

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{
		conn:   conn,
		driver: driver,
	}

	// Run migrations
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// detectDriver determines the database driver from the DSN.
func detectDriver(dsn string) string {
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		return "postgres"
	}
	return "sqlite3"
}

// Close closes the database connection.
func (db *DB) Close() error {
	if db == nil || db.conn == nil {
		return nil
	}
	return db.conn.Close()
}

// Driver returns the database driver name.
func (db *DB) Driver() string {
	if db == nil {
		return ""
	}
	return db.driver
}

// IsConfigured returns true if the database is configured and connected.
func (db *DB) IsConfigured() bool {
	return db != nil && db.conn != nil
}

// migrate runs database migrations.
func (db *DB) migrate() error {
	migrations := []string{
		db.migrateUsers(),
		db.migrateUserStats(),
		db.migrateRefreshTokens(),
		db.migrateLeaderboard(),
	}

	for _, migration := range migrations {
		if _, err := db.conn.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	// Run V2 migration for leaderboard (adds username and score columns)
	if err := db.migrateLeaderboardV2Columns(); err != nil {
		return fmt.Errorf("leaderboard v2 migration failed: %w", err)
	}

	return nil
}

// migrateLeaderboardV2Columns adds username and score columns to existing leaderboard_entries table.
func (db *DB) migrateLeaderboardV2Columns() error {
	// Check if the columns exist by trying to add them
	// SQLite and PostgreSQL both ignore duplicate column errors differently

	if db.driver == "postgres" {
		// PostgreSQL supports IF NOT EXISTS
		_, _ = db.conn.Exec(`ALTER TABLE leaderboard_entries ADD COLUMN IF NOT EXISTS username TEXT DEFAULT 'ANON'`)
		_, _ = db.conn.Exec(`ALTER TABLE leaderboard_entries ADD COLUMN IF NOT EXISTS score DOUBLE PRECISION`)
		// Update existing rows
		_, _ = db.conn.Exec(`UPDATE leaderboard_entries SET score = wpm * accuracy / 100 WHERE score IS NULL`)
		_, _ = db.conn.Exec(`UPDATE leaderboard_entries SET username = 'ANON' WHERE username IS NULL`)
	} else {
		// SQLite - check if column exists first
		rows, err := db.conn.Query(`PRAGMA table_info(leaderboard_entries)`)
		if err != nil {
			return nil // Table might not exist yet
		}
		defer rows.Close()

		hasUsername := false
		hasScore := false
		for rows.Next() {
			var cid int
			var name, colType string
			var notNull, pk int
			var dfltValue interface{}
			if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
				continue
			}
			if name == "username" {
				hasUsername = true
			}
			if name == "score" {
				hasScore = true
			}
		}

		if !hasUsername {
			_, _ = db.conn.Exec(`ALTER TABLE leaderboard_entries ADD COLUMN username TEXT DEFAULT 'ANON'`)
		}
		if !hasScore {
			_, _ = db.conn.Exec(`ALTER TABLE leaderboard_entries ADD COLUMN score REAL`)
		}
		// Update existing rows
		_, _ = db.conn.Exec(`UPDATE leaderboard_entries SET score = wpm * accuracy / 100 WHERE score IS NULL`)
		_, _ = db.conn.Exec(`UPDATE leaderboard_entries SET username = 'ANON' WHERE username IS NULL`)
	}

	return nil
}

// migrateUsers returns the SQL for creating the users table.
func (db *DB) migrateUsers() string {
	if db.driver == "postgres" {
		return `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			display_name TEXT,
			avatar_url TEXT,
			provider TEXT NOT NULL,
			provider_id TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_login_at TIMESTAMP,
			UNIQUE(provider, provider_id)
		);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_provider ON users(provider, provider_id);
		`
	}
	// SQLite
	return `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		display_name TEXT,
		avatar_url TEXT,
		provider TEXT NOT NULL,
		provider_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_login_at DATETIME
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_users_provider ON users(provider, provider_id);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`
}

// migrateUserStats returns the SQL for creating the user_stats table.
func (db *DB) migrateUserStats() string {
	if db.driver == "postgres" {
		return `
		CREATE TABLE IF NOT EXISTS user_stats (
			user_id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
			best_wpm DOUBLE PRECISION DEFAULT 0,
			best_accuracy DOUBLE PRECISION DEFAULT 0,
			best_time DOUBLE PRECISION DEFAULT 0,
			total_wpm DOUBLE PRECISION DEFAULT 0,
			total_accuracy DOUBLE PRECISION DEFAULT 0,
			total_time DOUBLE PRECISION DEFAULT 0,
			total_sessions INTEGER DEFAULT 0,
			last_session_date TIMESTAMP,
			letter_accuracy JSONB,
			letter_seek_time JSONB,
			bigram_seek_time JSONB,
			finger_stats JSONB,
			hand_stats JSONB,
			row_stats JSONB,
			error_substitution JSONB,
			sfb_stats JSONB,
			hand_alternations INTEGER DEFAULT 0,
			same_hand_runs INTEGER DEFAULT 0,
			rhythm_stats JSONB,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`
	}
	// SQLite - uses TEXT for JSON columns
	return `
	CREATE TABLE IF NOT EXISTS user_stats (
		user_id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
		best_wpm REAL DEFAULT 0,
		best_accuracy REAL DEFAULT 0,
		best_time REAL DEFAULT 0,
		total_wpm REAL DEFAULT 0,
		total_accuracy REAL DEFAULT 0,
		total_time REAL DEFAULT 0,
		total_sessions INTEGER DEFAULT 0,
		last_session_date DATETIME,
		letter_accuracy TEXT,
		letter_seek_time TEXT,
		bigram_seek_time TEXT,
		finger_stats TEXT,
		hand_stats TEXT,
		row_stats TEXT,
		error_substitution TEXT,
		sfb_stats TEXT,
		hand_alternations INTEGER DEFAULT 0,
		same_hand_runs INTEGER DEFAULT 0,
		rhythm_stats TEXT,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
}

// migrateRefreshTokens returns the SQL for creating the refresh_tokens table.
func (db *DB) migrateRefreshTokens() string {
	if db.driver == "postgres" {
		return `
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id TEXT PRIMARY KEY,
			user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
			token_hash TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			revoked BOOLEAN DEFAULT FALSE
		);
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash ON refresh_tokens(token_hash);
		`
	}
	// SQLite
	return `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id TEXT PRIMARY KEY,
		user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
		token_hash TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		revoked INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash ON refresh_tokens(token_hash);
	`
}
