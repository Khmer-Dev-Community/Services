package posts

import (
	"strings"

	"github.com/Khmer-Dev-Community/Services/api-service/lib/userclient"
)

// ToPostResponse converts a Post model to a PostResponse DTO.
func ToPostResponseNOComment(post Post) PostResponse {
	var tagsArray []string
	if post.Tags != "" {
		tagsArray = splitTagsString(post.Tags)
	}

	var authorResponse userclient.ClientUserResponseInfor
	// Check if the Author pointer is not nil.
	if post.Author != nil {
		var avatarURL string
		// Check if the AvatarURL pointer is not nil.
		if post.Author.AvatarURL != nil {
			avatarURL = *post.Author.AvatarURL
		}

		authorResponse = userclient.ClientUserResponseInfor{
			ID:        post.Author.ID,
			AvatarURL: avatarURL,
			FirstName: post.Author.FirstName,
			LastName:  post.Author.LastName,
			Username:  post.Author.Username,
			Likes:     uint(post.Author.Likes),
			Follower:  uint(post.Author.Follower),
			Following: uint(post.Author.Following),
		}
	}

	return PostResponse{
		ID:               post.ID,
		Title:            post.Title,
		Slug:             post.Slug,
		Description:      post.Description,
		ViewCount:        post.ViewCount,
		FeaturedImageURL: post.FeaturedImageURL,
		AuthorID:         post.AuthorID,
		Author:           authorResponse,

		Tags:      tagsArray,
		Status:    post.Status,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

// ToPostResponse converts a Post model to a PostResponse DTO.
func ToPostResponse(post Post) PostResponse {
	var tagsArray []string
	if post.Tags != "" {
		tagsArray = splitTagsString(post.Tags)
	}

	var authorResponse userclient.ClientUserResponseInfor
	if post.Author != nil {
		var avatarURL string
		if post.Author.AvatarURL != nil {
			avatarURL = *post.Author.AvatarURL
		}

		authorResponse = userclient.ClientUserResponseInfor{
			ID:        post.Author.ID,
			AvatarURL: avatarURL,
			FirstName: post.Author.FirstName,
			LastName:  post.Author.LastName,
			Username:  post.Author.Username,
			Likes:     uint(post.Author.Likes),
			Follower:  uint(post.Author.Follower),
			Following: uint(post.Author.Following),
		}
	}

	// Convert slices of GORM models to slices of DTOs
	commentResponses := make([]CommentResponse, len(post.Comments))
	for i, comment := range post.Comments {
		commentResponses[i] = ToCommentResponse(comment)
	}

	reactionResponses := make([]ReactionResponse, len(post.Reactions))
	for i, reaction := range post.Reactions {
		reactionResponses[i] = ToReactionResponse(reaction)
	}

	return PostResponse{
		ID:               post.ID,
		Title:            post.Title,
		Slug:             post.Slug,
		Description:      post.Description,
		ViewCount:        post.ViewCount,
		FeaturedImageURL: post.FeaturedImageURL,
		AuthorID:         post.AuthorID,
		Author:           authorResponse,
		Comments:         commentResponses,
		Reaction:         reactionResponses,
		Meta:             post.Meta,
		Link:             post.Link,
		Tags:             tagsArray,
		Status:           post.Status,
		CreatedAt:        post.CreatedAt,
		UpdatedAt:        post.UpdatedAt,
	}
}

// ToCommentResponse converts a Comment model to a CommentResponse DTO.
func ToCommentResponse(comment Comment) CommentResponse {
	var authorResponse userclient.ClientUserResponseInfor
	if comment.Author != nil {
		var avatarURL string
		if comment.Author.AvatarURL != nil {
			avatarURL = *comment.Author.AvatarURL
		}
		authorResponse = userclient.ClientUserResponseInfor{
			ID:        comment.Author.ID,
			AvatarURL: avatarURL,
			FirstName: comment.Author.FirstName,
			LastName:  comment.Author.LastName,
			Username:  comment.Author.Username,
		}
	}

	// Convert nested replies
	repliesResponses := make([]CommentResponse, len(comment.Replies))
	for i, reply := range comment.Replies {
		repliesResponses[i] = ToCommentResponse(reply)
	}

	return CommentResponse{
		ID:              comment.ID,
		PostID:          comment.PostID,
		AuthorID:        comment.AuthorID,
		ParentCommentID: comment.ParentCommentID,
		Content:         comment.Content,
		Upvotes:         comment.Upvotes,
		CreatedAt:       comment.CreatedAt,
		Author:          authorResponse, // The user data is assigned here
		Replies:         repliesResponses,
	}
}

// ToReactionResponse converts a Reaction model to a ReactionResponse DTO.
func ToReactionResponse(reaction Reaction) ReactionResponse {
	return ReactionResponse{
		ID:           reaction.ID,
		PostID:       reaction.PostID,
		UserID:       reaction.UserID,
		ReactionType: reaction.ReactionType,
		CreatedAt:    reaction.CreatedAt,
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
	if req.Meta != "" {
		post.Meta = req.Meta
	}
	if req.Link != "" {
		post.Link = req.Link
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
