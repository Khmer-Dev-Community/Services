package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil" // Used for backwards compatibility, io.ReadAll is preferred for Go 1.16+
	"log"       // For log.Printf in the caching goroutine
	"net/http"  // For HTTP status codes and methods
	"strconv"   // For string to int conversion
	"strings"   // For string manipulation (e.g., strings.TrimPrefix, strings.HasPrefix)
	"sync"      // For mutex in caching
	"time"      // For time-related operations (e.g., cache TTL)

	redis "github.com/Khmer-Dev-Community/Services/api-service/config" // Assuming this imports your Redis client setup

	"github.com/gin-gonic/gin"     // Gin framework import
	"github.com/golang-jwt/jwt/v4" // JWT library
)

// secretKey is used for signing and verifying JWT tokens.
// IMPORTANT: In a production environment, this should be loaded from
// environment variables or a secure configuration management system,
// not hardcoded.
var secretKey = []byte("ihuegrbnor7nou3hu3uh3uh3")

// UserDTO represents the structure of user data, typically used for
// transferring user information between layers or for JWT claims.
type UserDTO struct {
	ID        uint   `json:"id"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	RoleID    uint   `json:"roleId"`
	CompanyID uint   `json:"companyId"`
	Sex       string `json:"sex"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
	Token     string `json:"token"`
}

func AuthMiddlewareWithWhiteList(whitelist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Check if the current request path is in the whitelist.
		if IsWhitelisted(c.Request.URL.Path, whitelist) {
			c.Next() // If whitelisted, skip authentication and proceed to the next handler.
			return
		}
		var jwtToken string
		cookie, err := c.Request.Cookie("kdc.secure.token")
		if err == nil && cookie != nil {
			jwtToken = cookie.Value
			// Optionally set a header for compatibility if other parts of your system expect it.
			c.Request.Header.Set("kdc-x-token", jwtToken)
		} else {
			// If no cookie, try to get from Authorization header (Bearer token).
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				jwtToken = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				// If no token found in cookie or header, respond with Unauthorized.
				ErrorLog(map[string]interface{}{"path": c.Request.URL.Path}, "Unauthorized: No token found in cookie or Authorization header")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User authentication required: token not found"})
				return
			}
		}

		_, err = decryptToken(jwtToken)
		if err != nil {
			ErrorLog(map[string]interface{}{"error": err.Error(), "token_snippet": jwtToken[:min(len(jwtToken), 20)]}, "Error decrypting token or invalid token format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// 4. Parse and validate the full JWT token, including claims.
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil // Return the secret key for validation.
		})
		if err != nil || !token.Valid {
			ErrorLog(map[string]interface{}{"error": err.Error(), "token_valid": token.Valid}, "Error parsing or validating JWT token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// 5. Extract claims and verify user ID.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ErrorLog(map[string]interface{}{}, "Invalid token claims format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userIDFloat, ok := claims["id"].(float64)
		if !ok {
			ErrorLog(map[string]interface{}{"claims": claims}, "User ID not found or invalid type in token claims")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: User ID missing in token"})
			return
		}
		userID := uint(userIDFloat) // Convert from float64 to uint.

		// 6. Fetch full user data from Redis using the user ID from claims.
		redisKey := fmt.Sprintf("user:%d", userID)
		userDataJSON, err := redis.Get(redisKey)
		if err != nil {
			ErrorLog(map[string]interface{}{"user_id": userID, "error": err.Error()}, "Token expired or user data not found in Redis")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired or not found"})
			return
		}

		// 7. Unmarshal Redis user data into UserDTO.
		var storedUser UserDTO
		if err := json.Unmarshal([]byte(userDataJSON), &storedUser); err != nil {
			ErrorLog(map[string]interface{}{"user_id": userID, "error": err.Error()}, "Failed to decode user data from Redis")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user data"})
			return
		}
		if jwtToken != storedUser.Token {
			ErrorLog(map[string]interface{}{"user_id": userID, "token_match": false}, "Unauthorized: Account logged in on another device")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Your account is logged in on another device"})
			return
		}

		// 9. Update token expiration in Redis to keep the session alive.
		expiration := time.Minute * 30
		if err := redis.UpdateExpiration(redisKey, expiration); err != nil {
			ErrorLog(map[string]interface{}{"user_id": userID, "error": err.Error()}, "Failed to refresh token expiration in Redis")
			// Decide if this should block the request. Typically, this would be a warning.
			// However, adhering to your original logic, it aborts with a 500 error.
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token expiration"})
			return
		}
		c.Set("userID", storedUser.ID)
		c.Set("roleID", storedUser.RoleID)
		c.Set("companyID", storedUser.CompanyID)
		c.Set("username", storedUser.Username)
		c.Set("userDTO", storedUser) // Store the full DTO if needed
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
			bodyBytes, err := ioutil.ReadAll(c.Request.Body) // Read the original body bytes.
			if err != nil {
				ErrorLog(map[string]interface{}{"error": err.Error()}, "Failed to read request body in AuthMiddleware")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
				return
			}
			c.Request.Body.Close() // Close the original request body.

			var requestBody map[string]interface{}
			if len(bodyBytes) > 0 { // Only attempt to unmarshal if the body is not empty.
				if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
					ErrorLog(map[string]interface{}{"error": err.Error(), "body_snippet": string(bodyBytes[:min(len(bodyBytes), 50)])}, "Invalid request body JSON in AuthMiddleware")
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
					return
				}
			} else {
				requestBody = make(map[string]interface{}) // Initialize an empty map if body was empty.
			}

			// Add or update companyId and userId in the request body map.
			requestBody["companyId"] = storedUser.CompanyID
			requestBody["userId"] = storedUser.ID

			updatedBodyBytes, err := json.Marshal(requestBody) // Marshal the modified map back to JSON bytes.
			if err != nil {
				ErrorLog(map[string]interface{}{"error": err.Error()}, "Failed to marshal updated request body in AuthMiddleware")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request body"})
				return
			}

			// Reset the request body for subsequent handlers/controllers to read the modified version.
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(updatedBodyBytes))
			c.Request.ContentLength = int64(len(updatedBodyBytes))   // Set content length to the new body size.
			c.Request.Header.Set("Content-Type", "application/json") // Ensure Content-Type is correct.

			InfoLog(map[string]interface{}{"updated_body_snippet": string(updatedBodyBytes[:min(len(updatedBodyBytes), 100)])}, "Request body updated in AuthMiddleware")
		}

		c.Next() // Proceed to the next middleware or handler in the chain.
	}
}

// decryptToken parses and validates a JWT token string and extracts basic UserDTO information.
// It's a helper function for the authentication middleware.
func decryptToken(tokenString string) (*UserDTO, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil // Return the secret key for verification.
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims format")
	}

	// Extract claims and perform type assertions.
	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		return nil, errors.New("missing or invalid user ID in token claims")
	}
	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("missing or invalid username in token claims")
	}
	roleIDFloat, ok := claims["role_id"].(float64)
	if !ok {
		return nil, errors.New("missing or invalid role_id in token claims")
	}
	companyIDFloat, ok := claims["company_id"].(float64)
	if !ok {
		return nil, errors.New("missing or invalid company_id in token claims")
	}
	tokenFromClaims, ok := claims["token"].(string) // Assuming the token itself is stored in claims for comparison
	if !ok {
		// Log a warning if 'token' claim is missing, but don't fail unless critical for logic
		log.Printf("Warning: 'token' claim missing or invalid in JWT for user %v", userIDFloat)
		tokenFromClaims = "" // Default to empty string if not found
	}

	return &UserDTO{
		ID:        uint(userIDFloat),
		RoleID:    uint(roleIDFloat),
		Username:  username,
		CompanyID: uint(companyIDFloat),
		Token:     tokenFromClaims, // Include the token from claims for comparison
	}, nil
}

// IsWhitelisted checks if the given path matches any path in the whitelist.
// This supports basic prefix matching for simplicity.
func IsWhitelisted(path string, whitelist []string) bool {
	for _, wp := range whitelist {
		// Check for exact match or prefix match (e.g., "/swagger/" should match "/swagger/index.html")
		if path == wp || strings.HasPrefix(path, wp) {
			return true
		}
	}
	return false
}

// min is a helper function to get the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RolePermissions defines the structure for role permissions data.
type RolePermissions struct {
	Permissions []string `json:"permissions"`
	Roles       []string `json:"roles"` // Not used in current logic, but kept for completeness
}

// Global cache variables for permissions.
var (
	// ctx is a background context, consider passing context from Gin's c.Request.Context() for better traceability.
	ctx             = context.Background()
	permissionCache = make(map[string]RolePermissions)
	cacheMutex      sync.RWMutex
	cacheTTL        = 10 * time.Minute // Time-to-live for cached permissions.
)

// HasPermission is a Gin middleware that checks if the authenticated user
// (identified by roleID in context) has a specific required permission.
func HasPermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Retrieve role ID from Gin context, which should have been set by AuthMiddlewareWithWhiteList.
		roleIDAny, exists := c.Get("roleID")
		if !exists {
			ErrorLog(map[string]interface{}{"path": c.Request.URL.Path}, "Role ID not found in Gin context for permission check")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Role information missing"})
			return
		}
		roleID, ok := roleIDAny.(uint) // Assert the type to uint.
		if !!ok {
			ErrorLog(map[string]interface{}{"roleID_type": fmt.Sprintf("%T", roleIDAny)}, "Invalid Role ID type in Gin context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error: Invalid Role ID format"})
			return
		}
		roleIdStr := strconv.FormatUint(uint64(roleID), 10) // Convert uint to string for map key.

		// 2. Check in-memory cache first for role permissions.
		cacheMutex.RLock()
		rolePermissions, found := permissionCache[roleIdStr]
		cacheMutex.RUnlock()

		// 3. If permissions not found in cache, fetch from Redis.
		if !found {
			permissionsJSON, err := redis.Get(fmt.Sprintf("role:%d", roleID)) // Fetch permissions JSON from Redis.
			if err != nil {
				ErrorLog(map[string]interface{}{"role_id": roleID, "error": err.Error()}, "Error fetching permissions from Redis")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized: Error fetching permissions"})
				return
			}

			// Unmarshal permissions JSON into RolePermissions struct.
			if err := json.Unmarshal([]byte(permissionsJSON), &rolePermissions); err != nil {
				ErrorLog(map[string]interface{}{"role_id": roleID, "error": err.Error()}, "Error unmarshalling permissions from Redis")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized: Internal Server Error"})
				return
			}

			// 4. Cache the retrieved permissions.
			cacheMutex.Lock()
			permissionCache[roleIdStr] = rolePermissions
			cacheMutex.Unlock()

			// 5. Start a goroutine to remove the cached entry after TTL.
			go func(roleIdStr string) {
				time.Sleep(cacheTTL) // Wait for the cache TTL.
				cacheMutex.Lock()
				delete(permissionCache, roleIdStr) // Remove the entry from cache.
				cacheMutex.Unlock()
				InfoLog(map[string]interface{}{"role_id_str": roleIdStr}, "Permissions cache entry expired and removed")
			}(roleIdStr)
		}

		// 6. Check if the user's role has the required permission.
		hasPermission := false
		for _, perm := range rolePermissions.Permissions {
			if perm == requiredPermission {
				hasPermission = true
				break
			}
		}

		// 7. If permission is not found, respond with Forbidden.
		if !hasPermission {
			InfoLog(map[string]interface{}{"role_id": roleID, "required_permission": requiredPermission, "path": c.Request.URL.Path}, "Forbidden: User does not have required permission")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized: Forbidden"})
			return
		}

		c.Next()
	}
}
