package posts

import (
	"context"
	"errors"
	"fmt"

	"github.com/Khmer-Dev-Community/Services/api-service/utils"
	"github.com/gosimple/slug"
	"github.com/sirupsen/logrus"
)

var ErrUnauthorizedPostAction = errors.New("unauthorized to perform action on this post")
var ErrSlugGenerationFailed = errors.New("failed to generate unique slug")

// Assuming ErrPostNotFound is defined elsewhere, e.g., in a DTO or repository file.
// var ErrPostNotFound = errors.New("post not found") // If it's not, you might need to uncomment this or define it.

type PostService interface {
	CreatePost(ctx context.Context, req CreatePostRequest, authorID uint) (*PostResponse, error)
	GetPostByID(ctx context.Context, id uint) (*PostResponse, error)
	GetPostBySlug(ctx context.Context, slug string) (*PostResponse, error)
	UpdatePost(ctx context.Context, id uint, req UpdatePostRequest, userID uint) (*PostResponse, error)
	DeletePost(ctx context.Context, id uint, userID uint) error
	ListPosts(ctx context.Context, offset, limit int, status, tag string) ([]PostResponse, int64, error)
}

type postService struct {
	repo PostRepository // CORRECTED: Holds the PostRepository interface directly
}

// NewPostService creates a new instance of PostService.
func NewPostService(repo PostRepository) PostService { // CORRECTED: Accepts the PostRepository interface directly
	return &postService{repo: repo}
}

func (s *postService) CreatePost(ctx context.Context, req CreatePostRequest, authorID uint) (*PostResponse, error) {
	generatedSlug := slug.Make(req.Title)
	if generatedSlug == "" {
		utils.ErrorLog(map[string]interface{}{"title": req.Title}, "Failed to generate slug for post creation")
		return nil, ErrSlugGenerationFailed
	}

	if req.Title == "" || req.Description == "" {
		utils.WarnLog(nil, "Attempted to create post with empty title or description")
		return nil, errors.New("title and description cannot be empty")
	}

	postModel := ToPostModel(req, authorID, generatedSlug)
	postModel.Status = getDefaultStatus(req.Status) // Assuming getDefaultStatus is defined in this package

	err := s.repo.CreatePost(ctx, &postModel)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"post": postModel, "error": err.Error()}, "Failed to create post in repository")
		return nil, fmt.Errorf("service: failed to create post: %w", err)
	}

	utils.LoggerService(map[string]interface{}{"post_id": postModel.ID, "title": postModel.Title}, "Post created successfully", logrus.InfoLevel)
	response := ToPostResponse(postModel) // Assuming ToPostResponse is defined in this package
	return &response, nil
}

func (s *postService) GetPostByID(ctx context.Context, id uint) (*PostResponse, error) {
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		// Assuming ErrPostNotFound is accessible (e.g., from a common errors file or the repo package)
		if errors.Is(err, ErrPostNotFound) {
			utils.InfoLog(map[string]interface{}{"post_id": id}, "Post not found by ID")
			return nil, ErrPostNotFound
		}
		utils.ErrorLog(map[string]interface{}{"post_id": id, "error": err.Error()}, "Failed to retrieve post by ID from repository")
		return nil, fmt.Errorf("service: failed to get post by ID %d: %w", id, err)
	}

	err = s.repo.IncrementPostViewCount(ctx, id)
	if err != nil {
		utils.WarnLog(map[string]interface{}{"post_id": id, "error": err.Error()}, "Failed to increment view count for post")
		// Do not return error here, as the primary goal is to retrieve the post
	}

	utils.LoggerService(map[string]interface{}{"post_id": id, "title": post.Title}, "Post retrieved by ID successfully", logrus.InfoLevel)
	response := ToPostResponse(*post)
	return &response, nil
}

func (s *postService) GetPostBySlug(ctx context.Context, postSlug string) (*PostResponse, error) {
	post, err := s.repo.GetPostBySlug(ctx, postSlug)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			utils.InfoLog(map[string]interface{}{"slug": postSlug}, "Post not found by slug")
			return nil, ErrPostNotFound
		}
		utils.ErrorLog(map[string]interface{}{"slug": postSlug, "error": err.Error()}, "Failed to retrieve post by slug from repository")
		return nil, fmt.Errorf("service: failed to get post by slug '%s': %w", postSlug, err)
	}

	err = s.repo.IncrementPostViewCount(ctx, post.ID)
	if err != nil {
		utils.WarnLog(map[string]interface{}{"post_id": post.ID, "error": err.Error()}, "Failed to increment view count for post retrieved by slug")
	}

	utils.LoggerService(map[string]interface{}{"slug": postSlug, "post_id": post.ID}, "Post retrieved by slug successfully", logrus.InfoLevel)
	response := ToPostResponse(*post)
	return &response, nil
}

func (s *postService) UpdatePost(ctx context.Context, id uint, req UpdatePostRequest, userID uint) (*PostResponse, error) {
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			utils.InfoLog(map[string]interface{}{"post_id": id, "user_id": userID}, "Attempted to update non-existent post")
			return nil, ErrPostNotFound
		}
		utils.ErrorLog(map[string]interface{}{"post_id": id, "user_id": userID, "error": err.Error()}, "Failed to get post for update from repository")
		return nil, fmt.Errorf("service: failed to get post %d for update: %w", id, err)
	}

	if post.AuthorID != userID {
		utils.WarnLog(map[string]interface{}{"post_id": id, "author_id": post.AuthorID, "attempted_user_id": userID}, "Unauthorized attempt to update post")
		return nil, ErrUnauthorizedPostAction
	}

	var newSlug *string

	if req.Title != nil {
		generatedSlug := slug.Make(*req.Title)
		if generatedSlug == "" {
			utils.ErrorLog(map[string]interface{}{"post_id": id, "new_title": *req.Title}, "Failed to generate new slug for post update")
			return nil, ErrSlugGenerationFailed
		}
		newSlug = &generatedSlug
	}

	UpdatePostModel(post, req, newSlug) // Assuming UpdatePostModel is defined in this package

	err = s.repo.UpdatePost(ctx, post)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"post": post, "error": err.Error()}, "Failed to update post in repository")
		return nil, fmt.Errorf("service: failed to update post %d: %w", id, err)
	}

	utils.LoggerService(map[string]interface{}{"post_id": id, "title": post.Title, "user_id": userID}, "Post updated successfully", logrus.InfoLevel)
	response := ToPostResponse(*post)
	return &response, nil
}

func (s *postService) DeletePost(ctx context.Context, id uint, userID uint) error {
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			utils.InfoLog(map[string]interface{}{"post_id": id, "user_id": userID}, "Attempted to delete non-existent post")
			return ErrPostNotFound
		}
		utils.ErrorLog(map[string]interface{}{"post_id": id, "user_id": userID, "error": err.Error()}, "Failed to get post for deletion from repository")
		return fmt.Errorf("service: failed to get post %d for deletion: %w", id, err)
	}

	if post.AuthorID != userID {
		utils.WarnLog(map[string]interface{}{"post_id": id, "author_id": post.AuthorID, "attempted_user_id": userID}, "Unauthorized attempt to delete post")
		return ErrUnauthorizedPostAction
	}

	err = s.repo.DeletePost(ctx, id)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"post_id": id, "error": err.Error()}, "Failed to delete post in repository")
		return fmt.Errorf("service: failed to delete post %d: %w", id, err)
	}

	utils.LoggerService(map[string]interface{}{"post_id": id, "user_id": userID}, "Post deleted successfully", logrus.InfoLevel)
	return nil
}

func (s *postService) ListPosts(ctx context.Context, offset, limit int, status, tag string) ([]PostResponse, int64, error) {
	if limit <= 0 || limit > 100 {
		utils.WarnLog(map[string]interface{}{"requested_limit": limit, "defaulting_to": 100}, "Invalid limit for listing posts, defaulting to 100")
		limit = 100 // Set to a sensible default, maybe your default is different
	}
	if offset < 0 {
		utils.WarnLog(map[string]interface{}{"requested_offset": offset, "defaulting_to": 0}, "Invalid offset for listing posts, defaulting to 0")
		offset = 0
	}

	posts, err := s.repo.ListPosts(ctx, offset, limit, status, tag)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"offset": offset, "limit": limit, "status": status, "tag": tag, "error": err.Error()}, "Failed to list posts from repository")
		return nil, 0, fmt.Errorf("service: failed to list posts: %w", err)
	}

	totalCount, err := s.repo.CountPosts(ctx, status, tag)
	if err != nil {
		utils.ErrorLog(map[string]interface{}{"status": status, "tag": tag, "error": err.Error()}, "Failed to count posts from repository")
		return nil, 0, fmt.Errorf("service: failed to count posts: %w", err)
	}

	utils.LoggerService(map[string]interface{}{"offset": offset, "limit": limit, "count": len(posts), "total": totalCount}, "Posts listed successfully", logrus.InfoLevel)
	responses := ToPostResponseList(posts) // Assuming ToPostResponseList is defined in this package
	return responses, totalCount, nil
}
