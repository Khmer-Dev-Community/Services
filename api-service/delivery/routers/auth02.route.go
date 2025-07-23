package router

import (
	"encoding/json" // For health check/root route JSON
	"log"
	"net/http"

	// Import your swagger docs if you have them generated for mux
	_ "github.com/Khmer-Dev-Community/Services/api-service/cmd/docs"

	"github.com/gorilla/mux"

	"github.com/Khmer-Dev-Community/Services/api-service/auth02"           // The auth02 package (controller)
	"github.com/Khmer-Dev-Community/Services/api-service/pkg/oauth_config" // The new oauth_config package
)

// ClientAuthControllerWrapper is a wrapper for the client auth controller
type ClientAuthControllerWrapper struct {
	clientAuthController *auth02.ClientAuthController
}

// NewClientAuthControllerWrapper initializes the wrapper for ClientAuthController
func NewClientAuthControllerWrapper(cac *auth02.ClientAuthController) *ClientAuthControllerWrapper {
	return &ClientAuthControllerWrapper{
		clientAuthController: cac,
	}
}

// SetupRouterAuth02 initializes and configures the router for the client authentication.
// It now directly accepts the initialized ClientAuthController.
func SetupRouterAuth02(r *mux.Router, clientAuthCtrl *auth02.ClientAuthController) { // <--- CHANGED SIGNATURE

	clientAuthControllerWrapper := NewClientAuthControllerWrapper(clientAuthCtrl)
	githubClientID := "Ov23liBVXaZ0bV6B43Ut"
	githubClientSecret := "28f9c091b494fc8bc1fd2b795c77be27b322f2c4"
	githubRedirectURL := "http://localhost:3000/api/account/auth02/github/callback"
	if githubClientID == "" || githubClientSecret == "" || githubRedirectURL == "" {
		log.Println("WARNING: GitHub OAuth environment variables (GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET, GITHUB_REDIRECT_URL) are not fully set. GitHub login might not work.")
	}

	oauth_config.InitializeClientGitHubOAuthConfig(githubClientID, githubClientSecret, githubRedirectURL)

	// Define API routes
	api := r.PathPrefix("/api").Subrouter()

	// --- Client Authentication Routes ---
	clientAuthRouter := api.PathPrefix("/auth02").Subrouter()

	// Public routes (no authentication required)
	clientAuthRouter.HandleFunc("/register", clientAuthControllerWrapper.clientAuthController.RegisterClient).Methods("POST")
	clientAuthRouter.HandleFunc("/login", clientAuthControllerWrapper.clientAuthController.ClientLogin).Methods("POST")

	// GitHub OAuth routes
	clientAuthRouter.HandleFunc("/github/login", clientAuthControllerWrapper.clientAuthController.GitHubLoginRedirect).Methods("GET")
	clientAuthRouter.HandleFunc("/github/callback", clientAuthControllerWrapper.clientAuthController.GitHubCallback).Methods("GET")

	// Profile routes (likely authenticated)
	// You will need to add middleware here for authentication if not global
	clientAuthRouter.HandleFunc("/profile", clientAuthControllerWrapper.clientAuthController.GetClientProfile).Methods("GET")
	clientAuthRouter.HandleFunc("/profile", clientAuthControllerWrapper.clientAuthController.UpdateClientProfile).Methods("PUT")

	// Optional: Add a health check or root route if this `SetupRouterAuth02` is the primary router setup
	// These might belong in a more general router setup function, but keeping for reference.
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Telegram Service API is running!"})
	})

	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Assuming 'db' is available if this function were passed it,
		// otherwise, health check might need a different approach or depend on other services.
		// For now, a simple HTTP OK.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "up"})
	}).Methods("GET")
}
