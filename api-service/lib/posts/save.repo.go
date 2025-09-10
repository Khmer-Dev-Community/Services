package posts

import (
	"context"

	"github.com/Khmer-Dev-Community/Services/api-service/utils"
	"gorm.io/gorm"
)

// GormSavedPostRepository is a concrete GORM implementation of SavedPostRepository.
type GormSavedPostRepository struct {
	db *gorm.DB
}

// NewGormSavedPostRepository creates a new instance of GormSavedPostRepository.
func NewGormSavedPostRepository(db *gorm.DB) *GormSavedPostRepository {
	return &GormSavedPostRepository{db: db}
}

// Save creates a new SavedPost record.
func (r *GormSavedPostRepository) Saved(ctx context.Context, savedPost *SavedPost) error {
	return r.db.WithContext(ctx).Create(savedPost).Error
}
func (r *GormSavedPostRepository) Save(ctx context.Context, savedPost *SavedPost) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingSavedPost SavedPost
		// Check if the post is already saved by this user
		result := tx.Where("post_id = ? AND user_id = ?", savedPost.PostID, savedPost.UserID).First(&existingSavedPost)
		utils.InfoLog(result, "Check Post Save")
		if result.Error == gorm.ErrRecordNotFound {
			// Case 1: Post is not saved, so create a new record.
			if err := tx.Create(savedPost).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Unscoped().Delete(&existingSavedPost).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Delete removes a SavedPost record by user and post ID.
func (r *GormSavedPostRepository) Delete(ctx context.Context, postID, userID uint) error {
	utils.InfoLog(postID, "Excute Delete Save")
	return r.db.WithContext(ctx).
		Unscoped().
		Where("post_id = ? AND user_id = ?", postID, userID).
		Delete(&SavedPost{}).Error
}

func (r *GormSavedPostRepository) FindByUserID(ctx context.Context, userID string, viewerID uint, offset, limit int) ([]SavedPost, error) {
	var savedPosts []SavedPost
	if limit >= 25 {
		limit = 25
	}
	query := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Where("username=?", userID)
	query = query.
		Preload("Post", func(db *gorm.DB) *gorm.DB {
			return db.
				Order("client_content_post.created_at DESC").
				Preload("Author", func(db *gorm.DB) *gorm.DB {
					return db.Select("id", "avatar_url", "first_name", "last_name", "username", "likes", "follower", "following")
				}).
				Preload("Comments", func(db *gorm.DB) *gorm.DB {
					return db.Where("parent_comment_id IS NULL").
						Order("created_at desc").
						Preload("Author", func(db *gorm.DB) *gorm.DB {
							return db.Select("id", "avatar_url", "first_name", "last_name", "username")
						}).
						Preload("Replies", func(db *gorm.DB) *gorm.DB {
							return db.Order("created_at desc").
								Preload("Author", func(db *gorm.DB) *gorm.DB {
									return db.Select("id", "avatar_url", "first_name", "last_name", "username")
								}).
								Preload("Replies", func(db *gorm.DB) *gorm.DB {
									return db.Order("created_at desc").
										Preload("Author", func(db *gorm.DB) *gorm.DB {
											return db.Select("id", "avatar_url", "first_name", "last_name", "username")
										})
								})
						})
				}).
				Preload("Reactions", func(db *gorm.DB) *gorm.DB {
					return db.Where("client_reaction_post.user_id=?", viewerID)
				})
		})

	// Execute the query.
	result := query.Find(&savedPosts)
	return savedPosts, result.Error
}

func (r *GormSavedPostRepository) IsPostSaved(ctx context.Context, userID, postID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&SavedPost{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
