package service

import (
	"context"

	"github.com/faber-numeris/luciole-auth/authn/model"
)

// IUserService defines the interface for user business logic operations
type IUserService interface {
	// RegisterUser creates a new user account
	RegisterUser(ctx context.Context, req *model.User) (*model.User, error)

	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, id string) (*model.User, error)

	// GetUserByUsername retrieves a user by their username
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	// GetUserByEmail retrieves a user by their email
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	// UpdateUserProfile updates an existing user's profile
	UpdateUserProfile(ctx context.Context, userID string, req *model.User) (*model.User, error)

	// DeleteUser deactivates a user account
	DeleteUser(ctx context.Context, userID string) error

	// ListUsers retrieves a list of users with optional filtering
	ListUsers(ctx context.Context, params *ListUsersParams) ([]*model.User, error)
}

// ListUsersParams contains parameters for listing users at the service level
type ListUsersParams struct {
	Username          *string
	Email             *string
	CreatedStartRange *string
	CreatedEndRange   *string
	Active            bool
}
