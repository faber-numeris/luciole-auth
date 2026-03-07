package messaging

import (
	"context"

	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
)

type IMailService interface {
	SendConfirmationEmail(ctx context.Context, userConfirmation domain.UserConfirmation) error
}
