package httpapi

import (
	"context"
	"errors"
	"log/slog"

	api "github.com/faber-numeris/luciole-auth/authn/internal/adapters/inbound/httpapi/gen"
	inboundport "github.com/faber-numeris/luciole-auth/authn/internal/app/ports/inbound"
	"github.com/faber-numeris/luciole-auth/authn/internal/domain"
	"github.com/faber-numeris/luciole-auth/authn/internal/platform/mapper/generated"
)

type Handler struct {
	userService       inboundport.UserService
	hashingService    inboundport.HashingService
	encryptionService inboundport.EncryptionService
}

var converterImpl = generated.ConverterImpl{}

func NewHandler(
	userService inboundport.UserService,
	hashingService inboundport.HashingService,
	encryptionService inboundport.EncryptionService,
) *Handler {
	return &Handler{
		userService:       userService,
		hashingService:    hashingService,
		encryptionService: encryptionService,
	}
}

func (h *Handler) ConfirmUserRegistration(ctx context.Context, params api.ConfirmUserRegistrationParams) (api.ConfirmUserRegistrationRes, error) {
	err := h.userService.ConfirmUserRegistration(ctx, params.Token)
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

func (h *Handler) GetUserByID(ctx context.Context, params api.GetUserByIDParams) (api.GetUserByIDRes, error) {
	user, err := h.userService.GetUserByID(ctx, string(params.ID))
	if err != nil {
		errorResponse := &api.GetUserByIDInternalServerError{
			Error:   err.Error(),
			Message: "Could not retrieve user. Please try again later.",
			Details: api.OptErrorDetails{},
		}
		return errorResponse, nil
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
		return errorResponse, nil
	}

	return &apiUser, nil
}

func (h *Handler) GetUserProfile(ctx context.Context) (api.GetUserProfileRes, error) {
	userID := ctx.Value(UserIDKey)
	if userID == nil {
		return &api.GetUserProfileUnauthorized{
			Error:   "UNAUTHORIZED",
			Message: "User not authenticated",
			Details: api.OptErrorDetails{},
		}, nil
	}

	user, err := h.userService.GetUserByID(ctx, userID.(string))
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

func (h *Handler) LoginUser(ctx context.Context, req *api.LoginRequest) (api.LoginUserRes, error) {
	_, err := h.userService.VerifyPassword(ctx, req.Email, []byte(req.Password))
	if err != nil {
		slog.Error("Login failed", "email", req.Email, "error", err)
		return &api.LoginUserUnauthorized{
			Error:   "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
			Details: api.OptErrorDetails{},
		}, nil
	}

	user, err := h.userService.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return &api.LoginUserInternalServerError{
			Error:   "INTERNAL_ERROR",
			Message: "Could not retrieve user details",
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

func (h *Handler) LogoutUser(ctx context.Context) (api.LogoutUserRes, error) {
	return &api.LogoutUserNoContent{}, nil
}

func (h *Handler) RegisterUser(ctx context.Context, req *api.UserCreateRequest) (api.RegisterUserRes, error) {
	userModel, err := converterImpl.UserModelFromUserRequest(*req)
	if err != nil {
		errorResponse := &api.RegisterUserBadRequest{
			Error:   err.Error(),
			Message: "Invalid Register User Data",
			Details: api.OptErrorDetails{},
		}

		return errorResponse, nil
	}

	userResponse, err := h.userService.RegisterUser(ctx, &userModel, []byte(req.Password))
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return &api.RegisterUserConflict{
				Error:   "USER_ALREADY_EXISTS",
				Message: "User already exists",
				Details: api.OptErrorDetails{},
			}, nil
		}
		errorResponse := &api.RegisterUserInternalServerError{
			Error:   err.Error(),
			Message: "Could not register user. Please try again later.",
			Details: api.OptErrorDetails{},
		}
		return errorResponse, nil
	}

	apiUserResponse, err := converterImpl.UserModelToApiUser(*userResponse)
	if err != nil {
		errorResponse := &api.RegisterUserInternalServerError{
			Error:   err.Error(),
			Message: "Could not process user response.",
			Details: api.OptErrorDetails{},
		}
		return errorResponse, nil
	}

	return &apiUserResponse, nil
}

func (h *Handler) RequestPasswordReset(ctx context.Context, req *api.PasswordResetRequest) (api.RequestPasswordResetRes, error) {
	_, err := h.userService.RequestPasswordReset(ctx, req.Email)
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

func (h *Handler) ResetPassword(ctx context.Context, req *api.PasswordResetConfirm) (api.ResetPasswordRes, error) {
	err := h.userService.ResetPassword(ctx, req.Token, []byte(req.NewPassword))
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

func (h *Handler) UpdateUserProfile(ctx context.Context, req *api.UserUpdateRequest) (api.UpdateUserProfileRes, error) {
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

	updatedUser, err := h.userService.UpdateUserProfile(ctx, userID.(string), userModel)
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
