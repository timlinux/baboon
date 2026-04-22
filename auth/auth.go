// Package auth provides authentication and authorization services for Baboon.
// It supports OAuth2/OIDC authentication with Google, GitHub, Apple, and Microsoft.
package auth

import (
	"context"
	"time"

	"github.com/timlinux/baboon/database"
)

// Provider represents an OAuth provider.
type Provider string

const (
	ProviderGoogle    Provider = "google"
	ProviderGitHub    Provider = "github"
	ProviderApple     Provider = "apple"
	ProviderMicrosoft Provider = "microsoft"
)

// ValidProviders is the list of supported OAuth providers.
var ValidProviders = []Provider{
	ProviderGoogle,
	ProviderGitHub,
	ProviderApple,
	ProviderMicrosoft,
}

// IsValidProvider checks if a provider string is valid.
func IsValidProvider(p string) bool {
	for _, valid := range ValidProviders {
		if string(valid) == p {
			return true
		}
	}
	return false
}

// Config holds authentication configuration.
type Config struct {
	// JWTSecret is the secret used for signing JWTs (required if database is configured).
	JWTSecret string

	// BaseURL is the base URL of the application (for OAuth callbacks).
	BaseURL string

	// AccessTokenExpiry is how long access tokens are valid.
	AccessTokenExpiry time.Duration

	// RefreshTokenExpiry is how long refresh tokens are valid.
	RefreshTokenExpiry time.Duration

	// OAuth provider credentials
	GoogleClientID        string
	GoogleClientSecret    string
	GitHubClientID        string
	GitHubClientSecret    string
	AppleClientID         string
	AppleClientSecret     string
	MicrosoftClientID     string
	MicrosoftClientSecret string
}

// DefaultConfig returns the default authentication configuration.
func DefaultConfig() Config {
	return Config{
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour, // 7 days
	}
}

// Service provides authentication services.
type Service struct {
	config    Config
	db        *database.DB
	providers map[Provider]*OAuthProvider
}

// NewService creates a new authentication service.
func NewService(cfg Config, db *database.DB) (*Service, error) {
	s := &Service{
		config:    cfg,
		db:        db,
		providers: make(map[Provider]*OAuthProvider),
	}

	// Initialize OAuth providers
	if cfg.GoogleClientID != "" && cfg.GoogleClientSecret != "" {
		s.providers[ProviderGoogle] = NewGoogleProvider(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.BaseURL)
	}

	if cfg.GitHubClientID != "" && cfg.GitHubClientSecret != "" {
		s.providers[ProviderGitHub] = NewGitHubProvider(cfg.GitHubClientID, cfg.GitHubClientSecret, cfg.BaseURL)
	}

	if cfg.AppleClientID != "" && cfg.AppleClientSecret != "" {
		s.providers[ProviderApple] = NewAppleProvider(cfg.AppleClientID, cfg.AppleClientSecret, cfg.BaseURL)
	}

	if cfg.MicrosoftClientID != "" && cfg.MicrosoftClientSecret != "" {
		s.providers[ProviderMicrosoft] = NewMicrosoftProvider(cfg.MicrosoftClientID, cfg.MicrosoftClientSecret, cfg.BaseURL)
	}

	return s, nil
}

// IsEnabled returns true if authentication is enabled (database configured).
func (s *Service) IsEnabled() bool {
	return s.db != nil && s.db.IsConfigured()
}

// GetEnabledProviders returns a list of configured OAuth providers.
func (s *Service) GetEnabledProviders() []string {
	providers := make([]string, 0, len(s.providers))
	for p := range s.providers {
		providers = append(providers, string(p))
	}
	return providers
}

// HasProvider checks if a specific provider is configured.
func (s *Service) HasProvider(provider Provider) bool {
	_, ok := s.providers[provider]
	return ok
}

// GetProvider returns the OAuth provider configuration.
func (s *Service) GetProvider(provider Provider) (*OAuthProvider, bool) {
	p, ok := s.providers[provider]
	return p, ok
}

// GetAuthURL returns the OAuth authorization URL for a provider.
func (s *Service) GetAuthURL(provider Provider, state string) (string, error) {
	p, ok := s.providers[provider]
	if !ok {
		return "", ErrProviderNotConfigured
	}
	return p.GetAuthURL(state), nil
}

// HandleCallback handles the OAuth callback and returns user info.
func (s *Service) HandleCallback(ctx context.Context, provider Provider, code string) (*UserInfo, error) {
	p, ok := s.providers[provider]
	if !ok {
		return nil, ErrProviderNotConfigured
	}

	token, err := p.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return p.GetUserInfo(ctx, token)
}

// LoginUser handles user login, creating the user if needed and returning tokens.
func (s *Service) LoginUser(ctx context.Context, provider Provider, userInfo *UserInfo) (*database.User, *TokenPair, error) {
	if s.db == nil {
		return nil, nil, database.ErrNotConfigured
	}

	// Find or create the user
	user, _, err := s.db.Users().FindOrCreate(
		string(provider),
		userInfo.ProviderID,
		userInfo.Email,
		userInfo.Name,
		userInfo.Picture,
	)
	if err != nil {
		return nil, nil, err
	}

	// Update last login
	if err := s.db.Users().UpdateLastLogin(user.ID); err != nil {
		// Non-fatal, just log it
	}

	// Generate tokens
	tokens, err := s.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// GenerateTokenPair generates a new access token and refresh token pair.
func (s *Service) GenerateTokenPair(userID string) (*TokenPair, error) {
	if s.db == nil {
		return nil, database.ErrNotConfigured
	}

	// Generate access token (JWT)
	accessToken, err := GenerateAccessToken(userID, s.config.JWTSecret, s.config.AccessTokenExpiry)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, _, err := s.db.Tokens().Create(userID, s.config.RefreshTokenExpiry)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.config.AccessTokenExpiry.Seconds()),
	}, nil
}

// RefreshTokens validates a refresh token and generates a new token pair.
func (s *Service) RefreshTokens(refreshToken string) (*TokenPair, error) {
	if s.db == nil {
		return nil, database.ErrNotConfigured
	}

	// Validate the refresh token
	rt, err := s.db.Tokens().Validate(refreshToken)
	if err != nil {
		return nil, err
	}

	// Revoke the old refresh token
	if err := s.db.Tokens().Revoke(rt.ID); err != nil {
		// Non-fatal, continue with new tokens
	}

	// Generate new token pair
	return s.GenerateTokenPair(rt.UserID)
}

// ValidateAccessToken validates an access token and returns the user ID.
func (s *Service) ValidateAccessToken(token string) (string, error) {
	return ValidateAccessToken(token, s.config.JWTSecret)
}

// GetUser retrieves a user by ID.
func (s *Service) GetUser(userID string) (*database.User, error) {
	if s.db == nil {
		return nil, database.ErrNotConfigured
	}
	return s.db.Users().FindByID(userID)
}

// Logout revokes all refresh tokens for a user.
func (s *Service) Logout(userID string) error {
	if s.db == nil {
		return database.ErrNotConfigured
	}
	return s.db.Tokens().RevokeAllForUser(userID)
}

// DeleteUser removes a user and all their data.
func (s *Service) DeleteUser(userID string) error {
	if s.db == nil {
		return database.ErrNotConfigured
	}

	// Delete stats first (cascade should handle this, but be explicit)
	if err := s.db.Stats().Delete(userID); err != nil {
		// Non-fatal
	}

	// Revoke all tokens
	if err := s.db.Tokens().RevokeAllForUser(userID); err != nil {
		// Non-fatal
	}

	// Delete the user
	return s.db.Users().Delete(userID)
}

// TokenPair contains an access token and refresh token.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds until access token expires
}

// UserInfo contains user information from an OAuth provider.
type UserInfo struct {
	ProviderID string `json:"provider_id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Picture    string `json:"picture,omitempty"`
}
