package repository

import (
	"context"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/model"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *model.User, passwordHash string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, params *ListUsersParams) ([]*model.User, error)
	UpdatePassword(ctx context.Context, userID string, passwordHash []byte) error
}

type ListUsersParams struct {
	Email             *string
	CreatedStartRange *time.Time
	CreatedEndRange   *time.Time
	Active            bool
}

type IUserConfirmationRepository interface {
	CreateUserConfirmation(ctx context.Context, userID string, token string, expiresAt time.Time) (*model.UserConfirmation, error)
	GetUserConfirmationByToken(ctx context.Context, token string) (string, error)
	ConfirmUserRegistration(ctx context.Context, userID string) error
	DeleteUserConfirmation(ctx context.Context, userID string) error
}
