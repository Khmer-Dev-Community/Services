package dto

import (
	"errors"

	models "github.com/Khmer-Dev-Community/Services/api-service/lib/users/models"
)

// UserUpdateDTO is a struct for updating a user
type UserDTO struct {
	ID         int    `json:"id"`
	FirstName  string `json:"fname"`
	LastName   string `json:"lname"`
	Username   string `json:"username"`
	RoleID     int    `json:"roleId"`
	CompanyID  int    `json:"companyId"`
	DepamentId int    `json:"depamentId"`
	Position   string `json:"position"`
	Emial      string `json:"email"`
	Sex        string `json:"sex"`
	Phone      string `json:"phone"`
	Status     bool   `json:"status"`
	Token      string `json:"token"`
	Decription string `json:"decription"`
}

// UserCreateDTO is a struct for creating a user
type UserCreateDTO struct {
	FirstName  string `json:"fname" binding:"required"`
	LastName   string `json:"lname" binding:"required"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RoleID     int    `json:"roleId" binding:"required"`
	CompanyID  int    `json:"companyId"`
	DepamentId int    `json:"depamentId"`
	Position   string `json:"position"`
	Emial      string `json:"email"`
	Sex        string `json:"sex"`
	Phone      string `json:"phone"`
	Status     bool   `json:"status"`
	Token      string `json:"token"`
	Decription string `json:"decription"`
}

// UserUpdateDTO is a struct for updating a user
type UserUpdateDTO struct {
	ID         int    `json:"id"`
	FirstName  string `json:"fname"`
	LastName   string `json:"lname"`
	Username   string `json:"username"`
	RoleID     int    `json:"roleId"`
	CompanyID  int    `json:"companyId"`
	DepamentId int    `json:"depamentId"`
	Position   string `json:"position"`
	Emial      string `json:"email"`
	Sex        string `json:"sex"`
	Phone      string `json:"phone"`
	Status     bool   `json:"status"`
	Token      string `json:"token"`
	Decription string `json:"decription"`
}

func (dto *UserCreateDTO) Validate() error {
	if dto.FirstName == "" {
		return errors.New("first name is required")
	}
	if dto.LastName == "" {
		return errors.New("last name is required")
	}
	if dto.Username == "" {
		return errors.New("username is required")
	}
	if dto.Password == "" {
		return errors.New("password is required")
	}
	if dto.RoleID == 0 {
		return errors.New("role ID is required")
	}
	if dto.CompanyID < 0 {
		return errors.New("company ID must be a non-negative integer")
	}
	if dto.Sex != "" && dto.Sex != "male" && dto.Sex != "female" && dto.Sex != "other" {
		return errors.New("invalid sex value")
	}
	if dto.Phone != "" && len(dto.Phone) < 10 {
		return errors.New("phone number must be at least 10 digits")
	}
	return nil
}

// FromModel maps a User model to UserDTO
func (dto *UserDTO) FromModel(user *models.User) {
	dto.ID = int(user.ID)
	dto.FirstName = user.FirstName
	dto.LastName = user.LastName
	dto.Username = user.Username
	dto.RoleID = int(user.RoleID)
	dto.CompanyID = user.CompanyID
	dto.DepamentId = user.DepartmentId
	dto.Emial = user.Email
	dto.Sex = user.Sex
	dto.Phone = user.Phone
	dto.Status = user.Status
	dto.Token = user.Token
	dto.Decription = user.Decription
	dto.Position = user.Position
}

// ToModel maps UserDTO to a User model
func (dto *UserDTO) ToModel() *models.User {
	return &models.User{
		ID:           uint(dto.ID),
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Username:     dto.Username,
		RoleID:       uint(dto.RoleID),
		CompanyID:    dto.CompanyID,
		DepartmentId: dto.DepamentId,
		Email:        dto.Emial,
		Sex:          dto.Sex,
		Phone:        dto.Phone,
		Status:       dto.Status,
		Token:        dto.Token,
		Decription:   dto.Decription,
		Position:     dto.Position,
	}
}
