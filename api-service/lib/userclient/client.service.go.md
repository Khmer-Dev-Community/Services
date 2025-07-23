package userclient

import (
"context"
"encoding/json"
"errors"
"fmt"
"io/ioutil"
"net/http"
"time"

    "github.com/sirupsen/logrus"
    "golang.org/x/oauth2"

    "github.com/Khmer-Dev-Community/Services/api-service/auth02"
    "github.com/Khmer-Dev-Community/Services/api-service/utils" // For password hashing, and for logging

)

// ClientAuthService handles authentication-related business logic for client users.
type ClientAuthService struct {
clientUserRepo \*ClientUserRepository
}

// NewClientAuthService creates a new instance of ClientAuthService.
func NewClientAuthService(clientUserRepo *ClientUserRepository) *ClientAuthService {
return &ClientAuthService{clientUserRepo: clientUserRepo}
}

// RegisterClient handles new client user registration.
// It takes a ClientRegisterRequestDTO, processes it, and returns a ClientUser model.
func (s *ClientAuthService) RegisterClient(payload *ClientRegisterRequestDTO) (\*ClientUser, error) {
utils.LoggerService(payload.Username, fmt.Sprintf("Attempting to register new client: %s", payload.Username), logrus.InfoLevel)

    // Check for existing username or email to prevent duplicates
    _, errUsername := s.clientUserRepo.GetClientUserByUsername(payload.Username)
    if errUsername == nil {
    	utils.WarnLog(payload.Username, "Client registration failed: username already taken")
    	return nil, errors.New("username already taken")
    }

    _, errEmail := s.clientUserRepo.GetClientUserByEmail(payload.Email)
    if errEmail == nil {
    	utils.WarnLog(payload.Email, "Client registration failed: email already registered")
    	return nil, errors.New("email already registered")
    }

    // Hash the password
    hashedPassword, err := utils.EncryptPassword(payload.Password)
    if err != nil {
    	utils.ErrorLog(err, "Error hashing password during client registration")
    	return nil, errors.New("failed to process password")
    }
    hashedPasswordPtr := &hashedPassword

    // Map the DTO to the ClientUser model
    newUser := ToClientUserFromRegisterDTO(payload) // Direct call, no alias
    newUser.Password = hashedPasswordPtr            // Assign the hashed password
    newUser.Status = true                           // Ensure new users are active by default

    if err := s.clientUserRepo.CreateClientUser(newUser); err != nil {
    	utils.ErrorLog(err, "Failed to create new client user in database")
    	return nil, errors.New("failed to register user")
    }

    utils.InfoLog(newUser.ID, fmt.Sprintf("Client user registered successfully: %s (ID: %d)", newUser.Username, newUser.ID))
    return newUser, nil

}

// ClientLogin handles client user login with username/email and password.
func (s *ClientAuthService) ClientLogin(credentials *ClientLoginRequestDTO) (\*ClientUser, error) {
utils.LoggerService(credentials.Identifier, fmt.Sprintf("Attempting client login for identifier: %s", credentials.Identifier), logrus.InfoLevel)

    user, err := s.clientUserRepo.GetClientUserByUsername(credentials.Identifier)
    if err != nil {
    	user, err = s.clientUserRepo.GetClientUserByEmail(credentials.Identifier)
    	if err != nil {
    		utils.WarnLog(credentials.Identifier, "Client login failed: invalid credentials (user not found)")
    		return nil, errors.New("invalid credentials")
    	}
    }

    if !user.Status {
    	utils.WarnLog(user.ID, fmt.Sprintf("Client account inactive or blocked: %s", user.Username))
    	return nil, errors.New("account is inactive or blocked")
    }

    if user.Password == nil || *user.Password == "" {
    	utils.WarnLog(user.ID, fmt.Sprintf("Client %s tried password login, but account uses social login.", user.Username))
    	return nil, errors.New("this account uses social login, please use the GitHub login option")
    }

    if err := utils.ComparePassword(*user.Password, credentials.Password); err != nil {
    	utils.WarnLog(user.ID, fmt.Sprintf("Client login failed for %s: invalid password", user.Username))
    	return nil, errors.New("invalid credentials")
    }

    utils.InfoLog(user.ID, fmt.Sprintf("Client logged in successfully: %s", user.Username))
    return user, nil

}

// GenerateToken generates an application token (e.g., JWT) for a given ClientUser.
func (s *ClientAuthService) GenerateToken(user *ClientUser) (string, error) {
token := fmt.Sprintf("client*jwt_token_for_user*%d", user.ID) // Placeholder
utils.LoggerService(user.ID, fmt.Sprintf("Generated token for client user: %s", user.Username), logrus.InfoLevel)
return token, nil
}

// ExchangeGitHubCodeForToken exchanges the GitHub authorization code for an OAuth token.
func (s *ClientAuthService) ExchangeGitHubCodeForToken(code string) (*oauth2.Token, error) {
utils.LoggerService(code, "Exchanging GitHub code for token", logrus.DebugLevel)
token, err := auth02.ClientGithubOauthConfig.Exchange(context.Background(), code) // Access config from auth02
if err != nil {
utils.ErrorLog(err, fmt.Sprintf("Failed to exchange GitHub code for token: %v", err))
return nil, fmt.Errorf("could not exchange code for token: %w", err)
}
utils.LoggerService(token.AccessToken, "Successfully exchanged GitHub code for token", logrus.DebugLevel)
return token, nil
}

// GetGitHubUserData fetches user profile data from GitHub's API.
func (s *ClientAuthService) GetGitHubUserData(accessToken string) (*auth02.GitHubUser, error) { // Use auth02.GitHubUser
utils.LoggerService(nil, "Fetching GitHub user data", logrus.DebugLevel)
req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
if err != nil {
utils.ErrorLog(err, "Error creating GitHub API request")
return nil, fmt.Errorf("error creating GitHub API request: %w", err)
}
req.Header.Set("Authorization", "token "+accessToken)
req.Header.Set("Accept", "application/vnd.github.v3+json")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
    	utils.ErrorLog(err, "Error fetching user data from GitHub")
    	return nil, fmt.Errorf("error fetching user data from GitHub: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
    	body, _ := ioutil.ReadAll(resp.Body)
    	utils.ErrorLog(map[string]interface{}{"status": resp.StatusCode, "body": string(body)}, "GitHub API returned non-OK status")
    	return nil, fmt.Errorf("github API returned non-OK status: %d, body: %s", resp.StatusCode, body)
    }

    userData, err := ioutil.ReadAll(resp.Body)
    if err != nil {
    	utils.ErrorLog(err, "Error reading GitHub user data")
    	return nil, fmt.Errorf("error reading GitHub user data: %w", err)
    }

    var githubUser auth02.GitHubUser // Use auth02.GitHubUser
    if err := json.Unmarshal(userData, &githubUser); err != nil {
    	utils.ErrorLog(err, "Error parsing GitHub user data")
    	return nil, fmt.Errorf("error parsing GitHub user data: %w", err)
    }

    if githubUser.Email == "" {
    	utils.LoggerService(githubUser.Login, "GitHub email not directly available, attempting to fetch emails", logrus.DebugLevel)
    	emailsReq, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
    	if err == nil {
    		emailsReq.Header.Set("Authorization", "token "+accessToken)
    		emailsResp, err := client.Do(emailsReq)
    		if err == nil {
    			defer emailsResp.Body.Close()
    			if emailsResp.StatusCode == http.StatusOK {
    				emailsBody, _ := ioutil.ReadAll(emailsResp.Body)
    				var emails []struct {
    					Email    string `json:"email"`
    					Primary  bool   `json:"primary"`
    					Verified bool   `json:"verified"`
    				}
    				if json.Unmarshal(emailsBody, &emails) == nil {
    					for _, e := range emails {
    						if e.Primary && e.Verified {
    							githubUser.Email = e.Email
    							utils.LoggerService(githubUser.Email, fmt.Sprintf("Found primary verified GitHub email: %s", githubUser.Email), logrus.DebugLevel)
    							break
    						}
    					}
    				}
    			} else {
    				utils.WarnLog(map[string]interface{}{"status": emailsResp.StatusCode, "body": emailsResp.Body}, "GitHub emails API returned non-OK status")
    			}
    		} else {
    			utils.WarnLog(err, "Could not fetch GitHub user emails for client")
    		}
    	}
    }
    utils.LoggerService(githubUser.ID, fmt.Sprintf("Successfully fetched GitHub user data for: %s", githubUser.Login), logrus.DebugLevel)
    return &githubUser, nil

}

// HandleGitHubLoginOrRegister processes GitHub user data, either logging in an existing user
// or registering a new one if not found.
func (s *ClientAuthService) HandleGitHubLoginOrRegister(githubUser *auth02.GitHubUser) (\*ClientUser, error) { // Use auth02.GitHubUser
githubIDStr := fmt.Sprintf("%d", githubUser.ID)
utils.LoggerService(githubIDStr, fmt.Sprintf("Handling GitHub login/registration for ID: %s", githubIDStr), logrus.InfoLevel)

    user, err := s.clientUserRepo.GetClientUserByGitHubID(githubIDStr)
    if err == nil && user != nil {
    	utils.InfoLog(user.ID, fmt.Sprintf("Found existing client user with GitHub ID: %d (%s)", githubUser.ID, user.Username))
    	user.Name = githubUser.Name
    	user.AvatarURL = &githubUser.AvatarURL
    	if err := s.clientUserRepo.UpdateClientUser(user); err != nil {
    		utils.ErrorLog(err, fmt.Sprintf("Failed to update existing client user %s from GitHub", user.Username))
    		return nil, fmt.Errorf("failed to update existing client user from GitHub: %w", err)
    	}
    	return user, nil
    }

    if githubUser.Email != "" {
    	user, err := s.clientUserRepo.GetClientUserByEmail(githubUser.Email)
    	if err == nil && user != nil {
    		utils.InfoLog(user.ID, fmt.Sprintf("Found existing client user with email %s, linking to GitHub ID: %d", githubUser.Email, githubUser.ID))
    		user.GitHubID = &githubIDStr
    		user.Name = githubUser.Name
    		user.AvatarURL = &githubUser.AvatarURL
    		user.Password = nil
    		if err := s.clientUserRepo.UpdateClientUser(user); err != nil {
    			utils.ErrorLog(err, fmt.Sprintf("Failed to link existing client user %s with GitHub ID", user.Username))
    			return nil, fmt.Errorf("failed to link existing client user with GitHub ID: %w", err)
    		}
    		return user, nil
    	}
    }

    utils.InfoLog(githubUser.Login, fmt.Sprintf("Creating new client user from GitHub data: %s (%d)", githubUser.Login, githubUser.ID))
    newUser := &ClientUser{
    	Username:  githubUser.Login,
    	Email:     githubUser.Email,
    	FirstName: githubUser.Name,
    	LastName:  "",
    	Name:      githubUser.Name,
    	GitHubID:  &githubIDStr,
    	AvatarURL: &githubUser.AvatarURL,
    	Status:    true,
    	CreatedAt: time.Now(),
    	UpdatedAt: time.Now(),
    }
    if newUser.Email == "" {
    	dummyEmail := fmt.Sprintf("github-client-%d@example.com", githubUser.ID)
    	newUser.Email = dummyEmail
    	utils.WarnLog(githubUser.Login, fmt.Sprintf("GitHub email not provided for client user %s, using dummy email: %s", githubUser.Login, dummyEmail))
    }
    newUser.Password = nil

    if err := s.clientUserRepo.CreateClientUser(newUser); err != nil {
    	utils.ErrorLog(err, fmt.Sprintf("Failed to create new client user %s from GitHub", githubUser.Login))
    	return nil, fmt.Errorf("failed to create new client user from GitHub: %w", err)
    }

    utils.InfoLog(newUser.ID, fmt.Sprintf("New client user created from GitHub: %s (ID: %d)", newUser.Username, newUser.ID))
    return newUser, nil

}

// GetClientProfile retrieves a client user's profile by ID.
func (s *ClientAuthService) GetClientProfile(userID uint) (*ClientUserResponseDTO, error) {
utils.LoggerService(userID, fmt.Sprintf("Attempting to get client profile for ID: %d", userID), logrus.InfoLevel)
user, err := s.clientUserRepo.GetClientUserByID(userID)
if err != nil {
utils.WarnLog(userID, fmt.Sprintf("Failed to get client profile for ID %d: %v", userID, err))
return nil, err
}
utils.InfoLog(user.ID, fmt.Sprintf("Successfully retrieved client profile for ID: %d", userID))
return ToClientUserResponseDTO(user), nil
}

// UpdateClientProfile updates a client user's profile.
func (s *ClientAuthService) UpdateClientProfile(userID uint, dto *ClientUpdateRequestDTO) (\*ClientUserResponseDTO, error) {
utils.LoggerService(userID, fmt.Sprintf("Attempting to update client profile for ID: %d", userID), logrus.InfoLevel)
user, err := s.clientUserRepo.GetClientUserByID(userID)
if err != nil {
utils.WarnLog(userID, fmt.Sprintf("Failed to find client user with ID %d for update: %v", userID, err))
return nil, err
}

    UpdateClientUserFromDTO(user, dto)

    if err := s.clientUserRepo.UpdateClientUser(user); err != nil {
    	utils.ErrorLog(err, fmt.Sprintf("Failed to save updated client profile for ID %d", userID))
    	return nil, fmt.Errorf("failed to update client profile: %w", err)
    }

    utils.InfoLog(user.ID, fmt.Sprintf("Client profile updated successfully for ID: %d", userID))
    return ToClientUserResponseDTO(user), nil

}
