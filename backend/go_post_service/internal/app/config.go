package app

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr            string
	WriterDSN       string
	ReaderDSN       string
	JWTSecret       string
	CacheTTL        time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	RateLimitRPS    int
	RateLimitBurst  int
}

func LoadConfig() (Config, error) {
	cfg := Config{
		Addr:            getEnv("GO_POST_SERVICE_ADDR", ":8080"),
		JWTSecret:       getEnv("JWT_SECRET_KEY", "jwt-huxiang-secret-key-dev"),
		CacheTTL:        getDuration("POST_CACHE_TTL", 15*time.Second),
		ReadTimeout:     getDuration("POST_SERVICE_READ_TIMEOUT", 3*time.Second),
		WriteTimeout:    getDuration("POST_SERVICE_WRITE_TIMEOUT", 5*time.Second),
		MaxOpenConns:    getInt("POST_DB_MAX_OPEN_CONNS", 128),
		MaxIdleConns:    getInt("POST_DB_MAX_IDLE_CONNS", 32),
		ConnMaxLifetime: getDuration("POST_DB_CONN_MAX_LIFETIME", 30*time.Minute),
		RateLimitRPS:    getInt("POST_RATE_LIMIT_RPS", 300),
		RateLimitBurst:  getInt("POST_RATE_LIMIT_BURST", 600),
	}

	var err error
	cfg.WriterDSN, err = NormalizeMySQLDSN(getEnv("DATABASE_URL", ""))
	if err != nil {
		return Config{}, fmt.Errorf("normalize DATABASE_URL: %w", err)
	}

	reader := getEnv("READ_DATABASE_URL", "")
	if reader == "" {
		cfg.ReaderDSN = cfg.WriterDSN
	} else {
		cfg.ReaderDSN, err = NormalizeMySQLDSN(reader)
		if err != nil {
			return Config{}, fmt.Errorf("normalize READ_DATABASE_URL: %w", err)
		}
	}

	if cfg.WriterDSN == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func NormalizeMySQLDSN(raw string) (string, error) {
	if raw == "" {
		return "", nil
	}

	if !strings.Contains(raw, "://") {
		return withMySQLOptions(raw), nil
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	if u.Scheme != "mysql+pymysql" && u.Scheme != "mysql" {
		return "", fmt.Errorf("unsupported DSN scheme %q", u.Scheme)
	}

	host := u.Host
	if !strings.Contains(host, ":") {
		host = net.JoinHostPort(host, "3306")
	}

	dbName := strings.TrimPrefix(u.Path, "/")
	query := u.Query()
	if query.Get("charset") == "" {
		query.Set("charset", "utf8mb4")
	}
	query.Set("parseTime", "true")
	query.Set("loc", "Local")

	password, _ := u.User.Password()
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", u.User.Username(), password, host, dbName, query.Encode())
	return dsn, nil
}

func withMySQLOptions(dsn string) string {
	if strings.Contains(dsn, "parseTime=") {
		return dsn
	}
	joiner := "?"
	if strings.Contains(dsn, "?") {
		joiner = "&"
	}
	return dsn + joiner + "parseTime=true&loc=Local"
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
