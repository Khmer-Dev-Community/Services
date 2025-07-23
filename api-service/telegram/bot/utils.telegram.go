package bot

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"telegram-service/telegram/services"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func EscapeMarkdownV2(text string) string {
	text = strings.ReplaceAll(text, `\`, `\\`)
	replacer := strings.NewReplacer(
		`_`, `\_`,
		`[`, `\[`,
		`]`, `\]`,
		`(`, `\(`,
		`)`, `\)`,
		`~`, `\~`,
		`>`, `\>`,
		`#`, `\#`,
		`+`, `\+`,
		`-`, `\-`, // Escape hyphen (-)
		`=`, `\=`,
		`|`, `\|`,
		`{`, `\{`,
		`}`, `\}`,
		`.`, `\.`, // Important for decimals in numbers like 20.00
		`!`, `\!`,
	)
	return replacer.Replace(text)
}

func (b *Bot) GenerateDateRangeButtons(startOffsetDays int, numberOfButtons int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	now := time.Now()

	for i := 0; i < numberOfButtons; i++ {
		targetDate := now.AddDate(0, 0, -(startOffsetDays + i))

		dateString := targetDate.Format("02-01-2006") // Format as DD-MM-YYYY
		dateValue := targetDate.Format("02_01_2006")
		callbackData := fmt.Sprintf("report_date_%s", dateValue)

		button := tgbotapi.NewInlineKeyboardButtonData("ថ្ងៃទី"+dateString, callbackData)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
func (b *Bot) GenerateDateRangeButtonsBank(startOffsetDays int, numberOfButtons int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	now := time.Now()

	for i := 0; i < numberOfButtons; i++ {
		targetDate := now.AddDate(0, 0, -(startOffsetDays + i))

		dateString := targetDate.Format("02-01-2006") // Format as DD-MM-YYYY
		dateValue := targetDate.Format("02_01_2006")
		callbackData := fmt.Sprintf("report_bank_date_%s", dateValue)

		button := tgbotapi.NewInlineKeyboardButtonData("ថ្ងៃទី"+dateString, callbackData)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// ABA Pay (USD & KHR, 'paid by' type)
var abaPaidByRegex = regexp.MustCompile(`^([$៛])([\d,.]+) paid by (.+?) \(\*.+?\) on (.+?) via (.+?) at (.+?)\. Trx\. ID: (\d+), APV: (\d+)\.?$`)

// ACLEDA (Received KHR 'បានទទួល' type)
// Example: បានទទួល 25,000 រៀល ពី 095 933 996 OUN RATANA, ថ្ងៃទី១៥ កក្កដា ២០២៥ ០៧:៤២ព្រឹក, លេខយោង 51960552015, នៅ HENG ROTHPANHA.
var receivedKHRRegex = regexp.MustCompile(`^បានទទួល ([\d,]+) រៀល ពី (.+?), (.+?), លេខយោង (\d+), នៅ (.+?)\.?$`)

// ACLEDA (Received USD 'បានទទួល' type)
// Example: បានទទួល 3.00 ដុល្លារ ពី 071 7718 763 THOL REAKSMEY, ថ្ងៃទី១៤ កក្កដា ២០២៥ ០៧:០៩ព្រឹក, លេខយោង 51950404514, នៅ HENG ROTHPANHA.
var receivedUSDKhmerRegex = regexp.MustCompile(`^បានទទួល ([\d.]+) ដុល្លារ ពី (.+?), (.+?), លេខយោង (\d+), នៅ (.+?)\.?$`)

// Wing (Received KHR 'លោកអ្នកបានទទួលប្រាក់ចំនួន' type)
// Example: លោកអ្នកបានទទួលប្រាក់ចំនួន 26,000 រៀល ពីឈ្មោះ POLIN CHEA ធនាគារ ABA Bank តាមការស្កេន KHQR ថ្ងៃទី ១៤ កក្កដា ២០២៥ ម៉ោង ០៦:៤៦ព្រឹក (Hash. cd837ef1).
var receivedWingKHRRegex = regexp.MustCompile(`^លោកអ្នកបានទទួលប្រាក់ចំនួន ([\d,]+) រៀល\s*ពីឈ្មោះ (.+?) ធនាគារ (.+?) តាមការស្កេន\s*(.+?)\s*ថ្ងៃទី (.+?) \(\s*Hash\.\s*([a-f0-9]+)\)\.?$`)

// Wing (Received USD 'លោកអ្នកបានទទួលប្រាក់ចំនួន' type)
// Example: លោកអ្នកបានទទួលប្រាក់ចំនួន 6.00 ដុល្លារ ពីឈ្មោះ SOPHOIN PHAL ធនាគារ ABA Bank តាមការស្កេន KHQR ថ្ងៃទី ១៤ កក្កដា ២០២៥ ម៉ោង ០៥:១១ល្ងាច (Hash. 30cdfbea).
var receivedWingUSDRegex = regexp.MustCompile(`^លោកអ្នកបានទទួលប្រាក់ចំនួន ([\d.]+) ដុល្លារ\s*ពីឈ្មោះ (.+?) ធនាគារ (.+?) តាមការស្កេន\s*(.+?)\s*ថ្ងៃទី (.+?) \(\s*Hash\.\s*([a-f0-9]+)\)\.?$`)

// ParseTransactionText parses a raw transaction string into TransactionData.
func ParseTransactionText(rawText string) (TransactionData, error) {
	data := TransactionData{RawText: rawText}
	var matches []string
	var err error
	matches = abaPaidByRegex.FindStringSubmatch(rawText)
	if len(matches) > 0 {
		data.Currency = strings.TrimSpace(matches[1])             // Should be '$' or '៛'
		cleanAmountStr := strings.ReplaceAll(matches[2], ",", "") // Remove commas for KHR
		data.Amount, err = strconv.ParseFloat(cleanAmountStr, 64)
		if err != nil {
			return data, fmt.Errorf("ABA 'paid by' parser: failed to parse amount '%s': %w", matches[2], err)
		}
		data.Sender = strings.TrimSpace(matches[3])
		data.On = strings.TrimSpace(matches[4])
		data.Via = strings.TrimSpace(matches[5])      // e.g., "ABA PAY"
		data.Receiver = strings.TrimSpace(matches[6]) // e.g., "HENG ROTHPANHA"
		data.TrxID = strings.TrimSpace(matches[7])
		data.APV = strings.TrimSpace(matches[8])
		data.Bank = "ABA" // Explicitly setting the bank
		return data, nil
	}

	// 2. Try to parse ACLEDA "បានទទួល" messages (KHR)
	// Example: បានទទួល 25,000 រៀល ពី 095 933 996 OUN RATANA, ថ្ងៃទី១៥ កក្កដា ២០២៥ ០៧:៤២ព្រឹក, លេខយោង 51960552015, នៅ HENG ROTHPANHA.
	matches = receivedKHRRegex.FindStringSubmatch(rawText)
	if len(matches) > 0 {
		data.Currency = "៛"
		cleanAmountStr := strings.ReplaceAll(matches[1], ",", "")
		data.Amount, err = strconv.ParseFloat(cleanAmountStr, 64)
		if err != nil {
			return data, fmt.Errorf("ACLEDA KHR parser: failed to parse amount '%s': %w", matches[1], err)
		}
		data.Sender = strings.TrimSpace(matches[2]) // Includes number and name
		data.On = strings.TrimSpace(matches[3])
		data.APV = strings.TrimSpace(matches[4]) // លេខយោង is APV
		data.Receiver = strings.TrimSpace(matches[5])
		data.Via = "ACLEDA" // Explicitly setting for ACLEDA format
		data.Bank = "ACLEDA"
		data.TrxID = "" // ACLEDA often doesn't have a distinct TrxID, APV is used
		return data, nil
	}

	// 3. Try to parse ACLEDA "បានទទួល" messages (USD)
	// Example: បានទទួល 3.00 ដុល្លារ ពី 071 7718 763 THOL REAKSMEY, ថ្ងៃទី១៤ កក្កដា ២០២៥ ០៧:០៩ព្រឹក, លេខយោង 51950404514, នៅ HENG ROTHPANHA.
	matches = receivedUSDKhmerRegex.FindStringSubmatch(rawText)
	if len(matches) > 0 {
		data.Currency = "$"
		data.Amount, err = strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return data, fmt.Errorf("ACLEDA USD parser: failed to parse amount '%s': %w", matches[1], err)
		}
		data.Sender = strings.TrimSpace(matches[2]) // Includes number and name
		data.On = strings.TrimSpace(matches[3])
		data.APV = strings.TrimSpace(matches[4]) // លេខយោង is APV
		data.Receiver = strings.TrimSpace(matches[5])
		data.Via = "ACLEDA" // Explicitly setting for ACLEDA format
		data.Bank = "ACLEDA"
		data.TrxID = "" // ACLEDA often doesn't have a distinct TrxID, APV is used
		return data, nil
	}

	// 4. Try to parse Wing "លោកអ្នកបានទទួលប្រាក់ចំនួន" messages (KHR)
	// Example: លោកអ្នកបានទទួលប្រាក់ចំនួន 26,000 រៀល ពីឈ្មោះ POLIN CHEA ធនាគារ ABA Bank តាមការស្កេន KHQR ថ្ងៃទី ១៤ កក្កដា ២០២៥ ម៉ោង ០៦:៤៦ព្រឹក (Hash. cd837ef1).
	matches = receivedWingKHRRegex.FindStringSubmatch(rawText)
	if len(matches) > 0 {
		data.Currency = "៛"
		cleanAmountStr := strings.ReplaceAll(matches[1], ",", "")
		data.Amount, err = strconv.ParseFloat(cleanAmountStr, 64)
		if err != nil {
			return data, fmt.Errorf("Wing KHR parser: failed to parse amount '%s': %w", matches[1], err)
		}
		data.Sender = strings.TrimSpace(matches[2]) // POLIN CHEA
		data.Bank = strings.TrimSpace(matches[3])   // ABA Bank
		data.Via = strings.TrimSpace(matches[4])    // តាមការស្កេន KHQR
		data.On = strings.TrimSpace(matches[5])     // ថ្ងៃទី ១៤ កក្កដា ២០២៥ ម៉ោង ០៦:៤៦ព្រឹក
		data.TrxID = strings.TrimSpace(matches[6])  // Hash. cd837ef1
		data.APV = data.TrxID                       // Use TrxID as APV for Wing Hash
		data.Receiver = ""                          // Not present in this format
		return data, nil
	}

	// 5. Try to parse Wing "លោកអ្នកបានទទួលប្រាក់ចំនួន" messages (USD)
	// Example: លោកអ្នកបានទទួលប្រាក់ចំនួន 6.00 ដុល្លារ ពីឈ្មោះ SOPHOIN PHAL ធនាគារ ABA Bank តាមការស្កេន KHQR ថ្ងៃទី ១៤ កក្កដា ២០២៥ ម៉ោង ០៥:១១ល្ងាច (Hash. 30cdfbea).
	matches = receivedWingUSDRegex.FindStringSubmatch(rawText)
	if len(matches) > 0 {
		data.Currency = "$"
		data.Amount, err = strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return data, fmt.Errorf("Wing USD parser: failed to parse amount '%s': %w", matches[1], err)
		}
		data.Sender = strings.TrimSpace(matches[2]) // SOPHOIN PHAL
		data.Bank = strings.TrimSpace(matches[3])   // ABA Bank
		data.Via = strings.TrimSpace(matches[4])    // តាមការស្កេន KHQR
		data.On = strings.TrimSpace(matches[5])     // ថ្ងៃទី ១៤ កក្កដា ២០២៥ ម៉ោង ០៥:១១ល្ងាច
		data.TrxID = strings.TrimSpace(matches[6])  // Hash. 30cdfbea
		data.APV = data.TrxID                       // Use TrxID as APV for Wing Hash
		data.Receiver = ""                          // Not present in this format
		return data, nil
	}

	return data, fmt.Errorf("no matching transaction format found for: %s", rawText)
}

type ReportSummaryDTO struct {
	TotalUSD      float64 `json:"totalUSD"`
	TotalTrxUSD   int     `json:"totalTrxUSD"`
	TotalReil     float64 `json:"totalReil"`
	TotalTrxReil  int     `json:"totalTrxReil"`
	BankBreakdown map[string]map[string]struct {
		Total float64
		Count int
	} `json:"bankBreakdown"` // Added JSON tag for consistency
}

func FormatReportForTelegram(reportDTO *services.ReportSummaryDTO, date string) string {
	var sb strings.Builder

	// Start with the specific header for the bank summary section
	sb.WriteString(fmt.Sprintf("📊 ប្រតិបត្តិការណ៍សង្ខេប តាមធនាគារ \n-----------------------------------\n 🕑 `ថ្ងៃទី %s` \n", date))
	// Collect all unique transport names to ensure consistent ordering
	transportNames := make([]string, 0, len(reportDTO.BankBreakdown))
	for name := range reportDTO.BankBreakdown {
		transportNames = append(transportNames, name)
	}
	sort.Strings(transportNames) // Sort transport names alphabetically

	// Iterate through each transport name to build its summary
	for i, transportName := range transportNames {
		currencies := reportDTO.BankBreakdown[transportName]
		totalTrxForTransport := 0
		if usdSummary, ok := currencies["$"]; ok {
			totalTrxForTransport += usdSummary.Count
		}
		if khrSummary, ok := currencies["៛"]; ok {
			totalTrxForTransport += khrSummary.Count
		}
		cleanedName := formatTransportName(transportName)
		sb.WriteString(fmt.Sprintf("🏦 *`%s`* : សរុបប្រតិបត្តិការណ៍ %d\n", cleanedName, totalTrxForTransport))
		if usdSum, ok := currencies["$"]; ok {
			sb.WriteString(fmt.Sprintf("  $ (USD) : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n", usdSum.Total, usdSum.Count))
		} else {
			sb.WriteString(fmt.Sprintf("  $ (USD) : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n", 0.00, 0))
		}
		if khrSum, ok := currencies["៛"]; ok {
			sb.WriteString(fmt.Sprintf("  ៛ (KHR)  : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n", khrSum.Total, khrSum.Count))
		} else {
			sb.WriteString(fmt.Sprintf("  ៛ (KHR)  : *`%.2f`* | ប្រតិបត្តិការណ៍ : %d\n", 0.00, 0))
		}
		if i < len(transportNames)-1 {
			sb.WriteString("--------------------\n")
		}
	}

	return EscapeMarkdownV2(sb.String()) // Apply markdown escaping before returning
}

func formatTransportName(rawName string) string {
	// Remove trailing _bot or Bot (case-insensitive)
	re := regexp.MustCompile(`(?i)(_bot|bot)$`)
	cleanName := re.ReplaceAllString(rawName, "")
	return cleanName
}
