package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DSN        string

	// Auth
	JWTSecret   string
	AllowSignup bool

	// App
	AppEnv         string
	LogLevel       string
	MaxRequestSize int // MB

	// CORS
	CORSAllowedOrigins string
	CORSAllowedMethods string
	CORSAllowedHeaders string

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   int // seconds
}

var (
	instance *Config
	once     sync.Once
)

// Load returns singleton config instance
func Load() *Config {
	once.Do(func() {
		// Try to load .env file
		if err := godotenv.Load(); err != nil {
			fmt.Println("ℹ️ File .env not found, using system environment variables")
		}

		cfg := &Config{
			// Database
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBUser:     getEnv("DB_USER", "postgres"),
			DBPassword: getEnv("DB_PASSWORD", ""),
			DBName:     getEnv("DB_NAME", "gorders"),
			DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

			// Auth
			JWTSecret:   getEnv("JWT_SECRET", ""),
			AllowSignup: getEnv("ALLOWSIGNUP", "false") == "true",

			// App
			AppEnv:         getEnv("APP_ENV", "development"),
			LogLevel:       getEnv("LOG_LEVEL", "info"),
			MaxRequestSize: getEnvInt("MAX_REQUEST_SIZE", 10),

			// CORS
			CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
			CORSAllowedMethods: getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			CORSAllowedHeaders: getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization"),

			// Rate Limiting
			RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:   getEnvInt("RATE_LIMIT_WINDOW", 60),
		}

		// Build DSN
		cfg.DSN = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
		)

		// Validate critical fields
		if err := cfg.Validate(); err != nil {
			fmt.Printf("❌ Config validation failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✅ Config loaded and validated")
		instance = cfg
	})

	return instance
}

// getEnv returns environment variable or default
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getEnvInt returns environment variable as int or default
func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := parseInt(value); err == nil {
			return i
		}
	}
	return defaultVal
}

// parseInt converts string to int
func parseInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// Validate checks required configuration
func (c *Config) Validate() error {
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DBPort == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required and must be at least 32 characters")
	}
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters for security")
	}
	return nil
}
