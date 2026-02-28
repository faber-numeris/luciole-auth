package configuration

import _ "github.com/caarlos0/env/v11"

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
