package setting

import (
	"time"

	"gorm.io/gorm"
)

type VideoSnapSetting struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	ServerName       string         `gorm:"column:server_name" json:"server_name"`
	Rtmp             string         `gorm:"column:rtmp" json:"rtmp"`
	Duration         int            `json:"duration" gorm:"column:duration"`
	VideoSize        string         `json:"video_size" gorm:"column:video_size"`
	VideoType        string         `json:"video_type" gorm:"column:video_type"`
	BeforeStart      int            `json:"before_start" gorm:"column:before_start"`
	AfterEnd         int            `json:"after_end" gorm:"column:after_end"`
	OutputPath       string         `json:"output_path" gorm:"column:output_path"`
	ServiceAcccount  string         `json:"service_account" gorm:"column:service_account"`
	BucketName       string         `json:"bucket_name" gorm:"column:bucket_name"`
	AcceessDomain    string         `json:"access_domain" gorm:"column:access_domain"`
	DeleteLocalStore bool           `json:"delete_local_store" gorm:"column:delete_local_store"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (VideoSnapSetting) TableName() string {
	return "video_setting"
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&VideoSnapSetting{})
}
