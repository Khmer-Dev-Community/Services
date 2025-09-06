package auth02

import (
	// Still useful for marshaling/unmarshaling if needed explicitly
	"fmt"      // For logging in this file, Gin has its own logging.
	"net/http" // Still needed for http.StatusOK, http.StatusBadRequest etc.
	"time"

	"github.com/Khmer-Dev-Community/Services/api-service/config"
	redis "github.com/Khmer-Dev-Community/Services/api-service/config"
	clientauth_service "github.com/Khmer-Dev-Community/Services/api-service/lib/clientauth"
	"github.com/Khmer-Dev-Community/Services/api-service/lib/userclient"
	"github.com/Khmer-Dev-Community/Services/api-service/pkg/oauth_config"
	"github.com/Khmer-Dev-Community/Services/api-service/utils" // Your utility functions

	"github.com/gin-gonic/gin" // Gin framework
)

// ClientAuthController handles HTTP requests related to client authentication.
type ClientAuthController struct {
	clientAuthService *clientauth_service.ClientAuthService
}

// NewClientAuthController creates a new instance of ClientAuthController.
func NewClientAuthController(clientAuthService *clientauth_service.ClientAuthService) *ClientAuthController {
	return &ClientAuthController{clientAuthService: clientAuthService}
}

// RegisterClient handles the client registration request.
// @Summary Register a new client user
// @Description Registers a new client user with username, email, and password.
// @Tags Client Authentication
// @Accept json
// @Produce json
// @Param client body ClientRegisterPayload true "Client registration payload"
// @Success 201 {object} ClientAuthResponse "User successfully registered and logged in"
// @Failure 400 {object} gin.H "Invalid request payload"
// @Failure 409 {object} gin.H "Username or email already taken"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /client/register [post]
func (ctrl *ClientAuthController) RegisterClient(c *gin.Context) { // Changed signature
	var payload ClientRegisterPayload
	// Use c.ShouldBindJSON for robust JSON binding and validation
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.LoggerRequest(map[string]interface{}{"error": err.Error()}, "Client Registration", "Invalid request payload")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request payload: %v", err.Error())})
		return
	}

	newUser, err := ctrl.clientAuthService.RegisterClient((*userclient.ClientRegisterRequestDTO)(&payload))
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"error": err.Error()}, fmt.Sprintf("Client registration failed for %s", payload.Username))
		if err.Error() == "username already taken" || err.Error() == "email already registered" {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	token, err := ctrl.clientAuthService.GenerateToken(newUser)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"error": err.Error()}, fmt.Sprintf("Failed to generate token for client user %s", newUser.Username))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	userResponse := userclient.ToClientUserResponseDTO(newUser)
	response := ClientAuthResponse{
		Token: token,
		User:  userResponse,
	}

	c.JSON(http.StatusCreated, response) // Use c.JSON for success response
	utils.LoggerRequest(map[string]interface{}{"user_id": userResponse.ID}, "Client Registration", fmt.Sprintf("Client user registered and logged in: %s", userResponse.Username))
}

// ClientLogin handles the client login request.
// @Summary Log in a client user
// @Description Authenticates a client user with username/email and password.
// @Tags Client Authentication
// @Accept json
// @Produce json
// @Param client body ClientLoginCredentials true "Client login credentials"
// @Success 200 {object} ClientAuthResponse "User successfully logged in"
// @Failure 400 {object} gin.H "Invalid request payload"
// @Failure 401 {object} gin.H "Invalid credentials"
// @Failure 403 {object} gin.H "Account is inactive or blocked / This account uses social login"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /client/login [post]
func (ctrl *ClientAuthController) ClientLogin(c *gin.Context) { // Changed signature
	var credentials ClientLoginCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		utils.LoggerRequest(map[string]interface{}{"error": err.Error()}, "Client Login", "Invalid request payload")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request payload: %v", err.Error())})
		return
	}

	user, err := ctrl.clientAuthService.ClientLogin((*userclient.ClientLoginRequestDTO)(&credentials))
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"error": err.Error()}, fmt.Sprintf("Client login failed for %s", credentials.Identifier))
		if err.Error() == "invalid credentials" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		if err.Error() == "account is inactive or blocked" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Account is inactive or blocked"})
			return
		}
		if err.Error() == "this account uses social login, please use the GitHub login option" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to log in"})
		return
	}

	token, err := ctrl.clientAuthService.GenerateToken(user)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"error": err.Error()}, fmt.Sprintf("Failed to generate token for client user %s", user.Username))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	userResponse := userclient.ToClientUserResponseDTO(user)
	response := ClientAuthResponse{
		Token: token,
		User:  userResponse,
	}

	c.JSON(http.StatusOK, response) // Use c.JSON for success response
	utils.LoggerRequest(map[string]interface{}{"user_id": userResponse.ID}, "Client Login", fmt.Sprintf("Client user logged in: %s", userResponse.Username))
}

// GitHubLoginRedirect handles the initiation of GitHub OAuth login.
// @Summary Redirect to GitHub for OAuth login
// @Description Redirects the user to GitHub's authorization page.
// @Tags Client Authentication
// @Produce html
// @Success 307 {string} string "Redirect to GitHub"
// @Failure 500 {object} gin.H "GitHub OAuth not configured"
// @Router /client/github/login [get]
func (ctrl *ClientAuthController) GitHubLoginRedirect(c *gin.Context) { // Changed signature
	if oauth_config.ClientGithubOauthConfig == nil {
		utils.ErrorLog(nil, "GitHub OAuth configuration not initialized.")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "GitHub OAuth not configured"})
		return
	}
	url := oauth_config.ClientGithubOauthConfig.AuthCodeURL(oauth_config.ClientOauthStateString)
	utils.LoggerRequest(map[string]interface{}{"url": url}, "GitHub OAuth Redirect", "Redirecting client to GitHub for authorization")
	c.Redirect(http.StatusTemporaryRedirect, url) // Use c.Redirect
}

// GitHubCallback handles the callback from GitHub after successful authorization.
// @Summary Handle GitHub OAuth callback
// @Description Processes the callback from GitHub after user authorization. Exchanges code for token and logs in/registers user.
// @Tags Client Authentication
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from GitHub"
// @Param state query string true "State parameter for CSRF protection"
// @Success 200 {object} ClientAuthResponse "User successfully logged in or registered via GitHub"
// @Failure 400 {object} gin.H "Missing authorization code or state / Invalid state parameter"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /client/github/callback [get]
func (ctrl *ClientAuthController) GitHubCallback(c *gin.Context, cfg *config.GitConfig) { // Changed signature
	code := c.Query("code")   // Use c.Query
	state := c.Query("state") // Use c.Query

	if code == "" || state == "" {
		utils.LoggerRequest(map[string]interface{}{"code": code, "state": state}, "GitHub Callback", "Missing code or state in query parameters")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code or state"})
		return
	}

	if state != oauth_config.ClientOauthStateString {
		utils.WarnLog(map[string]interface{}{"state": state}, "GitHub callback: Invalid OAuth state parameter")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}

	token, err := ctrl.clientAuthService.ExchangeGitHubCodeForToken(code)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"error": err.Error()}, "GitHub callback: Failed to exchange code for token")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate with GitHub"})
		return
	}

	githubUser, err := ctrl.clientAuthService.GetGitHubUserData(token.AccessToken)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"error": err.Error()}, "GitHub callback: Failed to get GitHub user data")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user data from GitHub"})
		return
	}

	user, err := ctrl.clientAuthService.HandleGitHubLoginOrRegister(githubUser)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"error": err.Error()}, fmt.Sprintf("GitHub callback: Failed to handle login/register for GitHub user %s", githubUser.Login))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to process GitHub login/registration"})
		return
	}

	appToken, err := ctrl.clientAuthService.GenerateToken(user)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"error": err.Error()}, fmt.Sprintf("GitHub callback: Failed to generate app token for user %s", user.Username))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate application token"})
		return
	}

	userResponse := userclient.ToClientUserResponseDTO(user)
	response := ClientAuthResponse{
		Token: appToken,
		User:  userResponse,
	}

	expirationDuration := time.Hour * 24 // 24 hours
	key := fmt.Sprintf("user:%d", user.ID)
	user.Token = appToken

	if err := redis.SetWithExpiration(key, response, expirationDuration); err != nil {
		utils.ErrorLog(map[string]interface{}{"error": err.Error()}, "Failed to store user info in Redis")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user information in Redis"})
		return
	}

	// Set cookie using Gin's SetCookie
	c.SetCookie(
		"kdc.secure.token",                // name
		appToken,                          // value
		int(expirationDuration.Seconds()), // maxAge (seconds)
		"/",                               // path
		"",                                // domain (empty for current domain)
		c.Request.TLS != nil,              // secure (true if HTTPS, false for HTTP)
		false,                             // httpOnly
	)

	frontendRedirectURL := cfg.ClientEnd              // "http://localhost:8080/"   // Redirect to your frontend
	c.Redirect(http.StatusFound, frontendRedirectURL) // Use c.Redirect

	utils.LoggerRequest(map[string]interface{}{"user_id": userResponse.ID}, "GitHub Callback", fmt.Sprintf("Client user logged in/registered via GitHub: %s", userResponse.Username))
}

// GetClientProfile handles fetching a client user's profile.
// @Summary Get client user profile
// @Description Retrieves the profile of the authenticated client user.
// @Tags Client Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userclient.ClientUserResponseDTO "Client user profile"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Client user not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /client/profile [get]
func (ctrl *ClientAuthController) GetClientProfile(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		utils.ErrorLog(nil, "User ID not found in Gin context for GetClientProfile")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID missing from context"})
		return
	}

	userID, ok := userIDAny.(uint)
	if !ok {
		utils.ErrorLog(map[string]interface{}{"userID_type": fmt.Sprintf("%T", userIDAny)}, "Invalid User ID type in Gin context for GetClientProfile")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error: Invalid user ID format"})
		return
	}

	userDTO, err := ctrl.clientAuthService.GetClientProfile(userID)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"error": err.Error()}, fmt.Sprintf("Failed to retrieve profile for client ID %d", userID))
		if err.Error() == fmt.Sprintf("client user not found with ID: %d", userID) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Client user not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		return
	}
	utils.InfoLog(map[string]interface{}{"user_id": userID}, "Get Client Profile")

	utils.SuccessResponse(c, http.StatusOK, userDTO, "success")

}
func (ctrl *ClientAuthController) GetClientProfileByUsername(c *gin.Context) {
	username, exists := c.Params.Get("username")
	if !exists {
		utils.ErrorLog(nil, "User ID not found in Gin context for GetClientProfile")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID missing from context"})
		return
	}

	userDTO, err := ctrl.clientAuthService.GetClientProfileByUsername(username)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"error": err.Error()}, fmt.Sprintf("Failed to retrieve profile for client ID %d", username))
		if err.Error() == fmt.Sprintf("client user not found with ID: %d", username) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Client user not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		return
	}
	utils.InfoLog(map[string]interface{}{"user_id": username}, "Get Client Profile")

	utils.SuccessResponse(c, http.StatusOK, userDTO, "success")

}

// UpdateClientProfile handles updating a client user's profile.
// @Summary Update client user profile
// @Description Updates the profile of the authenticated client user.
// @Tags Client Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param client body userclient.ClientUpdateRequestDTO true "Client profile update data"
// @Success 200 {object} userclient.ClientUserResponseDTO "Client profile updated successfully"
// @Failure 400 {object} gin.H "Invalid request payload"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 404 {object} gin.H "Client user not found"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /client/profile [put]
func (ctrl *ClientAuthController) UpdateClientProfile(c *gin.Context) { // Changed signature
	// Retrieve userID from Gin context (set by AuthMiddleware)
	userIDAny, exists := c.Get("userID")
	if !exists {
		utils.ErrorLog(nil, "User ID not found in Gin context for UpdateClientProfile")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID missing from context"})
		return
	}
	userID, ok := userIDAny.(uint)
	if !ok {
		utils.ErrorLog(map[string]interface{}{"userID_type": fmt.Sprintf("%T", userIDAny)}, "Invalid User ID type in Gin context for UpdateClientProfile")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error: Invalid user ID format"})
		return
	}

	var dto userclient.ClientUpdateRequestDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.LoggerRequest(map[string]interface{}{"error": err.Error()}, "Update Client Profile", "Invalid request payload")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request payload: %v", err.Error())})
		return
	}

	updatedUserDTO, err := ctrl.clientAuthService.UpdateClientProfile(userID, &dto)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"error": err.Error()}, fmt.Sprintf("Failed to update profile for client ID %d", userID))
		if err.Error() == fmt.Sprintf("client user not found with ID: %d", userID) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Client user not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, updatedUserDTO) // Use c.JSON
	utils.LoggerRequest(map[string]interface{}{"user_id": userID}, "Update Client Profile", fmt.Sprintf("Updated profile for client ID %d", userID))
}

func (ctrl *ClientAuthController) ClientLogout(c *gin.Context, cfg *config.GitConfig) {
	c.SetCookie(
		"kdc.secure.token",
		"",
		-1,
		"/",                  // MUST match the path of the original cookie
		"",                   // MUST match the domain of the original cookie
		c.Request.TLS != nil, // MUST match the secure attribute
		false,                // MUST match the httpOnly attribute
	)
	frontendRedirectURL := cfg.ClientEnd + "home"
	c.Redirect(http.StatusFound, frontendRedirectURL)
	utils.InfoLog(map[string]interface{}{"logout": ""}, "success")
}
