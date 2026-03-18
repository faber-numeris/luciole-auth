package postgresadapter

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserConfirmationRepository_CreateUserConfirmation(t *testing.T) {
	sqlxDB, mock := setupTestDB(t)
	repo := NewUserConfirmationRepository(sqlxDB)
	ctx := context.Background()
	userID := "user-123"
	token := "token-123"
	expiresAt := time.Now().Add(time.Hour)

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "token", "expires_at", "confirmed_at", "created_at", "updated_at"}).
			AddRow("conf-123", userID, token, expiresAt, nil, time.Now(), time.Now())

		mock.ExpectQuery("INSERT INTO user_confirmations").
			WithArgs(userID, token, expiresAt).
			WillReturnRows(rows)

		res, err := repo.CreateUserConfirmation(ctx, userID, token, expiresAt)

		assert.NoError(t, err)
		assert.Equal(t, userID, res.UserID)
		assert.Equal(t, token, res.Token)
	})
}

func TestUserConfirmationRepository_GetUserConfirmationByToken(t *testing.T) {
	sqlxDB, mock := setupTestDB(t)
	repo := NewUserConfirmationRepository(sqlxDB)
	ctx := context.Background()
	token := "valid-token"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"user_id"}).AddRow("user-123")

		mock.ExpectQuery("SELECT user_id FROM user_confirmations").
			WithArgs(token).
			WillReturnRows(rows)

		userID, err := repo.GetUserConfirmationByToken(ctx, token)

		assert.NoError(t, err)
		assert.Equal(t, "user-123", userID)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT user_id FROM user_confirmations").
			WithArgs("invalid").
			WillReturnError(sql.ErrNoRows)

		userID, err := repo.GetUserConfirmationByToken(ctx, "invalid")

		assert.NoError(t, err)
		assert.Equal(t, "", userID)
	})
}

func TestUserConfirmationRepository_ConfirmUserRegistration(t *testing.T) {
	sqlxDB, mock := setupTestDB(t)
	repo := NewUserConfirmationRepository(sqlxDB)
	ctx := context.Background()
	userID := "user-123"

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE user_confirmations SET confirmed_at =").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.ConfirmUserRegistration(ctx, userID)

		assert.NoError(t, err)
	})
}

func TestUserConfirmationRepository_DeleteUserConfirmation(t *testing.T) {
	sqlxDB, mock := setupTestDB(t)
	repo := NewUserConfirmationRepository(sqlxDB)
	ctx := context.Background()
	userID := "user-123"

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM user_confirmations").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteUserConfirmation(ctx, userID)

		assert.NoError(t, err)
	})
}
