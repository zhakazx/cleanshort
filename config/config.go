package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// App settings
	Environment string
	Port        string
	BaseURL     string

	// Database
	DatabaseDSN string

	// JWT settings
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration

	// Rate limiting
	RateLimitAuth     int
	RateLimitRedirect int

	// Logging
	LogLevel string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		Environment:       getEnv("APP_ENV", "development"),
		Port:              getEnv("APP_PORT", "8080"),
		BaseURL:           getEnv("APP_BASE_URL", "http://localhost:8080"),
		DatabaseDSN:       getEnv("DB_DSN", "postgres://postgres:@localhost:5432/shortener?sslmode=disable"),
		JWTSecret:         getEnv("JWT_SECRET", "super-secret-change-in-production"),
		JWTAccessTTL:      parseDuration(getEnv("JWT_ACCESS_TTL", "15m")),
		JWTRefreshTTL:     parseDuration(getEnv("JWT_REFRESH_TTL", "168h")),
		RateLimitAuth:     parseInt(getEnv("RATE_LIMIT_AUTH", "5")),
		RateLimitRedirect: parseInt(getEnv("RATE_LIMIT_REDIRECT", "200")),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
	}

	// Validate required fields
	if cfg.JWTSecret == "super-secret-change-in-production" && cfg.Environment == "production" {
		log.Fatal("JWT_SECRET must be set in production")
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Invalid integer value: %s", s)
	}
	return i
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Fatalf("Invalid duration value: %s", s)
	}
	return d
}