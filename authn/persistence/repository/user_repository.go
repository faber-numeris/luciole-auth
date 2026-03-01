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

	// CreateUserConfirmation creates a new user confirmation record
	CreateUserConfirmation(ctx context.Context, userID string, token string, expiresAt time.Time) (string, error)

	// GetUserConfirmationByToken retrieves a user confirmation by token
	GetUserConfirmationByToken(ctx context.Context, token string) (string, error)

	// ConfirmUserRegistration confirms a user's registration
	ConfirmUserRegistration(ctx context.Context, userID string) error

	// DeleteUserConfirmation deletes a user confirmation record
	DeleteUserConfirmation(ctx context.Context, userID string) error

	// SetPasswordResetToken sets a password reset token for a user
	SetPasswordResetToken(ctx context.Context, userID string, token string, expiresAt time.Time) error

	// GetUserByPasswordResetToken retrieves a user by password reset token
	GetUserByPasswordResetToken(ctx context.Context, token string) (*model.User, error)

	// UpdatePassword updates the user's password
	UpdatePassword(ctx context.Context, userID string, passwordHash []byte) error
}

// ListUsersParams contains parameters for listing users
type ListUsersParams struct {
	Email             *string
	CreatedStartRange *time.Time
	CreatedEndRange   *time.Time
	Active            bool
}
