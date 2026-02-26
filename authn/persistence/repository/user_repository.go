package repository

import (
	"context"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/model"
)

// IUserRepository defines the interface for user data operations
type IUserRepository interface {
	// CreateUser creates a new user in the database
	CreateUser(ctx context.Context, user *model.User, passwordHash string) (*model.User, error)

	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, id string) (*model.User, error)

	// GetUserByEmail retrieves a user by their email
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *model.User) error

	// DeleteUser soft deletes a user by setting deleted_at
	DeleteUser(ctx context.Context, id string) error

	// ListUsers retrieves a list of users with optional filtering
	ListUsers(ctx context.Context, params *ListUsersParams) ([]*model.User, error)
}

// ListUsersParams contains parameters for listing users
type ListUsersParams struct {
	Email             *string
	CreatedStartRange *time.Time
	CreatedEndRange   *time.Time
	Active            bool
}
