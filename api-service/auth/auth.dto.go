// auth_dto.go
package auth

import (
	users "github.com/Khmer-Dev-Community/Services/api-service/lib/users/models"

	"golang.org/x/oauth2"
)

// LoginRequestDTO represents the structure of a login request
type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LogoutRequestDTO represents the structure of a logout request
type LogoutRequestDTO struct {
	Token string `json:"token"`
}

// GitHubCallback represents the data received from GitHub's callback
type GitHubCallback struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// GitHubUser represents the relevant data fetched from GitHub's /user API
type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// AuthResponse will be used for both normal login and GitHub login
type AuthResponse struct {
	Token string      `json:"token"`
	User  *users.User `json:"user"` // Assuming users.User is your internal user model
}

// Define your GitHub OAuth configuration globally or pass it
// This should ideally be loaded from environment variables
var (
	GithubOauthConfig *oauth2.Config
	OauthStateString  string // This should be generated dynamically per request in production
)
