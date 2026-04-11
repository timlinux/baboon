package auth

import "errors"

var (
	// ErrProviderNotConfigured is returned when attempting to use an unconfigured OAuth provider.
	ErrProviderNotConfigured = errors.New("OAuth provider not configured")

	// ErrInvalidToken is returned when a token is invalid.
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenExpired is returned when a token has expired.
	ErrTokenExpired = errors.New("token expired")

	// ErrInvalidState is returned when the OAuth state parameter doesn't match.
	ErrInvalidState = errors.New("invalid OAuth state")

	// ErrAuthRequired is returned when authentication is required but not provided.
	ErrAuthRequired = errors.New("authentication required")

	// ErrAuthDisabled is returned when auth features are called but auth is disabled.
	ErrAuthDisabled = errors.New("authentication is disabled")
)
