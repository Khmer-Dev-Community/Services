package main

import (
	"fmt"
	"log"
	"time"

	"telegram-service/config"
	"telegram-service/utils"

	"gorm.io/gorm"
)

var whitelist = map[string]bool{
	"/forward":                                 true,
	"/forward/":                                true,
	"/api/auth/login":                          true,
	"/api/auth/logout":                         true,
	"/api/video/videosnap":                     true,
	"/api/video/videostop":                     true,
	"/api/metric/query":                        true,
	"/api/metric/exportquery":                  true,
	"/api/metric/namespace":                    true,
	"/api/metric/namespace/sub":                true,
	"/api/v1/items/list":                       true,
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
