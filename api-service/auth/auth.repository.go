// user_repository.go
package auth

import (
	"errors"

	user "github.com/Khmer-Dev-Community/Services/api-service/lib/users/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (ur *UserRepository) GetUserByUsername(username string) (*user.User, error) {
	var user user.User
	if err := ur.DB.Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepository) Login(username string, password string) (*user.User, error) {
	var user user.User
	if err := r.DB.Preload("system_role").Where("username = ? AND password = ?", username, password).First(&user).Limit(1).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}
	return &user, nil
}
