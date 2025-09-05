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
