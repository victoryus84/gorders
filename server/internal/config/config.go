package config

import (
	"fmt"
	"os"
	"sync"
	"strings"

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
	Commit         string
	Version        string
	LogLevel       string
	MaxRequestSize int // MB

	// CORS
	CORSAllowedOrigins string
	CORSAllowedMethods string
	CORSAllowedHeaders string

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   int // seconds
	
	// Kafka
	KafkaAddr         string
    KafkaTopicPattern string
    
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

		// Creăm un helper local pentru curățare
        cleanEnv := func(key, fallback string) string {
            val := getEnv(key, fallback)
            return strings.TrimSpace(val) // Aici moare "huieta"
        }
		
		cfg := &Config{
			// Database
			DBHost:     cleanEnv("DB_HOST", "localhost"),
			DBPort:     cleanEnv("DB_PORT", "5432"),
			DBUser:     cleanEnv("DB_USER", "postgres"),
			DBPassword: cleanEnv("DB_PASSWORD", ""),
			DBName:     cleanEnv("DB_NAME", "gorders"),
			DBSSLMode:  cleanEnv("DB_SSLMODE", "disable"),

			// Auth
			JWTSecret:   cleanEnv("JWT_SECRET", ""),
			AllowSignup: getEnv("ALLOWSIGNUP", "false") == "true",

			// App
			AppEnv:         cleanEnv("APP_ENV", "development"),
			Commit:         cleanEnv("GIT_COMMIT", "unknown"),
			Version:        cleanEnv("VERSION", "1.0.0"),
			LogLevel:       getEnv("LOG_LEVEL", "info"),
			MaxRequestSize: getEnvInt("MAX_REQUEST_SIZE", 10),

		 	// Kafka
    		KafkaAddr:  	cleanEnv("KAFKA_BROKERS", ""),
    		KafkaTopicPattern: cleanEnv("KAFKA_TOPIC_PATTERN", ""),

			// CORS
			CORSAllowedOrigins: cleanEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
			CORSAllowedMethods: cleanEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			CORSAllowedHeaders: cleanEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization"),

			// Rate Limiting
			RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:   getEnvInt("RATE_LIMIT_WINDOW", 60),
		}

		// Build DSN
		cfg.DSN = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
		)

		// DEBUG: Scoate asta după ce te convingi că merge
        //fmt.Printf("🔍 DSN final: [%s]\n", cfg.DSN)

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
	if c.KafkaAddr == "" {
        return fmt.Errorf("KAFKA_BROKERS is required")
    }
	if c.KafkaTopicPattern == "" {
        return fmt.Errorf("KAFKA_TOPIC_PATTERN is missing in .env")
    }
    // Verificăm dacă ai pus %s în pattern
    if !strings.Contains(c.KafkaTopicPattern, "%s") {
        return fmt.Errorf("KAFKA_TOPIC_PATTERN must contain '%%s' (e.g., gorders.%%s.events)")
    }
	return nil
}

func (c *Config) GetTopic(domain string) string {
    if domain == "" {
        // Aici poți să faci log.Fatal sau să returnezi un topic de tip "garbage/dead-letter"
        return "internal.unknown.events" 
    }
    return fmt.Sprintf(c.KafkaTopicPattern, domain)
}