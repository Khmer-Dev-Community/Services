package models

import (
	"time"

	"gorm.io/gorm"
)

type ProxyList struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Address        string         `gorm:"type:varchar(250);column:proxy_address" json:"proxy_address"`
	Port           string         `gorm:"column:proxy_port" json:"proxy_port"`
	Session        string         `gorm:"column:proxy_seesion" json:"proxy_seesion"`
	Password       string         `gorm:"column:proxy_password" json:"proxy_password"`
	TimeOut        int            `gorm:"column:proxy_timeout" json:"proxy_timeout"`
	TragetLocation string         `gorm:"type:varchar(100);column:proxy_traget_location" json:"proxy_traget_location"`
	GroupId        int            `gorm:"column:proxy_group" json:"proxy_group"`
	Sort           int            `json:"sort" gorm:"column:sort"`
	CompanyID      int            `json:"companyId" gorm:"column:company_id"`
	Decription     string         `json:"description" gorm:"column:description"`
	Status         bool           `json:"status" gorm:"column:status"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type ProxySetting struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Port           int    `gorm:"column:proxy_port" json:"proxy_port"`
	ProxyType      string `gorm:"column:proxy_type" json:"proxy_type"`
	TimeOut        int    `gorm:"column:proxy_timeout" json:"proxy_timeout"`
	TragetLocation string `gorm:"type:varchar(100);column:proxy_traget_location" json:"proxy_traget_location"`
	TargetCity     string `gorm:"type:varchar(100);column:proxy_traget_city" json:"proxy_traget_city"`
	Decription     string `json:"description" gorm:"column:description"`
	Status         bool   `json:"status" gorm:"column:status"`
}

type ProxyTeamUser struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name" gorm:"column:name"`
}

func (ProxyList) TableName() string {
	return "proxy_list"
}

func MigrateProxy(db *gorm.DB) {
	db.AutoMigrate(&ProxyList{})
	db.AutoMigrate(&ProxySetting{})
	db.AutoMigrate(&ProxyTeamUser{})
}
