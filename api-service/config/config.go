package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config represents the overall configuration structure
type Config struct {
	Database    DatabaseConfig `yaml:"database"`
	SecretKey   string         `yaml:"secret_key"`
	RabbitMQURL string         `yaml:"RabbitMQURL"`
	Redis       ConfigRedis    `yaml:"redis"`
	Service     ServiceConfig  `yaml:"service"`
	//Telegram    TelegramConfig `yaml:"telegram"`
	Github GitConfig `yaml:"github"`
}

// telegram
type BotConfig struct {
	Name    string `yaml:"name"`
	Status  bool   `yaml:"status"`
	Debug   bool   `yaml:"debug"`
	Token   string `yaml:"token"`
	GroupID int64  `yaml:"groupId"`
	BackBtn bool   `yaml:"back_button"`
}
type TGAccountConfig struct {
	Name        string `yaml:"name"`
	Status      bool   `yaml:"status"`
	ApiHash     string `yaml:"apiHash"`
	AppID       int    `yaml:"appID"`
	PhoneNumber string `yaml:"phoneNumber"`
	SessionFile string `yaml:"sessionFile"`
	Debug       bool   `yaml:"debug"`
}

type TelegramConfig struct {
	Status  bool              `yaml:"status"`
	Bots    []BotConfig       `yaml:"bots"`
	Account []TGAccountConfig `yaml:"account"`
}

// Server
type ServiceConfig struct {
	Port     string `yaml:"port"`
	LogPtah  string `yaml:"logPtah"`
	TimeZone string `yaml:"timeZone"`
}

// Redis
type ConfigRedis struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
}

// DatabaseConfig represents database connection parameters
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Port     string `yaml:"port"`
	Driver   string `yaml:"driver"`
}

// LoadConfig loads the configuration from the YAML file
func LoadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// Connect establishes a connection to the database using the provided configuration
func Connect(cfg DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch cfg.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.Host, cfg.Username, cfg.Password, cfg.Name, cfg.Port)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	if err != nil {
		return nil, err
	}
	return db, nil
}

// InitDatabase initializes the database connection
func InitDatabase(filename string) *gorm.DB {
	config, err := LoadConfig(filename)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	db, err := Connect(config.Database)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	fmt.Println("Database connected")
	return db
}
