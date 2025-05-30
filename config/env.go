package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Configuration holds all the application config values
type Configuration struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	JWT         JWTConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	User     string
	Password string
	Server   string
	Port     string
	DBName   string
}

type JWTConfig struct {
	Secret string
	Expire int
}

var Config Configuration

func ConfigLoad() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	} else {
		log.Println("Successfully loaded .env file")
	}

	// Load configuration from environment variables
	Config = Configuration{
		Environment: getEnv("ENVIRONMENT", "DEV"),
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", ":8000"),
		},
		Database: DatabaseConfig{
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			Server:   getEnv("POSTGRES_SERVER", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			DBName:   getEnv("POSTGRES_DB", "ocrolus"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "ocrolus-secret-key"),
			Expire: getEnvAsInt("JWT_EXPIRE", 24),
		},
	}

	// Log loaded configuration for debugging
	logConfigValues()
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value := defaultValue
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		log.Printf("Warning: Invalid value for %s, using default: %v", key, err)
		return defaultValue
	}

	return value
}

// logConfigValues logs the loaded configuration for debugging
func logConfigValues() {
	log.Printf("Environment: %s", Config.Environment)
	log.Printf("Server Port: %s", Config.Server.Port)
	log.Printf("Database: Server=%s, User=%s, Port=%s, DB=%s",
		Config.Database.Server,
		Config.Database.User,
		Config.Database.Port,
		Config.Database.DBName,
	)
}
