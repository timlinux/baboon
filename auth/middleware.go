package auth

import (
	"context"
	"net/http"
)

// Middleware provides authentication middleware for HTTP handlers.
type Middleware struct {
	service *Service
}

// NewMiddleware creates a new authentication middleware.
func NewMiddleware(service *Service) *Middleware {
	return &Middleware{service: service}
}

// RequireAuth is a middleware that requires authentication.
// If the user is not authenticated, it returns 401 Unauthorized.
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.service.IsEnabled() {
			writeError(w, http.StatusServiceUnavailable, "authentication is disabled")
			return
		}

		token := extractAccessToken(r)
		if token == "" {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		userID, err := m.service.ValidateAccessToken(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth is a middleware that optionally extracts user info.
// If a valid token is present, it adds the user ID to the context.
// If not, it continues without authentication.
func (m *Middleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.service.IsEnabled() {
			next.ServeHTTP(w, r)
			return
		}

		token := extractAccessToken(r)
		if token != "" {
			if userID, err := m.service.ValidateAccessToken(token); err == nil {
				ctx := context.WithValue(r.Context(), userIDContextKey, userID)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}

// RequireAuthFunc is a middleware wrapper for http.HandlerFunc.
func (m *Middleware) RequireAuthFunc(next http.HandlerFunc) http.HandlerFunc {
	return m.RequireAuth(next).ServeHTTP
}

// OptionalAuthFunc is a middleware wrapper for http.HandlerFunc.
func (m *Middleware) OptionalAuthFunc(next http.HandlerFunc) http.HandlerFunc {
	return m.OptionalAuth(next).ServeHTTP
}

// GetUserID extracts the user ID from the request context.
// Returns empty string if no user is authenticated.
func GetUserID(r *http.Request) string {
	if id := r.Context().Value(userIDContextKey); id != nil {
		return id.(string)
	}
	return ""
}

// IsAuthenticated returns true if the request has a valid user ID in context.
func IsAuthenticated(r *http.Request) bool {
	return GetUserID(r) != ""
}
