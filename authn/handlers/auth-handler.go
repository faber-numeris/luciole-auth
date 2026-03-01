package handlers

import (
	"context"
	"log/slog"

	api2 "github.com/faber-numeris/luciole-auth/authn/api/gen"
	"github.com/faber-numeris/luciole-auth/authn/model"
	"github.com/faber-numeris/luciole-auth/authn/model/generated"
	"github.com/faber-numeris/luciole-auth/authn/service"
)

type IAuthnService = api2.Handler

type AuthnService struct {
	userService service.IUserService
}

var converterImpl = generated.ConverterImpl{}

func (a *AuthnService) ConfirmUserRegistration(ctx context.Context, params api2.ConfirmUserRegistrationParams) (api2.ConfirmUserRegistrationRes, error) {
	err := a.userService.ConfirmUserRegistration(ctx, params.Token)
	if err != nil {
		slog.Error("Failed to confirm registration", "error", err)
		return &api2.ConfirmUserRegistrationBadRequest{
			Error:   "INVALID_TOKEN",
			Message: "Invalid or expired confirmation token",
			Details: api2.OptErrorDetails{},
		}, nil
	}

	return &api2.MessageResponse{
		Message: "User registration confirmed successfully",
	}, nil
}

// GetUserByID retrieves a user by their ID
// assignees: rafaelsousa
func (a *AuthnService) GetUserByID(ctx context.Context, params api2.GetUserByIDParams) (api2.GetUserByIDRes, error) {
	user, err := a.userService.GetUserByID(ctx, string(params.ID))
	if err != nil {
		errorResponse := &api2.GetUserByIDInternalServerError{
			Error:   err.Error(),
			Message: "Could not retrieve user. Please try again later.",
			Details: api2.OptErrorDetails{},
		}
		return errorResponse, err
	}

	if user == nil {
		errorResponse := &api2.GetUserByIDNotFound{
			Error:   "USER_NOT_FOUND",
			Message: "User not found",
			Details: api2.OptErrorDetails{},
		}
		return errorResponse, nil
	}

	apiUser, err := converterImpl.UserModelToApiUser(*user)
	if err != nil {
		errorResponse := &api2.GetUserByIDInternalServerError{
			Error:   err.Error(),
			Message: "Could not process user response.",
			Details: api2.OptErrorDetails{},
		}
		return errorResponse, err
	}

	return &apiUser, nil
}

func (a *AuthnService) GetUserProfile(ctx context.Context) (api2.GetUserProfileRes, error) {
	userID := ctx.Value(UserIDKey)
	if userID == nil {
		return &api2.GetUserProfileUnauthorized{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
			Details: api2.OptErrorDetails{},
		}, nil
	}

	user, err := a.userService.GetUserByID(ctx, userID.(string))
	if err != nil {
		slog.Error("Failed to get user profile", "error", err)
		return &api2.GetUserProfileInternalServerError{
			Error:   "INTERNAL_ERROR",
			Message: "Could not retrieve user profile",
			Details: api2.OptErrorDetails{},
		}, nil
	}

	if user == nil {
		return &api2.GetUserProfileUnauthorized{
			Error:   "USER_NOT_FOUND",
			Message: "User not found",
			Details: api2.OptErrorDetails{},
		}, nil
	}

	apiUser, err := converterImpl.UserModelToApiUser(*user)
	if err != nil {
		return &api2.GetUserProfileInternalServerError{
			Error:   "INTERNAL_ERROR",
			Message: "Could not process user response",
			Details: api2.OptErrorDetails{},
		}, nil
	}

	return &apiUser, nil
}

func (a *AuthnService) LoginUser(ctx context.Context, req *api2.LoginRequest) (api2.LoginUserRes, error) {
	user, err := a.userService.VerifyPassword(ctx, req.Email, req.Password)
	if err != nil {
		slog.Error("Login failed", "email", req.Email, "error", err)
		return &api2.LoginUserUnauthorized{
			Error:   "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
			Details: api2.OptErrorDetails{},
		}, nil
	}

	apiUser, err := converterImpl.UserModelToApiUser(*user)
	if err != nil {
		return &api2.LoginUserInternalServerError{
			Error:   "INTERNAL_ERROR",
			Message: "Could not process user response",
			Details: api2.OptErrorDetails{},
		}, nil
	}

	return &api2.LoginResponse{
		AccessToken: user.ID,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		User:        apiUser,
	}, nil
}

func (a *AuthnService) LogoutUser(ctx context.Context) (api2.LogoutUserRes, error) {
	return &api2.LogoutUserNoContent{}, nil
}

func (a *AuthnService) RegisterUser(ctx context.Context, req *api2.UserCreateRequest) (api2.RegisterUserRes, error) {
	userModel, err := converterImpl.UserModelFromUserRequest(*req)
	if err != nil {
		errorResponse := &api2.RegisterUserBadRequest{
			Error:   err.Error(),
			Message: "Invalid Register User Data",
			Details: api2.OptErrorDetails{},
		}

		return errorResponse, err
	}

	userResponse, err := a.userService.RegisterUser(ctx, &userModel, req.Password)
	if err != nil {
		errorResponse := &api2.RegisterUserInternalServerError{
			Error:   err.Error(),
			Message: "Could not register user. Please try again later.",
			Details: api2.OptErrorDetails{},
		}
		return errorResponse, err
	}

	apiUserResponse, err := converterImpl.UserModelToApiUser(*userResponse)
	if err != nil {
		errorResponse := &api2.RegisterUserInternalServerError{
			Error:   err.Error(),
			Message: "Could not process user response.",
			Details: api2.OptErrorDetails{},
		}
		return errorResponse, err
	}

	return &apiUserResponse, nil
}

func (a *AuthnService) RequestPasswordReset(ctx context.Context, req *api2.PasswordResetRequest) (api2.RequestPasswordResetRes, error) {
	_, err := a.userService.RequestPasswordReset(ctx, req.Email)
	if err != nil {
		slog.Error("Password reset request failed", "email", req.Email, "error", err)
		return &api2.RequestPasswordResetNotFound{
			Error:   "USER_NOT_FOUND",
			Message: "If the email exists, a password reset link will be sent",
			Details: api2.OptErrorDetails{},
		}, nil
	}

	return &api2.MessageResponse{
		Message: "If the email exists, a password reset link will be sent",
	}, nil
}

func (a *AuthnService) ResetPassword(ctx context.Context, req *api2.PasswordResetConfirm) (api2.ResetPasswordRes, error) {
	err := a.userService.ResetPassword(ctx, req.Token, req.NewPassword)
	if err != nil {
		slog.Error("Password reset failed", "error", err)
		return &api2.ResetPasswordBadRequest{
			Error:   "INVALID_TOKEN",
			Message: "Invalid or expired password reset token",
			Details: api2.OptErrorDetails{},
		}, nil
	}

	return &api2.MessageResponse{
		Message: "Password reset successfully",
	}, nil
}

func (a *AuthnService) UpdateUserProfile(ctx context.Context, req *api2.UserUpdateRequest) (api2.UpdateUserProfileRes, error) {
	userID := ctx.Value(UserIDKey)
	if userID == nil {
		return &api2.UpdateUserProfileUnauthorized{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
			Details: api2.OptErrorDetails{},
		}, nil
	}

	userModel := &model.User{
		Profile: &model.UserProfile{},
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
		return &api2.UpdateUserProfileInternalServerError{
			Error:   "INTERNAL_ERROR",
			Message: "Could not update user profile",
			Details: api2.OptErrorDetails{},
		}, nil
	}

	apiUser, err := converterImpl.UserModelToApiUser(*updatedUser)
	if err != nil {
		return &api2.UpdateUserProfileInternalServerError{
			Error:   "INTERNAL_ERROR",
			Message: "Could not process user response",
			Details: api2.OptErrorDetails{},
		}, nil
	}

	return &apiUser, nil
}

func NewAuthnService(userService service.IUserService) IAuthnService {
	return &AuthnService{
		userService: userService,
	}
}
