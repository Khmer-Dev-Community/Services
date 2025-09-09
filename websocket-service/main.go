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
	"websocket-service/rabbitmq"
	"websocket-service/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	// Assuming InitConfig and config package are correctly defined elsewhere
	cfg, err := InitConfig("config/config.yml")
	if err != nil {
		logger.Fatalf("Initialization error: %v", err)
	}
	utils.InitializeLogger(cfg.Service.LogPtah)
	logger.Print("Initialization Database")

	// Initialize RabbitMQ first.
	if err := rabbitmq.InitializeRabbitMQ(cfg.RabbitMQURL); err != nil {
		logger.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitmq.RMQ.Close() // Defer the close after a successful connection

	// Create the Socket.IO server.
	socketIOServer := gateway.CreateServer()
	go func() {
		if err := socketIOServer.Serve(); err != nil {
			logger.Fatalf("Socket.IO server error: %s\n", err)
		}
	}()
	defer socketIOServer.Close()

	// Start the RabbitMQ consumer as a goroutine.
	// The gateway's consumer function will use the initialized rabbitmq.RMQ object.
	go func() {
		if err := gateway.StartRabbitMQConsumer(socketIOServer); err != nil {
			logger.Fatalf("Failed to start RabbitMQ consumer: %v", err)
		}
	}()

	mux := http.NewServeMux()
	// Handle Socket.IO connections on both /socket.io/ and /ws/ paths
	mux.Handle("/socket.io/", socketIOServer)
	mux.Handle("/ws/", socketIOServer)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	httpServer := &http.Server{
		Addr:    ":" + cfg.Service.Port,
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("      \n                   * \n.         * *A* *\n.        *A* **=** *A*\n        *\"\"\"* *|\"\"\"|* *\"\"\"*\n       *|***|* *|*+*|* *|***|*\n*********\"\"\"*___*//+\\\\*___*\"\"\"*********\n@@@@@@@@@@@@@@@@//   \\\\@@@@@@@@@@@@@@@@@\n###############||ព្រះពុទ្ធ||#################\nTTTTTTTTTTTTTTT||ព្រះធម័||TTTTTTTTTTTTTTTTT\nLLLLLLLLLLLLLL//ព្រះសង្ឃ\\\\LLLLLLLLLLLLLLLLL\n៚ សូមប្រោសប្រទានពរឱ្យប្រតិប័ត្តិការណ៍ប្រព្រឹត្តទៅជាធម្មតា ៚ \n៚ ជោគជ័យ   //  ៚សិរីសួរស្តី \\\\   ៚សុវត្តិភាព \n___________//___៚(♨️)៚__\\\\____________\n៚Application Service is Running Port: 80 ")
		log.Println("Application Service is Running Port: " + cfg.Service.Port)
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
func InitConfig(configPath string) (config.Config, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
