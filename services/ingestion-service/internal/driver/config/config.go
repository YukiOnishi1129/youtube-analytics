package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	// HTTP server configuration
	HTTPPort string
	
	// Database configuration
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	
	// Service configuration
	Environment string
	LogLevel    string
	
	// External services
	YouTubeAPIKey string
	
	// WebSub configuration
	WebSubCallbackURL string
	WebSubSecret      string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		// HTTP server
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		
		// Database
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", ""),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", ""),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		
		// Service
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		
		// External services
		YouTubeAPIKey: getEnv("YOUTUBE_API_KEY", ""),
		
		// WebSub
		WebSubCallbackURL: getEnv("WEBSUB_CALLBACK_URL", ""),
		WebSubSecret:      getEnv("WEBSUB_SECRET", ""),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as int with a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool gets an environment variable as bool with a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}