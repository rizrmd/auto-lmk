package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	LLM      LLMConfig
	WhatsApp WhatsAppConfig
	Security SecurityConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type LLMConfig struct {
	Provider    string
	APIKey      string
	Model       string
	ZAIEndpoint string // Z.AI specific endpoint
}

type WhatsAppConfig struct {
	SessionPath string
}

type SecurityConfig struct {
	JWTSecret string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "autolmk"),
			Password: getEnv("DB_PASSWORD", "autolmk_dev_password"),
			Name:     getEnv("DB_NAME", "autolmk"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		LLM: LLMConfig{
			Provider:    getEnv("LLM_PROVIDER", ""),
			APIKey:      getEnv("LLM_API_KEY", ""),
			Model:       getEnv("LLM_MODEL", ""),
			ZAIEndpoint: getEnv("ZAI_ENDPOINT", "https://api.z.ai/api/coding/paas/v4"),
		},
		WhatsApp: WhatsAppConfig{
			SessionPath: getEnv("WHATSAPP_SESSION_PATH", "./whatsapp_sessions"),
		},
		Security: SecurityConfig{
			JWTSecret: getEnv("JWT_SECRET", "change-this-in-production"),
		},
	}

	return cfg, nil
}

// DatabaseURL returns PostgreSQL connection string
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
