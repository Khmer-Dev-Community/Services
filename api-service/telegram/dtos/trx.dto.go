package dtos

import (
	"time"
)

type TransactionResponseDTO struct {
	ID              uint       `json:"id"`
	GroupID         int64      `json:"group_id"`
	TrxID           string     `json:"trx_id"`
	Amount          float64    `json:"amount"`
	Currency        string     `json:"currency"`
	Sender          string     `json:"sender"`
	APV             string     `json:"apv"`
	Transport       string     `json:"transport"`
	TransactionDate string     `json:"transaction_date"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	TxtRaw          string     `json:"txt_raw"`
	Bank            string     `json:"bank_bot"`
}

type TransactionCreateDTO struct {
	GroupID         int64   `json:"group_id" binding:"required"`
	TrxID           string  `json:"trx_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,gte=0"`
	Currency        string  `json:"currency" binding:"required"`
	Sender          string  `json:"sender" binding:"required"`
	Transport       string  `json:"transport" binding:"required"`        // This will receive result.Via
	APV             string  `json:"apv"`                                 // This will receive result.APV
	TransactionDate string  `json:"transaction_date" binding:"required"` // This will receive formatted date
	TxtRaw          string  `json:"txt_raw"`
	Bank            string  `json:"bank_bot"`
}

type TransactionUpdateDTO struct {
	GroupID         *int64   `json:"group_id"`
	TrxID           *string  `json:"trx_id"`
	Amount          *float64 `json:"amount"`
	Currency        *string  `json:"currency"`
	Sender          *string  `json:"sender"`
	Transport       *string  `json:"transport"`
	APV             *string  `json:"apv"`
	TransactionDate *string  `json:"transaction_date"`
	TxtRaw          string   `json:"txt_raw"`
	Bank            string   `json:"bank_bot"`
}

type TransactionFilter struct {
	GroupID         *int64     `json:"group_id"`
	TrxID           *string    `json:"trx_id"`
	CreatedAt       *time.Time `json:"created_at"`
	TransactionDate string     `json:"transaction_date"`
	StartDate       *time.Time `json:"from_date"`
	EndDate         *time.Time `json:"to_date"`
	Bank            string     `json:"bank_bot"`
}
