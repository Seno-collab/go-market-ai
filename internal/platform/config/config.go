package config

import (
	"fmt"
	"go-ai/pkg/logger"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultAccessSecret  = "your-access-secret-key"
	defaultRefreshSecret = "your-refresh-secret-key"
)

type Config struct {
	JwtAccessSecret     string `mapstructure:"JWT_SECRET"`
	JwtRefreshSecret    string `mapstructure:"JWT_REFRESH_SECRET"`
	JwtExpiresIn        int    `mapstructure:"JWT_EXPIRES_IN"`
	JwtRefreshExpiresIn int    `mapstructure:"JWT_REFRESH_EXPIRES_IN"`
	RedisHost           string `mapstructure:"REDIS_HOST"`
	RedisPassword       string `mapstructure:"REDIS_PASSWORD"`
	RedisPort           int    `mapstructure:"REDIS_PORT"`
	RedisDB             int    `mapstructure:"REDIS_DB"`
	DBName              string `mapstructure:"POSTGRES_DB"`
	DBHost              string `mapstructure:"POSTGRES_HOST"`
	DBPort              string `mapstructure:"POSTGRES_PORT"`
	DBUser              string `mapstructure:"POSTGRES_USER"`
	DBPassword          string `mapstructure:"POSTGRES_PASSWORD"`
	DBSSLMode           string `mapstructure:"db_sslmode"`
	ServerPort          string `mapstructure:"PORT"`
	ServerHost          string `mapstructure:"server_host"`
	Environment         string `mapstructure:"ENVIRONMENT"`
	MinioEndPoint       string `mapstructure:"MINIO_END_POINT"`
	MinioPort           string `mapstructure:"MINIO_PORT"`
	MinioAccessKey      string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey      string `mapstructure:"MINIO_SECRET_KEY"`
	Bucket              string `mapstructure:"MINIO_BUCKET"`
	MinioUseSSL         bool   `mapstructure:"MINIO_USE_SSL"`

	// PostgreSQL Connection Pool Settings
	DBMaxConns          int `mapstructure:"DB_MAX_CONNS"`
	DBMinConns          int `mapstructure:"DB_MIN_CONNS"`
	DBMaxConnLifetime   int `mapstructure:"DB_MAX_CONN_LIFETIME"`   // seconds
	DBMaxConnIdleTime   int `mapstructure:"DB_MAX_CONN_IDLE_TIME"` // seconds
	DBHealthCheckPeriod int `mapstructure:"DB_HEALTH_CHECK_PERIOD"` // seconds

	// Redis Connection Pool Settings
	RedisPoolSize     int `mapstructure:"REDIS_POOL_SIZE"`
	RedisMinIdleConns int `mapstructure:"REDIS_MIN_IDLE_CONNS"`
	RedisMaxRetries   int `mapstructure:"REDIS_MAX_RETRIES"`

	// Rate Limiting Settings
	RateLimitRequests int `mapstructure:"RATE_LIMIT_REQUESTS"` // requests per second
	RateLimitBurst    int `mapstructure:"RATE_LIMIT_BURST"`    // burst size

	// Timeout Settings
	RequestTimeout  int `mapstructure:"REQUEST_TIMEOUT"`  // seconds
	ShutdownTimeout int `mapstructure:"SHUTDOWN_TIMEOUT"` // seconds
}

func LoadConfig() (*Config, error) {
	logger := logger.NewLogger()
	logger.With().Str("component", "config").Logger()

	// Set default values
	setDefaults()

	// Set environment variable key replacer to convert dots to underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Enable automatic environment variable reading with highest priority
	viper.AutomaticEnv()

	// Ensure environment variables take precedence over config files
	viper.SetEnvPrefix("")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Warn().Err(err).Msg("error reading .env file")
		}
	} else {
		logger.Info().Msg(".env file loaded for backward compatibility")
	}

	// Create config struct
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		logger.Error().Err(err).Msg("unable to decode config")
		return nil, err
	}
	if cfg.JwtAccessSecret == defaultAccessSecret || cfg.JwtRefreshSecret == defaultRefreshSecret {
		logger.Warn().Msg("JWT secrets are using default values; set JWT_SECRET and JWT_REFRESH_SECRET")
		if cfg.IsProduction() {
			return nil, fmt.Errorf("jwt secrets must be set for production")
		}
	}

	return cfg, nil
}

func setDefaults() {
	// JWT defaults
	viper.SetDefault("JWT_SECRET", defaultAccessSecret)
	viper.SetDefault("JWT_REFRESH_SECRET", defaultRefreshSecret)
	viper.SetDefault("JWT_EXPIRES_IN", 3000)
	viper.SetDefault("JWT_REFRESH_EXPIRES_IN", 6480000)

	// Redis defaults
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_PASSWORD", "")

	// Database defaults
	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("POSTGRES_PORT", "5432")
	viper.SetDefault("POSTGRES_DB", "go-ai")
	viper.SetDefault("POSTGRES_USER", "postgres")
	viper.SetDefault("POSTGRES_PASSWORD", "")
	viper.SetDefault("DB_SSLMODE", "disable")

	// Server defaults
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("PORT", "8080")

	// Environment
	viper.SetDefault("ENVIRONMENT", "development")

	// Minio defaults
	viper.SetDefault("MINIO_END_POINT", "localhost")
	viper.SetDefault("MINIO_PORT", "9000")
	viper.SetDefault("MINIO_ACCESS_KEY", "minioadmin")
	viper.SetDefault("MINIO_SECRET_KEY", "minioadmin")
	viper.SetDefault("MINIO_USE_SSL", false)
	viper.SetDefault("MINIO_BUCKET", "uploads")

	// PostgreSQL Connection Pool defaults
	viper.SetDefault("DB_MAX_CONNS", 25)
	viper.SetDefault("DB_MIN_CONNS", 5)
	viper.SetDefault("DB_MAX_CONN_LIFETIME", 3600)   // 1 hour
	viper.SetDefault("DB_MAX_CONN_IDLE_TIME", 300)   // 5 minutes
	viper.SetDefault("DB_HEALTH_CHECK_PERIOD", 30)   // 30 seconds

	// Redis Connection Pool defaults
	viper.SetDefault("REDIS_POOL_SIZE", 10)
	viper.SetDefault("REDIS_MIN_IDLE_CONNS", 5)
	viper.SetDefault("REDIS_MAX_RETRIES", 3)

	// Rate Limiting defaults
	viper.SetDefault("RATE_LIMIT_REQUESTS", 20)
	viper.SetDefault("RATE_LIMIT_BURST", 50)

	// Timeout defaults
	viper.SetDefault("REQUEST_TIMEOUT", 30)
	viper.SetDefault("SHUTDOWN_TIMEOUT", 10)
}

// GetString returns a string value from config
func (c *Config) GetString(key string) string {
	return viper.GetString(key)
}

// GetInt returns an int value from config
func (c *Config) GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool returns a bool value from config
func (c *Config) GetBool(key string) bool {
	return viper.GetBool(key)
}

// IsDevelopment returns true if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
