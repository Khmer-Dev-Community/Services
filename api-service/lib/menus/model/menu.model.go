package menus

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID         int            `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"column:name" json:"name"`
	Path       string         `gorm:"column:path" json:"path"`
	Redirect   string         `gorm:"column:redirect" json:"redirect"`
	Hidden     bool           `gorm:"column:hidden" json:"hidden"`
	Title      string         `gorm:"column:title" json:"title"`
	Icon       string         `gorm:"column:icon" json:"icon"`
	IconColor  string         `gorm:"column:icon_color" json:"icon_color"`
	NoCache    bool           `gorm:"column:no_cache" json:"noCache"`   // Use snake_case for consistency
	TitleKey   string         `gorm:"column:title_key" json:"titleKey"` // Use snake_case for consistency
	Link       string         `gorm:"column:link" json:"link"`
	SubOf      int            `gorm:"column:subof" json:"subOf"` // Use snake_case for consistency
	Component  string         `gorm:"column:component" json:"component"`
	OrderNum   int            `gorm:"column:order_num" json:"orderNum"` // Use snake_case for consistency
	IsFrame    int            `gorm:"column:is_frame" json:"isFrame"`   // Use snake_case for consistency
	MenuType   string         `gorm:"column:menu_type" json:"menuType"` // Use snake_case for consistency
	Perms      string         `gorm:"column:perms" json:"perms"`
	CreateTime time.Time      `gorm:"column:created_time" json:"createTime"`
	AlwaysShow bool           `gorm:"column:always_show" json:"alwaysShow"`              // Use snake_case for consistency
	MenuStatus bool           `gorm:"column:menu_status" json:"menuStatus"`              // Use snake_case for consistency
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"` // Automatically set
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"` // Automatically set
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	APIURL     string         `gorm:"column:api_url" json:"apiURL"` // Use snake_case for consistency
}

func (m *Menu) BeforeCreate(tx *gorm.DB) (err error) {
	if m.CreateTime.IsZero() {
		m.CreateTime = time.Now()
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = time.Now()
	}
	return
}
func (m *Menu) BeforeUpdate(tx *gorm.DB) (err error) {
	if m.CreateTime.IsZero() {
		m.CreateTime = time.Now()
	}
	m.UpdatedAt = time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}
	return
}
func (m *Menu) TableName() string {
	return "system_menu"
}

func MigrateMenu(db *gorm.DB) {
	db.AutoMigrate(&Menu{})
}
