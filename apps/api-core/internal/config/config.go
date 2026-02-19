package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JWTConfig
	CORS        CORSConfig
	App         AppConfig
	Security    SecurityConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
}

type AppConfig struct {
	Name    string
	Version string
	Deploy  DeployConfig
}

type DeployConfig struct {
	CommitSHA       string
	CommitShortSHA  string
	CommitTitle     string
	CommitTimestamp string
}

type SecurityConfig struct {
	MaxLoginAttempts int           `yaml:"maxLoginAttempts" env:"SECURITY_MAX_LOGIN_ATTEMPTS" env-default:"5"`
	LockoutDuration  time.Duration `yaml:"lockoutDuration" env:"SECURITY_LOCKOUT_DURATION" env-default:"15m"`
}

func Load() *Config {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicitly bind environment variables for Kubernetes compatibility
	viper.BindEnv("database.mysql.host", "DATABASE_MYSQL_HOST")
	viper.BindEnv("database.mysql.port", "DATABASE_MYSQL_PORT")
	viper.BindEnv("database.mysql.user", "DATABASE_MYSQL_USER")
	viper.BindEnv("database.mysql.password", "DATABASE_MYSQL_PASSWORD")
	viper.BindEnv("database.mysql.database", "DATABASE_MYSQL_DATABASE", "DATABASE_MYSQL_SCHEMA")
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("security.max_login_attempts", "SECURITY_MAX_LOGIN_ATTEMPTS")
	viper.BindEnv("security.lockout_duration", "SECURITY_LOCKOUT_DURATION")

	profile := viper.GetString("PROFILE")
	if profile == "" {
		// Check for Spring Boot compatibility
		profile = viper.GetString("SPRING_PROFILES_ACTIVE")
	}
	if profile == "" {
		profile = "local"
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	if profile != "" {
		viper.SetConfigName("application-" + profile)
		if err := viper.MergeInConfig(); err != nil {
			log.Printf("Error reading profile config file: %v", err)
		}
	}

	cfg := &Config{
		Environment: viper.GetString("environment"),
		Server: ServerConfig{
			Port: viper.GetInt("server.port"),
		},
		Database: DatabaseConfig{
			Host:            viper.GetString("database.mysql.host"),
			Port:            viper.GetInt("database.mysql.port"),
			User:            viper.GetString("database.mysql.user"),
			Password:        viper.GetString("database.mysql.password"),
			Database:        viper.GetString("database.mysql.database"),
			MaxIdleConns:    viper.GetInt("database.mysql.max_idle_conns"),
			MaxOpenConns:    viper.GetInt("database.mysql.max_open_conns"),
			ConnMaxLifetime: viper.GetDuration("database.mysql.conn_max_lifetime"),
			ConnMaxIdleTime: viper.GetDuration("database.mysql.conn_max_idle_time"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("redis.host"),
			Port:     viper.GetInt("redis.port"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
		JWT: JWTConfig{
			Secret:             viper.GetString("jwt.secret"),
			AccessTokenExpiry:  time.Duration(viper.GetInt("jwt.access_token_expiry")) * time.Second,
			RefreshTokenExpiry: time.Duration(viper.GetInt("jwt.refresh_token_expiry")) * time.Second,
		},
		CORS: CORSConfig{
			AllowedOrigins: viper.GetStringSlice("cors.allowed_origins"),
		},
		App: AppConfig{
			Name:    viper.GetString("app.name"),
			Version: viper.GetString("app.version"),
			Deploy: DeployConfig{
				CommitSHA:       viper.GetString("app.deploy.commit_sha"),
				CommitShortSHA:  viper.GetString("app.deploy.commit_short_sha"),
				CommitTitle:     viper.GetString("app.deploy.commit_title"),
				CommitTimestamp: viper.GetString("app.deploy.commit_timestamp"),
			},
		},
	}

	if cfg.Environment == "" {
		cfg.Environment = profile
	}

	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}

	// Set defaults for database configuration if not provided
	if cfg.Database.Host == "" {
		cfg.Database.Host = "localhost"
	}

	if cfg.Database.Port == 0 {
		cfg.Database.Port = 3306
	}

	if cfg.Database.User == "" {
		cfg.Database.User = "rms"
	}

	if cfg.Database.Password == "" {
		cfg.Database.Password = "rms123"
	}

	if cfg.Database.Database == "" {
		cfg.Database.Database = "rms-core"
	}

	if cfg.Database.MaxIdleConns == 0 {
		cfg.Database.MaxIdleConns = 10
	}

	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 100
	}

	if cfg.Database.ConnMaxLifetime == 0 {
		cfg.Database.ConnMaxLifetime = time.Hour
	}

	if cfg.Database.ConnMaxIdleTime == 0 {
		cfg.Database.ConnMaxIdleTime = 10 * time.Minute
	}

	// Set defaults for Redis configuration if not provided
	if cfg.Redis.Host == "" {
		cfg.Redis.Host = "localhost"
	}

	if cfg.Redis.Port == 0 {
		cfg.Redis.Port = 6379
	}

	// Set defaults for JWT configuration if not provided
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "local-development-secret"
	}

	if cfg.JWT.AccessTokenExpiry == 0 {
		cfg.JWT.AccessTokenExpiry = 15 * time.Minute
	}

	if cfg.JWT.RefreshTokenExpiry == 0 {
		cfg.JWT.RefreshTokenExpiry = 7 * 24 * time.Hour
	}

	if len(cfg.CORS.AllowedOrigins) == 0 {
		cfg.CORS.AllowedOrigins = []string{"http://localhost:9000"}
	}

	if cfg.App.Name == "" {
		cfg.App.Name = "Resort Management System"
	}

	if cfg.App.Version == "" {
		cfg.App.Version = "1.0.0"
	}

	cfg.Security = SecurityConfig{
		MaxLoginAttempts: viper.GetInt("security.max_login_attempts"),
		LockoutDuration:  viper.GetDuration("security.lockout_duration"),
	}

	// Set defaults for Security configuration if not provided
	if cfg.Security.MaxLoginAttempts == 0 {
		cfg.Security.MaxLoginAttempts = 5
	}

	if cfg.Security.LockoutDuration == 0 {
		cfg.Security.LockoutDuration = 15 * time.Minute
	}

	return cfg
}
