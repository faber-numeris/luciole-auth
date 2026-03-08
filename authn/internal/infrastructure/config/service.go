package config

import "github.com/caarlos0/env/v11"

type IServiceConfig interface {
	Port() int
	AllowedOrigins() []string
}

type ServiceConfig struct {
	Port_           int      `env:"PORT" envDefault:"8080"`
	AllowedOrigins_ []string `env:"ALLOWED_ORIGINS" envDefault:"http://localhost:*,https://*"`
}

func (c *ServiceConfig) Port() int {
	return c.Port_
}

func (c *ServiceConfig) AllowedOrigins() []string {
	return c.AllowedOrigins_
}

func Load() (IAppConfig, error) {
	cfg, err := env.ParseAs[AppConfig]()
	return &cfg, err
}
