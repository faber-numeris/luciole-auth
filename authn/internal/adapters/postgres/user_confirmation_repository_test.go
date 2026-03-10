package postgresadapter

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/postgres/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserConfirmationRepository_CreateUserConfirmation(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserConfirmationRepository(querier)
		querier.EXPECT().CreateUserConfirmation(ctx, mock.Anything).
			Return(gen.UserConfirmation{ID: "123", UserID: "user-123", Token: "token-123"}, nil)

		res, err := repo.CreateUserConfirmation(ctx, "user-123", "token-123", now)

		assert.NoError(t, err)
		assert.Equal(t, "token-123", res.Token)
	})

	t.Run("error", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserConfirmationRepository(querier)
		querier.EXPECT().CreateUserConfirmation(ctx, mock.Anything).
			Return(gen.UserConfirmation{}, errors.New("db error"))

		res, err := repo.CreateUserConfirmation(ctx, "user-123", "token-123", now)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserConfirmationRepository_GetUserConfirmationByToken(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserConfirmationRepository(querier)
		querier.EXPECT().GetUserConfirmationByToken(ctx, "token-123").
			Return(gen.UserConfirmation{UserID: "user-123"}, nil)

		res, err := repo.GetUserConfirmationByToken(ctx, "token-123")

		assert.NoError(t, err)
		assert.Equal(t, "user-123", res)
	})

	t.Run("not found", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserConfirmationRepository(querier)
		querier.EXPECT().GetUserConfirmationByToken(ctx, "unknown").
			Return(gen.UserConfirmation{}, sql.ErrNoRows)

		res, err := repo.GetUserConfirmationByToken(ctx, "unknown")

		assert.NoError(t, err)
		assert.Equal(t, "", res)
	})

	t.Run("error", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserConfirmationRepository(querier)
		querier.EXPECT().GetUserConfirmationByToken(ctx, "error").
			Return(gen.UserConfirmation{}, errors.New("db error"))

		res, err := repo.GetUserConfirmationByToken(ctx, "error")

		assert.Error(t, err)
		assert.Equal(t, "", res)
	})
}

func TestUserConfirmationRepository_ConfirmUserRegistration(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserConfirmationRepository(querier)
		querier.EXPECT().ConfirmUserRegistration(ctx, "user-123").Return(nil)

		err := repo.ConfirmUserRegistration(ctx, "user-123")

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserConfirmationRepository(querier)
		querier.EXPECT().ConfirmUserRegistration(ctx, "user-123").Return(errors.New("db error"))

		err := repo.ConfirmUserRegistration(ctx, "user-123")

		assert.Error(t, err)
	})
}

func TestUserConfirmationRepository_DeleteUserConfirmation(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserConfirmationRepository(querier)
		querier.EXPECT().DeleteUserConfirmation(ctx, "user-123").Return(nil)

		err := repo.DeleteUserConfirmation(ctx, "user-123")

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserConfirmationRepository(querier)
		querier.EXPECT().DeleteUserConfirmation(ctx, "user-123").Return(errors.New("db error"))

		err := repo.DeleteUserConfirmation(ctx, "user-123")

		assert.Error(t, err)
	})
}
