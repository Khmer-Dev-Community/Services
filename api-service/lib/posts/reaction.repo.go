package posts

import (
	"context"

	"gorm.io/gorm"
)

type GormReactionRepository struct {
	db *gorm.DB
}
type ReactionRepository interface {
	CreateOrUpdate(ctx context.Context, reaction *Reaction) (*Reaction, error)
	Delete(ctx context.Context, postID, userID uint) error
	GetByPostAndUser(ctx context.Context, postID, userID uint) (*Reaction, error)
	GetCountByPostID(ctx context.Context, postID uint) (int64, error)
}

func NewGormReactionRepository(db *gorm.DB) ReactionRepository {
	return &GormReactionRepository{db: db}
}

// CreateOrUpdate handles both creating a new reaction or updating an existing one.
// This prevents a user from having multiple reactions on the same post.
func (r *GormReactionRepository) CreateOrUpdate(ctx context.Context, reaction *Reaction) (*Reaction, error) {
	var existingReaction Reaction
	result := r.db.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", reaction.PostID, reaction.UserID).
		First(&existingReaction)

	if result.Error == gorm.ErrRecordNotFound {
		// Reaction does not exist, create it
		if err := r.db.WithContext(ctx).Create(reaction).Error; err != nil {
			return nil, err
		}
		return reaction, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	// Reaction exists, update its type
	if err := r.db.WithContext(ctx).Model(&existingReaction).Update("reaction_type", reaction.ReactionType).Error; err != nil {
		return nil, err
	}
	return &existingReaction, nil
}

// Delete removes a user's reaction from a specific post.
func (r *GormReactionRepository) Delete(ctx context.Context, postID, userID uint) error {
	return r.db.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&Reaction{}).Error
}

// GetByPostAndUser retrieves a specific reaction for a post by a user.
func (r *GormReactionRepository) GetByPostAndUser(ctx context.Context, postID, userID uint) (*Reaction, error) {
	var reaction Reaction
	err := r.db.WithContext(ctx).
		Where("post_id = ? AND user_id = ?", postID, userID).
		First(&reaction).Error

	if err != nil {
		return nil, err
	}
	return &reaction, nil
}

// GetCountByPostID retrieves the total number of reactions for a post.
func (r *GormReactionRepository) GetCountByPostID(ctx context.Context, postID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&Reaction{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
