package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type AppEnvConfig struct {
	Port             int64 `env:"APP_PORT"`
	GoogleAuthConfig GoogleAuthEnvConfig
	DbConfig         DbEnvConfig
}

type GoogleAuthEnvConfig struct {
	ClientID     string   `env:"GAUTH_CLIENT_ID"`
	ClientSecret string   `env:"GAUTH_CLIENT_SECRET"`
	RedirectURL  string   `env:"GAUTH_REDIRECT_URL"`
	Scopes       []string `env:"GAUTH_SCOPES"`
}

type DbEnvConfig struct {
	Driver   string `env:"DB_DRIVER"`
	Host     string `env:"DB_HOST"`
	Port     int64  `env:"DB_PORT"`
	Catalog  string `env:"DB_CATALOG"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
}

func NewAppConfigFromEnv(ctx context.Context) (*AppEnvConfig, error) {
	c := &AppEnvConfig{}
	err := envconfig.Process(ctx, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
