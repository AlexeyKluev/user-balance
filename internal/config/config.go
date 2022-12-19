package config

import (
	"fmt"

	"github.com/vrischmann/envconfig"
)

// Config struct
type Config struct {
	Addr         string `envconfig:"default=0.0.0.0:8080"`
	IsProduction bool   `envconfig:"default=false"`
	// Redis        *redis.Config
	// Postgres     *postgres.Config
}

// InitConfig func
func InitConfig(prefix string) (*Config, error) {
	config := &Config{}
	if err := envconfig.InitWithPrefix(config, prefix); err != nil {
		return nil, fmt.Errorf("init config failed: %w", err)
	}

	return config, nil
}
