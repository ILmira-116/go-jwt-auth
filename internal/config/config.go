package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP       HTTPConfig
	JWT        JWTConfig
	Postgres   PostgresConfig
	WebhookURL string `env:"WEBHOOK_URL" env-default:"http://localhost:8081/webhook"`
	LogLevel   string `env:"LOG_LEVEL" env-default:"info"`
}

type JWTConfig struct {
	Secret          string `env:"JWT_SECRET" env-default:"supersecretkey"`
	AccessTokenTTL  int    `env:"ACCESS_TOKEN_TTL_MIN" env-default:"15"`
	RefreshTokenTTL int    `env:"REFRESH_TOKEN_TTL_HOURS" env-default:"168"`
}

type PostgresConfig struct {
	User     string `env:"POSTGRES_USER" env-default:"user"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"password"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	DBName   string `env:"POSTGRES_DB" env-default:"auth_db"`
	SSLMode  string `env:"POSTGRES_SSLMODE" env-default:"disable"`
}

type HTTPConfig struct {
	ServerHost         string        `env:"SERVER_HOST" env-default:"0.0.0.0"`
	ServerPort         string        `env:"SERVER_PORT" env-default:"8080"`
	ServerReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT" env-default:"5s"`
	ServerWriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" env-default:"10s"`
	ServerIdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT" env-default:"120s"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
