package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_DRIVER             string        `mapstructure:"DB_DRIVER"`
	DB_URL                string        `mapstructure:"DB_URL"`
	SYMMETRIC_KEY         string        `mapstructure:"SYMMETRIC_KEY"`
	LISTEN_ADDR           string        `mapstructure:"LISTEN_ADDR"`
	ACCESS_TOKEN_DURATION time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	var cfg *Config
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
