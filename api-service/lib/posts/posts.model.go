package posts

import (
	"time"

	"github.com/Khmer-Dev-Community/Services/api-service/lib/userclient"
	"gorm.io/gorm"
)

type ReactionType string

const (
	Like     ReactionType = "like"
	Love     ReactionType = "love"
	Laughing ReactionType = "laughing"
	// Add other reaction types as needed
)

type Post struct {
	ID               uint                   `gorm:"primaryKey"`
	Title            string                 `json:"title" gorm:"type:varchar(255);not null"`
	Slug             string                 `json:"slug" gorm:"type:varchar(255);uniqueIndex"`
	Meta             string                 `json:"meta" gorm:"type:varchar(300)"`
	Link             string                 `json:"link" gorm:"type:varchar(80)"`
	Description      string                 `json:"description" gorm:"type:text;not null"`
	ViewCount        uint                   `json:"view_count" gorm:"default:0"`
	FeaturedImageURL string                 `json:"featured_image_url" gorm:"type:varchar(255)"`
	AuthorID         uint                   `json:"author_id" gorm:"not null"`
	Author           *userclient.ClientUser `json:"author"`
	Tags             string                 `json:"tags" gorm:"type:text"`
	Status           string                 `json:"status" gorm:"type:varchar(50);default:'draft'"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`

	Comments      []Comment  `json:"discussion" gorm:"foreignKey:PostID"`
	Reactions     []Reaction `json:"reactions" gorm:"foreignKey:PostID"`
	ReactionCount int64      `json:"reaction_count" gorm:"-"`
}

type Comment struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	PostID          uint      `gorm:"not null" json:"post_id"`
	AuthorID        uint      `gorm:"not null" json:"author_id"`
	ParentCommentID *uint     `json:"parent_comment_id"` // Use a pointer for nullable foreign key
	Content         string    `gorm:"type:text;not null" json:"content"`
	Upvotes         int       `gorm:"default:0" json:"upvotes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt         `gorm:"index"`
	Author          *userclient.ClientUser `json:"author"`
	Replies         []Comment              `gorm:"foreignKey:ParentCommentID" json:"replies"`
}

type Reaction struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	PostID       uint         `gorm:"not null"`
	UserID       uint         `gorm:"not null"`
	ReactionType ReactionType `gorm:"type:varchar(20);not null"`
	CreatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (p *Post) TableName() string {
	return "client_content_post"
}
func (p *Comment) TableName() string {
	return "client_comment_post"
}
func (p *Reaction) TableName() string {
	return "client_reaction_post"
}

func MigrateClientPost(db *gorm.DB) {
	db.AutoMigrate(&Post{}, &Comment{}, &Reaction{})
}
