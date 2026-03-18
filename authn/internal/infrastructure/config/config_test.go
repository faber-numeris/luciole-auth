package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("ServiceConfig", func(t *testing.T) {
		cfg := &ServiceConfig{Port_: 9090, AllowedOrigins_: []string{"*"}}
		assert.Equal(t, 9090, cfg.Port())
		assert.Equal(t, []string{"*"}, cfg.AllowedOrigins())
	})

	t.Run("DatabaseConfig", func(t *testing.T) {
		cfg := &DatabaseConfig{
			DBHost_:     "host",
			DBPort_:     5432,
			DBUser_:     "user",
			DBPassword_: "pass",
			DBDBName_:   "name",
			DBSSLMode_:  "disable",
		}
		assert.Equal(t, "host", cfg.DBHost())
		assert.Equal(t, 5432, cfg.DBPort())
		assert.Equal(t, "user", cfg.DBUser())
		assert.Equal(t, "pass", cfg.DBPassword())
		assert.Equal(t, "name", cfg.DBName())
		assert.Equal(t, "disable", cfg.DBSSLMode())
	})

	t.Run("MailConfig", func(t *testing.T) {
		cfg := &MailConfig{
			MailFrom_:              "from",
			MailFromName_:          "name",
			SMTPHost_:              "host",
			SMTPPort_:              1025,
			ConfirmationURLFormat_: "fmt",
		}
		assert.Equal(t, "from", cfg.MailFrom())
		assert.Equal(t, "name", cfg.MailFromName())
		assert.Equal(t, "host", cfg.SMTPHost())
		assert.Equal(t, 1025, cfg.SMTPPort())
		assert.Equal(t, "fmt", cfg.ConfirmationURLFormat())
	})

	t.Run("LoadConfig success", func(t *testing.T) {
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_USER", "user")
		os.Setenv("DB_PASSWORD", "pass")
		os.Setenv("DB_NAME", "db")
		os.Setenv("DB_SSLMODE", "disable")
		os.Setenv("MAIL_FROM", "from")
		os.Setenv("MAIL_FROM_NAME", "name")
		os.Setenv("SMTP_HOST", "host")
		os.Setenv("CONFIRMATION_URL_FORMAT", "fmt")

		defer func() {
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("DB_NAME")
			os.Unsetenv("DB_SSLMODE")
			os.Unsetenv("MAIL_FROM")
			os.Unsetenv("MAIL_FROM_NAME")
			os.Unsetenv("SMTP_HOST")
			os.Unsetenv("CONFIRMATION_URL_FORMAT")
		}()

		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.Equal(t, "localhost", cfg.DBHost())
		assert.Equal(t, "from", cfg.MailFrom())
		assert.Equal(t, 8080, cfg.Port())
	})

	t.Run("LoadConfig error", func(t *testing.T) {
		// Ensure environment is clean of required variables
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("MAIL_FROM")
		os.Unsetenv("MAIL_FROM_NAME")
		os.Unsetenv("SMTP_HOST")
		os.Unsetenv("CONFIRMATION_URL_FORMAT")

		_, err := LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required environment variable")
	})
}
