package postgresadapter

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/postgres/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/faber-numeris/luciole-auth/authn/internal/mocks"
	"github.com/faber-numeris/luciole-auth/authn/internal/app/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRepository_CreateUser(t *testing.T) {
	ctx := context.Background()
	user := &domain.User{Email: "test@example.com"}
	passwordHash := "hashed"

	t.Run("success", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		querier.EXPECT().CreateUser(ctx, mock.Anything).Return(gen.User{ID: "123", Email: user.Email}, nil)

		res, err := repo.CreateUser(ctx, user, passwordHash)

		assert.NoError(t, err)
		assert.Equal(t, "123", res.ID)
	})

	t.Run("error", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		querier.EXPECT().CreateUser(ctx, mock.Anything).Return(gen.User{}, errors.New("db error"))

		res, err := repo.CreateUser(ctx, user, passwordHash)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_GetUserByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		querier.EXPECT().GetUser(ctx, "123").Return(gen.User{ID: "123"}, nil)

		res, err := repo.GetUserByID(ctx, "123")

		assert.NoError(t, err)
		assert.Equal(t, "123", res.ID)
	})

	t.Run("not found", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		querier.EXPECT().GetUser(ctx, "404").Return(gen.User{}, sql.ErrNoRows)

		res, err := repo.GetUserByID(ctx, "404")

		assert.NoError(t, err)
		assert.Nil(t, res)
	})

	t.Run("error", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		querier.EXPECT().GetUser(ctx, "500").Return(gen.User{}, errors.New("db error"))

		res, err := repo.GetUserByID(ctx, "500")

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		querier.EXPECT().GetUserByEmail(ctx, "test@example.com").Return(gen.User{Email: "test@example.com"}, nil)

		res, err := repo.GetUserByEmail(ctx, "test@example.com")

		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", res.Email)
	})

	t.Run("not found", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		querier.EXPECT().GetUserByEmail(ctx, "unknown").Return(gen.User{}, sql.ErrNoRows)

		res, err := repo.GetUserByEmail(ctx, "unknown")

		assert.NoError(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_UpdateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		user := &domain.User{ID: "123", Profile: &domain.UserProfile{FirstName: "New"}}
		querier.EXPECT().UpdateUser(ctx, mock.Anything).Return(gen.User{}, nil)

		err := repo.UpdateUser(ctx, user)

		assert.NoError(t, err)
	})

	t.Run("success without profile", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		user := &domain.User{ID: "123", Profile: nil}
		querier.EXPECT().UpdateUser(ctx, mock.Anything).Return(gen.User{}, nil)

		err := repo.UpdateUser(ctx, user)

		assert.NoError(t, err)
	})
}

func TestUserRepository_DeleteUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		querier.EXPECT().DeleteUser(ctx, "123").Return(nil)

		err := repo.DeleteUser(ctx, "123")

		assert.NoError(t, err)
	})
}

func TestUserRepository_ListUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		querier.EXPECT().ListUsers(ctx, mock.Anything).Return([]gen.User{{ID: "1"}, {ID: "2"}}, nil)

		res, err := repo.ListUsers(ctx, &ports.ListUsersParams{})

		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("error", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		querier.EXPECT().ListUsers(ctx, mock.Anything).Return(nil, errors.New("db error"))

		res, err := repo.ListUsers(ctx, &ports.ListUsersParams{})

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		querier := mocks.NewMockQuerier(t)
		repo := NewUserRepository(querier)
		querier.EXPECT().UpdatePassword(ctx, mock.Anything).Return(nil)

		err := repo.UpdatePassword(ctx, "123", []byte("hash"))

		assert.NoError(t, err)
	})
}
