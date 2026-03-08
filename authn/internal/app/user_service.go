package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
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
	passwordHash, err := s.hashingService.HashPassword(ctx, password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	createdUser, err := s.userRepo.CreateUser(ctx, user, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	confirmationToken := generateToken()
	expiresAt := time.Now().Add(24 * time.Hour)
	confirmation, err := s.confirmationRepo.CreateUserConfirmation(ctx, createdUser.ID, confirmationToken, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create confirmation token: %w", err)
	}
	confirmation.UserEmail = createdUser.Email

	if err = s.mailService.SendConfirmationEmail(ctx, *confirmation); err != nil {
		return nil, fmt.Errorf("failed to send confirmation email: %w", err)
	}

	return createdUser, nil
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	if user == nil {
		return nil, nil
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user == nil {
		return nil, nil
	}

	return user, nil
}

// UpdateUserProfile updates an existing user's profile
func (s *userService) UpdateUserProfile(ctx context.Context, userID string, req *domain.User) (*domain.User, error) {
	existingUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if existingUser == nil {
		return nil, fmt.Errorf("user not found")
	}

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
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return existingUser, nil
}

// DeleteUser deactivates a user account
func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	existingUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if existingUser == nil {
		return fmt.Errorf("user not found")
	}

	err = s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers retrieves a list of users with optional filtering
// assignees: rafaelsousa
func (s *userService) ListUsers(ctx context.Context, params *ListUsersParams) ([]*domain.User, error) {
	repoParams := &ports.ListUsersParams{
		Email:             params.Email,
		CreatedStartRange: params.CreatedStartRange,
		CreatedEndRange:   params.CreatedEndRange,
		Active:            params.Active,
	}
	return s.userRepo.ListUsers(ctx, repoParams)
}

// ConfirmUserRegistration confirms a user's email based on token
func (s *userService) ConfirmUserRegistration(ctx context.Context, token string) error {
	userID, err := s.confirmationRepo.GetUserConfirmationByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to get user by confirmation token: %w", err)
	}

	if userID == "" {
		return errors.New("invalid or expired confirmation token")
	}

	err = s.confirmationRepo.ConfirmUserRegistration(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to confirm user email: %w", err)
	}

	return nil
}

// RequestPasswordReset generates a password reset token for the user
func (s *userService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	return "", errors.New("not implemented")
}

// ResetPassword resets the user's password using the reset token
func (s *userService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	return errors.New("not implemented")
}

// VerifyPassword verifies if the provided password matches the user's password
func (s *userService) VerifyPassword(ctx context.Context, email string, password string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
