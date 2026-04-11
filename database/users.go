package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

// User represents a user in the database.
type User struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	DisplayName string     `json:"display_name,omitempty"`
	AvatarURL   string     `json:"avatar_url,omitempty"`
	Provider    string     `json:"provider"`
	ProviderID  string     `json:"provider_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

// UserRepository provides access to user data.
type UserRepository struct {
	db *DB
}

// Users returns the user repository.
func (db *DB) Users() *UserRepository {
	return &UserRepository{db: db}
}

// generateID creates a unique identifier.
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Create creates a new user in the database.
// If a user with the same provider/provider_id exists, returns the existing user.
func (r *UserRepository) Create(user *User) error {
	if r.db == nil || r.db.conn == nil {
		return ErrNotConfigured
	}

	if user.ID == "" {
		user.ID = generateID()
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, email, display_name, avatar_url, provider, provider_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	// SQLite uses ? placeholders, PostgreSQL uses $N
	if r.db.driver == "sqlite3" {
		query = `
			INSERT INTO users (id, email, display_name, avatar_url, provider, provider_id, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`
	}

	_, err := r.db.conn.Exec(query,
		user.ID, user.Email, user.DisplayName, user.AvatarURL,
		user.Provider, user.ProviderID, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		// Check if it's a duplicate error and return the existing user
		existing, findErr := r.FindByProvider(user.Provider, user.ProviderID)
		if findErr == nil {
			*user = *existing
			return nil
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByID finds a user by their ID.
func (r *UserRepository) FindByID(id string) (*User, error) {
	if r.db == nil || r.db.conn == nil {
		return nil, ErrNotConfigured
	}

	query := `SELECT id, email, display_name, avatar_url, provider, provider_id, created_at, updated_at, last_login_at FROM users WHERE id = $1`
	if r.db.driver == "sqlite3" {
		query = `SELECT id, email, display_name, avatar_url, provider, provider_id, created_at, updated_at, last_login_at FROM users WHERE id = ?`
	}

	var user User
	var displayName, avatarURL sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.conn.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &displayName, &avatarURL,
		&user.Provider, &user.ProviderID, &user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if displayName.Valid {
		user.DisplayName = displayName.String
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return &user, nil
}

// FindByEmail finds a user by their email address.
func (r *UserRepository) FindByEmail(email string) (*User, error) {
	if r.db == nil || r.db.conn == nil {
		return nil, ErrNotConfigured
	}

	query := `SELECT id, email, display_name, avatar_url, provider, provider_id, created_at, updated_at, last_login_at FROM users WHERE email = $1`
	if r.db.driver == "sqlite3" {
		query = `SELECT id, email, display_name, avatar_url, provider, provider_id, created_at, updated_at, last_login_at FROM users WHERE email = ?`
	}

	var user User
	var displayName, avatarURL sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.conn.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &displayName, &avatarURL,
		&user.Provider, &user.ProviderID, &user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	if displayName.Valid {
		user.DisplayName = displayName.String
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return &user, nil
}

// FindByProvider finds a user by their OAuth provider and provider ID.
func (r *UserRepository) FindByProvider(provider, providerID string) (*User, error) {
	if r.db == nil || r.db.conn == nil {
		return nil, ErrNotConfigured
	}

	query := `SELECT id, email, display_name, avatar_url, provider, provider_id, created_at, updated_at, last_login_at FROM users WHERE provider = $1 AND provider_id = $2`
	if r.db.driver == "sqlite3" {
		query = `SELECT id, email, display_name, avatar_url, provider, provider_id, created_at, updated_at, last_login_at FROM users WHERE provider = ? AND provider_id = ?`
	}

	var user User
	var displayName, avatarURL sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.conn.QueryRow(query, provider, providerID).Scan(
		&user.ID, &user.Email, &displayName, &avatarURL,
		&user.Provider, &user.ProviderID, &user.CreatedAt, &user.UpdatedAt, &lastLoginAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by provider: %w", err)
	}

	if displayName.Valid {
		user.DisplayName = displayName.String
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return &user, nil
}

// Update updates a user's information.
func (r *UserRepository) Update(user *User) error {
	if r.db == nil || r.db.conn == nil {
		return ErrNotConfigured
	}

	user.UpdatedAt = time.Now()

	query := `UPDATE users SET email = $1, display_name = $2, avatar_url = $3, updated_at = $4 WHERE id = $5`
	if r.db.driver == "sqlite3" {
		query = `UPDATE users SET email = ?, display_name = ?, avatar_url = ?, updated_at = ? WHERE id = ?`
	}

	result, err := r.db.conn.Exec(query, user.Email, user.DisplayName, user.AvatarURL, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateLastLogin updates the user's last login timestamp.
func (r *UserRepository) UpdateLastLogin(userID string) error {
	if r.db == nil || r.db.conn == nil {
		return ErrNotConfigured
	}

	now := time.Now()
	query := `UPDATE users SET last_login_at = $1, updated_at = $1 WHERE id = $2`
	if r.db.driver == "sqlite3" {
		query = `UPDATE users SET last_login_at = ?, updated_at = ? WHERE id = ?`
	}

	var err error
	if r.db.driver == "sqlite3" {
		_, err = r.db.conn.Exec(query, now, now, userID)
	} else {
		_, err = r.db.conn.Exec(query, now, userID)
	}
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// Delete removes a user from the database.
// This will cascade delete their stats and refresh tokens.
func (r *UserRepository) Delete(userID string) error {
	if r.db == nil || r.db.conn == nil {
		return ErrNotConfigured
	}

	query := `DELETE FROM users WHERE id = $1`
	if r.db.driver == "sqlite3" {
		query = `DELETE FROM users WHERE id = ?`
	}

	result, err := r.db.conn.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

// FindOrCreate finds a user by provider or creates them if they don't exist.
// This is the primary method used during OAuth login.
func (r *UserRepository) FindOrCreate(provider, providerID, email, displayName, avatarURL string) (*User, bool, error) {
	// Try to find existing user
	user, err := r.FindByProvider(provider, providerID)
	if err == nil {
		// User exists, update their info if changed
		needsUpdate := false
		if user.Email != email {
			user.Email = email
			needsUpdate = true
		}
		if user.DisplayName != displayName {
			user.DisplayName = displayName
			needsUpdate = true
		}
		if user.AvatarURL != avatarURL {
			user.AvatarURL = avatarURL
			needsUpdate = true
		}

		if needsUpdate {
			if err := r.Update(user); err != nil {
				return nil, false, err
			}
		}

		return user, false, nil
	}

	if err != ErrUserNotFound {
		return nil, false, err
	}

	// Create new user
	user = &User{
		Email:       email,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
		Provider:    provider,
		ProviderID:  providerID,
	}

	if err := r.Create(user); err != nil {
		return nil, false, err
	}

	return user, true, nil
}
