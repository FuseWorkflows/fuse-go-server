package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration settings
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTKey     string
	Port       string
	AIService  string
}

// NewConfig loads configuration settings from environment variables
func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		JWTKey:     os.Getenv("JWT_KEY"),
		Port:       os.Getenv("PORT"),
		AIService:  os.Getenv("AI_SERVICE"),
	}

	// Validate required environment variables
	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" || cfg.JWTKey == "" || cfg.Port == "" || cfg.AIService == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	// Parse port as integer
	portInt, err := strconv.Atoi(cfg.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid port value: %w", err)
	}
	cfg.Port = strconv.Itoa(portInt)

	return cfg, nil
}
