package mail

import (
	"github.com/faber-numeris/luciole-auth/authn/internal/ports/messaging"
)

// Verify at compile time that Mailpit implements IMailService
var _ messaging.IMailService = (*Mailpit)(nil)
