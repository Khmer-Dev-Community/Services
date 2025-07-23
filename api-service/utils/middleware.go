package utils

import (
	"fmt"
	"log"

	"github.com/Khmer-Dev-Community/Services/api-service/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CORSMiddleware adds CORS headers to the response
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func EncryptPassword(password string) (string, error) {

	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	secretKey := []byte(cfg.SecretKey)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+string(secretKey)), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hashed password: %v", err)
	}

	return string(hashedPassword), nil
}

func ComparePasswordss(hashedPassword, password string) error {
	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	secretKey := []byte(cfg.SecretKey)

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+string(secretKey)))
}
func ComparePassword(hashedPassword, password string) error {

	cfg, err := config.LoadConfig("config/config.yml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Print(cfg.SecretKey)
	secretKey := []byte(cfg.SecretKey)
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+string(secretKey)))
}
