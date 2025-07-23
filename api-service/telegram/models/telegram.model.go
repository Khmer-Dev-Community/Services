package models

import (
	"time"

	"gorm.io/gorm"
)

type TelegramMessage struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	ChatID      string    `gorm:"not null"`
	Text        string    `gorm:"not null"`
	MessageType string    `gorm:"not null"`
	SentAt      time.Time `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type TelegramAccount struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountName string     `gorm:"type:varchar(50);not null" json:"account_name"`
	AccountID   string     `gorm:"type:varchar(250);not null" json:"account_id"`
	PhoneNumber string     `gorm:"type:varchar(20);not null" json:"phone_number"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
type Transaction struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID         int64      `gorm:"type:bigint;not null" json:"group_id"`
	TrxID           string     `gorm:"type:varchar(50);not null" json:"trx_id"`
	Amount          float64    `gorm:"type:decimal(12,2);default:0" json:"amount"`
	Currency        string     `gorm:"type:varchar(30);not null" json:"currency"`
	Sender          string     `gorm:"type:varchar(150);not null" json:"sender"`
	Transport       string     `gorm:"type:varchar(150);not null" json:"transport"`
	APV             string     `gorm:"type:varchar(150);not null" json:"apv"`
	TransactionDate string     `gorm:"type:varchar(50);not null" json:"transaction_date"`
	TxtRaw          string     `gorm:"type:varchar(250);not null" json:"txt_raw"`
	Bank            string     `gorm:"type:varchar(50);" json:"bank_bot"`
	CreatedAt       *time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
type TelegramGroup struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID   int64      `gorm:"type:bigint;not null" json:"group_id"`
	GroupName string     `gorm:"type:varchar(250);not null" json:"group_name"`
	GameType  string     `gorm:"type:varchar(250);not null" json:"game_type"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (TelegramAccount) TableName() string {
	return "telegram_account"
}
func (TelegramMessage) TableName() string {
	return "telegram_message"
}
func (Transaction) TableName() string {
	return "telegram_transaction"
}
func (TelegramGroup) TableName() string {
	return "telegram_group"
}

// Migration table
func MigrateBotTable(db *gorm.DB) {
	//db.AutoMigrate(&TelegramMessage{})
	//db.AutoMigrate(&TelegramAccount{})
	db.AutoMigrate(&Transaction{})
	db.AutoMigrate(&TelegramGroup{})
}
