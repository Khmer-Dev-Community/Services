package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Khmer-Dev-Community/Services/api-service/config"
	"github.com/Khmer-Dev-Community/Services/api-service/delivery/rabbitmq"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"

	"gorm.io/gorm"
)

func InitConfigAndDatabase() (config.Config, *gorm.DB) {
	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db := config.InitDatabase("config/config.yml")

	if err := config.InitRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	// Initialize RabbitMQ first.
	if err := rabbitmq.InitializeRabbitMQ(cfg.RabbitMQURL); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitmq.RMQ.Close() // Defer the close after a successful connection
	utils.InitializeLogger(cfg.Service.LogPtah)

	loc, err := time.LoadLocation(cfg.Service.TimeZone)
	if err == nil {
		now := time.Now().In(loc)
		fmt.Printf("Current Time in %s: %s\n", cfg.Service.TimeZone, now)
		fmt.Printf("UTC Time: %s\n", time.Now().UTC())
	}

	return cfg, db
}
