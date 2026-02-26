package service

import (
	"context"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/model"
)

// IUserService defines the interface for user business logic operations
type IUserService interface {
	// RegisterUser creates a new user account
	RegisterUser(ctx context.Context, user *model.User, password string) (*model.User, error)

	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, id string) (*model.User, error)

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
	Email             *string
	CreatedStartRange *time.Time
	CreatedEndRange   *time.Time
	Active            bool
}
