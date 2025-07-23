package bot

import (
	"fmt"
	"log"
	"strings"
	"telegram-service/telegram/dtos"
	"telegram-service/telegram/services"
	"telegram-service/utils"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleCallbackQuery_(callbackQuery *tgbotapi.CallbackQuery) {
	var reportContent string
	var reportDate time.Time
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := b.botAPI.Request(callback); err != nil { // Use b.botAPI
		log.Printf("Failed to answer callback query: %v", err)
	}

	if callbackQuery.Message == nil {
		log.Printf("Callback query message is nil for query ID: %s", callbackQuery.ID)
		return
	}

	chatID := callbackQuery.Message.Chat.ID
	messageIDToDelete := callbackQuery.Message.MessageID
	callbackData := callbackQuery.Data
	username := callbackQuery.From.FirstName
	if callbackQuery.From.LastName != "" {
		username += " " + callbackQuery.From.LastName
	}

	escapedUsername := EscapeMarkdownV2(username)
	utils.InfoLog(callbackData, fmt.Sprintf("DEBUG: Received callback from %s with data: %s", escapedUsername, callbackData))

	deleteConfig := tgbotapi.NewDeleteMessage(chatID, messageIDToDelete)
	if _, err := b.botAPI.Request(deleteConfig); err != nil { // Use b.botAPI
		log.Printf("ERROR: Failed to delete button message (ID: %d) in chat %d after callback from %s: %v",
			messageIDToDelete, chatID, escapedUsername, err)
	} else {
		log.Printf("DEBUG: Successfully deleted button message (ID: %d) in chat %d.", messageIDToDelete, chatID)
	}
	if strings.HasPrefix(callbackData, "report_date_") {
		parts := strings.SplitN(callbackData, "_", 2)
		var reportCommandType string
		var dateValue string
		if len(parts) == 2 {
			reportCommandType = parts[0] // Assigns "report_date"
			dateValue = parts[1]         // Assigns "11_07_2025"
		} else {
			// This case handles malformed callback data (e.g., just "report_date")
			log.Printf("Invalid report_date command format: %s", callbackData)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ទ្រង់ទ្រាយបញ្ជាមិនត្រឹមត្រូវ។")
			if _, err := b.botAPI.Send(errMsg); err != nil {
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		formattedDateForParse := strings.ReplaceAll(dateValue, "_", "-")
		parsedTime, err := time.Parse("02-01-2006", formattedDateForParse)
		if err != nil {
			log.Printf("Error parsing date from callback data '%s': %v", formattedDateForParse, err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ទ្រង់ទ្រាយកាលបរិច្ឆេទមិនត្រឹមត្រូវ។")
			if _, err := b.botAPI.Send(errMsg); err != nil {
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		callbackData = reportCommandType
		reportDate = parsedTime
	}
	utils.InfoLog(callbackData, fmt.Sprintf("DEBUG: Received callback from %s with data: %s", callbackData, reportDate))
	switch callbackData {
	case "report_today":
		now := time.Now()
		reportDate = now
		transactionDTO := &dtos.TransactionFilter{
			CreatedAt: &now,
			GroupID:   &chatID, // Changed from &chatID to chatID if GroupID is int64
		}
		reportSummary, err := b.TransactionService.GetReports(transactionDTO)
		if err != nil {
			log.Printf("Error getting report for today: %v", err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ខ្ញុំមិនអាចបង្កើតរបាយការណ៍សម្រាប់ថ្ងៃនេះបានទេ។ សូមព្យាយាមម្តងទៀតនៅពេលក្រោយ")
			if _, err := b.botAPI.Send(errMsg); err != nil { // Use b.botAPI
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		reportContent = fmt.Sprintf(
			"📊 របាយការណ៍សរុបប្រតិបត្តិការណ៍\n"+
				"-----------------------------------\n"+
				"🕑 `ថ្ងៃនេះ %s `\n"+
				"$ (USD)សរុប : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n"+
				"៛ (KHR)សរុប : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n"+
				"-----------------------------------\n"+ // Added newline and made this a separate line for clarity
				"ប្រតិបត្តិការណ៍សង្ខេប តាមធនាគារ\n\n",
			reportDate.Format("2006/01/02"),
			reportSummary.TotalUSD,
			reportSummary.TotalTrxUSD,
			reportSummary.TotalReil,
			reportSummary.TotalTrxReil,
		)
		reportContent = EscapeMarkdownV2(reportContent) // Apply EscapeMarkdownV2
	case "report_yesterday":
		now := time.Now()
		todayStart := now.Truncate(24 * time.Hour)
		yesterdayStart := todayStart.Add(-24 * time.Hour)
		reportDate = yesterdayStart
		transactionDTO := &dtos.TransactionFilter{
			CreatedAt: &yesterdayStart,
			GroupID:   &chatID, // Changed from &chatID to chatID if GroupID is int64
		}
		reportSummary, err := b.TransactionService.GetReports(transactionDTO)
		if err != nil {
			log.Printf("Error getting report for today: %v", err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ខ្ញុំមិនអាចបង្កើតរបាយការណ៍សម្រាប់ថ្ងៃនេះបានទេ។ សូមព្យាយាមម្តងទៀតនៅពេលក្រោយ")
			if _, err := b.botAPI.Send(errMsg); err != nil { // Use b.botAPI
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		reportContent = fmt.Sprintf(
			"📊 `របាយការណ៍សរុបប្រតិបត្តិការណ៍ ` \n"+
				"-----------------------------------\n"+
				"🕑 `ម្សិលមិញ %s` \n"+
				"$ (USD)សរុប : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n៛ (KHR)សរុប : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d", reportDate.Format("2006/01/02"),
			reportSummary.TotalUSD,
			reportSummary.TotalTrxUSD,
			reportSummary.TotalReil,
			reportSummary.TotalTrxReil,
		)
		reportContent = EscapeMarkdownV2(reportContent)
	case "report_date":
		now := time.Now()
		reportDate = now
		transactionDTO := &dtos.TransactionFilter{
			CreatedAt: &now,
			GroupID:   &chatID, // Changed from &chatID to chatID if GroupID is int64
		}
		reportSummary, err := b.TransactionService.GetReports(transactionDTO)
		if err != nil {
			log.Printf("Error getting report for today: %v", err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ខ្ញុំមិនអាចបង្កើតរបាយការណ៍សម្រាប់ថ្ងៃនេះបានទេ។ សូមព្យាយាមម្តងទៀតនៅពេលក្រោយ")
			if _, err := b.botAPI.Send(errMsg); err != nil { // Use b.botAPI
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		reportContent = fmt.Sprintf(
			"📊 របាយការណ៍សរុបប្រតិបត្តិការណ៍\n"+
				"-----------------------------------\n"+
				"🕑 `ថ្ងៃទី %s `\n"+
				"$ (USD)សរុប : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n"+
				"៛ (KHR)សរុប : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n"+
				"-----------------------------------\n"+ // Added newline and made this a separate line for clarity
				"ប្រតិបត្តិការណ៍សង្ខេប តាមធនាគារ\n\n",
			reportDate.Format("2006/01/02"),
			reportSummary.TotalUSD,
			reportSummary.TotalTrxUSD,
			reportSummary.TotalReil,
			reportSummary.TotalTrxReil,
		)
		reportContent = EscapeMarkdownV2(reportContent) // Apply EscapeMarkdownV2

	case "report_bank_today":
		now := time.Now()
		reportDate = now
		transactionDTO := &dtos.TransactionFilter{
			CreatedAt: &now,
			GroupID:   &chatID, // Changed from &chatID to chatID if GroupID is int64
		}
		reportSummary, err := b.TransactionService.GetReportsByBankTransaction(transactionDTO)
		if err != nil {
			log.Printf("Error getting report for today: %v", err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ខ្ញុំមិនអាចបង្កើតរបាយការណ៍សម្រាប់ថ្ងៃនេះបានទេ។ សូមព្យាយាមម្តងទៀតនៅពេលក្រោយ")
			if _, err := b.botAPI.Send(errMsg); err != nil { // Use b.botAPI
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		reportContent = FormatReportForTelegram(reportSummary, reportDate.Format("2006/01/02"))

	default:
		reportContent = "សូមទោស, an unknown report type was selected."
		log.Printf("DEBUG: Unknown callback data received: %s", callbackData)
	}

	// Consolidated check for empty or default error message
	if reportContent == "" || strings.HasPrefix(reportContent, "សូមទោស, an unknown") { // Use strings.HasPrefix for robustness
		reportContent = "សូមទោស,មិនមាន របាយការណ៍ទេ"
		log.Printf("DEBUG: Report content is empty or default error message. Not sending main report, will send fallback.")
	}

	reportMsg := tgbotapi.NewMessage(chatID, reportContent)
	reportMsg.ParseMode = tgbotapi.ModeMarkdownV2 // Set ParseMode for MarkdownV2

	log.Printf("DEBUG: Attempting to send report message to chat %d with content length %d.", chatID, len(reportContent))

	if _, err := b.botAPI.Send(reportMsg); err != nil { // Use b.botAPI
		log.Printf("CRITICAL ERROR: Failed to send report message to chat %d: %v", chatID, err)
	} else {
		log.Printf("DEBUG: Report message successfully sent to chat %d.", chatID)
	}
}
func (b *Bot) HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {
	// Variables to be populated based on callback data
	var reportContent string
	var reportDate time.Time
	var reportType string // Stores the determined type of report (e.g., "report_date", "report_today")
	var (
		startDate time.Time
		endDate   time.Time
	)
	// 1. Answer the callback query (prevents "loading" state on button)
	callback := tgbotapi.NewCallback(callbackQuery.ID, "")
	if _, err := b.botAPI.Request(callback); err != nil {
		log.Printf("Failed to answer callback query: %v", err)
	}

	// 2. Basic validation: Ensure message is not nil
	if callbackQuery.Message == nil {
		log.Printf("Callback query message is nil for query ID: %s", callbackQuery.ID)
		return
	}

	// 3. Extract common callback data and log
	chatID := callbackQuery.Message.Chat.ID
	messageIDToDelete := callbackQuery.Message.MessageID
	callbackData := callbackQuery.Data // The raw data from the button
	username := callbackQuery.From.FirstName
	if callbackQuery.From.LastName != "" {
		username += " " + callbackQuery.From.LastName
	}
	escapedUsername := EscapeMarkdownV2(username)
	utils.InfoLog(callbackData, fmt.Sprintf("DEBUG: Received callback from %s with data: %s", escapedUsername, callbackData))

	// 4. Delete the button message to clean up chat
	deleteConfig := tgbotapi.NewDeleteMessage(chatID, messageIDToDelete)
	if _, err := b.botAPI.Request(deleteConfig); err != nil {
		log.Printf("ERROR: Failed to delete button message (ID: %d) in chat %d after callback from %s: %v",
			messageIDToDelete, chatID, escapedUsername, err)
	} else {
		log.Printf("DEBUG: Successfully deleted button message (ID: %d) in chat %d.", messageIDToDelete, chatID)
	}
	if strings.HasPrefix(callbackData, "report_date_") {
		log.Printf("TRACE: Inside 'report_date_' block. callbackData: '%s'", callbackData)
		dateValue := strings.TrimPrefix(callbackData, "report_date_")
		reportType = "report_date"
		log.Printf("TRACE: After trimming prefix. reportType: '%s', dateValue: '%s'", reportType, dateValue)
		formattedDateForParse := strings.ReplaceAll(dateValue, "_", "-")
		log.Printf("TRACE: Preparing for Parse. formattedDateForParse: '%s'", formattedDateForParse)
		var err error
		parsedTime, err := time.Parse("02-01-2006", formattedDateForParse)
		if err != nil {
			log.Printf("CRITICAL PARSE ERROR: Failed to parse '%s' (derived from '%s', raw callback '%s'): %v",
				formattedDateForParse, dateValue, callbackData, err)
			errMsg := tgbotapi.NewMessage(chatID, "សូមអភ័យទោស ទ្រង់ទ្រាយកាលបរិច្ឆេទមិនត្រឹមត្រូវ។")
			if _, sendErr := b.botAPI.Send(errMsg); sendErr != nil { // Use sendErr to avoid shadowing err
				log.Printf("Error sending Telegram error message: %v", sendErr)
			}
			return // IMPORTANT: Stop execution here if parsing failed
		}
		year, month, day := parsedTime.Date()
		reportDate = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		log.Printf("TRACE: Date parsed successfully. reportDate: '%s'", reportDate.Format("2006-01-02"))

	} else if strings.HasPrefix(callbackData, "report_bank_date") {
		dateValue := strings.TrimPrefix(callbackData, "report_bank_date_")
		reportType = "report_bank_date" // Assign new report type
		formattedDateForParse := strings.ReplaceAll(dateValue, "_", "-")
		var err error
		parsedTime, err := time.Parse("02-01-2006", formattedDateForParse)
		if err != nil {
			log.Printf("CRITICAL PARSE ERROR: Failed to parse 'report_bank_date_' '%s' (derived from '%s', raw callback '%s'): %v",
				formattedDateForParse, dateValue, callbackData, err)
			errMsg := tgbotapi.NewMessage(chatID, "សូមអភ័យទោស ទ្រង់ទ្រាយកាលបរិច្ឆេទមិនត្រឹមត្រូវ។")
			if _, sendErr := b.botAPI.Send(errMsg); sendErr != nil {
				log.Printf("Error sending Telegram error message: %v", sendErr)
			}
			return
		}
		year, month, day := parsedTime.Date()
		reportDate = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	} else if strings.HasPrefix(callbackData, "report_filter_date") {

		now := time.Now()
		reportType = "report_filter_date" // Assign new report type
		if callbackData == "report_filter_date_week" {
			// from date to last date last 1 week
			startOfDayNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			startDate = startOfDayNow.AddDate(0, 0, -6)
			endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

		} else if callbackData == "report_filter_date_month" {
			year, month, _ := now.Date()
			startDate = time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
			firstDayOfNextMonth := startDate.AddDate(0, 1, 0)
			endDate = firstDayOfNextMonth.Add(-time.Nanosecond)

			log.Printf("Filter applied: This month. Start: %s, End: %s",
				startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))

		} else if callbackData == "report_filter_date_lmonth" {
			year, month, _ := now.Date()
			firstDayOfCurrentMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
			startDate = firstDayOfCurrentMonth.AddDate(0, -1, 0)
			endDate = firstDayOfCurrentMonth.Add(-time.Nanosecond)

			log.Printf("Filter applied: Last month. Start: %s, End: %s",
				startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))

		}
		utils.InfoLog(reportType, fmt.Sprintf("Calculated dates: Start = %s, End = %s\n", startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05")))
		fmt.Printf("Calculated dates: Start = %s, End = %s\n",
			startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))
	} else if strings.HasPrefix(callbackData, "report_filter_bank") {

		now := time.Now()
		reportType = "report_filter_bank" // Assign new report type
		if callbackData == "report_filter_bank_week" {
			// from date to last date last 1 week
			startOfDayNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			startDate = startOfDayNow.AddDate(0, 0, -6)
			endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

		} else if callbackData == "report_filter_bank_month" {
			year, month, _ := now.Date()
			startDate = time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
			firstDayOfNextMonth := startDate.AddDate(0, 1, 0)
			endDate = firstDayOfNextMonth.Add(-time.Nanosecond)

			log.Printf("Filter applied: This month. Start: %s, End: %s",
				startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))

		} else if callbackData == "report_filter_bank_lmonth" {
			year, month, _ := now.Date()
			firstDayOfCurrentMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
			startDate = firstDayOfCurrentMonth.AddDate(0, -1, 0)
			endDate = firstDayOfCurrentMonth.Add(-time.Nanosecond)

			log.Printf("Filter applied: Last month. Start: %s, End: %s",
				startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))

		}
		utils.InfoLog(reportType, fmt.Sprintf("Calculated dates: Start = %s, End = %s\n", startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05")))
		fmt.Printf("Calculated dates: Start = %s, End = %s\n",
			startDate.Format("2006-01-02 15:04:05"), endDate.Format("2006-01-02 15:04:05"))

	} else { // This block handles other predefined callback types or unknown ones
		reportType = callbackData
		log.Printf("TRACE: Not a 'report_bank_' command. reportType set to: '%s'", reportType)
	}
	utils.InfoLog(reportType, fmt.Sprintf("DEBUG: Processing report type '%s' for date '%s'", reportType, reportDate.Format("2006-01-02")))

	var reportSummary *services.ReportSummaryDTO
	var serviceErr error   // Use a distinct variable for service errors
	var displayDate string // Date string to be displayed in the report header

	// 6. Centralized Report Generation Logic using a switch on `reportType`
	switch reportType {
	case "report_today":
		now := time.Now()
		reportDate = now
		utils.InfoLog(reportType, fmt.Sprintf("DEBUG: Processing report today '%s' for date '%s'", now))
		transactionDTO := &dtos.TransactionFilter{
			CreatedAt: &now,
			GroupID:   &chatID, // Changed from &chatID to chatID if GroupID is int64
		}
		reportSummary, err := b.TransactionService.GetReports(transactionDTO)
		if err != nil {
			log.Printf("Error getting report for today: %v", err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ខ្ញុំមិនអាចបង្កើតរបាយការណ៍សម្រាប់ថ្ងៃនេះបានទេ។ សូមព្យាយាមម្តងទៀតនៅពេលក្រោយ")
			if _, err := b.botAPI.Send(errMsg); err != nil { // Use b.botAPI
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		reportContent = fmt.Sprintf(
			"📊 របាយការណ៍ សរុបប្រតិបត្តិការណ៍\n"+
				"-----------------------------------\n"+
				"🕑 `ថ្ងៃនេះ %s `\n"+
				"$ (USD) : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n"+
				"៛ (KHR) : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n",
			reportDate.Format("2006/01/02"),
			reportSummary.TotalUSD,
			reportSummary.TotalTrxUSD,
			reportSummary.TotalReil,
			reportSummary.TotalTrxReil,
		)
		reportContent = EscapeMarkdownV2(reportContent) // Apply EscapeMarkdownV2
	case "report_yesterday":
		now := time.Now()
		todayStart := now.Truncate(24 * time.Hour)
		yesterdayStart := todayStart.Add(-24 * time.Hour)
		reportDate = yesterdayStart
		transactionDTO := &dtos.TransactionFilter{
			CreatedAt: &yesterdayStart,
			GroupID:   &chatID, // Changed from &chatID to chatID if GroupID is int64
		}
		reportSummary, err := b.TransactionService.GetReports(transactionDTO)
		if err != nil {
			log.Printf("Error getting report for today: %v", err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ខ្ញុំមិនអាចបង្កើតរបាយការណ៍សម្រាប់ថ្ងៃនេះបានទេ។ សូមព្យាយាមម្តងទៀតនៅពេលក្រោយ")
			if _, err := b.botAPI.Send(errMsg); err != nil { // Use b.botAPI
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		reportContent = fmt.Sprintf(
			"📊 `របាយការណ៍ សរុបប្រតិបត្តិការណ៍ ` \n"+
				"-----------------------------------\n"+
				"🕑 `ម្សិលមិញ %s` \n"+
				"$ (USD) : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n៛ (KHR) : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d", reportDate.Format("2006/01/02"),
			reportSummary.TotalUSD,
			reportSummary.TotalTrxUSD,
			reportSummary.TotalReil,
			reportSummary.TotalTrxReil,
		)
		reportContent = EscapeMarkdownV2(reportContent)
	case "report_date":
		transactionDTO := &dtos.TransactionFilter{
			CreatedAt: &reportDate,
			GroupID:   &chatID, // Changed from &chatID to chatID if GroupID is int64
		}
		reportSummary, err := b.TransactionService.GetReports(transactionDTO)
		if err != nil {
			log.Printf("Error getting report for today: %v", err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ខ្ញុំមិនអាចបង្កើតរបាយការណ៍សម្រាប់ថ្ងៃនេះបានទេ។ សូមព្យាយាមម្តងទៀតនៅពេលក្រោយ")
			if _, err := b.botAPI.Send(errMsg); err != nil { // Use b.botAPI
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		reportContent = fmt.Sprintf(
			"📊 របាយការណ៍ សរុបប្រតិបត្តិការណ៍\n"+
				"-----------------------------------\n"+
				"🕑 `ថ្ងៃទី %s `\n\n"+
				"$ (USD) : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n"+
				"៛ (KHR) : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n",
			reportDate.Format("2006/01/02"),
			reportSummary.TotalUSD,
			reportSummary.TotalTrxUSD,
			reportSummary.TotalReil,
			reportSummary.TotalTrxReil,
		)
		reportContent = EscapeMarkdownV2(reportContent) // Apply EscapeMarkdownV2
	case "report_filter_date":
		transactionDTO := &dtos.TransactionFilter{
			CreatedAt: &reportDate,
			StartDate: &startDate,
			EndDate:   &endDate,
			GroupID:   &chatID, // Changed from &chatID to chatID if GroupID is int64
		}
		reportSummary, err := b.TransactionService.GetFilterReports(transactionDTO)
		if err != nil {
			log.Printf("Error getting report for today: %v", err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ខ្ញុំមិនអាចបង្កើតរបាយការណ៍សម្រាប់ថ្ងៃនេះបានទេ។ សូមព្យាយាមម្តងទៀតនៅពេលក្រោយ")
			if _, err := b.botAPI.Send(errMsg); err != nil { // Use b.botAPI
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		reportContent = fmt.Sprintf(
			"📊 របាយការណ៍ សរុបប្រតិបត្តិការណ៍\n"+
				"-----------------------------------\n"+
				"🕑 ថ្ងៃទី: %s\n"+
				"$ (USD) : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n"+
				"៛ (KHR) : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n",

			startDate.Format("2006/01/02")+"-"+endDate.Format("2006/01/02"),
			reportSummary.TotalUSD,
			reportSummary.TotalTrxUSD,
			reportSummary.TotalReil,
			reportSummary.TotalTrxReil,
		)
		reportContent = EscapeMarkdownV2(reportContent) // Apply EscapeMarkdownV2
	case "report_bank_yesterday":
		now := time.Now()
		todayStart := now.Truncate(24 * time.Hour)
		yesterdayStart := todayStart.Add(-24 * time.Hour)
		reportDate = yesterdayStart
		transactionDTO := &dtos.TransactionFilter{
			CreatedAt: &yesterdayStart,
			GroupID:   &chatID, // Changed from &chatID to chatID if GroupID is int64
		}
		reportSummary, err := b.TransactionService.GetReportsByBankTransaction(transactionDTO)
		if err != nil {
			log.Printf("Error getting report for today: %v", err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ខ្ញុំមិនអាចបង្កើតរបាយការណ៍សម្រាប់ថ្ងៃនេះបានទេ។ សូមព្យាយាមម្តងទៀតនៅពេលក្រោយ")
			if _, err := b.botAPI.Send(errMsg); err != nil { // Use b.botAPI
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		reportContent = FormatReportForTelegram(reportSummary, reportDate.Format("2006/01/02"))
	case "report_bank_today":
		now := time.Now()
		reportDate = now
		transactionDTO := &dtos.TransactionFilter{
			CreatedAt: &now,
			GroupID:   &chatID, // Changed from &chatID to chatID if GroupID is int64
		}
		reportSummary, err := b.TransactionService.GetReportsByBankTransaction(transactionDTO)
		if err != nil {
			log.Printf("Error getting report for today: %v", err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ខ្ញុំមិនអាចបង្កើតរបាយការណ៍សម្រាប់ថ្ងៃនេះបានទេ។ សូមព្យាយាមម្តងទៀតនៅពេលក្រោយ")
			if _, err := b.botAPI.Send(errMsg); err != nil { // Use b.botAPI
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		reportContent = FormatReportForTelegram(reportSummary, reportDate.Format("2006/01/02"))
	case "report_bank_date":
		transactionDTO := &dtos.TransactionFilter{
			CreatedAt: &reportDate,
			GroupID:   &chatID, // Changed from &chatID to chatID if GroupID is int64
		}
		reportSummary, err := b.TransactionService.GetReportsByBankTransaction(transactionDTO)
		if err != nil {
			log.Printf("Error getting report for today: %v", err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ខ្ញុំមិនអាចបង្កើតរបាយការណ៍សម្រាប់ថ្ងៃនេះបានទេ។ សូមព្យាយាមម្តងទៀតនៅពេលក្រោយ")
			if _, err := b.botAPI.Send(errMsg); err != nil { // Use b.botAPI
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		reportContent = FormatReportForTelegram(reportSummary, reportDate.Format("2006/01/02"))
	case "report_filter_bank":
		transactionDTO := &dtos.TransactionFilter{
			CreatedAt: &reportDate,
			StartDate: &startDate,
			EndDate:   &endDate,
			GroupID:   &chatID, // Changed from &chatID to chatID if GroupID is int64
		}
		reportSummary, err := b.TransactionService.GetFilterReportsByBankTransaction(transactionDTO)
		if err != nil {
			log.Printf("Error getting report for today: %v", err)
			errMsg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "សូមអភ័យទោស ខ្ញុំមិនអាចបង្កើតរបាយការណ៍សម្រាប់ថ្ងៃនេះបានទេ។ សូមព្យាយាមម្តងទៀតនៅពេលក្រោយ")
			if _, err := b.botAPI.Send(errMsg); err != nil { // Use b.botAPI
				log.Printf("Error sending error message: %v", err)
			}
			return
		}
		reportContent = FormatReportForTelegram(reportSummary, startDate.Format("2006/01/02")+"-"+endDate.Format("2006/01/02"))
	default:
		reportContent = "សូមទោស, an unknown report type was selected."
		log.Printf("DEBUG: Unknown callback data received: %s", callbackData)
	}
	if serviceErr == nil && reportSummary != nil && reportContent == "" {
		reportContent = FormatReportForTelegram(reportSummary, displayDate)
	} else if reportContent == "" {
		// Fallback error if reportSummary is nil or an error occurred but no specific error message was assigned above.
		reportContent = "សូមទោស, មិនមានរបាយការណ៍សម្រាប់កាលបរិច្ឆេទនេះទេ។"
	}

	reportMsg := tgbotapi.NewMessage(chatID, reportContent)
	reportMsg.ParseMode = tgbotapi.ModeMarkdownV2 // Ensure MarkdownV2 is used for all reports
	reportMsg.DisableWebPagePreview = true        // Good practice for reports to avoid unexpected link previews

	log.Printf("DEBUG: Attempting to send report message to chat %d with content length %d.", chatID, len(reportContent))

	if _, err := b.botAPI.Send(reportMsg); err != nil {
		log.Printf("CRITICAL ERROR: Failed to send report message to chat %d: %v", chatID, err)
	} else {
		log.Printf("DEBUG: Report message successfully sent to chat %d.", chatID)
	}
}
