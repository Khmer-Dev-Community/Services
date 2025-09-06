package posts

import (
	"time"

	"github.com/Khmer-Dev-Community/Services/api-service/lib/userclient"
)

type PostResponse struct {
	ID               uint                               `json:"id"`
	Title            string                             `json:"title"`
	Slug             string                             `json:"slug"`
	Description      string                             `json:"description"`
	ViewCount        uint                               `json:"view_count"`
	FeaturedImageURL string                             `json:"featured_image_url"`
	AuthorID         uint                               `json:"author_id"`
	Author           userclient.ClientUserResponseInfor `json:"author"`
	Comments         []CommentResponse                  `json:"discussion"`
	Reaction         []ReactionResponse                 `json:"reaction"`
	Tags             []string                           `json:"tags"`
	Status           string                             `json:"status"`
	CreatedAt        time.Time                          `json:"created_at"`
	UpdatedAt        time.Time                          `json:"updated_at"`
}

type CreatePostRequest struct {
	Title            string   `json:"title" binding:"required"`
	Description      string   `json:"description" binding:"required"`
	FeaturedImageURL string   `json:"featured_image_url"`
	Tags             []string `json:"tags"`
	Status           string   `json:"status" binding:"omitempty,oneof=draft published archived"`
}

type UpdatePostRequest struct {
	Title            *string  `json:"title"`
	Slug             *string  `json:"slug"`
	Description      *string  `json:"description"`
	FeaturedImageURL *string  `json:"featured_image_url"`
	Tags             []string `json:"tags"`
	Status           *string  `json:"status" binding:"omitempty,oneof=draft published archived"`
}

type CommentResponse struct {
	ID              uint                               `json:"id"`
	PostID          uint                               `json:"post_id"`
	AuthorID        uint                               `json:"author_id"`
	ParentCommentID *uint                              `json:"parent_comment_id"`
	Content         string                             `json:"content"`
	Upvotes         int                                `json:"upvotes"`
	CreatedAt       time.Time                          `json:"created_at"`
	Author          userclient.ClientUserResponseInfor `json:"author"`
	Replies         []CommentResponse                  `json:"replies"` // Corrected: Nested replies are also a slice
}

type ReactionResponse struct {
	ID           uint         `json:"id"`
	PostID       uint         `json:"post_id"`
	UserID       uint         `json:"user_id"`
	ReactionType ReactionType `json:"reaction_type"`
	CreatedAt    time.Time    `json:"created_at"`
}
