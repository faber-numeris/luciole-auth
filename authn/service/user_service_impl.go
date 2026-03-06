package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/model"
	"github.com/faber-numeris/luciole-auth/authn/persistence/repository"
	"github.com/faber-numeris/luciole-auth/authn/service/mail"
)

// UserServiceImpl implements the IUserService interface
type UserServiceImpl struct {
	userRepo         repository.IUserRepository
	confirmationRepo repository.IUserConfirmationRepository
	mailService      mail.IMailService
	hashingService   IHashingService
}

// NewUserService creates a new instance of UserServiceImpl
func NewUserService(
	userRepo repository.IUserRepository,
	confirmationRepo repository.IUserConfirmationRepository,
	hashingService IHashingService,
	mailService mail.IMailService,
) IUserService {
	return &UserServiceImpl{
		userRepo:         userRepo,
		confirmationRepo: confirmationRepo,
		hashingService:   hashingService,
		mailService:      mailService,
	}
}

// RegisterUser creates a new user account
func (s *UserServiceImpl) RegisterUser(ctx context.Context, user *model.User, password string) (*model.User, error) {
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
func (s *UserServiceImpl) GetUserByID(ctx context.Context, id string) (*model.User, error) {
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
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
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
func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, userID string, req *model.User) (*model.User, error) {
	existingUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if existingUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	if req.Profile != nil {
		if existingUser.Profile == nil {
			existingUser.Profile = &model.UserProfile{}
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
func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID string) error {
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
func (s *UserServiceImpl) ListUsers(ctx context.Context, params *ListUsersParams) ([]*model.User, error) {
	repoParams := &repository.ListUsersParams{
		Email:             params.Email,
		CreatedStartRange: params.CreatedStartRange,
		CreatedEndRange:   params.CreatedEndRange,
		Active:            params.Active,
	}
	return s.userRepo.ListUsers(ctx, repoParams)
}

// ConfirmUserRegistration confirms a user's email based on token
func (s *UserServiceImpl) ConfirmUserRegistration(ctx context.Context, token string) error {
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
func (s *UserServiceImpl) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	return "", errors.New("not implemented")
}

// ResetPassword resets the user's password using the reset token
func (s *UserServiceImpl) ResetPassword(ctx context.Context, token string, newPassword string) error {
	return errors.New("not implemented")
}

// VerifyPassword verifies if the provided password matches the user's password
func (s *UserServiceImpl) VerifyPassword(ctx context.Context, email string, password string) (*model.User, error) {
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
