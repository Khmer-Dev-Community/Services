package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	redis "websocket-service/config"

	"strconv"
	"sync"

	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("ihuegrbnor7nou3hu3uh3uh3")

type UserDTO struct {
	ID        int    `json:"id"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	RoleID    int    `json:"roleId"`
	CompanyID int    `json:"companyId"`
	Sex       string `json:"sex"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	Token     string `json:"token"`
}

func AuthMiddlewareWithWhiteList(whitelist map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := whitelist[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				RespondWithError(w, http.StatusUnauthorized, "Authorization header is missing")
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			// Declare the variable outside
			var user *UserDTO
			user, err := decryptToken(tokenString) // Use assignment instead of short declaration
			if err != nil {
				ErrorLog("Error decrypting token", err.Error())
				RespondWithError(w, http.StatusUnauthorized, "Invalid token")
				return
			}
			r.Header.Set("token-x-userid", fmt.Sprintf("%d", user.ID))
			r.Header.Set("roleid", fmt.Sprintf("%d", user.RoleID))
			r.Header.Set("companyid", fmt.Sprintf("%d", user.CompanyID))
			// Append Body Data

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					ErrorLog("Unexpected signing method", "Error")
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})
			if err != nil {
				ErrorLog("Error parsing token", err.Error())
				RespondWithError(w, http.StatusUnauthorized, "Invalid token format")
				return
			}
			if !token.Valid {

				RespondWithError(w, http.StatusUnauthorized, "Invalid token")
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			userID := int(claims["id"].(float64))
			userData, err := redis.Get(fmt.Sprintf("user:%d", userID))

			if err != nil {
				RespondWithError(w, http.StatusUnauthorized, "Token expired or not found")
				return
			}

			if err := json.Unmarshal([]byte(userData), &user); err != nil {
				// Handle error if unable to decode JSON
				RespondWithError(w, http.StatusInternalServerError, "Failed to decode user data")
				return
			}
			// check request token & redis token
			if tokenString != user.Token {
				RespondWithError(w, http.StatusUnauthorized, "Your Account login on other device ")
				return
			}
			key := fmt.Sprintf("user:%d", userID)
			expiration := time.Minute * 30
			if err := redis.UpdateExpiration(key, expiration); err != nil {
				RespondWithError(w, http.StatusInternalServerError, "Failed to refresh token expiration")
				return
			}

			// append body
			if r.Method == http.MethodPost || r.Method == http.MethodPut {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					RespondWithError(w, http.StatusInternalServerError, "Failed to read request body")
					return
				}
				defer r.Body.Close() // Close the original body

				// Unmarshal the existing request body
				var requestBody map[string]interface{}
				if err := json.Unmarshal(body, &requestBody); err != nil {
					RespondWithError(w, http.StatusBadRequest, "Invalid request body")
					return
				}

				// Append companyId and userId to the request body
				requestBody["companyId"] = user.CompanyID
				requestBody["userId"] = user.ID

				// Marshal the updated request body
				updatedBody, err := json.Marshal(requestBody)
				if err != nil {
					RespondWithError(w, http.StatusInternalServerError, "Failed to marshal updated body")
					return
				}

				// Reset the request body to the updated body
				r.Body = ioutil.NopCloser(bytes.NewBuffer(updatedBody))
				r.Header.Set("Content-Type", "application/json")

				// Log the updated request body
				InfoLog(string(updatedBody), "Updated request body")
			}

			next.ServeHTTP(w, r)
		})
	}
}
func AuthMiddlewareWithWhiteList11(whitelist map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the request path is in the whitelist
			if _, ok := whitelist[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}

			// Retrieve the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				RespondWithError(w, http.StatusUnauthorized, "Authorization header is missing")
				return
			}

			// Parse the token from the Authorization header
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			var user *UserDTO
			user, err := decryptToken(tokenString) // Use assignment instead of short declaration
			if err != nil {
				ErrorLog("Error decrypting token", err.Error())
				RespondWithError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Set user information in the request header
			r.Header.Set("id", fmt.Sprintf("%d", user.ID))
			r.Header.Set("roleid", fmt.Sprintf("%d", user.RoleID))
			r.Header.Set("companyid", fmt.Sprintf("%d", user.CompanyID))

			// Validate the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					ErrorLog("Unexpected signing method", "Error")
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secretKey), nil
			})
			if err != nil || !token.Valid {
				RespondWithError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Handle claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			userID := int(claims["id"].(float64))
			userData, err := redis.Get(fmt.Sprintf("user:%d", userID))
			if err != nil {
				RespondWithError(w, http.StatusUnauthorized, "Token expired or not found")
				return
			}

			// Unmarshal user data from Redis
			if err := json.Unmarshal([]byte(userData), &user); err != nil {
				RespondWithError(w, http.StatusInternalServerError, "Failed to decode user data")
				return
			}

			// Check if the token matches the stored token
			if tokenString != user.Token {
				RespondWithError(w, http.StatusUnauthorized, "Your account is logged in on another device")
				return
			}

			// Refresh the token expiration
			key := fmt.Sprintf("user:%d", userID)
			expiration := time.Minute * 30
			if err := redis.UpdateExpiration(key, expiration); err != nil {
				RespondWithError(w, http.StatusInternalServerError, "Failed to refresh token expiration")
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func decryptToken(tokenString string) (*UserDTO, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	// Validate the token
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Extract user information
	userID, ok := claims["id"].(float64)
	if !ok {
		return nil, errors.New("missing or invalid user ID in token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("missing or invalid username in token claims")
	}

	// Populate and return the user DTO
	return &UserDTO{
		ID:     int(userID),                      // Convert from float64 to int
		RoleID: int(claims["role_id"].(float64)), // Convert from float64 to int

		Username:  username,                            // Extract as string
		CompanyID: int(claims["company_id"].(float64)), // Convert from float64 to int
		// Extract and set additional fields as needed
	}, nil
}

type RolePermissions struct {
	Permissions []string `json:"permissions"`
	Roles       []string `json:"roles"`
}

var (
	ctx             = context.Background()
	permissionCache = make(map[string]RolePermissions)
	cacheMutex      sync.RWMutex
	cacheTTL        = 10 * time.Minute // Cache time-to-live
)

func HasPermission(requiredPermission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleIdStr := r.Header.Get("Roleid") // Assuming role ID is passed in the header
			if roleIdStr == "" {
				RespondWithError(w, http.StatusUnauthorized, "Request Unauthorized")
				return
			}

			// Check in-memory cache first
			cacheMutex.RLock()
			rolePermissions, found := permissionCache[roleIdStr]
			cacheMutex.RUnlock()

			// Convert role ID from string to int
			roleId, err := strconv.Atoi(roleIdStr)
			if err != nil {
				RespondWithError(w, http.StatusUnauthorized, "Unauthorized, Invalid Role ID")
				return
			}

			// If not found in cache, fetch from Redis
			if !found {
				// Fetch user permissions from Redis
				permissions, err := redis.Get(fmt.Sprintf("role:%d", roleId))
				if err != nil {
					RespondWithError(w, http.StatusInternalServerError, "Unauthorized, Error fetching permissions")
					return
				}

				// Unmarshal permissions JSON into RolePermissions struct
				if err := json.Unmarshal([]byte(permissions), &rolePermissions); err != nil {
					log.Printf("Error unmarshalling permissions: %v", err)
					RespondWithError(w, http.StatusInternalServerError, "Unauthorized, Internal Server Error")
					return
				}

				// Cache the permissions
				cacheMutex.Lock()
				permissionCache[roleIdStr] = rolePermissions
				cacheMutex.Unlock()

				// Start a goroutine to remove the entry after TTL
				go func(roleIdStr string) {
					time.Sleep(cacheTTL)
					cacheMutex.Lock()
					delete(permissionCache, roleIdStr)
					cacheMutex.Unlock()
				}(roleIdStr)
			}

			// Check if user has the required permission
			hasPermission := false
			for _, perm := range rolePermissions.Permissions {
				if perm == requiredPermission {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				RespondWithError(w, http.StatusForbidden, "Unauthorized, Forbidden")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
