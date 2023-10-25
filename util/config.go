package util

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigDatabase struct {
	Environment       string        `env:"ENVIRONMENT" env-required:"true"`
	DatabaseURL       string        `env:"DATABASE_URL" env-required:"true"`
	MigrationURL      string        `env:"MIGRATION_URL" env-required:"true"`
	HTTPServerAddress string        `env:"HTTP_SERVER_ADDR" env-required:"true"`
	GRPCServerAddress string        `env:"GRPC_SERVER_ADDR" env-required:"true"`
	SecretKey         string        `env:"SECRET_KEY" env-required:"true"`
	TokenDuration     time.Duration `env:"TOKEN_DURATION" env-required:"true"`
	RefreshDuration   time.Duration `env:"REFRESH_DURATION" env-required:"true"`
}

func LoadConfig(path string) (*ConfigDatabase, error) {
	var cfg ConfigDatabase

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return &ConfigDatabase{}, err
	}

	return &cfg, nil
}
