package oauth_config

import (
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

// GitHubUser represents the relevant user data fetched from GitHub's API.
type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// Global variables for GitHub OAuth configuration for CLIENTS.
var (
	ClientGithubOauthConfig *oauth2.Config
	ClientOauthStateString  string
)

// InitializeClientGitHubOAuthConfig is a helper to set up the OAuth config.
// Call this from your main.go or config initialization.
func InitializeClientGitHubOAuthConfig(clientID, clientSecret, redirectURL string) {
	ClientGithubOauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"user:email", "read:user"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
	// Generate a secure random state string for CSRF protection
	ClientOauthStateString = "random-client-auth-state-" + fmt.Sprintf("%d", time.Now().UnixNano())
}
