package clientfollow

import "time"

type NotificationType string

const (
	Comment  NotificationType = "comment"
	Like     NotificationType = "like"
	Follower NotificationType = "follow"
	Post     NotificationType = "post"
	Message  NotificationType = "message"
)

type AccountFollow struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	FollowerID uint `gorm:"primaryKey"`
	FollowedID uint `gorm:"primaryKey"`
	CreatedAt  time.Time
}

type PostFollow struct {
	UserID    uint      `gorm:"primaryKey"`
	PostID    uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type Notification struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint `gorm:"index"`
	FromUserID uint
	PostID     uint             `gorm:"index"`
	Type       NotificationType `gorm:"size:50"`
	Message    string           `gorm:"type:text"`
	IsRead     bool             `gorm:"default:false"`
	CreatedAt  time.Time        `gorm:"default:CURRENT_TIMESTAMP"`
}
