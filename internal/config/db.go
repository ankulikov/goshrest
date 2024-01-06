package config

import "fmt"

type DbConfig struct {
	Driver   string `env:"DB_DRIVER"`
	Host     string `env:"DB_HOST"`
	Port     int64  `env:"DB_PORT"`
	Catalog  string `env:"DB_CATALOG"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
}

func (conf DbConfig) DSN() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		conf.Driver, conf.User, conf.Password, conf.Host, conf.Port, conf.Catalog)
}

func NewDbConfig(conf DbEnvConfig) DbConfig {
	return DbConfig(conf)
}
