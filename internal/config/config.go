package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ApiPort         string        `env:"API_PORT" env-default:"8080"`
	DiagPort        string        `env:"DIAG_PORT" env-default:"8081"`
	DBUrl           string        `env:"DB_URL" env-default:"postgres://user:password@localhost:5432/petstore?sslmode=disable"`
	HTTPTimeout     time.Duration `env:"HTTP_TIMEOUT" env-default:"5s"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" env-default:"30s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"30s"`
}

func Load() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("none exists .env file")
	}

	return &cfg, nil
}
