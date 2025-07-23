package employee

import (
	"time"

	"gorm.io/gorm"
)

// User represents the model for a user
type Employee struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	EmployeeID   string         `json:"empCard" gorm:"column:empCard"`
	EmployeeCard string         `json:"ci_card" gorm:"column:ci_card"`
	FirstName    string         `json:"fname" gorm:"column:fname"`
	LastName     string         `json:"lname" gorm:"column:lname"`
	Username     string         `json:"username" gorm:"column:username"`
	Sex          string         `json:"gender" gorm:"column:gender"`
	AccountName  string         `json:"account_name" gorm:"column:account_name"`
	Password     string         `json:"password" gorm:"column:password"`
	Position     string         `json:"position" gorm:"column:position"`
	RoleID       uint           `json:"roleId" gorm:"column:roleId"`
	CompanyID    int            `json:"companyId" gorm:"column:companyId"`
	DepartmentID int            `json:"departmentId" gorm:"column:departmentId"`
	Level        string         `json:"cardType" gorm:"column:cardType"`
	JoinDate     string         `json:"joinDate" gorm:"column:joinDate"`
	DOB          string         `json:"dob" gorm:"column:dob"`
	Language     string         `json:"default_lag" gorm:"column:default_lag"`
	Phone        string         `json:"phone" gorm:"column:phone"`
	Email        string         `json:"email" gorm:"column:email"`
	Balanace     string         `json:"balance" gorm:"column:balance"`
	Note         string         `json:"note" gorm:"column:note"`
	Status       string         `json:"status" gorm:"column:status"`
	Token        string         `json:"token" gorm:"column:token"`
	CreatedAt    time.Time      `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"column:updatedAt" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Employee) TableName() string {
	return "Employees"
}

// BeforeSave hook to handle setting RoleID
func (u *Employee) BeforeSave(tx *gorm.DB) (err error) {
	// If RoleID is 0 (not set), set it to 0
	if u.RoleID == 0 {
		u.RoleID = 0
	}
	return nil
}

// MigrateUsers automates the user table migration
func MigrateEmployee(db *gorm.DB) {
	db.AutoMigrate(&Employee{})
}
