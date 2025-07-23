package config

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CustomLogger struct {
	logger.Interface
}

func (c *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &CustomLogger{c.Interface.LogMode(level)}
}

func (c *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if len(data) > 0 && data[len(data)-1] != nil && data[len(data)-1].(*gorm.DB).Error != nil {
		log.Printf("[ERROR] SQL: %v | Args: %v\n", msg, data)
	}
}

func (c *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if len(data) > 0 && data[len(data)-1] != nil && data[len(data)-1].(*gorm.DB).Error != nil {
		log.Printf("[WARN] SQL: %v | Args: %v\n", msg, data)
	}
}

func (c *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if len(data) > 0 && data[len(data)-1] != nil && data[len(data)-1].(*gorm.DB).Error != nil {
		log.Printf("[ERROR] SQL: %v | Args: %v\n", msg, data)
	}
}

func (c *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// Only log errors
	if err != nil {
		log.Printf("[ERROR] SQL Execution Error: %v\n", err)
	}
}
