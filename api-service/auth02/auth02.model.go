package auth02

import (
	"github.com/Khmer-Dev-Community/Services/api-service/lib/userclient"
)

// ClientLoginCredentials represents the request body for client user login.
// This directly aliases the DTO type from the userclient package.
type ClientLoginCredentials userclient.ClientLoginRequestDTO

// ClientRegisterPayload represents the request body for client user registration.
// This directly aliases the DTO type from the userclient package.
type ClientRegisterPayload userclient.ClientRegisterRequestDTO

type ClientAuthResponse struct {
	Token string                            `json:"token"`
	User  *userclient.ClientUserResponseDTO `json:"user"`
}

// GitHubCallback represents the query parameters received from GitHub's OAuth callback.
// Keeping it here as it's a DTO specifically for the callback handler.
type GitHubCallback struct {
	Code  string `json:"code" form:"code"`
	State string `json:"state" form:"state"`
}
