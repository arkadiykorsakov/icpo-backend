package config

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	JWTSecret   string
	AccessTTL   time.Duration
	RefreshTTL  time.Duration

	CookieSecure       bool
	CookieSameSite     http.SameSite
	CookieDomain       string
	CORSAllowedOrigins []string
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

	appEnv := getEnv("APP_ENV", "production")
	secureDefault := strconv.FormatBool(appEnv != "development")
	if cfg.CookieSecure, err = strconv.ParseBool(getEnv("COOKIE_SECURE", secureDefault)); err != nil {
		return nil, fmt.Errorf("invalid COOKIE_SECURE: %w", err)
	}

	if cfg.CookieSameSite, err = parseSameSite(getEnv("COOKIE_SAMESITE", "strict")); err != nil {
		return nil, err
	}
	if cfg.CookieSameSite == http.SameSiteNoneMode && !cfg.CookieSecure {
		return nil, fmt.Errorf("COOKIE_SAMESITE=none requires COOKIE_SECURE=true")
	}

	cfg.CookieDomain = os.Getenv("COOKIE_DOMAIN")
	cfg.CORSAllowedOrigins = splitList(os.Getenv("CORS_ALLOWED_ORIGINS"))

	return cfg, nil
}

// LoadDatabaseURL возвращает только строку подключения к БД, не требуя остальных
// переменных (JWT_SECRET и др.). Нужен вспомогательным командам вроде cmd/seed.
func LoadDatabaseURL() string {
	_ = godotenv.Load()
	return buildDSN()
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

func parseSameSite(v string) (http.SameSite, error) {
	switch strings.ToLower(v) {
	case "strict":
		return http.SameSiteStrictMode, nil
	case "lax":
		return http.SameSiteLaxMode, nil
	case "none":
		return http.SameSiteNoneMode, nil
	default:
		return 0, fmt.Errorf("invalid COOKIE_SAMESITE %q: want strict|lax|none", v)
	}
}

func splitList(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
