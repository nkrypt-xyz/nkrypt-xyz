package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	MinIO       MinIOConfig       `mapstructure:"minio"`
	BlobStorage BlobStorageConfig `mapstructure:"blob_storage"`
	IAM         IAMConfig         `mapstructure:"iam"`
	Crypto      CryptoConfig      `mapstructure:"crypto"`
	Log         LogConfig         `mapstructure:"log"`
}

type ServerConfig struct {
	HTTP        HTTPConfig  `mapstructure:"http"`
	HTTPS       HTTPSConfig `mapstructure:"https"`
	ContextPath string      `mapstructure:"context_path"`
}

type HTTPConfig struct {
	Port int `mapstructure:"port"`
}

type HTTPSConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	Port         int    `mapstructure:"port"`
	KeyFile      string `mapstructure:"key_file"`
	CertFile     string `mapstructure:"cert_file"`
	CABundleFile string `mapstructure:"ca_bundle_file"`
}

type DatabaseConfig struct {
	URL             string        `mapstructure:"url"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type MinIOConfig struct {
	Endpoint   string `mapstructure:"endpoint"`
	AccessKey  string `mapstructure:"access_key"`
	SecretKey  string `mapstructure:"secret_key"`
	BucketName string `mapstructure:"bucket_name"`
	UseSSL     bool   `mapstructure:"use_ssl"`
}

type BlobStorageConfig struct {
	MaxFileSizeBytes int64 `mapstructure:"max_file_size_bytes"`
}

type IAMConfig struct {
	APIKeyLength            int           `mapstructure:"api_key_length"`
	SessionValidityDuration time.Duration `mapstructure:"session_validity_duration"`
	DefaultAdminUsername    string        `mapstructure:"default_admin_username"`
	DefaultAdminDisplayName string        `mapstructure:"default_admin_display_name"`
	DefaultAdminPassword    string        `mapstructure:"default_admin_password"`
}

type CryptoConfig struct {
	Argon2Memory      uint32 `mapstructure:"argon2_memory"`
	Argon2Iterations  uint32 `mapstructure:"argon2_iterations"`
	Argon2Parallelism uint8  `mapstructure:"argon2_parallelism"`
	Argon2SaltLength  uint32 `mapstructure:"argon2_salt_length"`
	Argon2KeyLength   uint32 `mapstructure:"argon2_key_length"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load loads configuration using Viper, following the precedence and defaults
// specified in 03-CONFIGURATION.md.
func Load() (*Config, error) {
	v := viper.New()

	// Defaults for application behavior (NOT external dependencies)
	v.SetDefault("server.http.port", 9041)
	v.SetDefault("server.https.enabled", false)
	v.SetDefault("server.https.port", 9443)
	v.SetDefault("server.context_path", "/")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.conn_max_lifetime", "5m")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("minio.bucket_name", "nkrypt-blobs")
	v.SetDefault("minio.use_ssl", false)
	v.SetDefault("blob_storage.max_file_size_bytes", int64(5368709120))
	v.SetDefault("iam.api_key_length", 128)
	v.SetDefault("iam.session_validity_duration", "168h")
	v.SetDefault("iam.default_admin_username", "admin")
	v.SetDefault("iam.default_admin_display_name", "Default Admin")
	// NOTE: No default for admin password - must be explicitly set!
	v.SetDefault("crypto.argon2_memory", 65536)
	v.SetDefault("crypto.argon2_iterations", 3)
	v.SetDefault("crypto.argon2_parallelism", 4)
	v.SetDefault("crypto.argon2_salt_length", 16)
	v.SetDefault("crypto.argon2_key_length", 32)
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")

	// NOTE: No defaults for external dependencies!
	// Database URL, Redis address, and MinIO endpoint MUST be provided
	// via environment variables or config file.

	// Config file lookup
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/nkrypt-xyz")
	v.AddConfigPath("/etc/nkrypt-xyz")

	// Environment variables
	v.SetEnvPrefix("NK")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	
	// Bind specific environment variables explicitly
	v.BindEnv("database.url", "NK_DATABASE_URL")
	v.BindEnv("redis.addr", "NK_REDIS_ADDR")
	v.BindEnv("minio.endpoint", "NK_MINIO_ENDPOINT")
	v.BindEnv("minio.access_key", "NK_MINIO_ACCESS_KEY")
	v.BindEnv("minio.secret_key", "NK_MINIO_SECRET_KEY")
	v.BindEnv("minio.bucket_name", "NK_MINIO_BUCKET_NAME")
	v.BindEnv("iam.default_admin_password", "NK_IAM_DEFAULT_ADMIN_PASSWORD")
	v.BindEnv("log.level", "NK_LOG_LEVEL")
	v.BindEnv("log.format", "NK_LOG_FORMAT")

	// Read config file (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate required external dependencies
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate ensures all required configuration values are provided.
// External dependencies (database, redis, minio) MUST be explicitly configured.
func (c *Config) Validate() error {
	var missing []string

	// Required: Database URL
	if c.Database.URL == "" {
		missing = append(missing, "NK_DATABASE_URL")
	}

	// Required: Redis address
	if c.Redis.Addr == "" {
		missing = append(missing, "NK_REDIS_ADDR")
	}

	// Required: MinIO endpoint
	if c.MinIO.Endpoint == "" {
		missing = append(missing, "NK_MINIO_ENDPOINT")
	}

	// Required: MinIO credentials
	if c.MinIO.AccessKey == "" {
		missing = append(missing, "NK_MINIO_ACCESS_KEY")
	}
	if c.MinIO.SecretKey == "" {
		missing = append(missing, "NK_MINIO_SECRET_KEY")
	}

	// Required: Default admin password (security - must be explicitly set)
	if c.IAM.DefaultAdminPassword == "" {
		missing = append(missing, "NK_IAM_DEFAULT_ADMIN_PASSWORD")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration: %s\n\nExternal dependencies and security credentials must be explicitly configured.\nSee .env.example for required environment variables.", strings.Join(missing, ", "))
	}

	return nil
}

