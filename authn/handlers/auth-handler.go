package handlers

import (
	"context"

	api "github.com/faber-numeris/luciole-auth/api/gen"
	"github.com/faber-numeris/luciole-auth/authn/model/generated"
	"github.com/faber-numeris/luciole-auth/authn/persistence/repository"
)

type IAuthnService = api.Handler

type AuthnService struct {
	userRepository repository.UserRepository
}

func (a *AuthnService) ConfirmUserRegistration(ctx context.Context, params api.ConfirmUserRegistrationParams) (api.ConfirmUserRegistrationRes, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthnService) GetUserProfile(ctx context.Context) (api.GetUserProfileRes, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthnService) LoginUser(ctx context.Context, req *api.LoginRequest) (api.LoginUserRes, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthnService) LogoutUser(ctx context.Context) (api.LogoutUserRes, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthnService) RegisterUser(ctx context.Context, req *api.RegisterRequest) (api.RegisterUserRes, error) {
	var converter generated.ConverterImpl

	userModel, err := converter.RegisterRequestToUserModel(*req)
	if err != nil {
		return api.RegisterRequest{}, err
	}

	a.userRepository.CreateUser(ctx, userModel)

}

func (a *AuthnService) RequestPasswordReset(ctx context.Context, req *api.PasswordResetRequest) (api.RequestPasswordResetRes, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthnService) ResetPassword(ctx context.Context, req *api.PasswordResetConfirm) (api.ResetPasswordRes, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthnService) UpdateUserProfile(ctx context.Context, req *api.ProfileUpdateRequest) (api.UpdateUserProfileRes, error) {
	//TODO implement me
	panic("implement me")
}

func NewAuthnService() IAuthnService {
	return &AuthnService{}
}
