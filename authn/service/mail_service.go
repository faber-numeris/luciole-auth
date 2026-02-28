package service

import (
	"context"
)

type IMailService interface {
	SendConfirmationEmail(ctx context.Context, email, code string) error
}
