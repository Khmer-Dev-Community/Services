package clientfollow

import "time"

type NotificationDTO struct {
	ID         uint      `json:"id"`
	FromUserID uint      `json:"from_user_id"`
	PostID     uint      `json:"post_id"`
	Type       string    `json:"type"`
	Message    string    `json:"message"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

type AccountFollowDTO struct {
	FollowerID uint      `json:"follower_id"`
	FollowedID uint      `json:"followed_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type PostFollowDTO struct {
	UserID    uint      `json:"user_id"`
	PostID    uint      `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}
