// internal/models/websocket_payloads.go
package gateway

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	// Register any custom validators here if needed
}

// WSJoinOrderChannelPayload represents the data for "ws:join_order_channel" event
type WSJoinOrderChannelPayload struct {
	Channel string `json:"channel" validate:"required,min=1"`
}

// TableSelectedPayload represents the data for "tableSelected" event
type TableSelectedPayload struct {
	TableID   float64 `json:"tableId" validate:"required,gt=0"`
	NewStatus string  `json:"newStatus" validate:"required,oneof=pending selected available"` // Example: status must be one of these
}

// BillSelectedPayload represents the data for "billSelected" event
type BillSelectedPayload struct {
	TableID   float64 `json:"tableId" validate:"required,gt=0"`
	InvoiceID float64 `json:"invoiceId" validate:"required,gt=0"`
	InvoiceNo float64 `json:"invoiceNo" validate:"required,gt=0"` // Note: You had InvoiceNo as float64 here, but string in orderSaved. Ensure consistency.
}

// ClientJoinBillPayload represents the data for "client:joinBill" event
type ClientJoinBillPayload struct {
	InvoiceNo string `json:"invoiceNo" validate:"required,min=1"`
}

// OrderSavedPayload represents the data for "orderSaved" event
type OrderSavedPayload struct {
	TableID      float64       `json:"tableId" validate:"required,gt=0"`
	InvoiceID    float64       `json:"invoiceId" validate:"required,gt=0"`
	InvoiceNo    string        `json:"invoiceNo" validate:"required,min=1"`         // Assuming InvoiceNo is a string
	UpdatedItems []interface{} `json:"updatedItems" validate:"required,min=1,dive"` // `dive` means validate elements if they were structs too
	ClientID     string        `json:"clientId" validate:"required,min=1"`          // Or `uuid` if it's a UUID
}

// OrderProcessedPayload represents the data for "orderProcessed" event
type OrderProcessedPayload struct {
	TableID   float64 `json:"tableId" validate:"required,gt=0"`
	InvoiceID float64 `json:"invoiceId" validate:"required,gt=0"`
	InvoiceNo string  `json:"invoiceNo" validate:"required,min=1"`
	Status    string  `json:"status" validate:"required,oneof=available"` // Should be 'available'
	ClientID  string  `json:"clientId" validate:"required,min=1"`
}
