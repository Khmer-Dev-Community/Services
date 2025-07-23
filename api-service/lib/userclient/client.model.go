package userclient

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ClientUser represents the model for a client user
type ClientUser struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	FirstName   string  `json:"fname" gorm:"column:first_name"`
	LastName    string  `json:"lname" gorm:"column:last_name"`
	Username    string  `json:"username" gorm:"column:username;uniqueIndex"`
	Password    *string `json:"-" gorm:"column:password"`
	Email       string  `json:"email" gorm:"column:email;uniqueIndex"`
	Phone       string  `json:"phone" gorm:"column:phone"`
	Sex         string  `json:"sex" gorm:"column:sex"`
	Status      bool    `json:"status" gorm:"column:status;default:true"`
	Token       string  `json:"token" gorm:"-"`
	Description string  `json:"description" gorm:"column:description"`

	// GitHub
	GitHubID  *string `json:"github_id,omitempty" gorm:"column:github_id;uniqueIndex"`
	AvatarURL *string `json:"avatar_url,omitempty" gorm:"column:avatar_url"`
	Name      string  `json:"name,omitempty" gorm:"column:name"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ClientUser) TableName() string {
	return "client_platform_user" // Distinct table name for client users
}

func MigrateClientUsers(db *gorm.DB) {
	db.AutoMigrate(&ClientUser{})
}

func (u *ClientUser) MarshalJSON() ([]byte, error) {
	type Alias ClientUser

	alias := &struct {
		*Alias
		Password interface{} `json:"password,omitempty"`
		Token    interface{} `json:"token,omitempty"`
		Role     interface{} `json:"Role,omitempty"`
	}{
		Alias: (*Alias)(u),
		Token: u.Token,
	}

	if u.Password != nil && *u.Password != "" {
		alias.Password = "********"
	} else {
		alias.Password = nil
	}

	return json.Marshal(alias)
}
