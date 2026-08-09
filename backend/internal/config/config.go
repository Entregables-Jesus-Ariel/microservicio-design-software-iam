// Package config loads backend runtime settings from the environment.
// Secret-bearing values are required: the process refuses to start with
// a default credential rather than running with a guessable one.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds every setting the backend needs to run.
type Config struct {
	DatabaseHost     string
	DatabasePort     int
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseSSLMode  string
	TokenSecret      string
	TokenTTL         time.Duration
	HTTPPort         int
	CORSOrigin       string
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	QueryTimeout     time.Duration
}

// Load reads configuration from the environment and validates it.
func Load() (Config, error) {
	host, err := requiredEnv("POSTGRES_HOST")
	if err != nil {
		return Config{}, err
	}
	name, err := requiredEnv("POSTGRES_DB")
	if err != nil {
		return Config{}, err
	}
	user, err := requiredEnv("POSTGRES_USER")
	if err != nil {
		return Config{}, err
	}
	password, err := requiredEnv("POSTGRES_PASSWORD")
	if err != nil {
		return Config{}, err
	}
	secret, err := requiredEnv("APP_TOKEN_SECRET")
	if err != nil {
		return Config{}, err
	}

	return Config{
		DatabaseHost:     host,
		DatabasePort:     intEnv("POSTGRES_PORT", 5432),
		DatabaseName:     name,
		DatabaseUser:     user,
		DatabasePassword: password,
		DatabaseSSLMode:  stringEnv("POSTGRES_SSLMODE", "disable"),
		TokenSecret:      secret,
		TokenTTL:         time.Duration(intEnv("APP_TOKEN_TTL_MINUTES", 15)) * time.Minute,
		HTTPPort:         intEnv("APP_HTTP_PORT", 8081),
		CORSOrigin:       stringEnv("APP_CORS_ORIGIN", "http://localhost:5173"),
		ReadTimeout:      time.Duration(intEnv("APP_READ_TIMEOUT_SECONDS", 10)) * time.Second,
		WriteTimeout:     time.Duration(intEnv("APP_WRITE_TIMEOUT_SECONDS", 30)) * time.Second,
		QueryTimeout:     time.Duration(intEnv("APP_DB_QUERY_TIMEOUT_SECONDS", 5)) * time.Second,
	}, nil
}

// ConnString builds the Postgres DSN used by the pgx driver.
func (c Config) ConnString() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.DatabaseHost, c.DatabasePort, c.DatabaseName, c.DatabaseUser, c.DatabasePassword, c.DatabaseSSLMode,
	)
}

func requiredEnv(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", name)
	}
	return value, nil
}

func stringEnv(name string, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func intEnv(name string, fallback int) int {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
