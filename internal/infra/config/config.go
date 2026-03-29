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
	BodyLimit    int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	RequestTTL   time.Duration
	AccessLog    AccessLogConfig
	AllowOrigins string
	AllowHeaders string
	AllowMethods string
	RateLimit    RateLimitConfig
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

type RateLimitConfig struct {
	Enabled bool
	Max     int
	Window  time.Duration
}

type AuthConfig struct {
	Enabled bool
	Mode    string // token | jwt | jwks | google
	Token   string
	Header  string
	Scheme  string
	JWT     JWTConfig
}

type JWTConfig struct {
	Issuer          string
	Audience        string
	JWKSURL         string
	PrivateKeyPath  string
	KeyID           string
	CacheTTL        time.Duration
	ClockSkew       time.Duration
	AllowedAlgs     []string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func Load() (Config, error) {
	loadEnv()

	cfg := Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "gofiber-hax"),
			Env:  getEnv("APP_ENV", "dev"),
		},
		HTTP: HTTPConfig{
			Addr:         getEnv("FIBER_ADDR", "5000"),
			Host:         getEnv("FIBER_HOST", "127.0.0.1"),
			BodyLimit:    getEnvInt("HTTP_BODY_LIMIT_MB", 4) * 1024 * 1024,
			ReadTimeout:  time.Duration(getEnvInt("HTTP_READ_TIMEOUT_SEC", 15)) * time.Second,
			WriteTimeout: time.Duration(getEnvInt("HTTP_WRITE_TIMEOUT_SEC", 15)) * time.Second,
			IdleTimeout:  time.Duration(getEnvInt("HTTP_IDLE_TIMEOUT_SEC", 30)) * time.Second,
			RequestTTL:   time.Duration(getEnvInt("HTTP_REQUEST_TIMEOUT_SEC", 15)) * time.Second,
			AllowOrigins: getEnv("FIBER_ALLOW_ORIGINS", "*"),
			AllowHeaders: getEnv("FIBER_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization,X-Request-ID"),
			AllowMethods: getEnv("FIBER_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			AccessLog: AccessLogConfig{
				Enabled:    getEnvBool("HTTP_ACCESS_LOG_ENABLED", true),
				Format:     getEnv("HTTP_ACCESS_LOG_FORMAT", `${ip} - - [${time}] "${method} ${url} ${protocol}" ${status} ${bytesSent}`+"\n"),
				TimeFormat: getEnv("HTTP_ACCESS_LOG_TIME_FORMAT", "02/Jan/2006:15:04:05 -0700"),
			},
			RateLimit: RateLimitConfig{
				Enabled: getEnvBool("HTTP_RATE_LIMIT_ENABLED", true),
				Max:     getEnvInt("HTTP_RATE_LIMIT_MAX", 100),
				Window:  time.Duration(getEnvInt("HTTP_RATE_LIMIT_WINDOW_SEC", 60)) * time.Second,
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
				DSN:         getEnv("MYSQL_DSN", "root:@tcp(127.0.0.1:3306)/app?parseTime=true"),
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

	authMode := strings.ToLower(getEnv("AUTH_MODE", "jwt"))
	jwtIssuer := getEnv("JWT_ISSUER", cfg.App.Name)
	jwtJWKS := getEnv("JWT_JWKS_URL", "")
	jwtAud := getEnv("JWT_AUDIENCE", cfg.App.Name)

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
			Issuer:         jwtIssuer,
			Audience:       jwtAud,
			JWKSURL:        jwtJWKS,
			PrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", "keys/jwt_private.pem"),
			KeyID:          getEnv("JWT_KEY_ID", "gofiber-hax-key-1"),
			CacheTTL:       time.Duration(getEnvInt("JWT_JWKS_TTL_SEC", 3600)) * time.Second,
			ClockSkew:      time.Duration(getEnvInt("JWT_CLOCK_SKEW_SEC", 60)) * time.Second,
			AllowedAlgs: []string{
				getEnv("JWT_ALG", "RS256"),
			},
			AccessTokenTTL:  time.Duration(getEnvInt("JWT_ACCESS_TTL_SEC", 900)) * time.Second,
			RefreshTokenTTL: time.Duration(getEnvInt("JWT_REFRESH_TTL_SEC", 604800)) * time.Second,
		},
	}

	if cfg.DB.Driver != "auto" && cfg.DB.Driver != "mongo" && cfg.DB.Driver != "mysql" && cfg.DB.Driver != "both" {
		return Config{}, fmt.Errorf("unsupported DB_DRIVER: %s", cfg.DB.Driver)
	}
	if isProduction(cfg.App.Env) && cfg.DB.MySQL.AutoMigrate {
		return Config{}, fmt.Errorf("MYSQL_AUTO_MIGRATE must be false in production; use the migrate command instead")
	}
	if cfg.Auth.Enabled {
		switch cfg.Auth.Mode {
		case "", "token":
			if cfg.Auth.Token == "" {
				return Config{}, fmt.Errorf("AUTH_MODE=token but AUTH_TOKEN is empty")
			}
		case "jwt":
			if cfg.Auth.JWT.Issuer == "" || cfg.Auth.JWT.Audience == "" {
				return Config{}, fmt.Errorf("AUTH_MODE=jwt requires JWT_ISSUER and JWT_AUDIENCE")
			}
			if cfg.Auth.JWT.PrivateKeyPath == "" {
				return Config{}, fmt.Errorf("AUTH_MODE=jwt requires JWT_PRIVATE_KEY_PATH")
			}
		case "jwks", "google":
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

func loadEnv() {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if isProduction(env) {
		_ = godotenv.Load(".env")
		return
	}
	_ = godotenv.Overload(".env")
}

func isProduction(env string) bool {
	env = strings.ToLower(strings.TrimSpace(env))
	return env == "prod" || env == "production"
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
