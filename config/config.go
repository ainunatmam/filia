package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Database configuration
	Database DatabaseConfig

	// Application configuration
	App AppConfig

	// JWT configuration
	JWT JWTConfig

	// Rate limiting configuration
	RateLimit RateLimitConfig
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string

	// Connection pool settings
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration

	// Query settings
	QueryTimeout time.Duration
}

// AppConfig holds application server configuration
type AppConfig struct {
	Port string
	Env  string // development, staging, production
}

// JWTConfig holds JWT-specific configuration
type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Max        int
	Expiration int // in seconds
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (optional in production)
	_ = godotenv.Load()

	cfg := &Config{
		Database:  loadDatabaseConfig(),
		App:       loadAppConfig(),
		JWT:       loadJWTConfig(),
		RateLimit: loadRateLimitConfig(),
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Username:        getEnvOrDefault("DATABASE_USERNAME", "root"),
		Password:        getEnvOrDefault("DATABASE_PASSWORD", ""),
		Host:            getEnvOrDefault("DATABASE_HOST", "127.0.0.1"),
		Port:            getEnvOrDefault("DATABASE_PORT", "3306"),
		Name:            getEnvOrDefault("DATABASE_NAME", "wallet_api"),
		MaxOpenConns:    getEnvIntOrDefault("DB_MAX_OPEN_CONNS", 100),
		MaxIdleConns:    getEnvIntOrDefault("DB_MAX_IDLE_CONNS", 25),
		ConnMaxLifetime: getEnvDurationOrDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		ConnMaxIdleTime: getEnvDurationOrDefault("DB_CONN_MAX_IDLE_TIME", 1*time.Minute),
		QueryTimeout:    getEnvDurationOrDefault("DB_QUERY_TIMEOUT", 30*time.Second),
	}
}

func loadAppConfig() AppConfig {
	return AppConfig{
		Port: getEnvOrDefault("APP_PORT", ":3000"),
		Env:  getEnvOrDefault("APP_ENV", "development"),
	}
}

func loadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:          getEnvOrDefault("JWT_SECRET", "your-secret-key-change-in-production"),
		ExpirationHours: getEnvIntOrDefault("JWT_EXPIRATION_HOURS", 24),
	}
}

func loadRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Max:        getEnvIntOrDefault("RATE_LIMIT_MAX", 100),
		Expiration: getEnvIntOrDefault("RATE_LIMIT_EXPIRATION", 60),
	}
}

// Validate validates required configuration fields
func (c *Config) Validate() error {
	// Validate database configuration
	if c.Database.Username == "" {
		return fmt.Errorf("DATABASE_USERNAME is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DATABASE_HOST is required")
	}
	if c.Database.Port == "" {
		return fmt.Errorf("DATABASE_PORT is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DATABASE_NAME is required")
	}

	// Validate JWT configuration
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.JWT.Secret == "your-secret-key-change-in-production" {
		if c.App.Env == "production" {
			return fmt.Errorf("JWT_SECRET must be changed in production environment")
		}
		fmt.Println("WARNING: Using default JWT_SECRET. Change this in production!")
	}

	return nil
}

// GetDSN returns the MySQL Data Source Name
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}

// GetJWTExpiration returns the JWT expiration duration
func (c *JWTConfig) GetJWTExpiration() time.Duration {
	return time.Duration(c.ExpirationHours) * time.Hour
}

// Helper functions for environment variable loading

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		// Try parsing as seconds (numeric string)
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
		// Try parsing as duration string (e.g., "5m")
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
