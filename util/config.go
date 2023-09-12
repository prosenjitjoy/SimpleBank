package util

import "github.com/spf13/viper"

// Config stores all configuration of the applications.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DatabaseURL   string `mapstructure:"DATABASE_URL"`
	ServerAddress string `mapstructure:"SERVER_ADDR"`
}

// LoadConfig read configuration from file or environment variables
func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
}
