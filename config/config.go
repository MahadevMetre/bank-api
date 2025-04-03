package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}

// Default server configuration values
const (
	DefaultReadTimeout    = 200
	DefaultWriteTimeout   = 200
	DefaultIdleTimeout    = 240
	DefaultMaxHeaderBytes = 65536 // 64 KB
)

// GetServerConfig returns the server configuration from environment variables
// with fallback to default values
func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		ReadTimeout:    time.Duration(getEnvInt("SERVER_READ_TIMEOUT", DefaultReadTimeout)) * time.Second,
		WriteTimeout:   time.Duration(getEnvInt("SERVER_WRITE_TIMEOUT", DefaultWriteTimeout)) * time.Second,
		IdleTimeout:    time.Duration(getEnvInt("SERVER_IDLE_TIMEOUT", DefaultIdleTimeout)) * time.Second,
		MaxHeaderBytes: getEnvInt("HEADER_MAX_HEADER_BYTES", DefaultMaxHeaderBytes),
	}
}

// getEnvInt safely retrieves an integer value from environment variables
// with a fallback default value
func getEnvInt(key string, defaultValue int) int {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		fmt.Printf("Warning: Invalid value for %s, using default: %d\n", key, defaultValue)
		return defaultValue
	}

	return value
}

// String returns a string representation of the server config
func (c *ServerConfig) String() string {
	return fmt.Sprintf(
		"ServerConfig{ReadTimeout: %v, WriteTimeout: %v, IdleTimeout: %v, MaxHeaderBytes: %d}",
		c.ReadTimeout,
		c.WriteTimeout,
		c.IdleTimeout,
		c.MaxHeaderBytes,
	)
}
