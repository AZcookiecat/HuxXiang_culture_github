package app

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName             string
	Environment         string
	Addr                string
	WriterDSN           string
	ReaderDSN           string
	JWTSecret           string
	CacheTTL            time.Duration
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	MaxOpenConns        int
	MaxIdleConns        int
	ConnMaxLifetime     time.Duration
	RateLimitRPS        int
	RateLimitBurst      int
	CORSAllowOrigins    []string
	EnableHTTPSRedirect bool
	LogJSON             bool
}

func LoadConfig() (Config, error) {
	loadEnvFileIfPresent()

	cfg := Config{
		AppName:             getEnv("GO_POST_SERVICE_NAME", "go-post-service"),
		Environment:         getEnv("GO_POST_SERVICE_ENV", "dev"),
		Addr:                getEnv("GO_POST_SERVICE_ADDR", ":8080"),
		JWTSecret:           getEnv("JWT_SECRET_KEY", "jwt-huxiang-secret-key-dev"),
		CacheTTL:            getDuration("POST_CACHE_TTL", 15*time.Second),
		ReadTimeout:         getDuration("POST_SERVICE_READ_TIMEOUT", 3*time.Second),
		WriteTimeout:        getDuration("POST_SERVICE_WRITE_TIMEOUT", 5*time.Second),
		MaxOpenConns:        getInt("POST_DB_MAX_OPEN_CONNS", 128),
		MaxIdleConns:        getInt("POST_DB_MAX_IDLE_CONNS", 32),
		ConnMaxLifetime:     getDuration("POST_DB_CONN_MAX_LIFETIME", 30*time.Minute),
		RateLimitRPS:        getInt("POST_RATE_LIMIT_RPS", 300),
		RateLimitBurst:      getInt("POST_RATE_LIMIT_BURST", 600),
		CORSAllowOrigins:    getCSV("POST_CORS_ALLOW_ORIGINS", []string{"*"}),
		EnableHTTPSRedirect: getBool("POST_ENABLE_HTTPS_REDIRECT", false),
		LogJSON:             getBool("POST_LOG_JSON", false),
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

func loadEnvFileIfPresent() {
	paths := []string{
		getEnv("GO_POST_SERVICE_ENV_FILE", ""),
		".env",
	}
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			_ = godotenv.Overload(path)
			return
		}
	}
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
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", u.User.Username(), password, host, dbName, query.Encode()), nil
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

func getBool(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func getCSV(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	if len(result) == 0 {
		return fallback
	}
	return result
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
