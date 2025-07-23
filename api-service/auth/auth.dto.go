// auth_dto.go
package auth

// LoginRequestDTO represents the structure of a login request
type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LogoutRequestDTO represents the structure of a logout request
type LogoutRequestDTO struct {
	Token string `json:"token"`
}
