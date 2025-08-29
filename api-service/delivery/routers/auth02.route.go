package router

import (
	"log"
	"net/http"

	_ "github.com/Khmer-Dev-Community/Services/api-service/cmd/docs"

	"github.com/gin-gonic/gin" // Gin framework import

	"github.com/Khmer-Dev-Community/Services/api-service/auth02" // The auth02 package (controller)
	"github.com/Khmer-Dev-Community/Services/api-service/pkg/oauth_config"
)

type ClientAuthControllerWrapper struct {
	clientAuthController *auth02.ClientAuthController
}

// NewClientAuthControllerWrapper initializes the wrapper for ClientAuthController
func NewClientAuthControllerWrapper(cac *auth02.ClientAuthController) *ClientAuthControllerWrapper {
	return &ClientAuthControllerWrapper{
		clientAuthController: cac,
	}
}

func SetupRouterAuth02(r *gin.Engine, clientAuthCtrl *auth02.ClientAuthController) {

	clientAuthControllerWrapper := NewClientAuthControllerWrapper(clientAuthCtrl)
	githubClientID := "Ov23liBVXaZ0bV6B43Ut"
	githubClientSecret := "28f9c091b494fc8bc1fd2b795c77be27b322f2c4"
	githubRedirectURL := "http://localhost:3000/api/auth02/github/callback"

	if githubClientID == "" || githubClientSecret == "" || githubRedirectURL == "" {
		log.Println("WARNING: GitHub OAuth environment variables (GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET, GITHUB_REDIRECT_URL) are not fully set. GitHub login might not work.")
	}

	oauth_config.InitializeClientGitHubOAuthConfig(githubClientID, githubClientSecret, githubRedirectURL)

	// Define API routes group (e.g., /api)
	api := r.Group("/api")
	clientAuthRouter := api.Group("/auth02") // Changed to Gin's Group method

	// Public routes (no authentication required)
	clientAuthRouter.POST("/register", clientAuthControllerWrapper.clientAuthController.RegisterClient) // Gin's POST method
	clientAuthRouter.POST("/login", clientAuthControllerWrapper.clientAuthController.ClientLogin)       // Gin's POST method

	// GitHub OAuth routes
	clientAuthRouter.GET("/github/login", clientAuthControllerWrapper.clientAuthController.GitHubLoginRedirect) // Gin's GET method
	clientAuthRouter.GET("/github/callback", clientAuthControllerWrapper.clientAuthController.GitHubCallback)   // Gin's GET method

	clientAuthRouter.GET("/profile", clientAuthControllerWrapper.clientAuthController.GetClientProfile)
	clientAuthRouter.PUT("/profile", clientAuthControllerWrapper.clientAuthController.UpdateClientProfile)

	api.GET("/health", func(c *gin.Context) { // Changed to Gin's GET method and handler signature
		c.JSON(http.StatusOK, gin.H{"status": "up", "message": "Telegram Service API is running!"}) // Use c.JSON
	})

	// Root route for Gin
	api.GET("/", func(c *gin.Context) { // Changed to Gin's GET method and handler signature
		c.JSON(http.StatusOK, gin.H{"message": "Telegram Service API is running!"}) // Use c.JSON
	})
}
