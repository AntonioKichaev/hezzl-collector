package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Postgres   `yaml:"postgres"`
		Redis      `yaml:"redis"`
		Nats       `yaml:"nats"`
		Clickhouse `yaml:"clickhouse"`
	}

	Postgres struct {
		MaxPoolSize  int    `yaml:"max_pool_size"`
		ConnAttempts int    `yaml:"conn_attempts"`
		URL          string `yaml:"url" env:"POSTGRES_URL"`
	}
	Redis struct {
		Addr     string `yaml:"addr" env:"REDIS_ADDR"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	}
	Nats struct {
		URL string `yaml:"url" env:"NATS_URL"`
	}
	Clickhouse struct {
		Addr       string `yaml:"addr" env:"CH_ADDR"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		NativePort int    `yaml:"native_port"`
		HttpPort   int    `yaml:"http_port"`
		DB         string `yaml:"db"`
	}
)

func New(configPath string) (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("config Readconfig %w", err)
	}

	err = cleanenv.UpdateEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("config UpdateEnv %w", err)
	}

	return &cfg, nil
}
