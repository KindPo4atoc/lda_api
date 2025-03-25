package apiserver

import (
	repository "github.com/lda_api/internal/app/repository"
)

type Config struct {
	BindAddr string             `toml:"bind_addr"`
	LogLevel string             `toml:"log_level"`
	DBConfig *repository.Config `toml:"database"`
}

func NewConfig() *Config {
	return &Config{
		DBConfig: repository.NewConfig(),
	}
}
