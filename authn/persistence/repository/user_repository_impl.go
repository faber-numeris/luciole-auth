package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/faber-numeris/luciole-auth/authn/model"
	"github.com/faber-numeris/luciole-auth/authn/model/generated"
	"github.com/faber-numeris/luciole-auth/authn/persistence/sqlc"
)

type SQLCUserRepository struct {
	querier sqlc.Querier
}

var conversion = generated.ConverterImpl{}

func NewSQLCUserRepository(querier sqlc.Querier) IUserRepository {
	return &SQLCUserRepository{
		querier: querier,
	}
}

func (r *SQLCUserRepository) CreateUser(ctx context.Context, user *model.User, passwordHash string) (*model.User, error) {
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

func (r *SQLCUserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
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

func (r *SQLCUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
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

func (r *SQLCUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	updateParams := sqlc.UpdateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: []byte{},
	}

	_, err := r.querier.UpdateUser(ctx, updateParams)
	return err
}

func (r *SQLCUserRepository) DeleteUser(ctx context.Context, id string) error {
	return r.querier.DeleteUser(ctx, id)
}

func (r *SQLCUserRepository) ListUsers(ctx context.Context, params *ListUsersParams) ([]*model.User, error) {
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

	result := make([]*model.User, len(sqlcUsers))
	for i, sqlcUser := range sqlcUsers {
		res, err := conversion.UserModelFromSQLC(sqlcUser)
		if err != nil {
			return nil, err
		}
		result[i] = &res
	}

	return result, nil
}
