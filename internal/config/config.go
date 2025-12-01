package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ApiPort         string        `env:"API_PORT"`
	DiagPort        string        `env:"DIAG_PORT"`
	HTTPTimeout     time.Duration `env:"HTTP_TIMEOUT"`
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
}

func Load() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("none exists .env file")
	}

	return &cfg, nil
}
