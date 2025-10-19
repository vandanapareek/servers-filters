package config

import (
	"fmt"
	"os"
	"strconv"
)

// Application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Cache    CacheConfig    `json:"cache"`
	Log      LogConfig      `json:"log"`
}

// Server configuration
type ServerConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// Database configuration
type DatabaseConfig struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}

// Cache configuration
type CacheConfig struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	URL     string `json:"url"`
	TTL     int    `json:"ttl"`
}

// Log configuration
type LogConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// load configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 30),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 30),
		},
		Database: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "sqlite3"),
			DSN:    getEnv("DB_DSN", "data/servers.db"),
		},
		Cache: CacheConfig{
			Enabled: getEnvAsBool("CACHE_ENABLED", false),
			Type:    getEnv("CACHE_TYPE", "redis"),
			URL:     getEnv("CACHE_URL", "redis://localhost:6379"),
			TTL:     getEnvAsInt("CACHE_TTL", 300), // 5 minutes
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	return config, nil
}

// To get server address
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// To get environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// To get environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// To get environment variable as boolean with a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
