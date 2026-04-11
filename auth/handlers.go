package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

// Handlers provides HTTP handlers for authentication endpoints.
type Handlers struct {
	service *Service
}

// NewHandlers creates new authentication handlers.
func NewHandlers(service *Service) *Handlers {
	return &Handlers{service: service}
}

// generateState creates a secure random state string for OAuth.
func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// HandleLogin initiates OAuth login for a provider.
// GET /api/auth/{provider}/login
func (h *Handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if !h.service.IsEnabled() {
		writeError(w, http.StatusServiceUnavailable, "authentication is disabled")
		return
	}

	provider := r.PathValue("provider")
	if !IsValidProvider(provider) {
		writeError(w, http.StatusBadRequest, "invalid provider")
		return
	}

	if !h.service.HasProvider(Provider(provider)) {
		writeError(w, http.StatusBadRequest, "provider not configured")
		return
	}

	state := generateState()

	// Store state in a short-lived cookie for validation
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		Secure:   r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
	})

	authURL, err := h.service.GetAuthURL(Provider(provider), state)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate auth URL")
		return
	}

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// HandleCallback handles the OAuth callback from a provider.
// GET /api/auth/{provider}/callback
func (h *Handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if !h.service.IsEnabled() {
		writeError(w, http.StatusServiceUnavailable, "authentication is disabled")
		return
	}

	provider := r.PathValue("provider")
	if !IsValidProvider(provider) {
		writeError(w, http.StatusBadRequest, "invalid provider")
		return
	}

	// Validate state
	state := r.URL.Query().Get("state")
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != state {
		writeError(w, http.StatusBadRequest, "invalid state")
		return
	}

	// Clear the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Check for error from provider
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		// Redirect to frontend with error
		http.Redirect(w, r, "/?auth_error="+errMsg, http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "missing authorization code")
		return
	}

	// Exchange code for user info
	userInfo, err := h.service.HandleCallback(r.Context(), Provider(provider), code)
	if err != nil {
		http.Redirect(w, r, "/?auth_error=oauth_failed", http.StatusTemporaryRedirect)
		return
	}

	// Login/create user and generate tokens
	user, tokens, err := h.service.LoginUser(r.Context(), Provider(provider), userInfo)
	if err != nil {
		http.Redirect(w, r, "/?auth_error=login_failed", http.StatusTemporaryRedirect)
		return
	}

	// Set cookies for tokens
	h.setAuthCookies(w, r, tokens)

	// Redirect to frontend with success
	http.Redirect(w, r, "/?auth_success=true&user_id="+user.ID, http.StatusTemporaryRedirect)
}

// HandleLogout logs out the user by revoking tokens.
// POST /api/auth/logout
func (h *Handlers) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if !h.service.IsEnabled() {
		writeError(w, http.StatusServiceUnavailable, "authentication is disabled")
		return
	}

	// Get user from context (set by middleware)
	userID := getUserIDFromRequest(r)
	if userID != "" {
		h.service.Logout(userID)
	}

	// Clear cookies
	h.clearAuthCookies(w)

	writeJSON(w, http.StatusOK, map[string]string{"status": "logged out"})
}

// HandleRefresh refreshes the access token using a refresh token.
// POST /api/auth/refresh
func (h *Handlers) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	if !h.service.IsEnabled() {
		writeError(w, http.StatusServiceUnavailable, "authentication is disabled")
		return
	}

	// Get refresh token from cookie
	refreshCookie, err := r.Cookie("baboon_refresh_token")
	if err != nil {
		writeError(w, http.StatusUnauthorized, "no refresh token")
		return
	}

	tokens, err := h.service.RefreshTokens(refreshCookie.Value)
	if err != nil {
		h.clearAuthCookies(w)
		writeError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	// Set new cookies
	h.setAuthCookies(w, r, tokens)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"access_token": tokens.AccessToken,
		"expires_in":   tokens.ExpiresIn,
	})
}

// HandleMe returns the current user's information.
// GET /api/auth/me
func (h *Handlers) HandleMe(w http.ResponseWriter, r *http.Request) {
	if !h.service.IsEnabled() {
		writeError(w, http.StatusServiceUnavailable, "authentication is disabled")
		return
	}

	userID := getUserIDFromRequest(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "user not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":           user.ID,
		"email":        user.Email,
		"display_name": user.DisplayName,
		"avatar_url":   user.AvatarURL,
		"provider":     user.Provider,
		"created_at":   user.CreatedAt,
	})
}

// HandleConfig returns auth configuration for the frontend.
// GET /api/auth/config (called from /api/config)
func (h *Handlers) GetAuthConfig() map[string]interface{} {
	return map[string]interface{}{
		"auth_enabled": h.service.IsEnabled(),
		"providers":    h.service.GetEnabledProviders(),
	}
}

// setAuthCookies sets the authentication cookies.
func (h *Handlers) setAuthCookies(w http.ResponseWriter, r *http.Request, tokens *TokenPair) {
	secure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"

	// Access token cookie (short-lived, accessible by JS for API calls)
	http.SetCookie(w, &http.Cookie{
		Name:     "baboon_access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		MaxAge:   tokens.ExpiresIn,
		HttpOnly: false, // Allow JS access for API calls
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})

	// Refresh token cookie (long-lived, HTTP-only for security)
	http.SetCookie(w, &http.Cookie{
		Name:     "baboon_refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/api/auth",
		MaxAge:   7 * 24 * 60 * 60, // 7 days
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}

// clearAuthCookies clears the authentication cookies.
func (h *Handlers) clearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "baboon_access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: false,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "baboon_refresh_token",
		Value:    "",
		Path:     "/api/auth",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// getUserIDFromRequest extracts the user ID from the request context.
func getUserIDFromRequest(r *http.Request) string {
	if id := r.Context().Value(userIDContextKey); id != nil {
		return id.(string)
	}
	return ""
}

// contextKey is a type for context keys to avoid collisions.
type contextKey string

const userIDContextKey contextKey = "userID"

// extractAccessToken extracts the access token from the request.
// It checks the Authorization header first, then falls back to cookies.
func extractAccessToken(r *http.Request) string {
	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Fall back to cookie
	if cookie, err := r.Cookie("baboon_access_token"); err == nil {
		return cookie.Value
	}

	return ""
}

// HandleDeleteAccount deletes the current user's account.
// DELETE /api/auth/account
func (h *Handlers) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	if !h.service.IsEnabled() {
		writeError(w, http.StatusServiceUnavailable, "authentication is disabled")
		return
	}

	userID := getUserIDFromRequest(r)
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	if err := h.service.DeleteUser(userID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete account")
		return
	}

	h.clearAuthCookies(w)

	writeJSON(w, http.StatusOK, map[string]string{"status": "account deleted"})
}
