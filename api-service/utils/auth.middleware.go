package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	redis "github.com/Khmer-Dev-Community/Services/api-service/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// A placeholder for your actual dependencies.
var secretKey = []byte("ihuegrbnor7nou3hu3uh3uh3")

// UserDTO is a placeholder for your user data transfer object.
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

func IsWhitelisted(path string, whitelist map[string]bool) bool {
	_, ok := whitelist[path]
	return ok
}

// AuthMiddlewareWithWhiteList is a Gin middleware for authenticating requests with a whitelist.
func AuthMiddlewareWithWhiteList(whitelist map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Check for whitelisted paths first
		fullRoute := c.FullPath()
		if IsWhitelisted(fullRoute, whitelist) || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// 2. Extract JWT from either a cookie or the Authorization header
		var jwtToken string
		if cookie, err := c.Cookie("kdc.secure.token"); err == nil && cookie != "" {
			jwtToken = cookie
		}
		/*else {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User authentication required: token not found"})
				return
			}
			jwtToken = strings.TrimPrefix(authHeader, "Bearer ")
		}*/

		// 3. Parse and validate JWT
		token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// 4. Extract claims and validate
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}
		userIDFloat, ok := claims["id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User ID missing in token"})
			return
		}
		userID := uint(userIDFloat)

		// 5. Validate token against Redis
		redisKey := fmt.Sprintf("user:%d", userID)
		userDataJSON, err := redis.Get(redisKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired or not found"})
			return
		}

		var storedUser UserDTO
		if err := json.Unmarshal([]byte(userDataJSON), &storedUser); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user data"})
			return
		}

		if jwtToken != storedUser.Token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Your account is logged in on another device"})
			return
		}

		// 6. Refresh token expiration in Redis
		_ = redis.UpdateExpiration(redisKey, 720*time.Hour)
		c.Set("userID", userID)
		c.Set("roleID", storedUser.RoleID)
		c.Set("companyID", storedUser.CompanyID)
		c.Set("username", storedUser.Username)
		c.Set("userDTO", storedUser)

		// 7. Inject companyId & userId into request body
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch || c.Request.Method == http.MethodDelete {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
				return
			}
			c.Request.Body.Close()

			var requestBody map[string]interface{}
			if len(bodyBytes) > 0 {
				if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
					return
				}
			} else {
				requestBody = make(map[string]interface{})
			}

			requestBody["companyId"] = storedUser.CompanyID
			requestBody["userId"] = storedUser.ID

			updatedBodyBytes, err := json.Marshal(requestBody)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal updated body"})
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(updatedBodyBytes))
			c.Request.ContentLength = int64(len(updatedBodyBytes))
		}
		c.Next()
	}
}

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
