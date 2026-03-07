package repository

import (
	"context"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
)

type IUserConfirmationRepository interface {
	CreateUserConfirmation(ctx context.Context, userID string, token string, expiresAt time.Time) (*domain.UserConfirmation, error)
	GetUserConfirmationByToken(ctx context.Context, token string) (string, error)
	ConfirmUserRegistration(ctx context.Context, userID string) error
	DeleteUserConfirmation(ctx context.Context, userID string) error
}
