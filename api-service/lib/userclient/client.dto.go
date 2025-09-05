package userclient

import (
	"time"
)

type ClientUserResponseInfor struct {
	ID        uint   `json:"id"`
	AvatarURL string `json:"avatar_url"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Likes     uint   `json:"likes"`
	Follower  uint   `json:"follower"`
	Following uint   `json:"following"`
}

// ClientUserResponseDTO represents the data response
type ClientUserResponseDTO struct {
	ID          uint    `json:"id"`
	FirstName   string  `json:"fname"`
	LastName    string  `json:"lname"`
	Username    string  `json:"username"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Sex         string  `json:"sex"`
	Status      bool    `json:"status"`
	Description string  `json:"description"`
	GitHubID    *string `json:"github_id,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Name        string  `json:"name,omitempty"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// DeletedAt is typically omitted from public responses
}

// ClientRegisterRequestDTO represents the data expected when a new client registers.
type ClientRegisterRequestDTO struct {
	FirstName   string `json:"fname" validate:"required"`
	LastName    string `json:"lname" validate:"required"`
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Password    string `json:"password" validate:"required,min=8"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"omitempty,e164"`
	Sex         string `json:"sex" validate:"omitempty,oneof=Male Female Other"`
	Description string `json:"description"`
}

// ClientLoginRequestDTO represents the data expected for a client login (username/email and password).
type ClientLoginRequestDTO struct {
	Identifier string `json:"identifier" validate:"required"` // Can be username or email
	Password   string `json:"password" validate:"required"`
}

// ClientUpdateRequestDTO represents data for updating a client's profile.
type ClientUpdateRequestDTO struct {
	FirstName   *string `json:"fname,omitempty"`
	LastName    *string `json:"lname,omitempty"`
	Username    *string `json:"username,omitempty"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Sex         *string `json:"sex,omitempty"`
	Status      *bool   `json:"status,omitempty"`
	Description *string `json:"description,omitempty"`
}
