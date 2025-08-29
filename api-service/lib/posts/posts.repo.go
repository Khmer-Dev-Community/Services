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
	ListPosts(ctx context.Context, offset, limit int, status, tag string) ([]Post, error)
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
	result := r.db.WithContext(ctx).First(&post, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, result.Error
	}
	return &post, nil
}

// GetPostBySlug implements PostRepository.GetPostBySlug
func (r *GormPostRepository) GetPostBySlug(ctx context.Context, slug string) (*Post, error) {
	var post Post
	result := r.db.WithContext(ctx).Where("slug = ?", slug).First(&post)
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
func (r *GormPostRepository) ListPosts(ctx context.Context, offset, limit int, status, tag string) ([]Post, error) {
	var posts []Post
	if limit > 50 {
		limit = 50
	}
	query := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at DESC")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}
	result := query.Find(&posts)
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
