package outboundport

import (
	"context"

	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
)

type Mailer interface {
	SendConfirmationEmail(ctx context.Context, userConfirmation domain.UserConfirmation) error
}
