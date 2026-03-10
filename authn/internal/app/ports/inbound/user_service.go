package inboundport

import (
	"context"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
)

// ListUsersParams contains parameters for listing users at the service level
type ListUsersParams struct {
	Email             *string
	CreatedStartRange *time.Time
	CreatedEndRange   *time.Time
	Active            bool
}

// UserService defines the interface for user business logic operations
type UserService interface {
	// RegisterUser creates a new user account
	RegisterUser(ctx context.Context, user *domain.User, password string) (*domain.User, error)

	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, id string) (*domain.User, error)

	// GetUserByEmail retrieves a user by their email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	// UpdateUserProfile updates an existing user's profile
	UpdateUserProfile(ctx context.Context, userID string, req *domain.User) (*domain.User, error)

	// DeleteUser deactivates a user account
	DeleteUser(ctx context.Context, userID string) error

	// ListUsers retrieves a list of users with optional filtering
	ListUsers(ctx context.Context, params *ListUsersParams) ([]*domain.User, error)

	// ConfirmUserRegistration confirms a user's email based on token
	ConfirmUserRegistration(ctx context.Context, token string) error

	// RequestPasswordReset generates a password reset token for the user
	RequestPasswordReset(ctx context.Context, email string) (string, error)

	// ResetPassword resets the user's password using the reset token
	ResetPassword(ctx context.Context, token string, newPassword string) error

	// VerifyPassword verifies if the provided password matches the user's password
	VerifyPassword(ctx context.Context, email string, password string) (*domain.User, error)
}
