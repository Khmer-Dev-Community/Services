package repository

import (
	"fmt"
	"log"
	"telegram-service/telegram/dtos"
	"telegram-service/telegram/mappers"
	"telegram-service/telegram/models" // Adjust this import path if your models are elsewhere
	"telegram-service/utils"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	GetByID(id uint) (*models.Transaction, error)
	Update(transaction *models.Transaction) error
	Delete(id uint) error
	List(offset, limit int) ([]models.Transaction, error)
	Count() (int64, error) // To get total count for pagination

	// Specific query methods for your bot's reporting needs
	GetTransactionsByDate(date string) ([]models.Transaction, error)
	GetTransactionsByDateRange(startDate, endDate string) ([]models.Transaction, error)
	SaveTrx(createDTO *dtos.TransactionCreateDTO) (*models.Transaction, error)
	SumTransactionsByDateAndGroupByCurrency(date dtos.TransactionFilter) (float64, int, float64, int, error)
	GetDetailedTransactionSummaryByDate(filter dtos.TransactionFilter) ([]CurrencySummary, error)
	SumTransactionsByDateAndGroupByCurrencyFilterDate(date dtos.TransactionFilter) (float64, int, float64, int, error)
	SumTransactionsByDateAndGroupByCurrencyFilter(date dtos.TransactionFilter) (float64, int, float64, int, error)
	GetDetailedBankTransactionSummaryByDate(filter dtos.TransactionFilter) ([]CurrencySummary, error)
	GetDetailedBankTransactionSummaryByDateFilter(filter dtos.TransactionFilter) ([]CurrencySummary, error)
}
type CurrencySummary struct {
	Bank      string `json:"bank_bot"`
	Transport string `json:"transport"`
	Currency  string
	Total     float64
	Count     int
}
type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}
func (r *transactionRepository) SaveTrx(createDTO *dtos.TransactionCreateDTO) (*models.Transaction, error) {
	trx := mappers.ToModelCreateMapper(createDTO)
	if err := r.db.Create(&trx).Error; err != nil {
		utils.LoggerRepository(err, "Trx: Insert failed")
		return nil, err
	}
	return trx, nil
}
func (r *transactionRepository) GetByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	result := r.db.First(&transaction, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // Return nil model and nil error if not found
		}
		return nil, result.Error // Return other database errors
	}
	return &transaction, nil
}

func (r *transactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

// Delete removes a transaction record from the database by its ID.
func (r *transactionRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Transaction{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no transaction found with ID %d to delete", id)
	}
	return nil
}

func (r *transactionRepository) List(offset, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := r.db.Offset(offset).Limit(limit).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

func (r *transactionRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Transaction{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *transactionRepository) GetTransactionsByDate(date string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := r.db.Where("transaction_date = ?", date).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

func (r *transactionRepository) GetTransactionsByDateRange(startDate, endDate string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	result := r.db.Where("transaction_date BETWEEN ? AND ?", startDate, endDate).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactions, nil
}

func (r *transactionRepository) SumTransactionsByDateAndGroupByCurrency(date dtos.TransactionFilter) (float64, int, float64, int, error) {
	var summaries []CurrencySummary
	var totalUSD float64
	var countUSD int // New variable for USD transaction count
	var totalReil float64
	var countReil int // New variable for Riel transaction count
	inputTimeLocal := *date.CreatedAt
	loc, err := time.LoadLocation("Asia/Phnom_Penh") // Or whatever your server's local timezone is
	if err != nil {
		log.Printf("Error loading location: %v, falling back to local.", err)
		loc = time.Local // Fallback to current system's local timezone
	}
	targetDateInLocal := time.Date(inputTimeLocal.Year(), inputTimeLocal.Month(), inputTimeLocal.Day(), 0, 0, 0, 0, loc)
	startOfDayUTC := targetDateInLocal.UTC()
	endOfDayUTC := targetDateInLocal.Add(24 * time.Hour).Add(-time.Nanosecond).UTC()
	result := r.db.Model(&models.Transaction{}).
		Select("currency, SUM(amount) as total, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startOfDayUTC, endOfDayUTC).
		Where("group_id = ?", date.GroupID).
		Group("currency").
		Order("currency ASC").
		Scan(&summaries)

	if result.Error != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to sum and count transactions by created_at date range: %w", result.Error)
	}

	// Iterate through the summaries to assign totals to USD and Riel
	for _, s := range summaries {
		switch s.Currency {
		case "$":
			totalUSD = s.Total
			countUSD = s.Count // Assign USD count
		case "៛":
			totalReil = s.Total
			countReil = s.Count // Assign Riel count
		}
	}

	return totalUSD, countUSD, totalReil, countReil, nil
}

func (r *transactionRepository) GetDetailedTransactionSummaryByDate(filter dtos.TransactionFilter) ([]CurrencySummary, error) {
	var summaries []CurrencySummary

	startOfDay := filter.CreatedAt.Truncate(24 * time.Hour)
	endOfDay := filter.CreatedAt.Add(24 * time.Hour).Add(-time.Nanosecond)

	result := r.db.Model(&models.Transaction{}).
		Select("transport, currency, SUM(amount) as total, COUNT(*) as count"). // Selects your desired aggregates
		Where("created_at >= ? AND created_at <= ?", startOfDay, endOfDay).
		Where("group_id = ?", filter.GroupID).
		Group("transport, currency").         // THIS IS THE KEY: Groups by both transport (via) and currency
		Order("transport ASC, currency ASC"). // Optional: order the results for consistent output
		Scan(&summaries)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get detailed transaction summary by date and transport: %w", result.Error)
	}

	return summaries, nil
}

func (r *transactionRepository) SumTransactionsByDateAndGroupByCurrencyFilterDate(date dtos.TransactionFilter) (float64, int, float64, int, error) {
	var summaries []CurrencySummary
	var totalUSD float64
	var countUSD int // New variable for USD transaction count
	var totalReil float64
	var countReil int // New variable for Riel transaction count

	startOfDay := date.StartDate
	endOfDay := date.EndDate

	result := r.db.Model(&models.Transaction{}).
		Select("currency, SUM(amount) as total, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startOfDay, endOfDay).
		Where("group_id = ?", date.GroupID).
		Group("currency").
		Order("currency ASC").
		Scan(&summaries)

	if result.Error != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to sum and count transactions by created_at date range: %w", result.Error)
	}

	// Iterate through the summaries to assign totals to USD and Riel
	for _, s := range summaries {
		switch s.Currency {
		case "$":
			totalUSD = s.Total
			countUSD = s.Count // Assign USD count
		case "៛":
			totalReil = s.Total
			countReil = s.Count // Assign Riel count
		}
	}

	return totalUSD, countUSD, totalReil, countReil, nil
}

func (r *transactionRepository) SumTransactionsByDateAndGroupByCurrencyFilter(date dtos.TransactionFilter) (float64, int, float64, int, error) {
	var summaries []CurrencySummary
	var totalUSD float64
	var countUSD int // New variable for USD transaction count
	var totalReil float64
	var countReil int // New variable for Riel transaction count

	startOfDay := date.StartDate
	endOfDay := date.EndDate

	result := r.db.Model(&models.Transaction{}).
		Select("currency, SUM(amount) as total, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startOfDay, endOfDay).
		Where("group_id = ?", date.GroupID).
		Group("currency").
		Order("currency ASC").
		Scan(&summaries)

	if result.Error != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to sum and count transactions by created_at date range: %w", result.Error)
	}

	// Iterate through the summaries to assign totals to USD and Riel
	for _, s := range summaries {
		switch s.Currency {
		case "$":
			totalUSD = s.Total
			countUSD = s.Count // Assign USD count
		case "៛":
			totalReil = s.Total
			countReil = s.Count // Assign Riel count
		}
	}

	return totalUSD, countUSD, totalReil, countReil, nil
}

func (r *transactionRepository) SumBankByDateAndGroupByCurrencyFilterDate(date dtos.TransactionFilter) (float64, int, float64, int, error) {
	var summaries []CurrencySummary
	var totalUSD float64
	var countUSD int // New variable for USD transaction count
	var totalReil float64
	var countReil int // New variable for Riel transaction count

	startOfDay := date.StartDate
	endOfDay := date.EndDate

	result := r.db.Model(&models.Transaction{}).
		Select("currency, SUM(amount) as total, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", startOfDay, endOfDay).
		Where("group_id = ?", date.GroupID).
		Group("currency").
		Order("currency ASC").
		Scan(&summaries)

	if result.Error != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to sum and count transactions by created_at date range: %w", result.Error)
	}

	// Iterate through the summaries to assign totals to USD and Riel
	for _, s := range summaries {
		switch s.Currency {
		case "$":
			totalUSD = s.Total
			countUSD = s.Count // Assign USD count
		case "៛":
			totalReil = s.Total
			countReil = s.Count // Assign Riel count
		}
	}

	return totalUSD, countUSD, totalReil, countReil, nil
}

func (r *transactionRepository) GetDetailedBankTransactionSummaryByDate(filter dtos.TransactionFilter) ([]CurrencySummary, error) {
	var summaries []CurrencySummary

	//startOfDay := filter.CreatedAt.Truncate(24 * time.Hour)
	//endOfDay := filter.CreatedAt.Add(24 * time.Hour).Add(-time.Nanosecond)
	inputTimeLocal := filter.CreatedAt
	loc, err := time.LoadLocation("Asia/Phnom_Penh") // Or whatever your server's local timezone is
	if err != nil {
		log.Printf("Error loading location: %v, falling back to local.", err)
		loc = time.Local // Fallback to current system's local timezone
	}
	targetDateInLocal := time.Date(inputTimeLocal.Year(), inputTimeLocal.Month(), inputTimeLocal.Day(), 0, 0, 0, 0, loc)
	startOfDayUTC := targetDateInLocal.UTC()
	endOfDayUTC := targetDateInLocal.Add(24 * time.Hour).Add(-time.Nanosecond).UTC()

	result := r.db.Model(&models.Transaction{}).
		Select("bank, currency, SUM(amount) as total, COUNT(*) as count"). // Selects your desired aggregates
		Where("created_at >= ? AND created_at <= ?", startOfDayUTC, endOfDayUTC).
		Where("group_id = ?", filter.GroupID).
		Group("bank, currency").         // THIS IS THE KEY: Groups by both transport (via) and currency
		Order("bank ASC, currency ASC"). // Optional: order the results for consistent output
		Scan(&summaries)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get detailed transaction summary by date and transport: %w", result.Error)
	}

	return summaries, nil
}

func (r *transactionRepository) GetDetailedBankTransactionSummaryByDateFilter(filter dtos.TransactionFilter) ([]CurrencySummary, error) {
	var summaries []CurrencySummary

	startOfDay := filter.StartDate
	endOfDay := filter.EndDate

	result := r.db.Model(&models.Transaction{}).
		Select("bank, currency, SUM(amount) as total, COUNT(*) as count"). // Selects your desired aggregates
		Where("created_at >= ? AND created_at <= ?", startOfDay, endOfDay).
		Where("group_id = ?", filter.GroupID).
		Group("bank, currency").         // THIS IS THE KEY: Groups by both transport (via) and currency
		Order("bank ASC, currency ASC"). // Optional: order the results for consistent output
		Scan(&summaries)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get detailed transaction summary by date and transport: %w", result.Error)
	}

	return summaries, nil
}
