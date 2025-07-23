package roles

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	RoleName   string         `gorm:"type:varchar(100);column:role_name" json:"roleName"` // Use snake_case for the column name
	RoleStatus int            `gorm:"column:role_status" json:"roleStatus"`               // Use snake_case for the column name
	RoleKey    string         `gorm:"type:varchar(100);column:role_key" json:"roleKey"`   // Use snake_case for the column name
	Sort       int            `json:"sort" gorm:"column:sort"`
	CompanyID  int            `json:"companyId" gorm:"column:company_id"`
	Decription string         `json:"description" gorm:"column:description"`
	Status     bool           `json:"status" gorm:"column:status"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"` // Automatically set creation time
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"` // Automatically set update time
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Role) TableName() string {
	return "system_role"
}

// MigrateUsers automates the user table migration
func MigrateRoles(db *gorm.DB) {
	db.AutoMigrate(&Role{})
}
