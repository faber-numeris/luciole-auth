package service

import (
	"context"
	"fmt"

	"github.com/faber-numeris/luciole-auth/authn/model"
	"github.com/faber-numeris/luciole-auth/authn/model/generated"
	"github.com/faber-numeris/luciole-auth/authn/persistence/repository"
)

// UserServiceImpl implements the IUserService interface
type UserServiceImpl struct {
	userRepo       repository.IUserRepository
	hashingService IHashingService
	converter      *generated.ConverterImpl
}

// NewUserService creates a new instance of UserServiceImpl
func NewUserService(userRepo repository.IUserRepository, hashingService IHashingService) IUserService {
	return &UserServiceImpl{
		userRepo:       userRepo,
		hashingService: hashingService,
		converter:      &generated.ConverterImpl{},
	}
}

// RegisterUser creates a new user account
func (s *UserServiceImpl) RegisterUser(ctx context.Context, user *model.User) (*model.User, error) {
	passwordHash, err := s.hashingService.HashPassword(ctx, user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	createdUser, err := s.userRepo.CreateUser(ctx, user, passwordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

// GetUserByID retrieves a user by their ID
// TODO: Implement GetUserByID method
// assignees: rafaelsousa
func (s *UserServiceImpl) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	panic("not implemented")
}

// GetUserByUsername retrieves a user by their username
// TODO: Implement GetUserByUsername method
// assignees: rafaelsousa
func (s *UserServiceImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	panic("not implemented")
}

// GetUserByEmail retrieves a user by their email
// TODO: Implement GetUserByEmail method
// assignees: rafaelsousa
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	panic("not implemented")
}

// UpdateUserProfile updates an existing user's profile
// TODO: Implement UpdateUserProfile method
// assignees: rafaelsousa
func (s *UserServiceImpl) UpdateUserProfile(ctx context.Context, userID string, req *model.User) (*model.User, error) {
	panic("not implemented")
}

// DeleteUser deactivates a user account
// TODO: Implement DeleteUser method
// assignees: rafaelsousa
func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID string) error {
	panic("not implemented")
}

// ListUsers retrieves a list of users with optional filtering
// TODO: Implement ListUsers method
// assignees: rafaelsousa
func (s *UserServiceImpl) ListUsers(ctx context.Context, params *ListUsersParams) ([]*model.User, error) {
	panic("not implemented")
}
