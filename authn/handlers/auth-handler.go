package handlers

import (
	"context"

	api "github.com/faber-numeris/luciole-auth/api/gen"
	"github.com/faber-numeris/luciole-auth/authn/model/generated"
	"github.com/faber-numeris/luciole-auth/authn/service"
)

type IAuthnService = api.Handler

type AuthnService struct {
	userService service.IUserService
}

// TODO: Implement ConfirmUserRegistration method
// assignees: rafaelsousa
func (a *AuthnService) ConfirmUserRegistration(ctx context.Context, params api.ConfirmUserRegistrationParams) (api.ConfirmUserRegistrationRes, error) {
	//TODO implement me
	panic("implement me")
}

// TODO: Implement GetUserProfile method
// assignees: rafaelsousa
func (a *AuthnService) GetUserProfile(ctx context.Context) (api.GetUserProfileRes, error) {
	//TODO implement me
	panic("implement me")
}

// TODO: Implement LoginUser method
// assignees: rafaelsousa
func (a *AuthnService) LoginUser(ctx context.Context, req *api.LoginRequest) (api.LoginUserRes, error) {
	//TODO implement me
	panic("implement me")
}

// TODO: Implement LogoutUser method
// assignees: rafaelsousa
func (a *AuthnService) LogoutUser(ctx context.Context) (api.LogoutUserRes, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthnService) RegisterUser(ctx context.Context, req *api.RegisterRequest) (api.RegisterUserRes, error) {
	var converter generated.ConverterImpl

	userModel, err := converter.RegisterRequestToUserModel(*req)
	if err != nil {
		errorResponse := &api.RegisterUserBadRequest{
			Error:   err.Error(),
			Message: "Invalid Register User Data",
			Details: api.OptErrorDetails{},
		}

		return errorResponse, err
	}

	userResponse, err := a.userService.RegisterUser(ctx, userModel)
	if err != nil {
		errorResponse := &api.RegisterUserInternalServerError{
			Error:   err.Error(),
			Message: "Could not register user. Please try again later.",
			Details: api.OptErrorDetails{},
		}
		return errorResponse, err
	}

	apiUserResponse, err := converter.UserToRegisterUserRes(*userResponse)
	if err != nil {
		errorResponse := &api.RegisterUserInternalServerError{
			Error:   err.Error(),
			Message: "Could not process user response.",
			Details: api.OptErrorDetails{},
		}
		return errorResponse, err
	}

	return apiUserResponse, nil

}

// TODO: Implement RequestPasswordReset method
// assignees: rafaelsousa
func (a *AuthnService) RequestPasswordReset(ctx context.Context, req *api.PasswordResetRequest) (api.RequestPasswordResetRes, error) {
	//TODO implement me
	panic("implement me")
}

// TODO: Implement ResetPassword method
// assignees: rafaelsousa
func (a *AuthnService) ResetPassword(ctx context.Context, req *api.PasswordResetConfirm) (api.ResetPasswordRes, error) {
	//TODO implement me
	panic("implement me")
}

// TODO: Implement UpdateUserProfile method
// assignees: rafaelsousa
func (a *AuthnService) UpdateUserProfile(ctx context.Context, req *api.ProfileUpdateRequest) (api.UpdateUserProfileRes, error) {
	//TODO implement me
	panic("implement me")
}

func NewAuthnService(userService service.IUserService) IAuthnService {
	return &AuthnService{
		userService: userService,
	}
}
