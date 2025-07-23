package bot

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"telegram-service/telegram/repository"
	"telegram-service/telegram/services"
	"telegram-service/utils"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotService struct {
}

type TransactionData struct {
	Currency string
	Amount   float64
	On       string
	Via      string
	TrxID    string
	APV      string
	RawText  string
	Sender   string
	Receiver string
	Bank     string
}
type Bot struct {
	Name  string
	Token string
	Debug bool
	//	GroupID            int64
	BackBtn            bool
	locationRequested  map[int64]time.Time
	ReplyCounts        map[int64]int
	Replies            map[int64][]string
	mu                 sync.Mutex
	TransactionService services.TransactionService
	botAPI             *tgbotapi.BotAPI
}
type TelegramBotService struct {
	repo *repository.TransactionRepository
}

func NewTelegramBotService(repo *repository.TransactionRepository) *TelegramBotService {

	return &TelegramBotService{repo: repo}
}

func NewBotService(name, token string, debug bool, trxService services.TransactionService) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}
	return &Bot{
		Name:  name,
		Token: token,
		Debug: debug,
		//GroupID:            groupID,
		locationRequested:  make(map[int64]time.Time),
		ReplyCounts:        make(map[int64]int),
		Replies:            make(map[int64][]string),
		TransactionService: trxService,
		botAPI:             botAPI,
	}, nil
}
func (b *Bot) Start() {
	log.Printf("Starting bot %s...", b.Name)
	// Initialize the bot
	botAPI, err := tgbotapi.NewBotAPI(b.Token)
	if err != nil {
		log.Printf("Failed to initialize bot %s: %v", b.Name, err)
		return
	}
	log.Printf("Bot %s initialized successfully", b.Name)
	// Set debug mode
	botAPI.Debug = b.Debug
	log.Printf("Authorized bot %s (%s)", b.Name, botAPI.Self.UserName)
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "ចាប់ផ្តើម ដំណើរការ រ៉ូបូត "},
		{Command: "report", Description: "📊 របាយការណ៍"},
		{Command: "bybank", Description: "🏛 របាយការណ៍ Bank"},
		{Command: "group", Description: "⚙️ ទាញយក លេខសំគាល់ក្រុម"},
		{Command: "menu", Description: "📓 បើកប៊ូតុង មីនុយរបាយការ"},
		{Command: "register", Description: "🔐 ចុះឈ្មោះ មើលរបាយការណ៍"},
		{Command: "miniapp", Description: "🌐 ប្រើប្រាស់ App"},
		// Add more commands here as needed
	}
	commandConfig := tgbotapi.NewSetMyCommands(commands...)
	_, err = botAPI.Request(commandConfig)
	if err != nil {
		log.Printf("Failed to set commands for bot %s: %v", b.Name, err)
	} else {
		log.Printf("Bot %s commands set successfully", b.Name)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := botAPI.GetUpdatesChan(u)
	log.Printf("Bot %s is now listening for updates...", b.Name)
	for update := range updates {
		if update.Message != nil {
			b.handleMessage(botAPI, update.Message, update)
		} else if update.CallbackQuery != nil { // <-- ADD THIS BLOCK
			b.HandleCallbackQuery(update.CallbackQuery)
		}

	}
}

// handleMessage handles regular messages
func (b *Bot) handleMessage(botAPI *tgbotapi.BotAPI, message *tgbotapi.Message, update tgbotapi.Update) {
	utils.InfoLog(message, "Income Message")
	if message.IsCommand() && message.Command() == "start" {
		b.initKeyBoard(botAPI, message)
		return
	}
	if message.IsCommand() && message.Command() == "group" {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("GroupID: %d", message.From.ID))
		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = removeKeyboard
		botAPI.Send(msg)
		return
	}
	// Check if the message is a reply to another message
	if message.ReplyToMessage != nil {
		b.mu.Lock()
		b.ReplyCounts[int64(message.ReplyToMessage.MessageID)]++
		b.Replies[int64(message.ReplyToMessage.MessageID)] = append(b.Replies[int64(message.ReplyToMessage.MessageID)], message.Text)
		b.mu.Unlock()
	}
	b.handleMessageTextDynamic(botAPI, message.Chat.ID, message)
	//b.HandleUpdate(botAPI, update)

}

func (b *Bot) initKeyBoard(botAPI *tgbotapi.BotAPI, message *tgbotapi.Message) bool {
	var Username = message.From.FirstName + " " + message.From.LastName

	var keyboardButtons []tgbotapi.KeyboardButton

	if len(keyboardButtons) == 0 {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("User: %s, ChatID: %d", Username, message.From.ID))
		removeKeyboard := tgbotapi.NewRemoveKeyboard(true)
		msg.ReplyMarkup = removeKeyboard
		botAPI.Send(msg)
		return false
	}
	var replyKeyboard tgbotapi.ReplyKeyboardMarkup
	if b.BackBtn {
		replyKeyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(keyboardButtons...),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("⬅️ Back"),
			),
		)
	} else {
		replyKeyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(keyboardButtons...),
		)
	}
	// Set properties after initialization
	replyKeyboard.ResizeKeyboard = true
	replyKeyboard.OneTimeKeyboard = true
	// Choose which keyboard to attach to the message
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("User: %s, ChatID: %d", Username, message.From.ID))
	msg.ReplyMarkup = replyKeyboard // Set the replyKeyboard here
	if _, err := botAPI.Send(msg); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	return true
}

func (b *Bot) handleMessageTextDynamic(botAPI *tgbotapi.BotAPI, botGroupID int64, message *tgbotapi.Message) {
	utils.InfoLog(message, fmt.Sprintf("IncomeMessage GroupID : [%d]", botGroupID))
	var Username = message.From.FirstName + " " + message.From.LastName
	var keyboardButtons []tgbotapi.KeyboardButton
	additionalButtons := []tgbotapi.KeyboardButton{}
	keyboardButtons = append(keyboardButtons, additionalButtons...)
	var replyKeyboard tgbotapi.ReplyKeyboardMarkup
	if b.BackBtn {
		replyKeyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(keyboardButtons...),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("⬅️ Back"),
			),
		)
	} else {
		replyKeyboard = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(keyboardButtons...),
		)
	}
	replyKeyboard.ResizeKeyboard = true
	replyKeyboard.OneTimeKeyboard = true

	if containsKeyword(message.Text, "⚙️ ចុះឈ្មោះ") || containsKeyword(message.Text, "/register") {
		msg := tgbotapi.NewMessage(message.Chat.ID, Username+"Your has registerd! ")
		msg.ReplyToMessageID = message.MessageID
		// Create inline keyboard button
		inlineBtn := tgbotapi.NewInlineKeyboardButtonURL("សូមទាកទងទៅកាន់", "t.me/PHEARUN_UM")
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(inlineBtn),
		)
		msg.ReplyMarkup = inlineKeyboard
		botAPI.Send(msg)
		return
	}
	if containsKeyword(message.Text, "👤 ទំនាក់ទំនង") || containsKeyword(message.Text, "/contact") {
		msg := tgbotapi.NewMessage(message.Chat.ID, "សូមទាក់ទង ដើម្បីពីភាក្សារ")
		msg.ReplyToMessageID = message.MessageID
		// Create inline keyboard button
		inlineBtn := tgbotapi.NewInlineKeyboardButtonURL("សូមទាកទងទៅកាន់", "t.me/PHEARUN_UM")
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(inlineBtn),
		)
		msg.ReplyMarkup = inlineKeyboard
		botAPI.Send(msg)
		return
	}
	if containsKeyword(message.Text, "🌐 MiniApp") || containsKeyword(message.Text, "miniapp") {
		msg := tgbotapi.NewMessage(message.Chat.ID, "ប្រើប្រាស់ កាន់តែងាយស្រួល ទាញរបាយការណ៍")
		msg.ReplyToMessageID = message.MessageID
		// Create inline keyboard button
		webAppInfo := tgbotapi.WebAppInfo{
			URL: fmt.Sprintf("https://track.igflexs.com?bot=Tg8899bot&group=%d", message.Chat.ID),
		}
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonWebApp("ប្រើប្រាស់ App", webAppInfo),
			),
		)
		msg.ReplyMarkup = inlineKeyboard
		botAPI.Send(msg)
		return
	}
	if strings.Contains(message.Text, "menu") { // Example trigger
		b.sendGroupReplyKeyboard(botAPI, message.Chat.ID, message.MessageID)
		return
	}
	if containsKeyword(message.Text, "📊 របាយការណ៍") || containsKeyword(message.Text, "/report") {
		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "សូមជ្រើសរើសប្រភេទរបាយការណ៍, "+message.From.FirstName+"!")
		replyMsg.ReplyToMessageID = message.MessageID

		var keyboardRows [][]tgbotapi.InlineKeyboardButton

		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🕑 ថ្ងៃនេះ", "report_today"),
			tgbotapi.NewInlineKeyboardButtonData("🕑 ម្សិលមិញ", "report_yesterday"),
		))

		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🕑 ៧ថ្ងៃចុងក្រោយ", "report_filter_date_week"),
		))

		remainingDaysKeyboard := b.GenerateDateRangeButtons(2, 6)

		var allDateRangeButtons []tgbotapi.InlineKeyboardButton
		for _, row := range remainingDaysKeyboard.InlineKeyboard {
			allDateRangeButtons = append(allDateRangeButtons, row...)
		}

		const buttonsPerRow = 2
		var groupedDateRangeRows [][]tgbotapi.InlineKeyboardButton

		for i := 0; i < len(allDateRangeButtons); i += buttonsPerRow {
			end := i + buttonsPerRow
			if end > len(allDateRangeButtons) {
				end = len(allDateRangeButtons)
			}

			currentRow := allDateRangeButtons[i:end]
			groupedDateRangeRows = append(groupedDateRangeRows, currentRow)
		}

		for _, row := range groupedDateRangeRows {
			keyboardRows = append(keyboardRows, row)
		}
		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🕑 ១ខែនេះ", "report_filter_date_month"),
			tgbotapi.NewInlineKeyboardButtonData("🕑 ១ខែមុន", "report_filter_date_lmonth"),
		))

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

		replyMsg.ReplyMarkup = inlineKeyboard

		if _, err := botAPI.Send(replyMsg); err != nil {
			log.Printf("Failed to send report options message: %v", err)
		}

		deleteConfig := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
		if _, err := botAPI.Request(deleteConfig); err != nil {
			log.Printf("Failed to delete user message (ID: %d) in chat %d: %v", message.MessageID, message.Chat.ID, err)
		}
		//b.initKeyBoard(botAPI, message)
		return
	}
	if containsKeyword(message.Text, "🏛 របាយការណ៍ Bank") || containsKeyword(message.Text, "🏛 ធនាគារ") || containsKeyword(message.Text, "/bybank") {
		replyMsg := tgbotapi.NewMessage(message.Chat.ID, "សូមជ្រើសរើសប្រភេទរបាយការណ៍, "+message.From.FirstName+"!")
		replyMsg.ReplyToMessageID = message.MessageID

		var keyboardRows [][]tgbotapi.InlineKeyboardButton

		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🕑 ថ្ងៃនេះ", "report_bank_today"),
			tgbotapi.NewInlineKeyboardButtonData("🕑 ម្សិលមិញ", "report_bank_yesterday"),
		))

		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🕑 ៧ថ្ងៃចុងក្រោយ", "report_filter_bank_week"),
		))

		remainingDaysKeyboard := b.GenerateDateRangeButtonsBank(2, 6)

		var allDateRangeButtons []tgbotapi.InlineKeyboardButton
		for _, row := range remainingDaysKeyboard.InlineKeyboard {
			allDateRangeButtons = append(allDateRangeButtons, row...)
		}

		const buttonsPerRow = 2
		var groupedDateRangeRows [][]tgbotapi.InlineKeyboardButton

		for i := 0; i < len(allDateRangeButtons); i += buttonsPerRow {
			end := i + buttonsPerRow
			if end > len(allDateRangeButtons) {
				end = len(allDateRangeButtons)
			}

			currentRow := allDateRangeButtons[i:end]
			groupedDateRangeRows = append(groupedDateRangeRows, currentRow)
		}

		for _, row := range groupedDateRangeRows {
			keyboardRows = append(keyboardRows, row)
		}
		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🕑 ១ខែនេះ", "report_filter_bank_month"),
			tgbotapi.NewInlineKeyboardButtonData("🕑 ១ខែមុន", "report_filter_bank_lmonth"),
		))

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

		replyMsg.ReplyMarkup = inlineKeyboard

		if _, err := botAPI.Send(replyMsg); err != nil {
			log.Printf("Failed to send report options message: %v", err)
		}

		deleteConfig := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
		if _, err := botAPI.Request(deleteConfig); err != nil {
			log.Printf("Failed to delete user message (ID: %d) in chat %d: %v", message.MessageID, message.Chat.ID, err)
		}
		//b.initKeyBoard(botAPI, message)
		return
	}
	transactions := []string{message.Text}
	parsedResults := []TransactionData{}
	for _, t := range transactions {
		data, err := ParseTransactionText(t)
		if err != nil {
			fmt.Printf("Error parsing text:\n%s\nError: %v\n\n", t, err)
			continue
		}
		parsedResults = append(parsedResults, data)
	}
	/*
		for _, result := range parsedResults {

			if result.TrxID != "" {
				transactionDTO := &dtos.TransactionCreateDTO{
					GroupID:         message.Chat.ID,
					TrxID:           result.TrxID,
					Amount:          result.Amount,
					Currency:        result.Currency,
					Sender:          result.Sender,
					Transport:       result.Via,   // <--- THIS MUST BE result.Via (e.g., "ABA PAY at SAN PHAN.")
					APV:             result.APV,   // <--- THIS MUST BE result.APV (e.g., "660217")
					TransactionDate: result.On,    // <--- THIS MUST BE transactionDateFormatted
					TxtRaw:          message.Text, // Optional: stores original message
				}
				//	successMsg, err := b.TransactionService.CreateTransaction(transactionDTO)
				//	if err != nil {
				//		log.Printf("CRITICAL ERROR: Failed to create transaction for GroupID %d: %v", message.Chat.ID, err)
				//		utils.ErrorLog(err.Error(), "CRITICAL ERROR: Failed")
				//		continue
				//	}

					successMsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(
						"✅ Transaction Recorded!\nAmount: %s %.2f\nSender: %s\nDate: %s",
						savedTrx.Currency, savedTrx.Amount, savedTrx.Sender, savedTrx.TransactionDate,
					))

				utils.InfoLog(successMsg, "Transaction Recorded")
				continue
			}
		}
	*/

}

func containsKeyword(text, keyword string) bool {
	return strings.Contains(text, keyword)
}

func (b *Bot) sendGroupReplyKeyboard(botAPI *tgbotapi.BotAPI, chatID int64, messageID int) {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📊 របាយការណ៍"),
			tgbotapi.NewKeyboardButton("🏛 ធនាគារ"),
			tgbotapi.NewKeyboardButton("⚙️ ចុះឈ្មោះ"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🌐 MiniApp"),
			tgbotapi.NewKeyboardButton("👤 ទំនាក់ទំនង"),
		),
	)
	keyboard.ResizeKeyboard = true  // Make the keyboard smaller
	keyboard.OneTimeKeyboard = true // Hide after one click (optional)

	msg := tgbotapi.NewMessage(chatID, "សូមជ្រើរើស មីនុយ:")
	msg.ReplyMarkup = keyboard
	msg.ReplyToMessageID = messageID // Reply to a specific message if appropriate

	if _, err := botAPI.Send(msg); err != nil {
		log.Printf("Failed to send reply keyboard in group %d: %v", chatID, err)
	}
}
