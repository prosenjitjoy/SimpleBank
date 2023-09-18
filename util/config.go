package util

import "github.com/ilyakaznacheev/cleanenv"

type ConfigDatabase struct {
	DatabaseURL   string `env:"DATABASE_URL" env-required:"true"`
	ServerAddress string `env:"SERVER_ADDR" env-required:"true"`
}

func LoadConfig(path string) (*ConfigDatabase, error) {
	var cfg ConfigDatabase

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return &ConfigDatabase{}, err
	}

	return &cfg, nil
}
