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
	var finalReaction *Reaction

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingReaction Reaction

		// Find the reaction
		result := tx.Where("post_id = ? AND user_id = ?", reaction.PostID, reaction.UserID).First(&existingReaction)

		if result.Error == gorm.ErrRecordNotFound {
			// Case 1: Reaction does not exist, so create it.
			if err := tx.Create(reaction).Error; err != nil {
				return err
			}

			// Atomically increment the post's reaction count.
			if err := tx.Model(&Post{}).
				Where("id = ?", reaction.PostID).
				Update("reaction_count", gorm.Expr("reaction_count + ?", 1)).Error; err != nil {
				return err
			}

			// Assign the created reaction to the outer variable
			finalReaction = reaction

		} else if result.Error != nil {
			// Case 2: An unexpected database error occurred.
			return result.Error

		} else {
			// Case 3: Reaction exists, so delete it.
			if err := tx.Delete(&existingReaction).Error; err != nil {
				return err
			}

			// Atomically decrement the post's reaction count.
			if err := tx.Model(&Post{}).
				Where("id = ?", existingReaction.PostID).
				Update("reaction_count", gorm.Expr("reaction_count - ?", 1)).Error; err != nil {
				return err
			}

			// Assign nil to the outer variable to indicate a deletion
			finalReaction = nil
		}

		return nil
	})

	// Now, return the values from the outer function
	return finalReaction, err
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
