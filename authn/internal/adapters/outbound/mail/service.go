package mail

import (
	outboundport "github.com/faber-numeris/luciole-auth/authn/internal/ports/outbound"
	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/config"
)

// NewService creates a new mailer service
func NewService() outboundport.Mailer {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	return NewMailpit(cfg)
}
