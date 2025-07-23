package services

import (
	"fmt"
	"log" // For logging if necessary inside service
	"strings"
	"telegram-service/telegram/dtos"
	"telegram-service/telegram/mappers"
	"telegram-service/telegram/repository" // Adjust this import path to your repository package
	"telegram-service/utils"
	"time" // For date calculations
)

// TransactionService defines the interface for transaction-related business logic.
type TransactionService interface {
	GetDailyReport(date string) (string, error)
	GetLast7DaysSummaryReport() (string, error)
	CreateTransaction(input *dtos.TransactionCreateDTO) (*dtos.TransactionResponseDTO, error)
	GetReports(input *dtos.TransactionFilter) (*ReportSummaryDTO, error)
	GetReportsByBankTransaction(input *dtos.TransactionFilter) (*ReportSummaryDTO, error)
	GetFilterReports(input *dtos.TransactionFilter) (*ReportSummaryDTO, error)
	GetFilterReportsByBankTransaction(input *dtos.TransactionFilter) (*ReportSummaryDTO, error)
	GetFilterReportsByBankTransport(input *dtos.TransactionFilter) (*ReportSummaryDTO, error)
}

// transactionService implements the TransactionService interface.
type transactionService struct {
	repo repository.TransactionRepository // Dependency on the repository layer
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) GetDailyReport(date string) (string, error) {
	queryDate := date
	if len(date) == 5 { // e.g., "02_01"
		currentYear := time.Now().Year()
		queryDate = fmt.Sprintf("%s_%d", date, currentYear) // Converts to "DD_MM_YYYY"
	} else if len(date) != 10 { // Expected "DD_MM_YYYY" or "DD_MM"
		return "", fmt.Errorf("invalid date format provided: %s. Expected 'DD_MM' or 'DD_MM_YYYY'", date)
	}
	transactions, err := s.repo.GetTransactionsByDate(queryDate)
	if err != nil {
		log.Printf("Error fetching daily transactions for %s: %v", queryDate, err)
		return "", fmt.Errorf("failed to retrieve daily report data")
	}
	if len(transactions) == 0 {
		return fmt.Sprintf("របាយការណ៍សម្រាប់ថ្ងៃ *%s*: មិនមានទិន្នន័យ", queryDate), nil
	}

	totalSales := 0.0
	totalTransactions := 0
	currencies := make(map[string]float64) // Sum sales by currency
	senders := make(map[string]int)        // Count transactions by sender
	transports := make(map[string]int)     // Count transactions by transport

	for _, trx := range transactions {
		totalSales += trx.Amount
		totalTransactions++
		currencies[trx.Currency] += trx.Amount
		senders[trx.Sender]++
		transports[trx.Transport]++
	}

	var report strings.Builder
	report.WriteString(fmt.Sprintf("របាយការណ៍សម្រាប់ថ្ងៃ *%s* :\n\n", queryDate))
	report.WriteString("\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\n") // Escaped hyphens for MarkdownV2
	report.WriteString(fmt.Sprintf("📊 **Total Sales:** $%.2f\n", totalSales))
	report.WriteString(fmt.Sprintf("📈 **Total Transactions:** %d\n", totalTransactions))
	// Detailed breakdown by currency
	if len(currencies) > 0 {
		report.WriteString("\n💰 **Sales by Currency:**\n")
		for curr, amt := range currencies {
			report.WriteString(fmt.Sprintf("    \\- %s: $%.2f\n", escapeMarkdownV2(curr), amt))
		}
	}
	// Detailed breakdown by sender
	if len(senders) > 0 {
		report.WriteString("\n👤 **Transactions by Sender:**\n")
		for sender, count := range senders {
			report.WriteString(fmt.Sprintf("    \\- %s: %d\n", escapeMarkdownV2(sender), count))
		}
	}

	// Detailed breakdown by transport
	if len(transports) > 0 {
		report.WriteString("\n🚚 **Transactions by Transport:**\n")
		for transport, count := range transports {
			report.WriteString(fmt.Sprintf("    \\- %s: %d\n", escapeMarkdownV2(transport), count))
		}
	}

	return report.String(), nil
}

// GetLast7DaysSummaryReport fetches and summarizes transaction data for the last 7 days.
func (s *transactionService) GetLast7DaysSummaryReport() (string, error) {
	now := time.Now()
	// endDate is today's date in "DD_MM_YYYY" format
	endDate := now.Format("02_01_2006")
	// startDate is 6 days ago (to include today, making it 7 days total) in "DD_MM_YYYY" format
	startDate := now.AddDate(0, 0, -6).Format("02_01_2006")

	transactions, err := s.repo.GetTransactionsByDateRange(startDate, endDate)
	if err != nil {
		log.Printf("Error fetching transactions for last 7 days (%s to %s): %v", startDate, endDate, err)
		return "", fmt.Errorf("failed to retrieve last 7 days report data")
	}

	if len(transactions) == 0 {
		return "របាយការណ៍ ៧ថ្ងៃចុងក្រោយ: មិនមានទិន្នន័យ", nil
	}

	totalSales := 0.0
	totalTransactions := 0
	currencies := make(map[string]float64)
	senders := make(map[string]int)
	transports := make(map[string]int)
	transactionsPerDay := make(map[string]int) // To count transactions per day

	for _, trx := range transactions {
		totalSales += trx.Amount
		totalTransactions++
		currencies[trx.Currency] += trx.Amount
		senders[trx.Sender]++
		transports[trx.Transport]++

		// Count transactions per day based on TransactionDate
		transactionsPerDay[trx.TransactionDate]++
	}
	var report strings.Builder
	report.WriteString(fmt.Sprintf("របាយការណ៍ *៧ថ្ងៃចុងក្រោយ* (%s \u2013 %s):\n\n", startDate, endDate)) // Using en dash for range
	report.WriteString("\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\-\\n")

	report.WriteString(fmt.Sprintf("📊 **Total Sales:** $%.2f\n", totalSales))
	report.WriteString(fmt.Sprintf("📈 **Total Transactions:** %d\n", totalTransactions))

	// Daily breakdown (iterate through the last 7 days to ensure all days are shown, even with 0 transactions)
	report.WriteString("\n🗓️ **Daily Transaction Count:**\n")
	for i := 6; i >= 0; i-- { // Iterate from startDate (6 days ago) up to today
		currentDate := now.AddDate(0, 0, -i).Format("02_01_2006")
		count, ok := transactionsPerDay[currentDate]
		if !ok {
			count = 0 // If no transactions for this day, count is 0
		}
		report.WriteString(fmt.Sprintf("    \\- %s: %d\n", currentDate, count))
	}

	// Detailed breakdown by currency
	if len(currencies) > 0 {
		report.WriteString("\n💰 **Sales by Currency:**\n")
		for curr, amt := range currencies {
			report.WriteString(fmt.Sprintf("    \\- %s: $%.2f\n", escapeMarkdownV2(curr), amt))
		}
	}

	// Detailed breakdown by sender
	if len(senders) > 0 {
		report.WriteString("\n👤 **Transactions by Sender:**\n")
		for sender, count := range senders {
			report.WriteString(fmt.Sprintf("    \\- %s: %d\n", escapeMarkdownV2(sender), count))
		}
	}

	// Detailed breakdown by transport
	if len(transports) > 0 {
		report.WriteString("\n🚚 **Transactions by Transport:**\n")
		for transport, count := range transports {
			report.WriteString(fmt.Sprintf("    \\- %s: %d\n", escapeMarkdownV2(transport), count))
		}
	}

	return report.String(), nil
}
func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(", "\\(", ")", "\\)",
		"~", "\\~", "`", "\\`", ">", "\\>", "#", "\\#", "+", "\\+",
		"-", "\\-", "=", "\\=", "|", "\\|", "{", "\\{", "}", "\\}",
		".", "\\.", "!", "\\!",
	)
	return replacer.Replace(text)
}

func (s *transactionService) CreateTransaction(input *dtos.TransactionCreateDTO) (*dtos.TransactionResponseDTO, error) {
	createdTable, err := s.repo.SaveTrx(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	responseDTO := mappers.TrxToResponseDTO(createdTable)
	utils.InfoLog(responseDTO, "CreateTable: Table created successfully via service")
	return responseDTO, nil
}

// ReportSummaryDTO represents the comprehensive summary report data.
type ReportSummaryDTO struct {
	Transport     string  `json:"transport"`
	Banker        string  `json:"bank_bot"`
	TotalUSD      float64 `json:"totalUSD"`
	TotalTrxUSD   int     `json:"totalTrxUSD"`
	TotalReil     float64 `json:"totalReil"`
	TotalTrxReil  int     `json:"totalTrxReil"`
	BankBreakdown map[string]map[string]struct {
		Total float64
		Count int
	} `json:"bankBreakdown"` // Added JSON tag for consistency
}

func (s *transactionService) GetReports(input *dtos.TransactionFilter) (*ReportSummaryDTO, error) {
	utils.InfoLog(input, "input *dtos.TransactionFilter")
	if input.CreatedAt == nil {
		return nil, fmt.Errorf("report date (CreatedAt) is required")
	}
	totalUSD, countUSD, totalReil, countReil, err := s.repo.SumTransactionsByDateAndGroupByCurrency(*input) // Dereference *input.CreatedAt if it's a pointer
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction summary from repository: %w", err)
	}
	reportSummary := &ReportSummaryDTO{
		TotalUSD:     totalUSD,
		TotalTrxUSD:  countUSD,
		TotalReil:    totalReil,
		TotalTrxReil: countReil,
	}
	return reportSummary, nil
}

func (s *transactionService) GetReportsByBankTransaction(input *dtos.TransactionFilter) (*ReportSummaryDTO, error) {
	//detailedSummary, err := s.repo.GetDetailedTransactionSummaryByDate(*input)
	detailedSummary, err := s.repo.GetDetailedBankTransactionSummaryByDate(*input)
	if err != nil {
		return nil, fmt.Errorf("failed to get detailed transaction summary from repository: %w", err)
	}

	reportDTO := &ReportSummaryDTO{
		BankBreakdown: make(map[string]map[string]struct {
			Total float64
			Count int
		}),
	}

	for _, sum := range detailedSummary { // Renamed 's' to 'sum' to avoid conflict with 's' for service
		// Accumulate overall totals
		switch sum.Currency {
		case "$":
			reportDTO.TotalUSD += sum.Total
			reportDTO.TotalTrxUSD += sum.Count
		case "៛":
			reportDTO.TotalReil += sum.Total
			reportDTO.TotalTrxReil += sum.Count
		}

		// Populate the BankBreakdown map
		if _, ok := reportDTO.BankBreakdown[sum.Bank]; !ok {
			reportDTO.BankBreakdown[sum.Bank] = make(map[string]struct {
				Total float64
				Count int
			})
		}
		reportDTO.BankBreakdown[sum.Bank][sum.Currency] = struct {
			Total float64
			Count int
		}{Total: sum.Total, Count: sum.Count}
	}

	return reportDTO, nil
}

func (s *transactionService) GetReportsByBankTransport(input *dtos.TransactionFilter) (*ReportSummaryDTO, error) {
	detailedSummary, err := s.repo.GetDetailedTransactionSummaryByDate(*input)
	if err != nil {
		return nil, fmt.Errorf("failed to get detailed transaction summary from repository: %w", err)
	}

	reportDTO := &ReportSummaryDTO{
		BankBreakdown: make(map[string]map[string]struct {
			Total float64
			Count int
		}),
	}

	for _, sum := range detailedSummary { // Renamed 's' to 'sum' to avoid conflict with 's' for service
		// Accumulate overall totals
		switch sum.Currency {
		case "$":
			reportDTO.TotalUSD += sum.Total
			reportDTO.TotalTrxUSD += sum.Count
		case "៛":
			reportDTO.TotalReil += sum.Total
			reportDTO.TotalTrxReil += sum.Count
		}

		// Populate the BankBreakdown map
		if _, ok := reportDTO.BankBreakdown[sum.Transport]; !ok {
			reportDTO.BankBreakdown[sum.Transport] = make(map[string]struct {
				Total float64
				Count int
			})
		}
		reportDTO.BankBreakdown[sum.Transport][sum.Currency] = struct {
			Total float64
			Count int
		}{Total: sum.Total, Count: sum.Count}
	}

	return reportDTO, nil
}
func (s *transactionService) GetFilterReports(input *dtos.TransactionFilter) (*ReportSummaryDTO, error) {
	if input.CreatedAt == nil {
		return nil, fmt.Errorf("report date (CreatedAt) is required")
	}
	totalUSD, countUSD, totalReil, countReil, err := s.repo.SumTransactionsByDateAndGroupByCurrencyFilterDate(*input) // Dereference *input.CreatedAt if it's a pointer
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction summary from repository: %w", err)
	}
	reportSummary := &ReportSummaryDTO{
		TotalUSD:     totalUSD,
		TotalTrxUSD:  countUSD,
		TotalReil:    totalReil,
		TotalTrxReil: countReil,
	}
	return reportSummary, nil
}

func (s *transactionService) GetFilterReportsByBankTransaction(input *dtos.TransactionFilter) (*ReportSummaryDTO, error) {
	detailedSummary, err := s.repo.GetDetailedBankTransactionSummaryByDateFilter(*input)
	if err != nil {
		return nil, fmt.Errorf("failed to get detailed transaction summary from repository: %w", err)
	}

	reportDTO := &ReportSummaryDTO{
		BankBreakdown: make(map[string]map[string]struct {
			Total float64
			Count int
		}),
	}

	for _, sum := range detailedSummary { // Renamed 's' to 'sum' to avoid conflict with 's' for service
		// Accumulate overall totals
		switch sum.Currency {
		case "$":
			reportDTO.TotalUSD += sum.Total
			reportDTO.TotalTrxUSD += sum.Count
		case "៛":
			reportDTO.TotalReil += sum.Total
			reportDTO.TotalTrxReil += sum.Count
		}

		// Populate the BankBreakdown map
		if _, ok := reportDTO.BankBreakdown[sum.Bank]; !ok {
			reportDTO.BankBreakdown[sum.Bank] = make(map[string]struct {
				Total float64
				Count int
			})
		}
		reportDTO.BankBreakdown[sum.Bank][sum.Currency] = struct {
			Total float64
			Count int
		}{Total: sum.Total, Count: sum.Count}
	}

	return reportDTO, nil
}

func (s *transactionService) GetFilterReportsByBankTransport(input *dtos.TransactionFilter) (*ReportSummaryDTO, error) {
	detailedSummary, err := s.repo.GetDetailedTransactionSummaryByDate(*input)
	if err != nil {
		return nil, fmt.Errorf("failed to get detailed transaction summary from repository: %w", err)
	}

	reportDTO := &ReportSummaryDTO{
		BankBreakdown: make(map[string]map[string]struct {
			Total float64
			Count int
		}),
	}

	for _, sum := range detailedSummary { // Renamed 's' to 'sum' to avoid conflict with 's' for service
		// Accumulate overall totals
		switch sum.Currency {
		case "$":
			reportDTO.TotalUSD += sum.Total
			reportDTO.TotalTrxUSD += sum.Count
		case "៛":
			reportDTO.TotalReil += sum.Total
			reportDTO.TotalTrxReil += sum.Count
		}

		// Populate the BankBreakdown map
		if _, ok := reportDTO.BankBreakdown[sum.Transport]; !ok {
			reportDTO.BankBreakdown[sum.Transport] = make(map[string]struct {
				Total float64
				Count int
			})
		}
		reportDTO.BankBreakdown[sum.Transport][sum.Currency] = struct {
			Total float64
			Count int
		}{Total: sum.Total, Count: sum.Count}
	}

	return reportDTO, nil
}
