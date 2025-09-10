package posts

import (
	"context"

	"github.com/Khmer-Dev-Community/Services/api-service/utils"
)

// SavedPostServiceImpl is the concrete implementation of SavedPostService.
type SavedPostServiceImpl struct {
	repo GormSavedPostRepository
}

// NewSavedPostService creates and returns a new instance of SavedPostServiceImpl.
func NewSavedPostService(repo GormSavedPostRepository) *SavedPostServiceImpl {
	return &SavedPostServiceImpl{repo: repo}
}

// SavePost handles the business logic for saving a post.
func (s *SavedPostServiceImpl) SavePost(ctx context.Context, userID, postID uint, userName string) error {
	savedPost := SavedPost{
		UserID:   userID,
		PostID:   postID,
		Username: userName,
	}
	utils.InfoLog(savedPost, "SavePost Request data")
	return s.repo.Save(ctx, &savedPost)
}

// UnsavePost removes a saved post record.
func (s *SavedPostServiceImpl) UnsavePost(ctx context.Context, userID, postID uint) error {
	return s.repo.Delete(ctx, userID, postID)
}

// GetSavedPosts fetches all posts saved by a user and maps them to DTOs.
func (s *SavedPostServiceImpl) GetSavedPosts(ctx context.Context, userID string, viewerID uint, offset, limit int) ([]SavePostData, error) {
	savedPosts, err := s.repo.FindByUserID(ctx, userID, viewerID, offset, limit)
	if err != nil {
		return nil, err
	}

	var responses []SavePostData
	for _, savedPost := range savedPosts {

		responses = append(responses, ToSavedPostData(savedPost))
	}

	return responses, nil
}

// IsPostSaved checks if a user has saved a specific post.
func (s *SavedPostServiceImpl) IsPostSaved(ctx context.Context, userID, postID uint) (bool, error) {
	return s.repo.IsPostSaved(ctx, userID, postID)
}
