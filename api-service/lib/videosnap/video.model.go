package videosnap

import (
	"time"

	"gorm.io/gorm"
)

type Videosnap struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ProcessId uint           `gorm:"column:process_id" json:"process_id"`
	Code      string         `gorm:"column:tran_id" json:"tran_id"`
	GameNo    string         `gorm:"column:game_no" json:"period"`
	Streamkey string         `gorm:"column:stream_key" json:"rtmpurl"`
	Rtmp      string         `gorm:"column:rtmp" json:"rtmp"`
	ImageURL  string         `gorm:"column:image_url" json:"image_url"`
	VideoURL  string         `gorm:"column:video_url" json:"video_url"`
	StorePath string         `gorm:"column:store_path" json:"store_path"`
	Status    bool           `gorm:"column:status" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Videosnap) TableName() string {
	return "video_records"
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&Videosnap{})
}
