package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github" // Specific GitHub endpoint
)

// Define your GitHub OAuth configurations
var (
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost/api/account/auth/github/callback", // Must match the one registered on GitHub
		ClientID:     "Ov23liBVXaZ0bV6B43Ut",
		ClientSecret: "e57c914b57ffdfda1f8af96357bfdceeddb8370d",
		Scopes:       []string{"user:email"}, // Request email access
		Endpoint:     githubOAuth.Endpoint,
	}
	// A random string to prevent CSRF attacks. Store this in a session or cookie.
	// For simplicity, using a hardcoded string here, but in production, generate dynamically.
	oauthStateString = "random-state-string"
)

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/auth/github/login", handleGitHubLogin)
	http.HandleFunc("/auth/github/callback", handleGitHubCallback)

	fmt.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Go backend! Try to login with GitHub via your Vue app.")
}

// handleGitHubLogin redirects to GitHub's authorization page
func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleGitHubCallback receives the callback from GitHub
func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	// 1. Verify the state parameter to prevent CSRF
	state := r.FormValue("state")
	if state != oauthStateString {
		http.Error(w, "Invalid state parameter", http.StatusInternalServerError)
		return
	}

	// 2. Exchange the authorization code for an access token
	code := r.FormValue("code")
	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Could not exchange code for token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Use the access token to fetch user data from GitHub API
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		http.Error(w, "Error creating GitHub API request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "token "+token.AccessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error fetching user data from GitHub: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	userData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading GitHub user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse user data (example: just print it for now)
	var githubUser map[string]interface{}
	if err := json.Unmarshal(userData, &githubUser); err != nil {
		http.Error(w, "Error parsing GitHub user data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("GitHub User Data: %+v\n", githubUser)

	// 4. (Crucial) Implement your application's user management logic:
	//    a. Check if a user with this GitHub ID (githubUser["id"]) already exists in your database.
	//    b. If yes, log them in (e.g., generate a JWT token for your frontend).
	//    c. If no, create a new user account in your database using the data from GitHub (e.g., githubUser["login"] for username, githubUser["email"] if available and requested in scope).
	//       Then log them in.

	// Example: Redirect back to the frontend with some user info or a token
	// In a real app, you'd generate a JWT here and redirect to your frontend with it.
	// For now, let's just show success.
	// Note: You need to configure CORS if your frontend is on a different origin.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "GitHub login successful!", "user_id": fmt.Sprintf("%v", githubUser["id"])})

	// To redirect to Vue frontend after successful login:
	// frontendRedirectURL := "http://localhost:8080/dashboard" // Or where your Vue app handles post-login
	// http.Redirect(w, r, frontendRedirectURL, http.StatusTemporaryRedirect)
}
