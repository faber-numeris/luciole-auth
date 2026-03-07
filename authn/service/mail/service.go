package mail

import (
	"context"

	"github.com/faber-numeris/luciole-auth/authn/model"
)

type IMailService interface {
	SendConfirmationEmail(ctx context.Context, userConfirmation model.UserConfirmation) error
}
