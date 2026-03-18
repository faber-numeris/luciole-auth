package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
)

type IAppConfig interface {
	IServiceConfig
	IDatabaseConfig
	IMailConfig
}

type AppConfig struct {
	ServiceConfig
	DatabaseConfig
	MailConfig
}

var GeneralConfig IAppConfig
var configErr error
var once sync.Once

func LoadConfig() (IAppConfig, error) {
	once.Do(func() {
		cfg, err := env.ParseAs[AppConfig]()
		if err != nil {
			configErr = fmt.Errorf("failed to parse environment variables: %w", err)
			return
		}
		GeneralConfig = &cfg
	})
	return GeneralConfig, configErr
}
