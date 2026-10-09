package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL: buildDSN(),
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func buildDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		getEnv("POSTGRES_USER", "postgres"),
		getEnv("POSTGRES_PASSWORD", ""),
		getEnv("POSTGRES_HOST", "localhost"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_DB", "postgres"),
		getEnv("POSTGRES_SSLMODE", "disable"),
	)
}
