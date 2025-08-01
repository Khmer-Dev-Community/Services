package posts

import (
	"strings"
)

// ToPostResponse converts a Post model to a PostResponse DTO.
func ToPostResponse(post Post) PostResponse {
	var tagsArray []string
	if post.Tags != "" {
		tagsArray = splitTagsString(post.Tags)
	}
	return PostResponse{
		ID:               post.ID,
		Title:            post.Title,
		Slug:             post.Slug,
		Description:      post.Description,
		ViewCount:        post.ViewCount,
		FeaturedImageURL: post.FeaturedImageURL,
		AuthorID:         post.AuthorID,
		Tags:             tagsArray,
		Status:           post.Status,
		CreatedAt:        post.CreatedAt,
		UpdatedAt:        post.UpdatedAt,
	}
}

// ToPostResponseList converts a slice of Post models to a slice of PostResponse DTOs.
func ToPostResponseList(posts []Post) []PostResponse {
	responses := make([]PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = ToPostResponse(post)
	}
	return responses
}

// ToPostModel converts a CreatePostRequest DTO to a Post model.
// The slug is now expected to be provided by the service layer.
func ToPostModel(req CreatePostRequest, authorID uint, generatedSlug string) Post {
	return Post{
		Title:            req.Title,
		Slug:             generatedSlug, // Slug is provided by the service
		Description:      req.Description,
		FeaturedImageURL: req.FeaturedImageURL,
		AuthorID:         authorID,
		Tags:             joinTagsArray(req.Tags),
		Status:           getDefaultStatus(req.Status),
	}
}

// UpdatePostModel updates an existing Post model from an UpdatePostRequest DTO.
// The updated slug is now expected to be provided by the service layer if the title changes.
func UpdatePostModel(post *Post, req UpdatePostRequest, updatedSlug *string) {
	if req.Title != nil {
		post.Title = *req.Title
		// If title changed and an updatedSlug is provided by the service, use it
		if updatedSlug != nil {
			post.Slug = *updatedSlug
		}
	}
	if req.Description != nil {
		post.Description = *req.Description
	}
	if req.FeaturedImageURL != nil {
		post.FeaturedImageURL = *req.FeaturedImageURL
	}
	if req.Tags != nil {
		post.Tags = joinTagsArray(req.Tags)
	}
	if req.Status != nil {
		post.Status = *req.Status
	}
}

// splitTagsString converts a comma-separated string to a slice of strings.
func splitTagsString(tags string) []string {
	if tags == "" {
		return []string{}
	}
	parts := []string{}
	for _, part := range strings.Split(tags, ",") {
		parts = append(parts, strings.TrimSpace(part))
	}
	return parts
}

// joinTagsArray converts a slice of strings to a comma-separated string.
func joinTagsArray(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	trimmedTags := make([]string, len(tags))
	for i, tag := range tags {
		trimmedTags[i] = strings.TrimSpace(tag)
	}
	return strings.Join(trimmedTags, ",")
}

// getDefaultStatus returns "draft"
func getDefaultStatus(status string) string {
	if status == "" {
		return "draft"
	}
	return status
}
