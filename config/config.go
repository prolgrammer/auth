package config

import (
	"auth/config/pg"
	"time"
)

type (
	Config struct {
		TokenConfig `mapstructure:"token_config"`
		HTTP        `mapstructure:"http"`
		JWT         `mapstructure:"jwt"`
		PG          pg.Config `mapstructure:"pg"`
	}

	TokenConfig struct {
		AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
		RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
	}

	HTTP struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	JWT struct {
		SecretKey string `mapstructure:"secret_key"`
	}
)
