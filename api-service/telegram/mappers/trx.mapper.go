package mappers

import (
	"telegram-service/telegram/dtos"
	"telegram-service/telegram/models"
)

func TrxToResponseDTO(t *models.Transaction) *dtos.TransactionResponseDTO {
	if t == nil {
		return nil
	}
	return &dtos.TransactionResponseDTO{
		ID:              t.ID,
		GroupID:         t.GroupID, // Correctly converts uint to string
		TrxID:           t.TrxID,
		Amount:          t.Amount,
		Currency:        t.Currency,
		Sender:          t.Sender,
		Transport:       t.Transport,
		APV:             t.APV,
		TransactionDate: t.TransactionDate,
		TxtRaw:          t.TxtRaw,
		Bank:            t.Bank,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}

func ToModelCreateMapper(d *dtos.TransactionCreateDTO) *models.Transaction {
	if d == nil {
		return nil
	}

	return &models.Transaction{
		GroupID:         d.GroupID, // Convert uint64 to uint
		TrxID:           d.TrxID,
		Amount:          d.Amount,
		Currency:        d.Currency,
		Sender:          d.Sender,
		Transport:       d.Transport,
		APV:             d.APV,
		TransactionDate: d.TransactionDate,
		TxtRaw:          d.TxtRaw,
		Bank:            d.Bank,
	}
}

func UpdateModel(d *dtos.TransactionUpdateDTO, t *models.Transaction) {
	if d == nil || t == nil {
		return // Do nothing if DTO or model is nil
	}

	if d.GroupID != nil {
		t.GroupID = *d.GroupID

	}
	if d.TrxID != nil {
		t.TrxID = *d.TrxID
	}
	if d.Amount != nil {
		t.Amount = *d.Amount
	}
	if d.Currency != nil {
		t.Currency = *d.Currency
	}
	if d.Sender != nil {
		t.Sender = *d.Sender
	}
	if d.Transport != nil {
		t.Transport = *d.Transport
	}
	if d.TransactionDate != nil {
		t.TransactionDate = *d.TransactionDate
	}
}
