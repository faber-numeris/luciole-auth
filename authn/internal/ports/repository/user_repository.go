package repository

import (
	"context"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *domain.User, passwordHash string) (*domain.User, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, params *ListUsersParams) ([]*domain.User, error)
	UpdatePassword(ctx context.Context, userID string, passwordHash []byte) error
}

type ListUsersParams struct {
	Email             *string
	CreatedStartRange *time.Time
	CreatedEndRange   *time.Time
	Active            bool
}
