package postgresadapter

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/postgres/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/app/ports"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/faber-numeris/luciole-auth/authn/internal/platform/mapper/generated"
)

var confirmationConversion = generated.ConverterImpl{}

type userConfirmationRepository struct {
	querier gen.Querier
}

func NewUserConfirmationRepository(querier gen.Querier) ports.UserConfirmationRepository {
	return &userConfirmationRepository{
		querier: querier,
	}
}

func (r *userConfirmationRepository) CreateUserConfirmation(ctx context.Context, userID string, token string, expiresAt time.Time) (*domain.UserConfirmation, error) {
	sqlcConfirmation, err := r.querier.CreateUserConfirmation(ctx, gen.CreateUserConfirmationParams{
		Userid:    userID,
		Token:     token,
		Expiresat: expiresAt,
	})
	if err != nil {
		return nil, err
	}

	result, err := confirmationConversion.UserConfirmationModelFromSQLC(sqlcConfirmation)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *userConfirmationRepository) GetUserConfirmationByToken(ctx context.Context, token string) (string, error) {
	confirmation, err := r.querier.GetUserConfirmationByToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return confirmation.UserID, nil
}

func (r *userConfirmationRepository) ConfirmUserRegistration(ctx context.Context, userID string) error {
	return r.querier.ConfirmUserRegistration(ctx, userID)
}

func (r *userConfirmationRepository) DeleteUserConfirmation(ctx context.Context, userID string) error {
	return r.querier.DeleteUserConfirmation(ctx, userID)
}
