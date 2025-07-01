package config

import (
	"auth/config/pg"
	"github.com/spf13/viper"

	"fmt"
	"os"
	"time"
)

type (
	Config struct {
		TokenConfig `mapstructure:"token_config"`
		HTTP        `mapstructure:"http"`
		JWT         `mapstructure:"jwt"`
		PG          pg.Config `mapstructure:"pg"`
		Cookie      `mapstructure:"cookie"`
	}

	TokenConfig struct {
		AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
		RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
	}

	HTTP struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}

	JWT struct {
		SecretKey string `mapstructure:"secret_key"`
	}

	Cookie struct {
		Domain   string `mapstructure:"domain"`
		Path     string `mapstructure:"path"`
		Secure   bool   `mapstructure:"secure"`
		HttpOnly bool   `mapstructure:"http_only"`
		SameSite string `mapstructure:"same_site"`
	}
)

func New() (*Config, error) {
	cfg := Config{}
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config/")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	for _, k := range v.AllKeys() {
		anyValue := v.Get(k)
		str, ok := anyValue.(string)
		if !ok {
			continue
		}

		replaced := os.ExpandEnv(str)
		v.Set(k, replaced)
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling file: %w", err))
	}

	return &cfg, nil
}
