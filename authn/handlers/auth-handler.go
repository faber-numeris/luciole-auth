package handlers

import (
	"context"

	api2 "github.com/faber-numeris/luciole-auth/authn/api/gen"
	"github.com/faber-numeris/luciole-auth/authn/model/generated"
	"github.com/faber-numeris/luciole-auth/authn/service"
)

type IAuthnService = api2.Handler

type AuthnService struct {
	userService service.IUserService
}

var converterImpl = generated.ConverterImpl{}

// TODO: Implement ConfirmUserRegistration method
// assignees: rafaelsousa
func (a *AuthnService) ConfirmUserRegistration(ctx context.Context, params api2.ConfirmUserRegistrationParams) (api2.ConfirmUserRegistrationRes, error) {
	//TODO implement me
	panic("implement me")
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

// TODO: Implement GetUserProfile method
// assignees: rafaelsousa
func (a *AuthnService) GetUserProfile(ctx context.Context) (api2.GetUserProfileRes, error) {
	//TODO implement me
	panic("implement me")
}

// TODO: Implement LoginUser method
// assignees: rafaelsousa
func (a *AuthnService) LoginUser(ctx context.Context, req *api2.LoginRequest) (api2.LoginUserRes, error) {
	//TODO implement me
	panic("implement me")
}

// TODO: Implement LogoutUser method
// assignees: rafaelsousa
func (a *AuthnService) LogoutUser(ctx context.Context) (api2.LogoutUserRes, error) {
	//TODO implement me
	panic("implement me")
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

// TODO: Implement RequestPasswordReset method
// assignees: rafaelsousa
func (a *AuthnService) RequestPasswordReset(ctx context.Context, req *api2.PasswordResetRequest) (api2.RequestPasswordResetRes, error) {
	//TODO implement me
	panic("implement me")
}

// TODO: Implement ResetPassword method
// assignees: rafaelsousa
func (a *AuthnService) ResetPassword(ctx context.Context, req *api2.PasswordResetConfirm) (api2.ResetPasswordRes, error) {
	//TODO implement me
	panic("implement me")
}

// TODO: Implement UpdateUserProfile method
// assignees: rafaelsousa
func (a *AuthnService) UpdateUserProfile(ctx context.Context, req *api2.UserUpdateRequest) (api2.UpdateUserProfileRes, error) {
	//TODO implement me
	panic("implement me")
}

func NewAuthnService(userService service.IUserService) IAuthnService {
	return &AuthnService{
		userService: userService,
	}
}
