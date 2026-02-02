package configuration

type IDatabaseConfig interface {
	DBHost() string
	DBPort() int
	DBUser() string
	DBPassword() string
	DBName() string
	DBSSLMode() string
}

var _ IDatabaseConfig = (*DatabaseConfig)(nil)

type DatabaseConfig struct {
	DBHost_     string `env:"DB_HOST,required"`
	DBPort_     int    `env:"DB_PORT,required"`
	DBUser_     string `env:"DB_USER,required"`
	DBPassword_ string `env:"DB_PASSWORD,required"`
	DBDBName_   string `env:"DB_NAME,required"`
	DBSSLMode_  string `env:"DB_SSLMODE,required"`
}

func (d DatabaseConfig) DBHost() string {
	return d.DBHost_
}

func (d DatabaseConfig) DBPort() int {
	return d.DBPort_
}

func (d DatabaseConfig) DBUser() string {
	return d.DBUser_
}

func (d DatabaseConfig) DBPassword() string {
	return d.DBPassword_
}

func (d DatabaseConfig) DBName() string {
	return d.DBDBName_
}

func (d DatabaseConfig) DBSSLMode() string {
	return d.DBSSLMode_
}
