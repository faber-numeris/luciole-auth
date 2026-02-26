package service

import (
	"context"
	"fmt"

	"github.com/faber-numeris/luciole-auth/authn/model"
	"github.com/faber-numeris/luciole-auth/authn/persistence/repository"
)

// UserServiceImpl implements the IUserService interface
type UserServiceImpl struct {
	userRepo       repository.IUserRepository
	hashingService IHashingService
}

// NewUserService creates a new instance of UserServiceImpl
func NewUserService(userRepo repository.IUserRepository, hashingService IHashingService) IUserService {
	return &UserServiceImpl{
		userRepo:       userRepo,
		hashingService: hashingService,
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
