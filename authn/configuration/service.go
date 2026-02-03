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

func Load() (IAppConfig, error) {
	cfg, err := env.ParseAs[AppConfig]()
	return &cfg, err
}
