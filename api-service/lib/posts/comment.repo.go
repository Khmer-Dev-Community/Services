package posts

import (
	"context"

	"github.com/Khmer-Dev-Community/Services/api-service/utils"
	"gorm.io/gorm"
)

type GormCommentRepository struct {
	db *gorm.DB
}
type CommentRepository interface {
	CreateComment(ctx context.Context, post *Comment) (*Comment, error)
	GetCommentsByPostID(ctx context.Context, postID uint) ([]*Comment, error)
	GetCommentByID(ctx context.Context, commentID uint) (*Comment, error)
	DeleteComment(ctx context.Context, commentID uint) error
}

func NewGormCommentRepository(db *gorm.DB) CommentRepository {
	return &GormCommentRepository{db: db}
}

func (r *GormCommentRepository) CreateComment(ctx context.Context, post *Comment) (*Comment, error) {
	err := r.db.Create(post).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}

	utils.LoggerRepository(post, "Execute")
	return post, nil
}

func (r *GormCommentRepository) GetCommentsByPostID(ctx context.Context, postID uint) ([]*Comment, error) {
	var comments []*Comment
	err := r.db.WithContext(ctx).
		Preload("Replies"). // Preload replies
		Where("post_id = ? AND parent_comment_id IS NULL", postID).
		Find(&comments).Error

	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *GormCommentRepository) GetCommentByID(ctx context.Context, commentID uint) (*Comment, error) {
	var comment Comment
	err := r.db.WithContext(ctx).First(&comment, commentID).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// DeleteComment removes a comment from the database by its ID.
func (r *GormCommentRepository) DeleteComment(ctx context.Context, commentID uint) error {
	return r.db.WithContext(ctx).Delete(&Comment{}, commentID).Error
}
