package posts

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID               uint   `gorm:"primaryKey"`
	Title            string `json:"title" gorm:"type:varchar(255);not null"`
	Slug             string `json:"slug" gorm:"type:varchar(255);uniqueIndex"`
	Description      string `json:"description" gorm:"type:text;not null"`
	ViewCount        uint   `json:"view_count" gorm:"default:0"`
	FeaturedImageURL string `json:"featured_image_url" gorm:"type:varchar(255)"`
	AuthorID         uint   `json:"author_id" gorm:"not null"`
	Tags             string `json:"tags" gorm:"type:text"`
	Status           string `json:"status" gorm:"type:varchar(50);default:'draft'"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type Comment struct {
	gorm.Model
	PostID  uint   `gorm:"not null"`
	UserID  uint   `gorm:"not null"`
	Content string `gorm:"type:text;not null"`
}

type Reaction struct {
	gorm.Model
	PostID uint   `gorm:"index"`
	UserID uint   `gorm:"index"`
	Type   string `gorm:"type:varchar(20)"` // e.g., "like", "heart", "upvote", "downvote"
}

func (p *Post) TableName() string {
	return "client_content_post" // Distinct table name for client users
}

func MigrateClientPost(db *gorm.DB) {
	db.AutoMigrate(&Post{})
}
