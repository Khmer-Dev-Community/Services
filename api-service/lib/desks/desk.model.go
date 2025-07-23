package desks

import (
	"telegram-service/lib/videosnap/setting"
	"time"

	"gorm.io/gorm"
)

type DeskSetting struct {
	ID               uint                     `json:"id" gorm:"primaryKey"`
	DeskNo           string                   `gorm:"column:desk_no" json:"desk_no"`
	DeskStreamKey    string                   `gorm:"column:desk_streamkey" json:"desk_streamkey"`
	DeskStreamServer uint                     `gorm:"column:desk_stream_server" json:"desk_stream_server"`
	StreamServer     setting.VideoSnapSetting `json:"server" gorm:"foreignKey:DeskStreamServer;references:ID;constraint:OnDelete:SET NULL,OnUpdate:SET NULL"`
	CreatedAt        time.Time                `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time                `gorm:"column:updated_at" json:"updatedAt"`
}

func (DeskSetting) TableName() string {
	return "desk_setting"
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&DeskSetting{})
}
