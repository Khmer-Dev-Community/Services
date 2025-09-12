package clientfollow

import "errors"

// NotificationService defines business logic for notifications.
type NotificationService interface {
	GetUserNotifications(userID uint) ([]Notification, error)
	MarkNotificationAsRead(notificationID uint) error
}

// FollowService defines business logic for follower relationships.
type FollowService interface {
	ToggleFollow(followerID, followedID uint) (bool, error)
	GetUserFollowers(userID uint) ([]AccountFollow, error)
	GetUserFollowing(userID uint) ([]AccountFollow, error)
}

// GormNotificationService is a concrete implementation of NotificationService.
type GormNotificationService struct {
	repo NotificationRepository
}

// NewGormNotificationService creates a new instance of the notification service.
func NewGormNotificationService(r NotificationRepository) NotificationService {
	return &GormNotificationService{repo: r}
}

// GormFollowService is a concrete implementation of FollowService.
type GormFollowService struct {
	repo FollowRepository
}

// NewGormFollowService creates a new instance of the follow service.
func NewGormFollowService(r FollowRepository) FollowService {
	return &GormFollowService{repo: r}
}

// NOTIFICATION SERVICE FUNCTIONS

// GetUserNotifications retrieves all notifications for a user.
func (s *GormNotificationService) GetUserNotifications(userID uint) ([]Notification, error) {
	// Business logic: You could add filtering, sorting, or pagination here
	return s.repo.GetNotificationsByUserID(userID)
}

// MarkNotificationAsRead marks a single notification as read.
func (s *GormNotificationService) MarkNotificationAsRead(notificationID uint) error {
	// Business logic: You could check user ownership here before calling the repo
	return s.repo.MarkAsRead(notificationID)
}

// FOLLOW SERVICE FUNCTIONS

// ToggleFollow toggles the follow status between two users.
func (s *GormFollowService) ToggleFollow(followerID, followedID uint) (bool, error) {
	// Business logic: Prevent a user from following themselves
	if followerID == followedID {
		return false, errors.New("cannot follow yourself")
	}

	// Check if the follower relationship already exists
	isFollowing, err := s.repo.IsFollowing(followerID, followedID)
	if err != nil {
		return false, err
	}

	if isFollowing {
		// If they are already following, unfollow them
		err := s.repo.UnfollowUser(followerID, followedID)
		return false, err // Return false because the user is now unfollowed
	} else {
		// If they are not following, create a new follow record
		follow := &AccountFollow{
			FollowerID: followerID,
			FollowedID: followedID,
		}
		err := s.repo.FollowUser(follow)
		return true, err // Return true because the user is now followed
	}
}

// GetUserFollowers retrieves all followers for a given user.
func (s *GormFollowService) GetUserFollowers(userID uint) ([]AccountFollow, error) {
	return s.repo.GetFollowers(userID)
}

// GetUserFollowing retrieves all users a given user is following.
func (s *GormFollowService) GetUserFollowing(userID uint) ([]AccountFollow, error) {
	return s.repo.GetFollowing(userID)
}
