package postgresadapter

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	outboundport "github.com/faber-numeris/luciole-auth/authn/internal/ports/outbound"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")
	return sqlxDB, mock
}

func TestUserRepository_CreateUser(t *testing.T) {
	sqlxDB, mock := setupTestDB(t)
	repo := NewUserRepository(sqlxDB)
	ctx := context.Background()
	user := &domain.User{Email: "test@example.com"}
	passwordHash := []byte("hashed")

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "first_name", "last_name", "locale", "timezone", "created_at", "updated_at", "deleted_at"}).
			AddRow("123", user.Email, passwordHash, "", "", "", "", time.Now(), time.Now(), nil)

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Email, passwordHash).
			WillReturnRows(rows)

		res, err := repo.CreateUser(ctx, user, passwordHash)

		assert.NoError(t, err)
		assert.Equal(t, user.Email, res.Email)
		assert.Equal(t, "123", res.ID)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(user.Email, passwordHash).
			WillReturnError(errors.New("db error"))

		res, err := repo.CreateUser(ctx, user, passwordHash)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_GetUserByID(t *testing.T) {
	sqlxDB, mock := setupTestDB(t)
	repo := NewUserRepository(sqlxDB)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "first_name", "last_name", "locale", "timezone", "created_at", "updated_at", "deleted_at"}).
			AddRow("123", "test@example.com", []byte("hash"), "", "", "", "", time.Now(), time.Now(), nil)

		mock.ExpectQuery("SELECT .* FROM users WHERE id = ?").
			WithArgs("123").
			WillReturnRows(rows)

		res, err := repo.GetUserByID(ctx, "123")

		assert.NoError(t, err)
		assert.Equal(t, "123", res.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM users WHERE id = ?").
			WithArgs("404").
			WillReturnError(sql.ErrNoRows)

		res, err := repo.GetUserByID(ctx, "404")

		assert.NoError(t, err)
		assert.Nil(t, res)
	})
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	sqlxDB, mock := setupTestDB(t)
	repo := NewUserRepository(sqlxDB)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "first_name", "last_name", "locale", "timezone", "created_at", "updated_at", "deleted_at"}).
			AddRow("123", "test@example.com", []byte("hash"), "", "", "", "", time.Now(), time.Now(), nil)

		mock.ExpectQuery("SELECT .* FROM users WHERE email = ?").
			WithArgs("test@example.com").
			WillReturnRows(rows)

		res, err := repo.GetUserByEmail(ctx, "test@example.com")

		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", res.Email)
	})
}

func TestUserRepository_UpdateUser(t *testing.T) {
	sqlxDB, mock := setupTestDB(t)
	repo := NewUserRepository(sqlxDB)
	ctx := context.Background()
	user := &domain.User{ID: "123", Email: "new@example.com"}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET").
			WithArgs(user.Email, "", "", "", "", user.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUser(ctx, user)

		assert.NoError(t, err)
	})
}

func TestUserRepository_DeleteUser(t *testing.T) {
	sqlxDB, mock := setupTestDB(t)
	repo := NewUserRepository(sqlxDB)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET deleted_at = NOW").
			WithArgs("123").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteUser(ctx, "123")

		assert.NoError(t, err)
	})
}

func TestUserRepository_ListUsers(t *testing.T) {
	sqlxDB, mock := setupTestDB(t)
	repo := NewUserRepository(sqlxDB)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "first_name", "last_name", "locale", "timezone", "created_at", "updated_at", "deleted_at"}).
			AddRow("1", "1@ex.com", []byte("h"), "", "", "", "", time.Now(), time.Now(), nil).
			AddRow("2", "2@ex.com", []byte("h"), "", "", "", "", time.Now(), time.Now(), nil)

		mock.ExpectQuery("SELECT .* FROM users").
			WillReturnRows(rows)

		res, err := repo.ListUsers(ctx, &outboundport.ListUsersParams{Active: true})

		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	sqlxDB, mock := setupTestDB(t)
	repo := NewUserRepository(sqlxDB)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET password_hash =").
			WithArgs([]byte("hash"), "123").
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdatePassword(ctx, "123", []byte("hash"))

		assert.NoError(t, err)
	})
}
