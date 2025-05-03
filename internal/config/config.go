package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Nats struct {
		URL string `env:"NATS_URL,notEmpty"`
	}

	Logger struct {
		Level       string `env:"LOGGER_LEVEL" env-default:"info"`
		FieldMaxLen int    `env:"LOG_FIELD_MAX_LEN" envDefault:"2000"`
	}
}

func Load() (Config, error) {
	var config Config

	if err := env.Parse(&config); err != nil {
		return Config{}, fmt.Errorf("env.Parse: %w", err)
	}

	return config, nil
}
