package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	JWTSecret   string
	AccessTTL   time.Duration
	RefreshTTL  time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL: buildDSN(),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	var err error
	if cfg.AccessTTL, err = parseDuration("JWT_ACCESS_TTL", "15m"); err != nil {
		return nil, err
	}
	if cfg.RefreshTTL, err = parseDuration("JWT_REFRESH_TTL", "720h"); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// buildDSN собирает строку подключения. Логин и пароль экранируются через
// url.UserPassword, поэтому спецсимволы (@ : / # %) в пароле не ломают DSN.
func buildDSN() string {
	u := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			getEnv("POSTGRES_USER", "postgres"),
			getEnv("POSTGRES_PASSWORD", ""),
		),
		Host: net.JoinHostPort(
			getEnv("POSTGRES_HOST", "localhost"),
			getEnv("POSTGRES_PORT", "5432"),
		),
		Path:     getEnv("POSTGRES_DB", "postgres"),
		RawQuery: url.Values{"sslmode": {getEnv("POSTGRES_SSLMODE", "disable")}}.Encode(),
	}
	return u.String()
}

func parseDuration(key, def string) (time.Duration, error) {
	d, err := time.ParseDuration(getEnv(key, def))
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	if d <= 0 {
		return 0, fmt.Errorf("invalid %s: must be positive", key)
	}
	return d, nil
}
