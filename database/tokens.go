package database

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

// RefreshToken represents a refresh token in the database.
type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	Revoked   bool      `json:"revoked"`
}

// TokenRepository provides access to refresh token data.
type TokenRepository struct {
	db *DB
}

// Tokens returns the token repository.
func (db *DB) Tokens() *TokenRepository {
	return &TokenRepository{db: db}
}

// generateToken creates a secure random token and returns both the token and its hash.
func generateToken() (token string, hash string) {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	token = hex.EncodeToString(bytes)

	hashBytes := sha256.Sum256([]byte(token))
	hash = hex.EncodeToString(hashBytes[:])

	return token, hash
}

// hashToken creates a SHA256 hash of a token.
func hashToken(token string) string {
	hashBytes := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hashBytes[:])
}

// Create creates a new refresh token for a user.
// Returns the raw token (to be sent to the client) and the token record.
func (r *TokenRepository) Create(userID string, expiresIn time.Duration) (string, *RefreshToken, error) {
	if r.db == nil || r.db.conn == nil {
		return "", nil, ErrNotConfigured
	}

	token, tokenHash := generateToken()
	now := time.Now()
	expiresAt := now.Add(expiresIn)

	rt := &RefreshToken{
		ID:        generateID(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: now,
		Revoked:   false,
	}

	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at, revoked)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	if r.db.driver == "sqlite3" {
		query = `
			INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at, revoked)
			VALUES (?, ?, ?, ?, ?, ?)
		`
	}

	_, err := r.db.conn.Exec(query, rt.ID, rt.UserID, rt.TokenHash, rt.ExpiresAt, rt.CreatedAt, rt.Revoked)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return token, rt, nil
}

// Validate validates a refresh token and returns the associated token record.
// Returns an error if the token is invalid, expired, or revoked.
func (r *TokenRepository) Validate(token string) (*RefreshToken, error) {
	if r.db == nil || r.db.conn == nil {
		return nil, ErrNotConfigured
	}

	tokenHash := hashToken(token)

	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	if r.db.driver == "sqlite3" {
		query = `
			SELECT id, user_id, token_hash, expires_at, created_at, revoked
			FROM refresh_tokens
			WHERE token_hash = ?
		`
	}

	var rt RefreshToken
	var revoked int
	if r.db.driver == "sqlite3" {
		err := r.db.conn.QueryRow(query, tokenHash).Scan(
			&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.CreatedAt, &revoked,
		)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid token")
		}
		if err != nil {
			return nil, fmt.Errorf("failed to validate token: %w", err)
		}
		rt.Revoked = revoked != 0
	} else {
		err := r.db.conn.QueryRow(query, tokenHash).Scan(
			&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.CreatedAt, &rt.Revoked,
		)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid token")
		}
		if err != nil {
			return nil, fmt.Errorf("failed to validate token: %w", err)
		}
	}

	if rt.Revoked {
		return nil, fmt.Errorf("token has been revoked")
	}

	if time.Now().After(rt.ExpiresAt) {
		return nil, fmt.Errorf("token has expired")
	}

	return &rt, nil
}

// Revoke revokes a specific refresh token.
func (r *TokenRepository) Revoke(tokenID string) error {
	if r.db == nil || r.db.conn == nil {
		return ErrNotConfigured
	}

	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE id = $1`
	if r.db.driver == "sqlite3" {
		query = `UPDATE refresh_tokens SET revoked = 1 WHERE id = ?`
	}

	_, err := r.db.conn.Exec(query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

// RevokeAllForUser revokes all refresh tokens for a user.
func (r *TokenRepository) RevokeAllForUser(userID string) error {
	if r.db == nil || r.db.conn == nil {
		return ErrNotConfigured
	}

	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = $1`
	if r.db.driver == "sqlite3" {
		query = `UPDATE refresh_tokens SET revoked = 1 WHERE user_id = ?`
	}

	_, err := r.db.conn.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}

	return nil
}

// CleanupExpired removes expired and revoked tokens from the database.
// This should be called periodically to clean up old tokens.
func (r *TokenRepository) CleanupExpired() error {
	if r.db == nil || r.db.conn == nil {
		return ErrNotConfigured
	}

	query := `DELETE FROM refresh_tokens WHERE expires_at < $1 OR revoked = TRUE`
	if r.db.driver == "sqlite3" {
		query = `DELETE FROM refresh_tokens WHERE expires_at < ? OR revoked = 1`
	}

	_, err := r.db.conn.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to cleanup tokens: %w", err)
	}

	return nil
}
