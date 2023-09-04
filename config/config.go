package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP HTTPServerConfig `yaml:"http_server"`
	PG   PGConfig         `yaml:"-"`
}

type HTTPServerConfig struct {
	Port         string        `yaml:"port"          env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout"  env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"5s"`
}

type PGConfig struct {
	Address  string `env:"POSTGRES_ADDR"     env-required:"true"`
	DBName   string `env:"POSTGRES_DB"       env-required:"true"`
	Username string `env:"POSTGRES_USER"     env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	SSLMode  string `env:"POSTGRES_SSLMODE"  env-default:"false"`
}

func New(configPath string) (Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return Config{}, fmt.Errorf("cleanenv - read config err: %w", err)
	}
	if err := cleanenv.ReadEnv(&cfg.PG); err != nil {
		return Config{}, fmt.Errorf("cleanenv - read env err: %w", err)
	}
	return cfg, nil
}
