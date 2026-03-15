package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App  AppConfig
	HTTP HTTPConfig
	DB   DBConfig
	Log  LogConfig
	Auth AuthConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type HTTPConfig struct {
	Addr         string
	Host         string
	AccessLog    AccessLogConfig
	AlowOrigins  string
	AllowHeaders string
	AllowMethods string
}

type DBConfig struct {
	Driver string // auto | mongo | mysql | both
	Mongo  MongoConfig
	MySQL  MySQLConfig
}

type MongoConfig struct {
	URI  string
	DB   string
	Host string
	Port string
	User string
	Pass string
}

type MySQLConfig struct {
	DSN         string
	ReplicaDSN  string
	AutoMigrate bool
}

type LogConfig struct {
	Level  string // debug | info | warn | error
	Format string // json | pretty | text
}

type AccessLogConfig struct {
	Enabled    bool
	Format     string
	TimeFormat string
}

type AuthConfig struct {
	Enabled bool
	Mode    string // token | jwt | google
	Token   string
	Header  string
	Scheme  string
	JWT     JWTConfig
}

type JWTConfig struct {
	Issuer      string
	Audience    string
	JWKSURL     string
	CacheTTL    time.Duration
	ClockSkew   time.Duration
	AllowedAlgs []string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "gofiber-hax"),
			Env:  getEnv("APP_ENV", "dev"),
		},
		HTTP: HTTPConfig{
			Addr:         getEnv("FIBER_ADDR", "5000"),
			Host:         getEnv("FIBER_HOST", "127.0.0.1"),
			AlowOrigins:  getEnv("FIBER_ALLOW_ORIGINS", "*"),
			AllowHeaders: getEnv("FIBER_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization"),
			AllowMethods: getEnv("FIBER_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			AccessLog: AccessLogConfig{
				Format:     getEnv("HTTP_ACCESS_LOG_FORMAT", `${ip} - - [${time}] "${method} ${url} ${protocol}" ${status} ${bytesSent}`+"\n"),
				TimeFormat: getEnv("HTTP_ACCESS_LOG_TIME_FORMAT", "02/Jan/2006:15:04:05 -0700"),
			},
		},
		DB: DBConfig{
			Driver: strings.ToLower(getEnv("DB_DRIVER", "auto")),
			Mongo: MongoConfig{
				URI:  getEnv("MONGO_URI", "mongodb://localhost:27017"),
				DB:   getEnv("MONGO_DB", "example_db"),
				Host: getEnv("MONGO_HOST", "localhost"),
				Port: getEnv("MONGO_PORT", "27017"),
				User: getEnv("MONGO_USER", ""),
				Pass: getEnv("MONGO_PASS", ""),
			},
			MySQL: MySQLConfig{
				DSN:         getEnv("MYSQL_DSN", "root:password@tcp(127.0.0.1:3306)/app?parseTime=true"),
				ReplicaDSN:  getEnv("MYSQL_REPLICA_DSN", ""),
				AutoMigrate: getEnvBool("MYSQL_AUTO_MIGRATE", false),
			},
		},
		Log: LogConfig{
			Level:  strings.ToLower(getEnv("LOG_LEVEL", "debug")),
			Format: strings.ToLower(getEnv("LOG_FORMAT", "json")),
		},
		Auth: AuthConfig{},
	}

	authMode := strings.ToLower(getEnv("AUTH_MODE", "token"))
	jwtIssuer := getEnv("JWT_ISSUER", "")
	jwtJWKS := getEnv("JWT_JWKS_URL", "")
	jwtAud := getEnv("JWT_AUDIENCE", "")

	if authMode == "google" {
		if jwtIssuer == "" {
			jwtIssuer = "https://accounts.google.com"
		}
		if jwtJWKS == "" {
			jwtJWKS = "https://www.googleapis.com/oauth2/v3/certs"
		}
		if jwtAud == "" {
			jwtAud = getEnv("GOOGLE_CLIENT_ID", "")
		}
	}

	cfg.Auth = AuthConfig{
		Enabled: getEnvBool("AUTH_ENABLED", false),
		Mode:    authMode,
		Token:   getEnv("AUTH_TOKEN", ""),
		Header:  getEnv("AUTH_HEADER", "Authorization"),
		Scheme:  getEnv("AUTH_SCHEME", "Bearer"),
		JWT: JWTConfig{
			Issuer:      jwtIssuer,
			Audience:    jwtAud,
			JWKSURL:     jwtJWKS,
			CacheTTL:    time.Duration(getEnvInt("JWT_JWKS_TTL_SEC", 3600)) * time.Second,
			ClockSkew:   time.Duration(getEnvInt("JWT_CLOCK_SKEW_SEC", 60)) * time.Second,
			AllowedAlgs: []string{"RS256"},
		},
	}

	if cfg.DB.Driver != "auto" && cfg.DB.Driver != "mongo" && cfg.DB.Driver != "mysql" && cfg.DB.Driver != "both" {
		return Config{}, fmt.Errorf("unsupported DB_DRIVER: %s", cfg.DB.Driver)
	}
	if cfg.Auth.Enabled {
		switch cfg.Auth.Mode {
		case "", "token":
			if cfg.Auth.Token == "" {
				return Config{}, fmt.Errorf("AUTH_MODE=token but AUTH_TOKEN is empty")
			}
		case "jwt", "google":
			if cfg.Auth.JWT.JWKSURL == "" || cfg.Auth.JWT.Issuer == "" {
				return Config{}, fmt.Errorf("AUTH_MODE=%s requires JWT_ISSUER and JWT_JWKS_URL", cfg.Auth.Mode)
			}
			if cfg.Auth.Mode == "google" && cfg.Auth.JWT.Audience == "" {
				return Config{}, fmt.Errorf("AUTH_MODE=google requires GOOGLE_CLIENT_ID or JWT_AUDIENCE")
			}
		default:
			return Config{}, fmt.Errorf("unsupported AUTH_MODE: %s", cfg.Auth.Mode)
		}
	}
	return cfg, nil
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if val == "" {
		return def
	}
	switch val {
	case "1", "true", "t", "yes", "y", "on":
		return true
	case "0", "false", "f", "no", "n", "off":
		return false
	default:
		return def
	}
}

func getEnvInt(key string, def int) int {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return def
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return n
}
