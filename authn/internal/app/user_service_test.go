package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/faber-numeris/luciole-auth/authn/internal/app"
	inboundport "github.com/faber-numeris/luciole-auth/authn/internal/app/ports/inbound"
	outboundport "github.com/faber-numeris/luciole-auth/authn/internal/app/ports/outbound"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/faber-numeris/luciole-auth/authn/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_RegisterUser(t *testing.T) {
	ctx := context.Background()
	user := &domain.User{
		Email: "test@example.com",
	}
	password := []byte("password123")
	passwordHash := []byte("hashed_password")

	t.Run("success", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		confRepo := mocks.NewMockUserConfirmationRepository(t)
		mailService := mocks.NewMockMailer(t)
		hashingService := mocks.NewMockHashingService(t)
		service := app.NewUserService(userRepo, confRepo, hashingService, mailService)

		hashingService.EXPECT().HashPassword(ctx, password).Return(passwordHash, nil)
		userRepo.EXPECT().CreateUser(ctx, user, passwordHash).Return(&domain.User{ID: "user-123", Email: user.Email}, nil)
		confRepo.EXPECT().CreateUserConfirmation(ctx, "user-123", mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).
			Return(&domain.UserConfirmation{UserID: "user-123", Token: "token-123"}, nil)
		mailService.EXPECT().SendConfirmationEmail(ctx, mock.MatchedBy(func(c domain.UserConfirmation) bool {
			return c.UserID == "user-123" && c.UserEmail == "test@example.com"
		})).Return(nil)

		createdUser, err := service.RegisterUser(ctx, user, password)

		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, "user-123", createdUser.ID)
	})

	t.Run("hashing failure", func(t *testing.T) {
		hashingService := mocks.NewMockHashingService(t)
		service := app.NewUserService(nil, nil, hashingService, nil)

		hashingService.EXPECT().HashPassword(ctx, password).Return(nil, errors.New("hash error"))

		createdUser, err := service.RegisterUser(ctx, user, password)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to hash password")
		assert.Nil(t, createdUser)
	})

	t.Run("create user failure", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		hashingService := mocks.NewMockHashingService(t)
		service := app.NewUserService(userRepo, nil, hashingService, nil)

		hashingService.EXPECT().HashPassword(ctx, password).Return(passwordHash, nil)
		userRepo.EXPECT().CreateUser(ctx, user, passwordHash).Return(nil, errors.New("db error"))

		createdUser, err := service.RegisterUser(ctx, user, password)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create user in repository")
		assert.Nil(t, createdUser)
	})

	t.Run("confirmation creation failure", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		confRepo := mocks.NewMockUserConfirmationRepository(t)
		hashingService := mocks.NewMockHashingService(t)
		service := app.NewUserService(userRepo, confRepo, hashingService, nil)

		hashingService.EXPECT().HashPassword(ctx, password).Return(passwordHash, nil)
		userRepo.EXPECT().CreateUser(ctx, user, passwordHash).Return(&domain.User{ID: "user-123", Email: user.Email}, nil)
		confRepo.EXPECT().CreateUserConfirmation(ctx, "user-123", mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).
			Return(nil, errors.New("conf error"))

		createdUser, err := service.RegisterUser(ctx, user, password)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create confirmation token")
		assert.Nil(t, createdUser)
	})

	t.Run("mail failure", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		confRepo := mocks.NewMockUserConfirmationRepository(t)
		mailService := mocks.NewMockMailer(t)
		hashingService := mocks.NewMockHashingService(t)
		service := app.NewUserService(userRepo, confRepo, hashingService, mailService)

		hashingService.EXPECT().HashPassword(ctx, password).Return(passwordHash, nil)
		userRepo.EXPECT().CreateUser(ctx, user, passwordHash).Return(&domain.User{ID: "user-123", Email: user.Email}, nil)
		confRepo.EXPECT().CreateUserConfirmation(ctx, "user-123", mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).
			Return(&domain.UserConfirmation{UserID: "user-123", Token: "token-123"}, nil)
		mailService.EXPECT().SendConfirmationEmail(ctx, mock.Anything).Return(errors.New("mail error"))

		createdUser, err := service.RegisterUser(ctx, user, password)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to send confirmation email")
		assert.Nil(t, createdUser)
	})
}

func TestUserService_GetUserByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		expectedUser := &domain.User{ID: "user-123"}
		userRepo.EXPECT().GetUserByID(ctx, "user-123").Return(expectedUser, nil)

		user, err := service.GetUserByID(ctx, "user-123")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("not found", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserByID(ctx, "non-existent").Return(nil, nil)

		user, err := service.GetUserByID(ctx, "non-existent")

		assert.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("error", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserByID(ctx, "error-id").Return(nil, errors.New("db error"))

		user, err := service.GetUserByID(ctx, "error-id")

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserService_GetUserByEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		expectedUser := &domain.User{Email: "test@example.com"}
		userRepo.EXPECT().GetUserByEmail(ctx, "test@example.com").Return(expectedUser, nil)

		user, err := service.GetUserByEmail(ctx, "test@example.com")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("not found", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserByEmail(ctx, "unknown").Return(nil, nil)

		user, err := service.GetUserByEmail(ctx, "unknown")

		assert.NoError(t, err)
		assert.Nil(t, user)
	})

	t.Run("error", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserByEmail(ctx, "error").Return(nil, errors.New("db error"))

		user, err := service.GetUserByEmail(ctx, "error")

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserService_ConfirmUserRegistration(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		confRepo := mocks.NewMockUserConfirmationRepository(t)
		service := app.NewUserService(nil, confRepo, nil, nil)
		confRepo.EXPECT().GetUserConfirmationByToken(ctx, "valid-token").Return("user-123", nil)
		confRepo.EXPECT().ConfirmUserRegistration(ctx, "user-123").Return(nil)

		err := service.ConfirmUserRegistration(ctx, "valid-token")

		assert.NoError(t, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		confRepo := mocks.NewMockUserConfirmationRepository(t)
		service := app.NewUserService(nil, confRepo, nil, nil)
		confRepo.EXPECT().GetUserConfirmationByToken(ctx, "invalid-token").Return("", nil)

		err := service.ConfirmUserRegistration(ctx, "invalid-token")

		assert.Error(t, err)
		assert.Equal(t, "invalid or expired confirmation token", err.Error())
	})

	t.Run("token retrieval error", func(t *testing.T) {
		confRepo := mocks.NewMockUserConfirmationRepository(t)
		service := app.NewUserService(nil, confRepo, nil, nil)
		confRepo.EXPECT().GetUserConfirmationByToken(ctx, "error-token").Return("", errors.New("db error"))

		err := service.ConfirmUserRegistration(ctx, "error-token")

		assert.Error(t, err)
	})

	t.Run("confirmation error", func(t *testing.T) {
		confRepo := mocks.NewMockUserConfirmationRepository(t)
		service := app.NewUserService(nil, confRepo, nil, nil)
		confRepo.EXPECT().GetUserConfirmationByToken(ctx, "valid-token").Return("user-123", nil)
		confRepo.EXPECT().ConfirmUserRegistration(ctx, "user-123").Return(errors.New("db error"))

		err := service.ConfirmUserRegistration(ctx, "valid-token")

		assert.Error(t, err)
	})
}

func TestUserService_UpdateUserProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("success with all fields", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		existingUser := &domain.User{
			ID: "user-123",
			Profile: &domain.UserProfile{
				FirstName: "OldFirstName",
				LastName:  "OldLastName",
				Locale:    "en",
				Timezone:  "UTC",
			},
		}
		updateReq := &domain.User{
			Profile: &domain.UserProfile{
				FirstName: "NewFirstName",
				LastName:  "NewLastName",
				Locale:    "fr",
				Timezone:  "Europe/Paris",
			},
		}

		userRepo.EXPECT().GetUserByID(ctx, "user-123").Return(existingUser, nil)
		userRepo.EXPECT().UpdateUser(ctx, mock.MatchedBy(func(u *domain.User) bool {
			return u.ID == "user-123" &&
				u.Profile.FirstName == "NewFirstName" &&
				u.Profile.LastName == "NewLastName" &&
				u.Profile.Locale == "fr" &&
				u.Profile.Timezone == "Europe/Paris"
		})).Return(nil)

		updatedUser, err := service.UpdateUserProfile(ctx, "user-123", updateReq)

		assert.NoError(t, err)
		assert.Equal(t, "NewFirstName", updatedUser.Profile.FirstName)
	})

	t.Run("success with missing profile in existing", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		existingUser := &domain.User{ID: "user-123", Profile: nil}
		updateReq := &domain.User{Profile: &domain.UserProfile{FirstName: "New"}}

		userRepo.EXPECT().GetUserByID(ctx, "user-123").Return(existingUser, nil)
		userRepo.EXPECT().UpdateUser(ctx, mock.Anything).Return(nil)

		updatedUser, err := service.UpdateUserProfile(ctx, "user-123", updateReq)

		assert.NoError(t, err)
		assert.Equal(t, "New", updatedUser.Profile.FirstName)
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserByID(ctx, "unknown").Return(nil, nil)

		updatedUser, err := service.UpdateUserProfile(ctx, "unknown", &domain.User{})

		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	})

	t.Run("get user error", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserByID(ctx, "error").Return(nil, errors.New("db error"))

		updatedUser, err := service.UpdateUserProfile(ctx, "error", &domain.User{})

		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	})

	t.Run("update user error", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserByID(ctx, "user-123").Return(&domain.User{ID: "user-123"}, nil)
		userRepo.EXPECT().UpdateUser(ctx, mock.Anything).Return(errors.New("db error"))

		updatedUser, err := service.UpdateUserProfile(ctx, "user-123", &domain.User{})

		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserByID(ctx, "user-123").Return(&domain.User{ID: "user-123"}, nil)
		userRepo.EXPECT().DeleteUser(ctx, "user-123").Return(nil)

		err := service.DeleteUser(ctx, "user-123")

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserByID(ctx, "unknown").Return(nil, nil)

		err := service.DeleteUser(ctx, "unknown")

		assert.Error(t, err)
	})

	t.Run("get user error", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserByID(ctx, "error").Return(nil, errors.New("db error"))

		err := service.DeleteUser(ctx, "error")

		assert.Error(t, err)
	})

	t.Run("delete error", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserByID(ctx, "user-123").Return(&domain.User{ID: "user-123"}, nil)
		userRepo.EXPECT().DeleteUser(ctx, "user-123").Return(errors.New("db error"))

		err := service.DeleteUser(ctx, "user-123")

		assert.Error(t, err)
	})
}

func TestUserService_ListUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		params := &inboundport.ListUsersParams{Active: true}
		expectedUsers := []*domain.User{{ID: "1"}, {ID: "2"}}
		userRepo.EXPECT().ListUsers(ctx, mock.MatchedBy(func(p *outboundport.ListUsersParams) bool {
			return p.Active == true
		})).Return(expectedUsers, nil)

		users, err := service.ListUsers(ctx, params)

		assert.NoError(t, err)
		assert.Len(t, users, 2)
	})

	t.Run("error", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().ListUsers(ctx, mock.Anything).Return(nil, errors.New("db error"))

		users, err := service.ListUsers(ctx, &inboundport.ListUsersParams{})

		assert.Error(t, err)
		assert.Nil(t, users)
	})
}

func TestUserService_VerifyPassword(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		hashingService := mocks.NewMockHashingService(t)
		service := app.NewUserService(userRepo, nil, hashingService, nil)
		password := []byte("password")
		passwordHash := []byte("hashed_password")

		expectedCreds := &domain.UserCredentials{Email: "test@example.com", PasswordHash: passwordHash}
		userRepo.EXPECT().GetUserCredentials(ctx, "test@example.com").Return(expectedCreds, nil)
		hashingService.EXPECT().HashPassword(ctx, password).Return(passwordHash, nil)

		user, err := service.VerifyPassword(ctx, "test@example.com", password)

		assert.NoError(t, err)
		assert.Equal(t, expectedCreds, user)
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserCredentials(ctx, "unknown@example.com").Return(nil, nil)

		user, err := service.VerifyPassword(ctx, "unknown@example.com", []byte("password"))

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("db error", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository(t)
		service := app.NewUserService(userRepo, nil, nil, nil)
		userRepo.EXPECT().GetUserCredentials(ctx, "error").Return(nil, errors.New("db error"))

		user, err := service.VerifyPassword(ctx, "error", []byte("password"))

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserService_RequestPasswordReset(t *testing.T) {
	service := app.NewUserService(nil, nil, nil, nil)
	_, err := service.RequestPasswordReset(context.Background(), "email")
	assert.Error(t, err)
	assert.Equal(t, "not implemented", err.Error())
}

func TestUserService_ResetPassword(t *testing.T) {
	service := app.NewUserService(nil, nil, nil, nil)
	err := service.ResetPassword(context.Background(), "token", []byte("new"))
	assert.Error(t, err)
	assert.Equal(t, "not implemented", err.Error())
}
