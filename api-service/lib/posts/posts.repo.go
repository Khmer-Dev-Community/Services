package posts

import (
	"context"
	"errors"

	"github.com/Khmer-Dev-Community/Services/api-service/utils"
	"gorm.io/gorm" // Make sure you have GORM installed: go get gorm.io/gorm
)

var ErrPostNotFound = errors.New("post not found")

type PostRepository interface {
	CreatePost_(ctx context.Context, post *Post) error
	CreatePost(ctx context.Context, post *Post) (*Post, error)
	GetPostByID(ctx context.Context, id uint) (*Post, error)
	GetPostBySlug(ctx context.Context, slug string) (*Post, error)
	UpdatePost(ctx context.Context, post *Post) error
	DeletePost(ctx context.Context, id uint) error
	ListPosts(ctx context.Context, offset, limit int, status, tag string, username string, viewerId uint) ([]Post, error)
	CountPosts(ctx context.Context, status, tag string) (int64, error)
	IncrementPostViewCount(ctx context.Context, id uint) error
}

// GormPostRepository is an implementation of PostRepository using GORM.
type GormPostRepository struct {
	db *gorm.DB // GORM database client
}

// NewGormPostRepository creates a new GormPostRepository instance.
func NewGormPostRepository(db *gorm.DB) PostRepository {
	return &GormPostRepository{db: db}
}

// CreatePost implements PostRepository.CreatePost
func (r *GormPostRepository) CreatePost_(ctx context.Context, post *Post) error {
	result := r.db.WithContext(ctx).Create(post)
	return result.Error
}
func (r *GormPostRepository) CreatePost(ctx context.Context, post *Post) (*Post, error) {
	err := r.db.Create(post).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}

	utils.LoggerRepository(post, "Execute")
	return post, nil
}

// GetPostByID implements PostRepository.GetPostByID
func (r *GormPostRepository) GetPostByID(ctx context.Context, id uint) (*Post, error) {
	var post Post

	// Build the base query and apply the Preload
	query := r.db.WithContext(ctx).
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "avatar_url", "first_name", "last_name", "username", "likes", "follower", "following")
		})

	// Execute the query
	result := query.First(&post, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, result.Error
	}
	return &post, nil
}

// GetPostBySlug implements PostRepository.GetPostBySlug
func (r *GormPostRepository) GetPostBySlugNOComment(ctx context.Context, slug string) (*Post, error) {
	var post Post

	// Build the base query with the Where clause and Preload
	query := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "avatar_url", "first_name", "last_name", "username", "likes", "follower", "following")
		})

	// Execute the query
	result := query.First(&post)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, result.Error
	}
	return &post, nil
}
func (r *GormPostRepository) GetPostBySlug(ctx context.Context, slug string) (*Post, error) {
	var post Post

	query := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "avatar_url", "first_name", "last_name", "username", "likes", "follower", "following")
		}).
		// Preload top-level comments (level 1)
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Where("parent_comment_id IS NULL").
				Order("created_at desc").
				Preload("Author", func(db *gorm.DB) *gorm.DB {
					return db.Select("id", "avatar_url", "first_name", "last_name", "username")
				}).
				// Preload replies to top-level comments (level 2)
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
		})

	result := query.First(&post)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, result.Error
	}
	return &post, nil
}

// UpdatePost implements PostRepository.UpdatePost
func (r *GormPostRepository) UpdatePost(ctx context.Context, post *Post) error {
	result := r.db.WithContext(ctx).Save(post)
	return result.Error
}

// DeletePost implements PostRepository.DeletePost
func (r *GormPostRepository) DeletePost(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&Post{}, id)
	if result.RowsAffected == 0 {
		return ErrPostNotFound // Or specific error if record didn't exist to delete
	}
	return result.Error
}

// ListPosts implements PostRepository.ListPosts
func (r *GormPostRepository) ListPosts(ctx context.Context, offset, limit int, status, tag string, account string, vewerId uint) ([]Post, error) {
	var posts []Post
	if limit >= 25 {
		limit = 25
	}
	utils.WarnLog(vewerId, " viewer Id")
	// Build the query and apply all conditions first
	query := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
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
			return db.Where("client_reaction_post.user_id=?", vewerId)
		})

	// Apply the conditional WHERE clauses to the query builder
	if status != "" {
		query = query.Where("client_content_post.status = ?", status)
	}
	if tag != "" {
		query = query.Where("client_content_post.tags LIKE ?", "%"+tag+"%")
	}
	if account != "" {
		query = query.Joins("JOIN client_platform_user ON client_platform_user.id = client_content_post.author_id").
			Where("client_platform_user.username = ?", account)
	}

	// Now, apply the final Select and Find
	result := query.Find(&posts)
	//	Select("*, (SELECT COUNT(*) FROM client_reaction_post WHERE client_reaction_post.post_id = client_content_post.id) AS reaction_count").
	//	Find(&posts)

	return posts, result.Error
}

// CountPosts implements PostRepository.CountPosts
func (r *GormPostRepository) CountPosts(ctx context.Context, status, tag string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&Post{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}

	result := query.Count(&count)
	return count, result.Error
}

// IncrementPostViewCount implements PostRepository.IncrementPostViewCount
func (r *GormPostRepository) IncrementPostViewCount(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Model(&Post{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrPostNotFound
	}
	return nil
}
