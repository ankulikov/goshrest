package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewGAuthConfig(conf GoogleAuthEnvConfig) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  conf.RedirectURL,
		Scopes:       conf.Scopes,
	}
}
