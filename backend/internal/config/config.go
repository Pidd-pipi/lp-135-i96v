package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
)

// Config 应用配置，全部由环境变量注入。
type Config struct {
	AppEnv     string `env:"APP_ENV" envDefault:"development"`
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	DBHost     string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPort     string `env:"DB_PORT" envDefault:"3306"`
	DBUser     string `env:"DB_USER" envDefault:"root"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"root"`
	DBName     string `env:"DB_NAME" envDefault:"givetrack"`

	DBMaxOpenConns    int `env:"DB_MAX_OPEN_CONNS" envDefault:"50"`
	DBMaxIdleConns    int `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
	DBConnMaxLifetime int `env:"DB_CONN_MAX_LIFETIME_MINUTES" envDefault:"60"`
	DBRetryCount      int `env:"DB_CONNECT_RETRY_COUNT" envDefault:"10"`
	DBRetryInterval   int `env:"DB_CONNECT_RETRY_INTERVAL_SECONDS" envDefault:"3"`

	JWTSecret string `env:"JWT_SECRET" envDefault:"change_me_to_a_long_random_string"`
	JWTExpire int    `env:"JWT_EXPIRE_HOURS" envDefault:"72"`

	CORSAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`

	AuthRateLimit      int `env:"AUTH_RATE_LIMIT" envDefault:"10"`
	AuthRateWindowSecs int `env:"AUTH_RATE_WINDOW_SECONDS" envDefault:"60"`
}

// Load 从环境变量解析并校验配置。
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.AppEnv == "production" {
		if len(cfg.JWTSecret) < 32 || cfg.JWTSecret == "change_me_to_a_long_random_string" {
			return nil, fmt.Errorf("JWT_SECRET must be at least 32 chars and must not use the default value in production")
		}
		if cfg.CORSAllowedOrigins == "*" {
			return nil, fmt.Errorf("CORS_ALLOWED_ORIGINS must be an explicit origin in production")
		}
	}
	if cfg.DBMaxOpenConns <= 0 || cfg.DBMaxIdleConns < 0 || cfg.DBMaxIdleConns > cfg.DBMaxOpenConns {
		return nil, fmt.Errorf("invalid db pool config")
	}
	return cfg, nil
}

// DSN 返回 MySQL 连接串。
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// CORSOrigins 解析 CORS 允许来源列表。
func (c *Config) CORSOrigins() []string {
	parts := strings.Split(c.CORSAllowedOrigins, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
