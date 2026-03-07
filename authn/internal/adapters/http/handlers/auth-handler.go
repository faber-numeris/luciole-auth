package handlers

import (
	"context"
	"log/slog"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/http/api/gen"
	"github.com/faber-numeris/luciole-auth/authn/internal/application"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/faber-numeris/luciole-auth/authn/internal/platform/mapper/generated"
)

type IAuthnService = api.Handler

type AuthnService struct {
	userService application.IUserService
}

var converterImpl = generated.ConverterImpl{}

func (a *AuthnService) ConfirmUserRegistration(ctx context.Context, params api.ConfirmUserRegistrationParams) (api.ConfirmUserRegistrationRes, error) {
	err := a.userService.ConfirmUserRegistration(ctx, params.Token)
	if err != nil {
		slog.Error("Failed to confirm registration", "error", err)
		return &api.ConfirmUserRegistrationBadRequest{
			Error:   "INVALID_TOKEN",
			Message: "Invalid or expired confirmation token",
			Details: api.OptErrorDetails{},
		}, nil
	}

	return &api.MessageResponse{
		Message: "User registration confirmed successfully",
	}, nil
}

// GetUserByID retrieves a user by their ID
// assignees: rafaelsousa
func (a *AuthnService) GetUserByID(ctx context.Context, params api.GetUserByIDParams) (api.GetUserByIDRes, error) {
	user, err := a.userService.GetUserByID(ctx, string(params.ID))
	if err != nil {
		errorResponse := &api.GetUserByIDInternalServerError{
			Error:   err.Error(),
			Message: "Could not retrieve user. Please try again later.",
			Details: api.OptErrorDetails{},
		}
		return errorResponse, err
	}

	if user == nil {
		errorResponse := &api.GetUserByIDNotFound{
			Error:   "USER_NOT_FOUND",
			Message: "User not found",
			Details: api.OptErrorDetails{},
		}
		return errorResponse, nil
	}

	apiUser, err := converterImpl.UserModelToApiUser(*user)
	if err != nil {
		errorResponse := &api.GetUserByIDInternalServerError{
			Error:   err.Error(),
			Message: "Could not process user response.",
			Details: api.OptErrorDetails{},
		}
		return errorResponse, err
	}

	return &apiUser, nil
}

func (a *AuthnService) GetUserProfile(ctx context.Context) (api.GetUserProfileRes, error) {
	userID := ctx.Value(UserIDKey)
	if userID == nil {
		return &api.GetUserProfileUnauthorized{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
			Details: api.OptErrorDetails{},
		}, nil
	}

	user, err := a.userService.GetUserByID(ctx, userID.(string))
	if err != nil {
		slog.Error("Failed to get user profile", "error", err)
		return &api.GetUserProfileInternalServerError{
			Error:   "INTERNAL_ERROR",
			Message: "Could not retrieve user profile",
			Details: api.OptErrorDetails{},
		}, nil
	}

	if user == nil {
		return &api.GetUserProfileUnauthorized{
			Error:   "USER_NOT_FOUND",
			Message: "User not found",
			Details: api.OptErrorDetails{},
		}, nil
	}

	apiUser, err := converterImpl.UserModelToApiUser(*user)
	if err != nil {
		return &api.GetUserProfileInternalServerError{
			Error:   "INTERNAL_ERROR",
			Message: "Could not process user response",
			Details: api.OptErrorDetails{},
		}, nil
	}

	return &apiUser, nil
}

func (a *AuthnService) LoginUser(ctx context.Context, req *api.LoginRequest) (api.LoginUserRes, error) {
	user, err := a.userService.VerifyPassword(ctx, req.Email, req.Password)
	if err != nil {
		slog.Error("Login failed", "email", req.Email, "error", err)
		return &api.LoginUserUnauthorized{
			Error:   "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
			Details: api.OptErrorDetails{},
		}, nil
	}

	apiUser, err := converterImpl.UserModelToApiUser(*user)
	if err != nil {
		return &api.LoginUserInternalServerError{
			Error:   "INTERNAL_ERROR",
			Message: "Could not process user response",
			Details: api.OptErrorDetails{},
		}, nil
	}

	return &api.LoginResponse{
		AccessToken: user.ID,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		User:        apiUser,
	}, nil
}

func (a *AuthnService) LogoutUser(ctx context.Context) (api.LogoutUserRes, error) {
	return &api.LogoutUserNoContent{}, nil
}

func (a *AuthnService) RegisterUser(ctx context.Context, req *api.UserCreateRequest) (api.RegisterUserRes, error) {
	userModel, err := converterImpl.UserModelFromUserRequest(*req)
	if err != nil {
		errorResponse := &api.RegisterUserBadRequest{
			Error:   err.Error(),
			Message: "Invalid Register User Data",
			Details: api.OptErrorDetails{},
		}

		return errorResponse, err
	}

	userResponse, err := a.userService.RegisterUser(ctx, &userModel, req.Password)
	if err != nil {
		errorResponse := &api.RegisterUserInternalServerError{
			Error:   err.Error(),
			Message: "Could not register user. Please try again later.",
			Details: api.OptErrorDetails{},
		}
		return errorResponse, err
	}

	apiUserResponse, err := converterImpl.UserModelToApiUser(*userResponse)
	if err != nil {
		errorResponse := &api.RegisterUserInternalServerError{
			Error:   err.Error(),
			Message: "Could not process user response.",
			Details: api.OptErrorDetails{},
		}
		return errorResponse, err
	}

	return &apiUserResponse, nil
}

func (a *AuthnService) RequestPasswordReset(ctx context.Context, req *api.PasswordResetRequest) (api.RequestPasswordResetRes, error) {
	_, err := a.userService.RequestPasswordReset(ctx, req.Email)
	if err != nil {
		slog.Error("Password reset request failed", "email", req.Email, "error", err)
		return &api.RequestPasswordResetNotFound{
			Error:   "USER_NOT_FOUND",
			Message: "If the email exists, a password reset link will be sent",
			Details: api.OptErrorDetails{},
		}, nil
	}

	return &api.MessageResponse{
		Message: "If the email exists, a password reset link will be sent",
	}, nil
}

func (a *AuthnService) ResetPassword(ctx context.Context, req *api.PasswordResetConfirm) (api.ResetPasswordRes, error) {
	err := a.userService.ResetPassword(ctx, req.Token, req.NewPassword)
	if err != nil {
		slog.Error("Password reset failed", "error", err)
		return &api.ResetPasswordBadRequest{
			Error:   "INVALID_TOKEN",
			Message: "Invalid or expired password reset token",
			Details: api.OptErrorDetails{},
		}, nil
	}

	return &api.MessageResponse{
		Message: "Password reset successfully",
	}, nil
}

func (a *AuthnService) UpdateUserProfile(ctx context.Context, req *api.UserUpdateRequest) (api.UpdateUserProfileRes, error) {
	userID := ctx.Value(UserIDKey)
	if userID == nil {
		return &api.UpdateUserProfileUnauthorized{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
			Details: api.OptErrorDetails{},
		}, nil
	}

	userModel := &domain.User{
		Profile: &domain.UserProfile{},
	}

	if firstName, ok := req.FirstName.Get(); ok && firstName != "" {
		userModel.Profile.FirstName = firstName
	}
	if lastName, ok := req.LastName.Get(); ok && lastName != "" {
		userModel.Profile.LastName = lastName
	}
	if locale, ok := req.Locale.Get(); ok && locale != "" {
		userModel.Profile.Locale = locale
	}
	if timezone, ok := req.Timezone.Get(); ok && timezone != "" {
		userModel.Profile.Timezone = timezone
	}

	updatedUser, err := a.userService.UpdateUserProfile(ctx, userID.(string), userModel)
	if err != nil {
		slog.Error("Failed to update user profile", "error", err)
		return &api.UpdateUserProfileInternalServerError{
			Error:   "INTERNAL_ERROR",
			Message: "Could not update user profile",
			Details: api.OptErrorDetails{},
		}, nil
	}

	apiUser, err := converterImpl.UserModelToApiUser(*updatedUser)
	if err != nil {
		return &api.UpdateUserProfileInternalServerError{
			Error:   "INTERNAL_ERROR",
			Message: "Could not process user response",
			Details: api.OptErrorDetails{},
		}, nil
	}

	return &apiUser, nil
}

func NewAuthnService(userService application.IUserService) IAuthnService {
	return &AuthnService{
		userService: userService,
	}
}
