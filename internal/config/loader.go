package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func Load(configPath string) (*Config, error) {
	cfg := &Config{}

	if configPath != "" {
		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		if err := cleanenv.ReadEnv(cfg); err != nil {
			return nil, fmt.Errorf("failed to read environment variables: %w", err)
		}
	}

	return cfg, nil
}

func LoadDefault() (*Config, error) {
	return Load("")
}
