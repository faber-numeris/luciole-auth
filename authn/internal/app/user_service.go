package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/app/ports"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
)

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

// ListUsersParams contains parameters for listing users at the service level
type ListUsersParams struct {
	Email             *string
	CreatedStartRange *time.Time
	CreatedEndRange   *time.Time
	Active            bool
}

// userService implements the UserService interface
type userService struct {
	userRepo         ports.UserRepository
	confirmationRepo ports.UserConfirmationRepository
	mailService      ports.Mailer
	hashingService   HashingService
}

// NewUserService creates a new instance of UserService
func NewUserService(
	userRepo ports.UserRepository,
	confirmationRepo ports.UserConfirmationRepository,
	hashingService HashingService,
	mailService ports.Mailer,
) UserService {
	return &userService{
		userRepo:         userRepo,
		confirmationRepo: confirmationRepo,
		hashingService:   hashingService,
		mailService:      mailService,
	}
}

// RegisterUser creates a new user account
func (s *userService) RegisterUser(ctx context.Context, user *domain.User, password string) (*domain.User, error) {
	slog.Info("Registering new user", "user", user.Email)
	passwordHash, err := s.hashingService.HashPassword(ctx, password)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	slog.Debug("Password hashed successfully")

	createdUser, err := s.userRepo.CreateUser(ctx, user, passwordHash)
	if err != nil {
		slog.Error("Failed to create user in repository", "error", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	slog.Info("User created successfully in repository", "userID", createdUser.ID)

	confirmationToken := generateToken()
	expiresAt := time.Now().Add(24 * time.Hour)
	confirmation, err := s.confirmationRepo.CreateUserConfirmation(ctx, createdUser.ID, confirmationToken, expiresAt)
	if err != nil {
		slog.Error("Failed to create confirmation token", "error", err)
		return nil, fmt.Errorf("failed to create confirmation token: %w", err)
	}
	slog.Info("Confirmation token created successfully", "token", confirmationToken)
	confirmation.UserEmail = createdUser.Email

	if err = s.mailService.SendConfirmationEmail(ctx, *confirmation); err != nil {
		slog.Error("Failed to send confirmation email", "error", err)
		return nil, fmt.Errorf("failed to send confirmation email: %w", err)
	}
	slog.Info("Confirmation email sent successfully", "email", createdUser.Email)

	return createdUser, nil
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	slog.Info("Getting user by ID", "id", id)
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		slog.Error("Failed to get user by ID", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	if user == nil {
		slog.Warn("User not found by ID", "id", id)
		return nil, nil
	}

	slog.Info("User found by ID", "id", id)
	return user, nil
}

// GetUserByEmail retrieves a user by their email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	slog.Info("Getting user by email", "email", email)
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		slog.Error("Failed to get user by email", "email", email, "error", err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user == nil {
		slog.Warn("User not found by email", "email", email)
		return nil, nil
	}

	slog.Info("User found by email", "email", email)
	return user, nil
}

// UpdateUserProfile updates an existing user's profile
func (s *userService) UpdateUserProfile(ctx context.Context, userID string, req *domain.User) (*domain.User, error) {
	slog.Info("Updating user profile", "userID", userID)
	existingUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		slog.Error("Failed to get user for update", "userID", userID, "error", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if existingUser == nil {
		slog.Warn("User not found for update", "userID", userID)
		return nil, fmt.Errorf("user not found")
	}
	slog.Debug("User found for update", "userID", userID)

	if req.Profile != nil {
		if existingUser.Profile == nil {
			existingUser.Profile = &domain.UserProfile{}
		}
		if req.Profile.FirstName != "" {
			existingUser.Profile.FirstName = req.Profile.FirstName
		}
		if req.Profile.LastName != "" {
			existingUser.Profile.LastName = req.Profile.LastName
		}
		if req.Profile.Locale != "" {
			existingUser.Profile.Locale = req.Profile.Locale
		}
		if req.Profile.Timezone != "" {
			existingUser.Profile.Timezone = req.Profile.Timezone
		}
	}

	err = s.userRepo.UpdateUser(ctx, existingUser)
	if err != nil {
		slog.Error("Failed to update user in repository", "userID", userID, "error", err)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	slog.Info("User profile updated successfully", "userID", userID)

	return existingUser, nil
}

// DeleteUser deactivates a user account
func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	slog.Info("Deleting user", "userID", userID)
	existingUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		slog.Error("Failed to get user for deletion", "userID", userID, "error", err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	if existingUser == nil {
		slog.Warn("User not found for deletion", "userID", userID)
		return fmt.Errorf("user not found")
	}
	slog.Debug("User found for deletion", "userID", userID)

	err = s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		slog.Error("Failed to delete user in repository", "userID", userID, "error", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	slog.Info("User deleted successfully", "userID", userID)

	return nil
}

// ListUsers retrieves a list of users with optional filtering
func (s *userService) ListUsers(ctx context.Context, params *ListUsersParams) ([]*domain.User, error) {
	slog.Info("Listing users", "params", params)
	repoParams := &ports.ListUsersParams{
		Email:             params.Email,
		CreatedStartRange: params.CreatedStartRange,
		CreatedEndRange:   params.CreatedEndRange,
		Active:            params.Active,
	}
	users, err := s.userRepo.ListUsers(ctx, repoParams)
	if err != nil {
		slog.Error("Failed to list users from repository", "error", err)
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	slog.Info("Users listed successfully", "count", len(users))
	return users, nil
}

// ConfirmUserRegistration confirms a user's email based on token
func (s *userService) ConfirmUserRegistration(ctx context.Context, token string) error {
	slog.Info("Confirming user registration")
	userID, err := s.confirmationRepo.GetUserConfirmationByToken(ctx, token)
	if err != nil {
		slog.Error("Failed to get user by confirmation token", "error", err)
		return fmt.Errorf("failed to get user by confirmation token: %w", err)
	}

	if userID == "" {
		slog.Warn("Invalid or expired confirmation token used")
		return errors.New("invalid or expired confirmation token")
	}
	slog.Debug("User confirmation found", "userID", userID)

	err = s.confirmationRepo.ConfirmUserRegistration(ctx, userID)
	if err != nil {
		slog.Error("Failed to confirm user email", "userID", userID, "error", err)
		return fmt.Errorf("failed to confirm user email: %w", err)
	}
	slog.Info("User email confirmed successfully", "userID", userID)

	return nil
}

// RequestPasswordReset generates a password reset token for the user
func (s *userService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	slog.Info("Requesting password reset", "email", email)
	return "", errors.New("not implemented")
}

// ResetPassword resets the user's password using the reset token
func (s *userService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	slog.Info("Resetting password")
	return errors.New("not implemented")
}

// VerifyPassword verifies if the provided password matches the user's password
func (s *userService) VerifyPassword(ctx context.Context, email string, password string) (*domain.User, error) {
	slog.Info("Verifying password", "email", email)
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		slog.Error("Failed to get user for password verification", "email", email, "error", err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user == nil {
		slog.Warn("Invalid credentials: user not found", "email", email)
		return nil, errors.New("invalid credentials")
	}
	slog.Debug("User found for password verification", "email", email)

	return user, nil
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
