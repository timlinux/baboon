package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
)

// OAuthProvider represents a configured OAuth2 provider.
type OAuthProvider struct {
	Name        string
	Config      *oauth2.Config
	UserInfoURL string
	ParseUser   func(body []byte) (*UserInfo, error)
}

// GetAuthURL returns the OAuth authorization URL.
func (p *OAuthProvider) GetAuthURL(state string) string {
	return p.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Exchange exchanges an authorization code for an access token.
func (p *OAuthProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.Config.Exchange(ctx, code)
}

// GetUserInfo retrieves user information from the OAuth provider.
func (p *OAuthProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	client := p.Config.Client(ctx, token)
	resp, err := client.Get(p.UserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	return p.ParseUser(body)
}

// NewGoogleProvider creates a Google OAuth provider.
func NewGoogleProvider(clientID, clientSecret, baseURL string) *OAuthProvider {
	return &OAuthProvider{
		Name: "google",
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  baseURL + "/api/auth/google/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		UserInfoURL: "https://www.googleapis.com/oauth2/v3/userinfo",
		ParseUser: func(body []byte) (*UserInfo, error) {
			var data struct {
				Sub     string `json:"sub"`
				Email   string `json:"email"`
				Name    string `json:"name"`
				Picture string `json:"picture"`
			}
			if err := json.Unmarshal(body, &data); err != nil {
				return nil, err
			}
			return &UserInfo{
				ProviderID: data.Sub,
				Email:      data.Email,
				Name:       data.Name,
				Picture:    data.Picture,
			}, nil
		},
	}
}

// NewGitHubProvider creates a GitHub OAuth provider.
func NewGitHubProvider(clientID, clientSecret, baseURL string) *OAuthProvider {
	return &OAuthProvider{
		Name: "github",
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  baseURL + "/api/auth/github/callback",
			Scopes:       []string{"user:email", "read:user"},
			Endpoint:     github.Endpoint,
		},
		UserInfoURL: "https://api.github.com/user",
		ParseUser: func(body []byte) (*UserInfo, error) {
			var data struct {
				ID        int64  `json:"id"`
				Email     string `json:"email"`
				Name      string `json:"name"`
				Login     string `json:"login"`
				AvatarURL string `json:"avatar_url"`
			}
			if err := json.Unmarshal(body, &data); err != nil {
				return nil, err
			}

			name := data.Name
			if name == "" {
				name = data.Login
			}

			return &UserInfo{
				ProviderID: fmt.Sprintf("%d", data.ID),
				Email:      data.Email,
				Name:       name,
				Picture:    data.AvatarURL,
			}, nil
		},
	}
}

// NewAppleProvider creates an Apple OAuth provider.
// Note: Apple Sign-In requires additional configuration and a private key.
func NewAppleProvider(clientID, clientSecret, baseURL string) *OAuthProvider {
	return &OAuthProvider{
		Name: "apple",
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  baseURL + "/api/auth/apple/callback",
			Scopes:       []string{"name", "email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://appleid.apple.com/auth/authorize",
				TokenURL: "https://appleid.apple.com/auth/token",
			},
		},
		// Apple doesn't have a user info endpoint - info comes from ID token
		UserInfoURL: "",
		ParseUser: func(body []byte) (*UserInfo, error) {
			// Apple returns user info in the ID token, not from a user info endpoint
			// This is handled specially in the callback
			var data struct {
				Sub   string `json:"sub"`
				Email string `json:"email"`
				Name  struct {
					FirstName string `json:"firstName"`
					LastName  string `json:"lastName"`
				} `json:"name"`
			}
			if err := json.Unmarshal(body, &data); err != nil {
				return nil, err
			}

			name := ""
			if data.Name.FirstName != "" || data.Name.LastName != "" {
				name = data.Name.FirstName + " " + data.Name.LastName
			}

			return &UserInfo{
				ProviderID: data.Sub,
				Email:      data.Email,
				Name:       name,
			}, nil
		},
	}
}

// NewMicrosoftProvider creates a Microsoft OAuth provider.
func NewMicrosoftProvider(clientID, clientSecret, baseURL string) *OAuthProvider {
	return &OAuthProvider{
		Name: "microsoft",
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  baseURL + "/api/auth/microsoft/callback",
			Scopes:       []string{"openid", "profile", "email", "User.Read"},
			Endpoint:     microsoft.AzureADEndpoint("common"),
		},
		UserInfoURL: "https://graph.microsoft.com/v1.0/me",
		ParseUser: func(body []byte) (*UserInfo, error) {
			var data struct {
				ID                string `json:"id"`
				Mail              string `json:"mail"`
				UserPrincipalName string `json:"userPrincipalName"`
				DisplayName       string `json:"displayName"`
			}
			if err := json.Unmarshal(body, &data); err != nil {
				return nil, err
			}

			email := data.Mail
			if email == "" {
				email = data.UserPrincipalName
			}

			return &UserInfo{
				ProviderID: data.ID,
				Email:      email,
				Name:       data.DisplayName,
			}, nil
		},
	}
}
