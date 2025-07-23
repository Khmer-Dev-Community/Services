// auth_service.go
package auth

import (
	"errors"
	"fmt"
	user "telegram-service/lib/users/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AuthService handles authentication-related operations
type AuthService struct {
	userRepository *UserRepository
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userRepository *UserRepository) *AuthService {
	return &AuthService{userRepository: userRepository}
}

func (s *AuthService) Login(credentials *LoginCredentials) (*user.User, error) {
	user, err := s.userRepository.GetUserByUsername(credentials.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	return user, nil
}
func (s *AuthService) GenerateToken(user *user.User) (string, error) {
	// Define the claims for the token
	claims := jwt.MapClaims{
		"sub":        user.ID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":        time.Now().Unix(),
		"username":   user.Username,
		"id":         user.ID,
		"role_id":    user.RoleID,
		"company_id": user.CompanyID,
	}

	// Generate the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	tokenString, err := token.SignedString([]byte("ihuegrbnor7nou3hu3uh3uh3"))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}
