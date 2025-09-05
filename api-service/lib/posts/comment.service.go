package posts

import (
	"context"
	"fmt"

	"github.com/Khmer-Dev-Community/Services/api-service/utils"
)

// CommentService defines the interface for our comment-related business logic.
type CommentService interface {
	CreateComment(ctx context.Context, comment *Comment) (*Comment, error)
	GetCommentsByPostID(ctx context.Context, postID uint) ([]*Comment, error)
	DeleteComment(ctx context.Context, userID uint, commentID uint) error
}

// GormCommentService is the concrete implementation of the CommentService interface.
type GormCommentService struct {
	repo CommentRepository
}

// NewGormCommentService creates a new instance of the CommentService.
func NewGormCommentService(repo CommentRepository) CommentService {
	return &GormCommentService{repo: repo}
}

// CreateComment handles the creation of a new comment.
func (s *GormCommentService) CreateComment(ctx context.Context, comment *Comment) (*Comment, error) {
	utils.InfoLog(comment, "CreateComment")
	if comment.Content == "" {
		err := fmt.Errorf("comment content cannot be empty")
		utils.InfoLog(err, "CreateComment")
		return nil, err
	}

	createdComment, err := s.repo.CreateComment(ctx, comment)
	if err != nil {
		utils.ErrorLog(err.Error(), "CreateComment")
		return nil, err
	}

	return createdComment, nil
}

// GetCommentsByPostID retrieves all top-level comments and their nested replies for a given post.
func (s *GormCommentService) GetCommentsByPostID(ctx context.Context, postID uint) ([]*Comment, error) {
	utils.InfoLog(fmt.Sprintf("Fetching comments for post ID: %d", postID), "GetCommentsByPostID")

	comments, err := s.repo.GetCommentsByPostID(ctx, postID)
	if err != nil {
		utils.ErrorLog(err, "GetCommentsByPostID")
		return nil, err
	}

	return comments, nil
}

// DeleteComment handles the deletion of a comment.
func (s *GormCommentService) DeleteComment(ctx context.Context, userID uint, commentID uint) error {
	utils.InfoLog(fmt.Sprintf("Attempting to delete comment ID: %d by user: %d", commentID, userID), "DeleteComment")
	comment, err := s.repo.GetCommentByID(ctx, commentID)
	if err != nil {
		utils.ErrorLog(err, "DeleteComment")
		return fmt.Errorf("comment not found or an error occurred")
	}
	if comment.AuthorID != userID {
		err := fmt.Errorf("user is not authorized to delete this comment")
		utils.WarnLog(err, "DeleteComment")
		return err
	}
	if err := s.repo.DeleteComment(ctx, commentID); err != nil {
		utils.ErrorLog(err, "DeleteComment")
		return err
	}
	utils.InfoLog(fmt.Sprintf("Successfully deleted comment ID: %d", commentID), "DeleteComment")
	return nil
}
