package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	TokenEncKey        []byte
	SessionSecret      []byte
	Port               string
	SyncIntervalMin    int
	AppEnv             string
	LogLevel           string
	LogFormat          string
	DBPath             string
	FrontendURL        string
	GmailUnitsPerMin   int
	GmailConcurrency   int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		GoogleRedirectURL: envOr("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/callback"),
		Port:              envOr("PORT", "8080"),
		SyncIntervalMin:   envIntOr("SYNC_INTERVAL_MINUTES", 5),
		AppEnv:            envOr("APP_ENV", "dev"),
		LogLevel:          envOr("LOG_LEVEL", "info"),
		LogFormat:         envOr("LOG_FORMAT", "json"),
		DBPath:            envOr("DB_PATH", "./data/jobtrack.db"),
		FrontendURL:       envOr("FRONTEND_URL", defaultFrontendURL(envOr("APP_ENV", "dev"))),
		GmailUnitsPerMin:  envIntOr("GMAIL_UNITS_PER_MIN", 2500),
		GmailConcurrency:  envIntOr("GMAIL_CONCURRENCY", 3),
	}

	var missing []string
	cfg.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	if cfg.GoogleClientID == "" {
		missing = append(missing, "GOOGLE_CLIENT_ID")
	}
	cfg.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	if cfg.GoogleClientSecret == "" {
		missing = append(missing, "GOOGLE_CLIENT_SECRET")
	}

	tokenKeyB64 := os.Getenv("TOKEN_ENC_KEY")
	if tokenKeyB64 == "" {
		missing = append(missing, "TOKEN_ENC_KEY")
	} else {
		key, err := decodeKey32("TOKEN_ENC_KEY", tokenKeyB64)
		if err != nil {
			return nil, err
		}
		cfg.TokenEncKey = key
	}

	sessionB64 := os.Getenv("SESSION_SECRET")
	if sessionB64 == "" {
		missing = append(missing, "SESSION_SECRET")
	} else {
		secret, err := base64.StdEncoding.DecodeString(sessionB64)
		if err != nil {
			return nil, fmt.Errorf("SESSION_SECRET must be valid base64: %w", err)
		}
		if len(secret) < 32 {
			return nil, fmt.Errorf("SESSION_SECRET must decode to at least 32 bytes, got %d", len(secret))
		}
		cfg.SessionSecret = secret
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s (copy .env.example to .env)", strings.Join(missing, ", "))
	}

	if cfg.AppEnv != "dev" && cfg.AppEnv != "prod" {
		return nil, fmt.Errorf("APP_ENV must be dev or prod, got %q", cfg.AppEnv)
	}
	if cfg.LogFormat != "json" && cfg.LogFormat != "console" {
		return nil, fmt.Errorf("LOG_FORMAT must be json or console, got %q", cfg.LogFormat)
	}

	if cfg.GmailUnitsPerMin < 1 {
		return nil, fmt.Errorf("GMAIL_UNITS_PER_MIN must be >= 1, got %d", cfg.GmailUnitsPerMin)
	}
	if cfg.GmailConcurrency < 1 {
		return nil, fmt.Errorf("GMAIL_CONCURRENCY must be >= 1, got %d", cfg.GmailConcurrency)
	}

	return cfg, nil
}

func decodeKey32(name, b64 string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("%s must be valid base64: %w", name, err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("%s must decode to exactly 32 bytes for AES-256-GCM, got %d", name, len(key))
	}
	return key, nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envIntOr(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func (c *Config) ListenAddr() string {
	if isTruthy(os.Getenv("BIND_ALL")) {
		return ":" + c.Port
	}
	return "127.0.0.1:" + c.Port
}

func isTruthy(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func (c *Config) IsDev() bool {
	return c.AppEnv == "dev"
}

func defaultFrontendURL(appEnv string) string {
	if appEnv == "dev" {
		return "http://localhost:5173"
	}
	return ""
}

// PostLoginRedirectURL sends the browser back to the SPA after OAuth (Vite in dev, same host when embedded).
func (c *Config) PostLoginRedirectURL() string {
	base := strings.TrimRight(c.FrontendURL, "/")
	if base == "" {
		return "/board"
	}
	return base + "/board"
}

func (c *Config) FrontendConnectURL(query string) string {
	base := strings.TrimRight(c.FrontendURL, "/")
	path := "/connect"
	if query != "" {
		path += "?" + query
	}
	if base == "" {
		return path
	}
	return base + path
}
