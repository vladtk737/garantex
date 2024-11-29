package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		GRPC
		Log
		PG
	}

	GRPC struct {
		Port int `env-required:"true" env:"GRPC_PORT"`
	}

	Log struct {
		Level       string   `env-required:"true" env:"LOG_LEVEL" mapstructure:"level"`
		Encoding    string   `mapstructure:"encoding"`
		OutputPaths []string `mapstructure:"output"`
		ErrorOutput []string `mapstructure:"errorOutput"`
	}

	PG struct {
		Host          string `env-required:"true" env:"PG_HOST"`
		Port          uint16 `env-required:"true" env:"PG_PORT"`
		User          string `env-required:"true" env:"PG_USER"`
		Password      string `env-required:"true" env:"PG_PASSWORD"`
		UserAdmin     string `env-required:"true" env:"PG_USER_ADMIN"`
		PasswordAdmin string `env-required:"true" env:"PG_PASSWORD_ADMIN"`
		Database      string `env-required:"true" env:"PG_DATABASE"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
