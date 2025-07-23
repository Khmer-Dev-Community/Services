package userclient

import (
	"fmt"

	"github.com/Khmer-Dev-Community/Services/api-service/utils"
	"gorm.io/gorm"
)

// ClientUserRepository handles database operations for ClientUser.
type ClientUserRepository struct {
	db *gorm.DB
}

// NewClientUserRepository creates a new instance of ClientUserRepository.
func NewClientUserRepository(db *gorm.DB) *ClientUserRepository {
	return &ClientUserRepository{db: db}
}

// GetClientUserByUsername retrieves a ClientUser by their username.
func (r *ClientUserRepository) GetClientUserByUsername(username string) (*ClientUser, error) {
	user := &ClientUser{}
	result := r.db.Where("username = ?", username).First(user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.InfoLog(username, fmt.Sprintf("Client user not found with username: %s", username))
			return nil, fmt.Errorf("client user not found with username: %s", username)
		}
		utils.ErrorLog(result.Error, fmt.Sprintf("Error getting client user by username %s: %v", username, result.Error))
		return nil, fmt.Errorf("database error: %w", result.Error)
	}
	utils.LoggerRepository(user.ID, fmt.Sprintf("Client user retrieved by username: %s", username))
	return user, nil
}

// GetClientUserByGitHubID retrieves a ClientUser by their GitHub ID.
func (r *ClientUserRepository) GetClientUserByGitHubID(githubID string) (*ClientUser, error) {
	user := &ClientUser{}
	result := r.db.Where("github_id = ?", githubID).First(user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.InfoLog(githubID, fmt.Sprintf("Client user not found with GitHub ID: %s", githubID))
			return nil, fmt.Errorf("client user not found with GitHub ID: %s", githubID)
		}
		utils.ErrorLog(result.Error, fmt.Sprintf("Error getting client user by GitHub ID %s: %v", githubID, result.Error))
		return nil, fmt.Errorf("database error: %w", result.Error)
	}
	utils.LoggerRepository(user.ID, fmt.Sprintf("Client user retrieved by GitHub ID: %s", githubID))
	return user, nil
}

// GetClientUserByEmail retrieves a ClientUser by their email.
func (r *ClientUserRepository) GetClientUserByEmail(email string) (*ClientUser, error) {
	user := &ClientUser{}
	result := r.db.Where("email = ?", email).First(user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.InfoLog(email, fmt.Sprintf("Client user not found with email: %s", email))
			return nil, fmt.Errorf("client user not found with email: %s", email)
		}
		utils.ErrorLog(result.Error, fmt.Sprintf("Error getting client user by email %s: %v", email, result.Error))
		return nil, fmt.Errorf("database error: %w", result.Error)
	}
	utils.LoggerRepository(user.ID, fmt.Sprintf("Client user retrieved by email: %s", email))
	return user, nil
}

// CreateClientUser creates a new ClientUser in the database.
func (r *ClientUserRepository) CreateClientUser(user *ClientUser) error {
	if user.Status == false {
		user.Status = true
	}
	result := r.db.Create(user)
	if result.Error != nil {
		utils.ErrorLog(result.Error, fmt.Sprintf("Failed to create client user %s: %v", user.Username, result.Error))
		return fmt.Errorf("failed to create client user: %w", result.Error)
	}
	utils.LoggerRepository(user.ID, fmt.Sprintf("New client user created: %s (ID: %d)", user.Username, user.ID))
	return nil
}

// UpdateClientUser updates an existing ClientUser in the database.
func (r *ClientUserRepository) UpdateClientUser(user *ClientUser) error {
	result := r.db.Save(user) // Save updates all fields. For partial updates, use db.Model(&user).Updates(map[string]interface{}) or selectively update fields.
	if result.Error != nil {
		utils.ErrorLog(result.Error, fmt.Sprintf("Failed to update client user %s (ID: %d): %v", user.Username, user.ID, result.Error))
		return fmt.Errorf("failed to update client user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		utils.WarnLog(user.ID, fmt.Sprintf("No client user found to update with ID: %d", user.ID))
	} else {
		utils.LoggerRepository(user.ID, fmt.Sprintf("Client user updated: %s (ID: %d)", user.Username, user.ID))
	}
	return nil
}

// GetClientUserByID retrieves a ClientUser by their ID. (Often useful for profile fetching)
func (r *ClientUserRepository) GetClientUserByID(id uint) (*ClientUser, error) {
	user := &ClientUser{}
	result := r.db.First(user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.InfoLog(id, fmt.Sprintf("Client user not found with ID: %d", id))
			return nil, fmt.Errorf("client user not found with ID: %d", id)
		}
		utils.ErrorLog(result.Error, fmt.Sprintf("Error getting client user by ID %d: %v", id, result.Error))
		return nil, fmt.Errorf("database error: %w", result.Error)
	}
	utils.LoggerRepository(user.ID, fmt.Sprintf("Client user retrieved by ID: %d", id))
	return user, nil
}

// DeleteClientUser marks a ClientUser as deleted (soft delete).
func (r *ClientUserRepository) DeleteClientUser(id uint) error {
	user := &ClientUser{}
	result := r.db.Delete(user, id)
	if result.Error != nil {
		utils.ErrorLog(result.Error, fmt.Sprintf("Failed to delete client user with ID: %d: %v", id, result.Error))
		return fmt.Errorf("failed to delete client user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		utils.WarnLog(id, fmt.Sprintf("Client user with ID %d not found for deletion", id))
		return fmt.Errorf("client user with ID %d not found for deletion", id)
	}
	utils.LoggerRepository(id, fmt.Sprintf("Client user soft-deleted: ID %d", id))
	return nil
}
