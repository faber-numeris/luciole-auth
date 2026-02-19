package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/model"
	"github.com/faber-numeris/luciole-auth/authn/persistence/sqlc"
)

// SQLCUserRepository implements UserRepository using sqlc generated components
type SQLCUserRepository struct {
	querier sqlc.Querier
}

// NewSQLCUserRepository creates a new instance of SQLCUserRepository
func NewSQLCUserRepository(querier sqlc.Querier) UserRepository {
	return &SQLCUserRepository{
		querier: querier,
	}
}

// CreateUser creates a new user in the database
func (r *SQLCUserRepository) CreateUser(ctx context.Context, user *model.User, passwordHash string) (*model.User, error) {
	// Convert domain model to sqlc parameter
	createParams := sqlc.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: []byte(passwordHash),
	}

	// Execute the create query
	sqlcUser, err := r.querier.CreateUser(ctx, createParams)
	if err != nil {
		return nil, err
	}

	// Convert sqlc result back to domain model
	result := &model.User{}
	result.FromSQLC(sqlcUser)
	return result, nil
}

// GetUserByID retrieves a user by their ID
func (r *SQLCUserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	sqlcUser, err := r.querier.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	result := &model.User{}
	result.FromSQLC(sqlcUser)
	return result, nil
}

// GetUserByUsername retrieves a user by their username
func (r *SQLCUserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	sqlcUser, err := r.querier.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	result := &model.User{}
	result.FromSQLC(sqlcUser)
	return result, nil
}

// GetUserByEmail retrieves a user by their email
func (r *SQLCUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	sqlcUser, err := r.querier.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	result := &model.User{}
	result.FromSQLC(sqlcUser)
	return result, nil
}

// UpdateUser updates an existing user
func (r *SQLCUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	// Convert domain model to sqlc parameter
	updateParams := sqlc.UpdateUserParams{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: []byte{}, // This would need to be handled separately for password updates
	}

	_, err := r.querier.UpdateUser(ctx, updateParams)
	return err
}

// DeleteUser soft deletes a user by setting deleted_at
func (r *SQLCUserRepository) DeleteUser(ctx context.Context, id string) error {
	return r.querier.DeleteUser(ctx, id)
}

// ListUsers retrieves a list of users with optional filtering
func (r *SQLCUserRepository) ListUsers(ctx context.Context, params *ListUsersParams) ([]*model.User, error) {
	// Convert domain parameters to sqlc parameters
	sqlcParams := sqlc.ListUsersParams{
		Active: params.Active,
	}

	if params.Username != nil {
		sqlcParams.Username = params.Username
	}
	if params.Email != nil {
		sqlcParams.Email = params.Email
	}
	if params.CreatedStartRange != nil {
		if parsedTime, err := time.Parse(time.RFC3339, *params.CreatedStartRange); err == nil {
			sqlcParams.CreatedStartRange = &parsedTime
		}
	}
	if params.CreatedEndRange != nil {
		if parsedTime, err := time.Parse(time.RFC3339, *params.CreatedEndRange); err == nil {
			sqlcParams.CreatedEndRange = &parsedTime
		}
	}

	sqlcUsers, err := r.querier.ListUsers(ctx, sqlcParams)
	if err != nil {
		return nil, err
	}

	// Convert sqlc results to domain models
	result := make([]*model.User, len(sqlcUsers))
	for i, sqlcUser := range sqlcUsers {
		result[i] = &model.User{}
		result[i].FromSQLC(sqlcUser)
	}

	return result, nil
}
