package httpapi

import (
	"context"
	"errors"
	"testing"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/inbound/httpapi/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/faber-numeris/luciole-auth/authn/internal/mocks"
	"github.com/faber-numeris/luciole-auth/authn/internal/platform/mapper/generated"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_ConfirmUserRegistration(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().ConfirmUserRegistration(ctx, "token-123").Return(nil)

		res, err := handler.ConfirmUserRegistration(ctx, api.ConfirmUserRegistrationParams{Token: "token-123"})

		assert.NoError(t, err)
		assert.IsType(t, &api.MessageResponse{}, res)
	})

	t.Run("error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().ConfirmUserRegistration(ctx, "bad-token").Return(errors.New("invalid"))

		res, err := handler.ConfirmUserRegistration(ctx, api.ConfirmUserRegistrationParams{Token: "bad-token"})

		assert.NoError(t, err)
		assert.IsType(t, &api.ConfirmUserRegistrationBadRequest{}, res)
	})
}

func TestHandler_RegisterUser(t *testing.T) {
	ctx := context.Background()
	req := &api.UserCreateRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	t.Run("success", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().RegisterUser(ctx, mock.Anything, "password").
			Return(&domain.User{ID: "123", Email: "test@example.com"}, nil)

		res, err := handler.RegisterUser(ctx, req)

		assert.NoError(t, err)
		assert.IsType(t, &api.User{}, res)
	})

	t.Run("error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().RegisterUser(ctx, mock.Anything, "password").
			Return(nil, errors.New("db error"))

		res, err := handler.RegisterUser(ctx, req)

		assert.NoError(t, err)
		assert.IsType(t, &api.RegisterUserInternalServerError{}, res)
	})

	t.Run("converter error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().RegisterUser(ctx, mock.Anything, "password").
			Return(&domain.User{ID: "123"}, nil)

		patches := gomonkey.ApplyMethod(&converterImpl, "UserModelToApiUser", func(_ *generated.ConverterImpl, _ domain.User) (api.User, error) {
			return api.User{}, errors.New("conv error")
		})
		defer patches.Reset()

		res, err := handler.RegisterUser(ctx, req)
		assert.NoError(t, err)
		assert.IsType(t, &api.RegisterUserInternalServerError{}, res)
	})
}

func TestHandler_LoginUser(t *testing.T) {
	ctx := context.Background()
	req := &api.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	t.Run("success", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().VerifyPassword(ctx, "test@example.com", []byte("password")).
			Return(&domain.UserCredentials{Email: "test@example.com"}, nil)
		userService.EXPECT().GetUserByEmail(ctx, "test@example.com").
			Return(&domain.User{ID: "123", Email: "test@example.com"}, nil)

		res, err := handler.LoginUser(ctx, req)

		assert.NoError(t, err)
		assert.IsType(t, &api.LoginResponse{}, res)
	})

	t.Run("unauthorized", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().VerifyPassword(ctx, "test@example.com", []byte("password")).
			Return(nil, errors.New("invalid"))

		res, err := handler.LoginUser(ctx, req)

		assert.NoError(t, err)
		assert.IsType(t, &api.LoginUserUnauthorized{}, res)
	})

	t.Run("converter error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().VerifyPassword(ctx, "test@example.com", []byte("password")).
			Return(&domain.UserCredentials{Email: "test@example.com"}, nil)
		userService.EXPECT().GetUserByEmail(ctx, "test@example.com").
			Return(&domain.User{ID: "123"}, nil)

		patches := gomonkey.ApplyMethod(&converterImpl, "UserModelToApiUser", func(_ *generated.ConverterImpl, _ domain.User) (api.User, error) {
			return api.User{}, errors.New("conv error")
		})
		defer patches.Reset()

		res, err := handler.LoginUser(ctx, req)
		assert.NoError(t, err)
		assert.IsType(t, &api.LoginUserInternalServerError{}, res)
	})
}

func TestHandler_GetUserByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().GetUserByID(ctx, "123").
			Return(&domain.User{ID: "123", Email: "test@example.com"}, nil)

		res, err := handler.GetUserByID(ctx, api.GetUserByIDParams{ID: "123"})

		assert.NoError(t, err)
		assert.IsType(t, &api.User{}, res)
	})

	t.Run("not found", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().GetUserByID(ctx, "404").Return(nil, nil)

		res, err := handler.GetUserByID(ctx, api.GetUserByIDParams{ID: "404"})

		assert.NoError(t, err)
		assert.IsType(t, &api.GetUserByIDNotFound{}, res)
	})

	t.Run("error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().GetUserByID(ctx, "500").Return(nil, errors.New("db error"))

		res, err := handler.GetUserByID(ctx, api.GetUserByIDParams{ID: "500"})

		assert.NoError(t, err)
		assert.IsType(t, &api.GetUserByIDInternalServerError{}, res)
	})

	t.Run("converter error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().GetUserByID(ctx, "123").
			Return(&domain.User{ID: "123"}, nil)

		patches := gomonkey.ApplyMethod(&converterImpl, "UserModelToApiUser", func(_ *generated.ConverterImpl, _ domain.User) (api.User, error) {
			return api.User{}, errors.New("conv error")
		})
		defer patches.Reset()

		res, err := handler.GetUserByID(ctx, api.GetUserByIDParams{ID: "123"})
		assert.NoError(t, err)
		assert.IsType(t, &api.GetUserByIDInternalServerError{}, res)
	})
}

func TestHandler_UpdateUserProfile(t *testing.T) {
	req := &api.UserUpdateRequest{
		FirstName: api.NewOptString("New"),
	}

	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "123")
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().UpdateUserProfile(ctx, "123", mock.Anything).
			Return(&domain.User{ID: "123", Profile: &domain.UserProfile{FirstName: "New"}}, nil)

		res, err := handler.UpdateUserProfile(ctx, req)

		assert.NoError(t, err)
		assert.IsType(t, &api.User{}, res)
	})

	t.Run("unauthorized", func(t *testing.T) {
		handler := NewHandler(nil, nil, nil)
		res, err := handler.UpdateUserProfile(context.Background(), req)
		assert.NoError(t, err)
		assert.IsType(t, &api.UpdateUserProfileUnauthorized{}, res)
	})

	t.Run("service error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "123")
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().UpdateUserProfile(ctx, "123", mock.Anything).
			Return(nil, errors.New("error"))

		res, err := handler.UpdateUserProfile(ctx, req)
		assert.NoError(t, err)
		assert.IsType(t, &api.UpdateUserProfileInternalServerError{}, res)
	})

	t.Run("converter error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "123")
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().UpdateUserProfile(ctx, "123", mock.Anything).
			Return(&domain.User{ID: "123"}, nil)

		patches := gomonkey.ApplyMethod(&converterImpl, "UserModelToApiUser", func(_ *generated.ConverterImpl, _ domain.User) (api.User, error) {
			return api.User{}, errors.New("conv error")
		})
		defer patches.Reset()

		res, err := handler.UpdateUserProfile(ctx, req)
		assert.NoError(t, err)
		assert.IsType(t, &api.UpdateUserProfileInternalServerError{}, res)
	})
}

func TestHandler_GetUserProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "123")
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().GetUserByID(ctx, "123").
			Return(&domain.User{ID: "123", Email: "test@example.com"}, nil)

		res, err := handler.GetUserProfile(ctx)

		assert.NoError(t, err)
		assert.IsType(t, &api.User{}, res)
	})

	t.Run("unauthorized", func(t *testing.T) {
		handler := NewHandler(nil, nil, nil)
		res, err := handler.GetUserProfile(context.Background())
		assert.NoError(t, err)
		assert.IsType(t, &api.GetUserProfileUnauthorized{}, res)
	})

	t.Run("user not found", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "123")
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().GetUserByID(ctx, "123").Return(nil, nil)

		res, err := handler.GetUserProfile(ctx)
		assert.NoError(t, err)
		assert.IsType(t, &api.GetUserProfileUnauthorized{}, res)
	})

	t.Run("service error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "123")
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().GetUserByID(ctx, "123").Return(nil, errors.New("error"))

		res, err := handler.GetUserProfile(ctx)
		assert.NoError(t, err)
		assert.IsType(t, &api.GetUserProfileInternalServerError{}, res)
	})

	t.Run("converter error", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), UserIDKey, "123")
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().GetUserByID(ctx, "123").
			Return(&domain.User{ID: "123"}, nil)

		patches := gomonkey.ApplyMethod(&converterImpl, "UserModelToApiUser", func(_ *generated.ConverterImpl, _ domain.User) (api.User, error) {
			return api.User{}, errors.New("conv error")
		})
		defer patches.Reset()

		res, err := handler.GetUserProfile(ctx)
		assert.NoError(t, err)
		assert.IsType(t, &api.GetUserProfileInternalServerError{}, res)
	})
}

func TestHandler_RequestPasswordReset(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().RequestPasswordReset(ctx, "test@example.com").Return("token", nil)

		res, err := handler.RequestPasswordReset(ctx, &api.PasswordResetRequest{Email: "test@example.com"})

		assert.NoError(t, err)
		assert.IsType(t, &api.MessageResponse{}, res)
	})

	t.Run("error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().RequestPasswordReset(ctx, "test@example.com").Return("", errors.New("not found"))

		res, err := handler.RequestPasswordReset(ctx, &api.PasswordResetRequest{Email: "test@example.com"})

		assert.NoError(t, err)
		assert.IsType(t, &api.RequestPasswordResetNotFound{}, res)
	})
}

func TestHandler_ResetPassword(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().ResetPassword(ctx, "token", []byte("new-pass")).Return(nil)

		res, err := handler.ResetPassword(ctx, &api.PasswordResetConfirm{Token: "token", NewPassword: "new-pass"})

		assert.NoError(t, err)
		assert.IsType(t, &api.MessageResponse{}, res)
	})

	t.Run("error", func(t *testing.T) {
		userService := mocks.NewMockUserService(t)
		handler := NewHandler(userService, nil, nil)
		userService.EXPECT().ResetPassword(ctx, "token", []byte("new-pass")).Return(errors.New("invalid"))

		res, err := handler.ResetPassword(ctx, &api.PasswordResetConfirm{Token: "token", NewPassword: "new-pass"})

		assert.NoError(t, err)
		assert.IsType(t, &api.ResetPasswordBadRequest{}, res)
	})
}

func TestHandler_LogoutUser(t *testing.T) {
	handler := NewHandler(nil, nil, nil)
	res, err := handler.LogoutUser(context.Background())
	assert.NoError(t, err)
	assert.IsType(t, &api.LogoutUserNoContent{}, res)
}
