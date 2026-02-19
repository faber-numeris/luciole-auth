package repository

import (
	"context"

	"github.com/faber-numeris/luciole-auth/authn/model"
)

// IUserRepository defines the interface for user data operations
type IUserRepository interface {
	// CreateUser creates a new user in the database
	CreateUser(ctx context.Context, user *model.User, passwordHash string) (*model.User, error)

	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, id string) (*model.User, error)

	// GetUserByUsername retrieves a user by their username
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

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
	Username          *string
	Email             *string
	CreatedStartRange *string
	CreatedEndRange   *string
	Active            bool
}
