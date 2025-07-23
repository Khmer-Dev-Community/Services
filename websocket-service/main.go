package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"websocket-service/config"
	"websocket-service/gateway"
	"websocket-service/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	// Assuming InitConfig and config package are correctly defined elsewhere
	cfg, err := InitConfig("config/config.yml")
	logger.Print("Initialization Database") // This log seems misplaced, should be after successful DB init
	if err != nil {
		logger.Fatalf("Initialization error: %v", err)
	}
	utils.InitializeLogger(cfg.Service.LogPtah)
	socketIOServer := gateway.CreateServer()
	go func() {
		if err := socketIOServer.Serve(); err != nil {
			logger.Fatalf("Socket.IO server error: %s\n", err)
		}
	}()
	defer socketIOServer.Close()

	mux := http.NewServeMux()
	// Handle Socket.IO connections on both /socket.io/ and /ws/ paths
	mux.Handle("/socket.io/", socketIOServer)
	mux.Handle("/ws/", socketIOServer)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	httpServer := &http.Server{
		Addr:    ":" + cfg.Service.Port, // Dynamically use port from config
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("      \n                   *       \n.         *        *        *\n.        ***     **=**     ***\n        *\"\"\"*   *|***|*   *\"\"\"*\n       *|***|*  *|*+*|*  *|***|*\n**********\"\"\"*___*//+\\\\*___*\"\"\"*********\n@@@@@@@@@@@@@@@@//   \\\\@@@@@@@@@@@@@@@@@\n###############||ព្រះពុទ្ធ||#################\nTTTTTTTTTTTTTTT||ព្រះធម័||TTTTTTTTTTTTTTTTT\n------------- -//ព្រះសង្ឃ\\\\----------------\n៚ សូមប្រោសប្រទានពរឱ្យប្រតិប័ត្តិការណ៍ប្រព្រឹត្តទៅជាធម្មតា ៚ \n ៚ ជោគជ័យ      ៚សិរីសួរស្តី       ៚សុវត្តិភាព \n_________________________________________\n៚  Application Service is Running Port: " + cfg.Service.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not listen on %s: %v\n", httpServer.Addr, err)
		}
	}()

	<-stop
	logger.Info("Shutting down the server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatalf("HTTP server shutdown error: %v", err)
	}
	logger.Info("Server stopped.")
}

// InitConfig is assumed to be defined in your config package
// and handles loading configuration from a YAML file.
func InitConfig(configPath string) (config.Config, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
