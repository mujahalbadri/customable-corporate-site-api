package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "your_user"),
			Password: getEnv("DB_PASSWORD", "your_password"),
			DBName:   getEnv("DB_NAME", "your_db_name"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "your_jwt_secret_key"),
			ExpiresIn: 24 * time.Hour,
		},
	}

	// Validate critical configurations
	if config.JWT.Secret == "your_jwt_secret_key" {
		log.Println("Warning: Using default JWT secret key. Please set JWT_SECRET in environment variables for better security.")
	}

	if config.Server.Mode != "development" && config.Server.Mode != "production" {
		log.Fatalf("Invalid SERVER_MODE: %s. Must be 'development' or 'production'.", config.Server.Mode)
	}

	if config.Database.Host == "your_db_host" || config.Database.User == "your_user" || config.Database.DBName == "your_db_name" {
		log.Fatal("Database configuration is incomplete.")
	}

	// Log loaded configuration (excluding sensitive info)
	log.Printf("Configuration loaded: Server Mode=%s", config.Server.Mode)
	log.Printf("Server will start on port: %s", config.Server.Port)
	log.Printf("Database Host: %s, Port: %s, User: %s, DBName: %s, SSLMode: %s",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.DBName, config.Database.SSLMode)

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
