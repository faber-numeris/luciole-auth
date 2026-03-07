package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/persistence/postgres/sqlc"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/faber-numeris/luciole-auth/authn/internal/platform/mapper/generated"
	"github.com/faber-numeris/luciole-auth/authn/internal/ports/repository"
)

type SQLCUserRepository struct {
	querier sqlc.Querier
}

var conversion = generated.ConverterImpl{}

func NewSQLCUserRepository(querier sqlc.Querier) *SQLCUserRepository {
	return &SQLCUserRepository{
		querier: querier,
	}
}

func (r *SQLCUserRepository) CreateUser(ctx context.Context, user *domain.User, passwordHash string) (*domain.User, error) {
	createParams := sqlc.CreateUserParams{
		Email:        user.Email,
		PasswordHash: []byte(passwordHash),
	}

	sqlcUser, err := r.querier.CreateUser(ctx, createParams)
	if err != nil {
		return nil, err
	}

	result, err := conversion.UserModelFromSQLC(sqlcUser)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *SQLCUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	sqlcUser, err := r.querier.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	result, err := conversion.UserModelFromSQLC(sqlcUser)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *SQLCUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	sqlcUser, err := r.querier.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	result, err := conversion.UserModelFromSQLC(sqlcUser)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *SQLCUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	var firstName, lastName, locale, timezone string
	if user.Profile != nil {
		firstName = user.Profile.FirstName
		lastName = user.Profile.LastName
		locale = user.Profile.Locale
		timezone = user.Profile.Timezone
	}

	updateParams := sqlc.UpdateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: []byte{},
		FirstName:    firstName,
		LastName:     lastName,
		Locale:       locale,
		Timezone:     timezone,
	}

	_, err := r.querier.UpdateUser(ctx, updateParams)
	return err
}

func (r *SQLCUserRepository) DeleteUser(ctx context.Context, id string) error {
	return r.querier.DeleteUser(ctx, id)
}

func (r *SQLCUserRepository) ListUsers(ctx context.Context, params *repository.ListUsersParams) ([]*domain.User, error) {
	sqlcParams := sqlc.ListUsersParams{
		Active:            params.Active,
		Email:             params.Email,
		CreatedStartRange: params.CreatedStartRange,
		CreatedEndRange:   params.CreatedEndRange,
	}

	sqlcUsers, err := r.querier.ListUsers(ctx, sqlcParams)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.User, len(sqlcUsers))
	for i, sqlcUser := range sqlcUsers {
		res, err := conversion.UserModelFromSQLC(sqlcUser)
		if err != nil {
			return nil, err
		}
		result[i] = &res
	}

	return result, nil
}

func (r *SQLCUserRepository) UpdatePassword(ctx context.Context, userID string, passwordHash []byte) error {
	return r.querier.UpdatePassword(ctx, sqlc.UpdatePasswordParams{
		Userid:       userID,
		Passwordhash: passwordHash,
	})
}
