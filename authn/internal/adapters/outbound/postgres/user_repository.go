package postgresadapter

import (
	"context"
	"database/sql"
	"errors"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/outbound/postgres/gen"
	outboundport "github.com/faber-numeris/luciole-auth/authn/internal/app/ports/outbound"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/faber-numeris/luciole-auth/authn/internal/platform/mapper/generated"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type userRepository struct {
	querier gen.Querier
}

var conversion = generated.ConverterImpl{}

func NewUserRepository(querier gen.Querier) outboundport.UserRepository {
	return &userRepository{
		querier: querier,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User, passwordHash string) (*domain.User, error) {
	createParams := gen.CreateUserParams{
		Email:        user.Email,
		PasswordHash: []byte(passwordHash),
	}

	sqlcUser, err := r.querier.CreateUser(ctx, createParams)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) && pgerrcode.IsIntegrityConstraintViolation(pgerr.Code) {
			return nil, domain.ErrUserAlreadyExists
		}
		return nil, err
	}

	result, err := conversion.UserModelFromSQLC(sqlcUser)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
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

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
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

func (r *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	var firstName, lastName, locale, timezone string
	if user.Profile != nil {
		firstName = user.Profile.FirstName
		lastName = user.Profile.LastName
		locale = user.Profile.Locale
		timezone = user.Profile.Timezone
	}

	updateParams := gen.UpdateUserParams{
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

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	return r.querier.DeleteUser(ctx, id)
}

func (r *userRepository) ListUsers(ctx context.Context, params *outboundport.ListUsersParams) ([]*domain.User, error) {
	sqlcParams := gen.ListUsersParams{
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

func (r *userRepository) UpdatePassword(ctx context.Context, userID string, passwordHash []byte) error {
	return r.querier.UpdatePassword(ctx, gen.UpdatePasswordParams{
		Userid:       userID,
		Passwordhash: passwordHash,
	})
}

func (r *userRepository) GetUserCredentials(ctx context.Context, email string) (*domain.UserCredentials, error) {
	sqlcUser, err := r.querier.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	result, err := conversion.UserCredentialsModelFromSQLC(sqlcUser)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
