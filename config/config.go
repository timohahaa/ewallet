package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		PG     `yaml:"postgres"`
		Server `yaml:"server"`
	}
	PG struct {
		URL          string `yaml:"url" env:"PG_URL" env-required:"true"`
		ConnPoolSize int    `yaml:"maxConnPoolSize" env:"PG_MAX_POOL_SIZE"`
	}
	Server struct {
		Port string `yaml:"port" env:"HTTP_SERVER_PORT"`
	}
)

func NewConfig(filePath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(filePath, cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	//	err = cleanenv.UpdateEnv(cfg)
	//	if err != nil {
	//		return nil, fmt.Errorf("error updating env: %w", err)
	//	}

	return cfg, nil
}
