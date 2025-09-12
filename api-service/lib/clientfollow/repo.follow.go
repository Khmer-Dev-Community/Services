package clientfollow

import "gorm.io/gorm"

// NotificationRepository defines the methods for interacting with notifications.
type NotificationRepository interface {
	CreateNotification(notification *Notification) error
	GetNotificationsByUserID(userID uint) ([]Notification, error)
	MarkAsRead(notificationID uint) error
	DeleteNotification(notificationID uint) error
}

// FollowRepository defines the methods for managing follow relationships.
type FollowRepository interface {
	FollowUser(follow *AccountFollow) error
	UnfollowUser(followerID, followedID uint) error
	IsFollowing(followerID, followedID uint) (bool, error)
	GetFollowers(userID uint) ([]AccountFollow, error)
	GetFollowing(userID uint) ([]AccountFollow, error)
}

// GormNotificationRepository is a concrete implementation of NotificationRepository using GORM.
type GormNotificationRepository struct {
	DB *gorm.DB
}

// NewGormNotificationRepository creates a new instance of the notification repository.
func NewGormNotificationRepository(db *gorm.DB) NotificationRepository {
	return &GormNotificationRepository{DB: db}
}

// GormFollowRepository is a concrete implementation of FollowRepository using GORM.
type GormFollowRepository struct {
	DB *gorm.DB
}

// NewGormFollowRepository creates a new instance of the follow repository.
func NewGormFollowRepository(db *gorm.DB) FollowRepository {
	return &GormFollowRepository{DB: db}
}

// GORM NOTIFICATION REPOSITORY METHODS

func (r *GormNotificationRepository) CreateNotification(notification *Notification) error {
	result := r.DB.Create(notification)
	return result.Error
}

func (r *GormNotificationRepository) GetNotificationsByUserID(userID uint) ([]Notification, error) {
	var notifications []Notification
	result := r.DB.Where("user_id = ?", userID).Find(&notifications)
	return notifications, result.Error
}

func (r *GormNotificationRepository) MarkAsRead(notificationID uint) error {
	result := r.DB.Model(&Notification{}).Where("id = ?", notificationID).Update("is_read", true)
	return result.Error
}

func (r *GormNotificationRepository) DeleteNotification(notificationID uint) error {
	result := r.DB.Delete(&Notification{}, notificationID)
	return result.Error
}

// GORM FOLLOW REPOSITORY METHODS

func (r *GormFollowRepository) FollowUser(follow *AccountFollow) error {
	result := r.DB.Create(follow)
	return result.Error
}

func (r *GormFollowRepository) UnfollowUser(followerID, followedID uint) error {
	result := r.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).Delete(&AccountFollow{})
	return result.Error
}

func (r *GormFollowRepository) IsFollowing(followerID, followedID uint) (bool, error) {
	var count int64
	err := r.DB.Model(&AccountFollow{}).Where("follower_id = ? AND followed_id = ?", followerID, followedID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormFollowRepository) GetFollowers(userID uint) ([]AccountFollow, error) {
	var followers []AccountFollow
	result := r.DB.Where("followed_id = ?", userID).Find(&followers)
	return followers, result.Error
}

func (r *GormFollowRepository) GetFollowing(userID uint) ([]AccountFollow, error) {
	var following []AccountFollow
	result := r.DB.Where("follower_id = ?", userID).Find(&following)
	return following, result.Error
}
