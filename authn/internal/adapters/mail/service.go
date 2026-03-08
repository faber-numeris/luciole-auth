package mail

import (
	"github.com/faber-numeris/luciole-auth/authn/internal/app/ports"
	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/config"
)

// NewService creates a new mailer service
func NewService(cfg config.IMailConfig) ports.Mailer {
	return NewMailpit(cfg)
}
