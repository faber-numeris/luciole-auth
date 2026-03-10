package mail

import (
	outboundport "github.com/faber-numeris/luciole-auth/authn/internal/app/ports/outbound"
	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/config"
)

// NewService creates a new mailer service
func NewService(cfg config.IMailConfig) outboundport.Mailer {
	return NewMailpit(cfg)
}
