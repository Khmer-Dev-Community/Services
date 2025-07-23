package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Khmer-Dev-Community/Services/api-service/config"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"

	"gorm.io/gorm"
)

var whitelist = map[string]bool{

	"/api/auth/login":             true,
	"/api/auth/logout":            true,
	"/api/auth02/register":        true,
	"/api/auth02/login":           true,
	"/api/auth02/github/login":    true,
	"/api/auth02/github/callback": true,
	"/api/auth02/profile":         true,

	"/api/swagger/index.html":                  true,
	"/swagger/index.html":                      true,
	"/swagger/swagger-ui-bundle.js":            true,
	"/swagger/swagger-ui.css":                  true,
	"/swagger/swagger-ui-standalone-preset.js": true,
	"/swagger/doc.json":                        true,
	"/swagger/favicon-32x32.png":               true,
	"/swagger/favicon-16x16.png":               true,
}

func InitConfigAndDatabase() (config.Config, *gorm.DB) {
	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db := config.InitDatabase("config/config.yml")

	if err := config.InitRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	utils.InitializeLogger(cfg.Service.LogPtah)

	loc, err := time.LoadLocation(cfg.Service.TimeZone)
	if err == nil {
		now := time.Now().In(loc)
		fmt.Printf("Current Time in %s: %s\n", cfg.Service.TimeZone, now)
		fmt.Printf("UTC Time: %s\n", time.Now().UTC())
	}

	return cfg, db
}
