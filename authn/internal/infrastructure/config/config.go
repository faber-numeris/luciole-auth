package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
	_ "github.com/caarlos0/env/v11"
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
var once sync.Once

func LoadConfig() IAppConfig {
	once.Do(func() {
		cfg, err := env.ParseAs[AppConfig]()
		if err != nil {
			panic(fmt.Errorf("failed to parse environment variables: %w", err))
		}
		GeneralConfig = &cfg
	})
	return GeneralConfig
}
