package permission

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID                int       `gorm:"primaryKey" json:"id"`
	RoleID            int       `gorm:"column:role_id" json:"roleId"`
	MenuID            int       `gorm:"column:menu_id" json:"menuId"` // changed to lowercase
	MenuCheckStrictly bool      `gorm:"column:menu_check_strictly" json:"menuCheckStrictly"`
	RoleKey           string    `gorm:"column:role_key" json:"roleKey"`
	RoleName          string    `gorm:"column:role_name" json:"roleName"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"` // No pointer, GORM handles it automatically
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"` // No pointer, GORM handles it automatically
	//DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Permission) TableName() string {
	return "system_menu_role"
}
func (m *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	// Set CreatedAt to current time if it is zero (not needed if autoCreateTime is used)
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now() // Not necessary if you are using autoCreateTime
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = time.Now() // Not necessary if you are using autoUpdateTime
	}
	return
}

func (m *Permission) BeforeUpdate(tx *gorm.DB) (err error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now() // Not necessary if you are using autoCreateTime
	}
	// Set UpdatedAt to current time
	m.UpdatedAt = time.Now() // Not necessary if you are using autoUpdateTime
	return
}

func MigrateRoleAccesses(db *gorm.DB) {
	db.AutoMigrate(&Permission{})
}
