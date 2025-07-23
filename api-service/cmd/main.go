package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	telegrambot "telegram-service/telegram/bot"
	modelInit "telegram-service/telegram/models"
	telegramRepo "telegram-service/telegram/repository"
	telegramService "telegram-service/telegram/services"
	"telegram-service/utils"
	// Ensure these are imported if your InitConfigAndDatabase relies on them
)

// Ensure these functions are defined either globally in main.go
// or imported from their respective packages.
// For example, if InitConfigAndDatabase is in config/config.go, you'd import config
// and call config.InitConfigAndDatabase().
//
// Assuming your InitConfigAndDatabase and StartHTTPServer are still defined elsewhere
// in main.go or correctly imported:
//
// func InitConfigAndDatabase() (*config.Config, *gorm.DB) { /* ... your implementation ... */ }
// func StartHTTPServer(port string, handler http.Handler) error { /* ... your implementation ... */ }

func main() {
	loginBotFlag := flag.Bool("login-bot", false, "Set to true to initiate interactive login for Telegram userbots.")
	loginPhoneFlag := flag.String("phone", "", "Specify phone number of the bot to login (requires --login-bot).")
	flag.Parse()

	currentWorkingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}
	log.Printf("Application started from working directory: %s", currentWorkingDir)

	cfg, db := InitConfigAndDatabase() // Assuming this is defined and returns config and db

	// Corrected logic for determining absoluteSessionDir
	projectRoot := currentWorkingDir
	if filepath.Base(currentWorkingDir) == "cmd" {
		projectRoot = filepath.Dir(currentWorkingDir)
	}
	absoluteSessionDir := filepath.Join(projectRoot, "sessions")
	if err := os.MkdirAll(absoluteSessionDir, 0755); err != nil {
		log.Fatalf("Failed to create session directory %s: %v", absoluteSessionDir, err)
	}
	log.Printf("Session files will be stored in: %s", absoluteSessionDir)

	for i := range cfg.Telegram.Account {
		acc := &cfg.Telegram.Account[i]
		if acc.SessionFile != "" {
			acc.SessionFile = filepath.Join(absoluteSessionDir, filepath.Base(acc.SessionFile))
			log.Printf("Userbot %s session file path resolved to: %s", acc.Name, acc.SessionFile)
		}
	}

	if *loginBotFlag {
		log.Println("Initiating Telegram userbot login mode (interactive)...")
		tgRepo := telegramRepo.NewTransactionRepository(db)
		tgService := telegramService.NewTransactionService(tgRepo)

		foundBotToLogin := false
		for _, acc := range cfg.Telegram.Account {
			if !acc.Status {
				continue
			}
			if *loginPhoneFlag != "" && acc.PhoneNumber != *loginPhoneFlag {
				continue
			}

			log.Printf("Attempting interactive login for userbot: %s (phone: %s)", acc.Name, acc.PhoneNumber)
			userBot, err := telegrambot.NewTgService(
				acc.Name,
				acc.AppID,
				acc.ApiHash,
				acc.PhoneNumber,
				acc.SessionFile,
				acc.Debug,
				true,
				tgService,
			)
			if err != nil {
				log.Printf("ERROR: Failed to create userbot service for %s: %v", acc.Name, err)
				continue
			}

			err = userBot.Start()
			if err != nil {
				log.Printf("ERROR: Userbot %s login failed: %v", acc.Name, err)
			} else {
				log.Printf("SUCCESS: Userbot %s logged in and session saved.", acc.Name)
			}
			foundBotToLogin = true

			if *loginPhoneFlag != "" {
				break
			}
		}

		if !foundBotToLogin {
			if *loginPhoneFlag != "" {
				log.Fatalf("No active userbot found with phone number: %s", *loginPhoneFlag)
			} else {
				log.Println("No active userbots found to login. Check your config.yaml and ensure status: true for accounts.")
			}
		}

		log.Println("Login process complete. Exiting.")
		return
	}

	log.Println("Starting normal application mode (non-interactive)...")
	tgRepo := telegramRepo.NewTransactionRepository(db)
	tgService := telegramService.NewTransactionService(tgRepo)
	modelInit.MigrateBotTable(db)

	activeUserBots := make(map[string]*telegrambot.BotAccount)
	var wg sync.WaitGroup

	if cfg.Telegram.Status {
		for _, acc := range cfg.Telegram.Account {
			if !acc.Status {
				log.Printf("Skipping userbot %s (status: false)", acc.Name)
				continue
			}
			log.Printf("Initializing Telegram userbot: %s", acc.Name)
			userBot, err := telegrambot.NewTgService(
				acc.Name,
				acc.AppID,
				acc.ApiHash,
				acc.PhoneNumber,
				acc.SessionFile,
				acc.Debug,
				false,
				tgService,
			)
			if err != nil {
				log.Printf("ERROR: Failed to create userbot service for %s: %v", acc.Name, err)
				continue
			}
			go userBot.Start()
			activeUserBots[acc.Name] = userBot
		}
	} else {
		log.Println("Telegram integration is disabled in config.")
	}
	log.Println("Attempting to initialize services...")
	wg.Wait()
	time.Sleep(2 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		utils.WarnLog("Received signal: %s", sig.String())
		cancel()
	}()

	// --- ADDED DIAGNOSTIC LOGS HERE ---

	services := InitServices(db, &cfg, activeUserBots) // Check if cfg is passed by reference here, previously it was &cfg
	log.Println("Services initialized. Attempting to initialize router...")
	router := InitRoutes(cfg, services)
	log.Println("Router initialized. Launching HTTP server goroutine...")
	// --- END ADDED DIAGNOSTIC LOGS ---

	go func() {
		log.Printf("HTTP server goroutine started. Binding to port %s...", cfg.Service.Port) // More specific log
		log.Println("      \n                   * \n.         * *A* *\n.        *A* **=** *A*\n        *\"\"\"* *|\"\"\"|* *\"\"\"*\n       *|***|* *|*+*|* *|***|*\n*********\"\"\"*___*//+\\\\*___*\"\"\"*********\n@@@@@@@@@@@@@@@@//   \\\\@@@@@@@@@@@@@@@@@\n###############||ព្រះពុទ្ធ||#################\nTTTTTTTTTTTTTTT||ព្រះធម័||TTTTTTTTTTTTTTTTT\nLLLLLLLLLLLLLL//ព្រះសង្ឃ\\\\LLLLLLLLLLLLLLLLL\n៚ សូមប្រោសប្រទានពរឱ្យប្រតិប័ត្តិការណ៍ប្រព្រឹត្តទៅជាធម្មតា ៚ \n៚ ជោគជ័យ   //  ៚សិរីសួរស្តី \\\\   ៚សុវត្តិភាព \n___________//___៚(♨️)៚__\\\\____________\n៚Application Service is Running Port: 80 ")
		if err := StartHTTPServer(cfg.Service.Port, router.(http.Handler)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
		log.Printf("HTTP server successfully started on port %s", cfg.Service.Port) // This log on successful start
	}()

	<-ctx.Done()
	utils.InfoLog("Shutting down service...", "")
}
