package configuration

import "github.com/caarlos0/env/v11"

type IServiceConfig interface {
	Port() int
}

type ServiceConfig struct {
	Port_ int `env:"PORT" envDefault:"8080"`
}

func (c *ServiceConfig) Port() int {
	return c.Port_
}

type IMailConfig interface {
	MailHost() string
	MailPort() int
	MailUsername() string
	MailPassword() string
	MailFrom() string
}

type MailConfig struct {
	MailHost_     string `env:"MAIL_HOST,required"`
	MailPort_     int    `env:"MAIL_PORT" envDefault:"1025"`
	MailUsername_ string `env:"MAIL_USERNAME"`
	MailPassword_ string `env:"MAIL_PASSWORD"`
	MailFrom_     string `env:"MAIL_FROM" envDefault:"noreply@example.com"`
}

func (c MailConfig) MailHost() string {
	return c.MailHost_
}

func (c MailConfig) MailPort() int {
	return c.MailPort_
}

func (c MailConfig) MailUsername() string {
	return c.MailUsername_
}

func (c MailConfig) MailPassword() string {
	return c.MailPassword_
}

func (c MailConfig) MailFrom() string {
	return c.MailFrom_
}

func Load() (IAppConfig, error) {
	cfg, err := env.ParseAs[AppConfig]()
	return &cfg, err
}
