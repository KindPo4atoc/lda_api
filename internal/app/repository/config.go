package repository

// конфиг для БД, подтягивается из apiserver.toml
type Config struct {
	DatabaseURL string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}
