package department

import (
	"time"

	"gorm.io/gorm"
)

// User represents the model for a user
type Department struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Code      string         `json:"dept_code"`
	Name      string         `json:"name"`
	CompanyID int            `json:"companyId"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Department) TableName() string {
	return "Departments"
}

// MigrateUsers automates the user table migration
func MigrateDepartment(db *gorm.DB) {
	db.AutoMigrate(&Department{})
}
