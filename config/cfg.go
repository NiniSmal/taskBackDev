package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Port            int64         `yaml:"port"`
	Postgres        string        `yaml:"postgres"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
	Email           string        `yaml:"email"`
	Password        string        `yaml:"password"`
}

func GetConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
