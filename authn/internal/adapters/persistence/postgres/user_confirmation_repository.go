package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/persistence/postgres/sqlc"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/faber-numeris/luciole-auth/authn/internal/platform/mapper/generated"
	"github.com/faber-numeris/luciole-auth/authn/internal/ports/repository"
)

var confirmationConversion = generated.ConverterImpl{}

type SQLCUserConfirmationRepository struct {
	querier sqlc.Querier
}

func NewSQLCUserConfirmationRepository(querier sqlc.Querier) repository.IUserConfirmationRepository {
	return &SQLCUserConfirmationRepository{
		querier: querier,
	}
}

func (r *SQLCUserConfirmationRepository) CreateUserConfirmation(ctx context.Context, userID string, token string, expiresAt time.Time) (*domain.UserConfirmation, error) {
	sqlcConfirmation, err := r.querier.CreateUserConfirmation(ctx, sqlc.CreateUserConfirmationParams{
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

func (r *SQLCUserConfirmationRepository) GetUserConfirmationByToken(ctx context.Context, token string) (string, error) {
	confirmation, err := r.querier.GetUserConfirmationByToken(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return confirmation.UserID, nil
}

func (r *SQLCUserConfirmationRepository) ConfirmUserRegistration(ctx context.Context, userID string) error {
	return r.querier.ConfirmUserRegistration(ctx, userID)
}

func (r *SQLCUserConfirmationRepository) DeleteUserConfirmation(ctx context.Context, userID string) error {
	return r.querier.DeleteUserConfirmation(ctx, userID)
}
