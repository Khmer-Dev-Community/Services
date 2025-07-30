package auth02

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	redis "github.com/Khmer-Dev-Community/Services/api-service/config"
	clientauth_service "github.com/Khmer-Dev-Community/Services/api-service/lib/clientauth" // Aliased service import
	"github.com/Khmer-Dev-Community/Services/api-service/lib/userclient"                    // Renamed package
	"github.com/Khmer-Dev-Community/Services/api-service/pkg/oauth_config"                  // NEW IMPORT for shared OAuth types
	"github.com/Khmer-Dev-Community/Services/api-service/utils"
)

// ClientAuthController handles HTTP requests related to client authentication.
type ClientAuthController struct {
	clientAuthService *clientauth_service.ClientAuthService // Refers to the service in `lib/clientauth`
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
// @Failure 400 {string} string "Invalid request payload"
// @Failure 409 {string} string "Username or email already taken"
// @Failure 500 {string} string "Internal server error"
// @Router /client/register [post]
func (ctrl *ClientAuthController) RegisterClient(w http.ResponseWriter, r *http.Request) {
	var payload ClientRegisterPayload // From auth02 package
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.LoggerRequest(err.Error(), "Client Registration", "Invalid request payload")
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	newUser, err := ctrl.clientAuthService.RegisterClient((*userclient.ClientRegisterRequestDTO)(&payload)) // Use userclient DTO
	if err != nil {
		utils.ErrorLog(err, fmt.Sprintf("Client registration failed for %s", payload.Username))
		if err.Error() == "username already taken" || err.Error() == "email already registered" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	token, err := ctrl.clientAuthService.GenerateToken(newUser)
	if err != nil {
		utils.ErrorLog(err, fmt.Sprintf("Failed to generate token for client user %s", newUser.Username))
		http.Error(w, "Failed to generate authentication token", http.StatusInternalServerError)
		return
	}

	userResponse := userclient.ToClientUserResponseDTO(newUser) // Use userclient DTO
	response := ClientAuthResponse{                             // From auth02 package
		Token: token,
		User:  userResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.ErrorLog(err, "Failed to encode response during client registration")
	}
	utils.LoggerRequest(userResponse.ID, "Client Registration", fmt.Sprintf("Client user registered and logged in: %s", userResponse.Username))
}

// ClientLogin handles the client login request.
// @Summary Log in a client user
// @Description Authenticates a client user with username/email and password.
// @Tags Client Authentication
// @Accept json
// @Produce json
// @Param client body ClientLoginCredentials true "Client login credentials"
// @Success 200 {object} ClientAuthResponse "User successfully logged in"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 403 {string} string "Account is inactive or blocked / This account uses social login"
// @Failure 500 {string} string "Internal server error"
// @Router /client/login [post]
func (ctrl *ClientAuthController) ClientLogin(w http.ResponseWriter, r *http.Request) {
	var credentials ClientLoginCredentials // From auth02 package
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		utils.LoggerRequest(err.Error(), "Client Login", "Invalid request payload")
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	user, err := ctrl.clientAuthService.ClientLogin((*userclient.ClientLoginRequestDTO)(&credentials)) // Use userclient DTO
	if err != nil {
		utils.ErrorLog(err, fmt.Sprintf("Client login failed for %s", credentials.Identifier))
		if err.Error() == "invalid credentials" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		if err.Error() == "account is inactive or blocked" {
			http.Error(w, "Account is inactive or blocked", http.StatusForbidden)
			return
		}
		if err.Error() == "this account uses social login, please use the GitHub login option" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to log in", http.StatusInternalServerError)
		return
	}

	token, err := ctrl.clientAuthService.GenerateToken(user)
	if err != nil {
		utils.ErrorLog(err, fmt.Sprintf("Failed to generate token for client user %s", user.Username))
		http.Error(w, "Failed to generate authentication token", http.StatusInternalServerError)
		return
	}

	userResponse := userclient.ToClientUserResponseDTO(user) // Use userclient DTO
	response := ClientAuthResponse{                          // From auth02 package
		Token: token,
		User:  userResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		utils.ErrorLog(err, "Failed to encode response during client login")
	}
	utils.LoggerRequest(userResponse.ID, "Client Login", fmt.Sprintf("Client user logged in: %s", userResponse.Username))
}

// GitHubLoginRedirect handles the initiation of GitHub OAuth login.
// @Summary Redirect to GitHub for OAuth login
// @Description Redirects the user to GitHub's authorization page.
// @Tags Client Authentication
// @Produce html
// @Success 307 {string} string "Redirect to GitHub"
// @Failure 500 {string} string "GitHub OAuth not configured"
// @Router /client/github/login [get]
func (ctrl *ClientAuthController) GitHubLoginRedirect(w http.ResponseWriter, r *http.Request) {
	if oauth_config.ClientGithubOauthConfig == nil { // Access from oauth_config package
		utils.ErrorLog(nil, "GitHub OAuth configuration not initialized.")
		http.Error(w, "GitHub OAuth not configured", http.StatusInternalServerError)
		return
	}
	url := oauth_config.ClientGithubOauthConfig.AuthCodeURL(oauth_config.ClientOauthStateString) // Access from oauth_config package
	utils.LoggerRequest(url, "GitHub OAuth Redirect", "Redirecting client to GitHub for authorization")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
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
// @Failure 400 {string} string "Missing authorization code or state / Invalid state parameter"
// @Failure 500 {string} string "Internal server error"
// @Router /client/github/callback [get]
func (ctrl *ClientAuthController) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")
	state := query.Get("state")

	if code == "" || state == "" {
		utils.LoggerRequest(nil, "GitHub Callback", "Missing code or state in query parameters")
		http.Error(w, "Missing authorization code or state", http.StatusBadRequest)
		return
	}

	if state != oauth_config.ClientOauthStateString {
		utils.WarnLog(state, "GitHub callback: Invalid OAuth state parameter")
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}
	token, err := ctrl.clientAuthService.ExchangeGitHubCodeForToken(code)
	if err != nil {
		utils.ErrorLog(err, "GitHub callback: Failed to exchange code for token")
		http.Error(w, "Failed to authenticate with GitHub", http.StatusInternalServerError)
		return
	}
	githubUser, err := ctrl.clientAuthService.GetGitHubUserData(token.AccessToken) // Service returns oauth_config.GitHubUser
	if err != nil {
		utils.ErrorLog(err, "GitHub callback: Failed to get GitHub user data")
		http.Error(w, "Failed to get user data from GitHub", http.StatusInternalServerError)
		return
	}
	user, err := ctrl.clientAuthService.HandleGitHubLoginOrRegister(githubUser) // Service expects oauth_config.GitHubUser
	if err != nil {
		utils.ErrorLog(err, fmt.Sprintf("GitHub callback: Failed to handle login/register for GitHub user %s", githubUser.Login))
		http.Error(w, "Failed to process GitHub login/registration", http.StatusInternalServerError)
		return
	}
	appToken, err := ctrl.clientAuthService.GenerateToken(user)
	if err != nil {
		utils.ErrorLog(err, fmt.Sprintf("GitHub callback: Failed to generate app token for user %s", user.Username))
		http.Error(w, "Failed to generate application token", http.StatusInternalServerError)
		return
	}
	userResponse := userclient.ToClientUserResponseDTO(user) // Use userclient DTO
	response := ClientAuthResponse{                          // From auth02 package
		Token: appToken,
		User:  userResponse,
	}
	expiration := time.Hour * 24 // 1 nimutes
	expirationTime := time.Now().Add(time.Hour * 24 * 7)
	key := fmt.Sprintf("user:%d", user.ID)
	user.Token = appToken
	if err := redis.SetWithExpiration(key, response, expiration); err != nil {
		utils.ErrorLog(err, "Failed Store User Infor in Redis")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to store user information in Redis")
		return
	}
	cookie := http.Cookie{
		Name:     "kdc.secure.token", // Name of your authentication cookie
		Value:    appToken,           // The JWT you generated
		Expires:  expirationTime,     // When the cookie expires
		HttpOnly: true,               // Prevents client-side JavaScript access <-- CRITICAL FOR SECURITY
		Secure:   r.TLS != nil,       // Only send over HTTPS in production <-- CRITICAL FOR SECURITY
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	frontendRedirectURL := "http://localhost/auth/callback"
	http.Redirect(w, r, frontendRedirectURL, http.StatusFound)

	utils.LoggerRequest(userResponse.ID, "GitHub Callback", fmt.Sprintf("Client user logged in/registered via GitHub: %s", userResponse.Username))
}

// GetClientProfile handles fetching a client user's profile.
// @Summary Get client user profile
// @Description Retrieves the profile of the authenticated client user.
// @Tags Client Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} userclient.ClientUserResponseDTO "Client user profile"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Client user not found"
// @Failure 500 {string} string "Internal server error"
// @Router /client/profile [get]
func (ctrl *ClientAuthController) GetClientProfile(w http.ResponseWriter, r *http.Request) {
	userID := uint(1)
	userDTO, err := ctrl.clientAuthService.GetClientProfile(userID) // Returns userclient.ClientUserResponseDTO
	if err != nil {
		utils.ErrorLog(err, fmt.Sprintf("Failed to retrieve profile for client ID %d", userID))
		if err.Error() == fmt.Sprintf("client user not found with ID: %d", userID) {
			http.Error(w, "Client user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve profile", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userDTO); err != nil {
		utils.ErrorLog(err, "Failed to encode response during GetClientProfile")
	}
	utils.LoggerRequest(userID, "Get Client Profile", fmt.Sprintf("Retrieved profile for client ID %d", userID))
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
// @Failure 400 {string} string "Invalid request payload"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Client user not found"
// @Failure 500 {string} string "Internal server error"
// @Router /client/profile [put]
func (ctrl *ClientAuthController) UpdateClientProfile(w http.ResponseWriter, r *http.Request) {
	// !!! IMPORTANT: This is a placeholder. You need to implement actual authentication
	// and extract the userID from the authenticated user's context/token.
	// For now, it's hardcoded to 1 for testing.
	userID := uint(1)

	var dto userclient.ClientUpdateRequestDTO // Use userclient DTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.LoggerRequest(err.Error(), "Update Client Profile", "Invalid request payload")
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	updatedUserDTO, err := ctrl.clientAuthService.UpdateClientProfile(userID, &dto) // Returns userclient.ClientUserResponseDTO
	if err != nil {
		utils.ErrorLog(err, fmt.Sprintf("Failed to update profile for client ID %d", userID))
		if err.Error() == fmt.Sprintf("client user not found with ID: %d", userID) {
			http.Error(w, "Client user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedUserDTO); err != nil {
		utils.ErrorLog(err, "Failed to encode response during UpdateClientProfile")
	}
	utils.LoggerRequest(userID, "Update Client Profile", fmt.Sprintf("Updated profile for client ID %d", userID))
}
