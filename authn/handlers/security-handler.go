package handlers

import (
	"context"
	"log/slog"

	api2 "github.com/faber-numeris/luciole-auth/authn/api/gen"
)

type ISecurityService = api2.SecurityHandler

type SecurityService struct {
	api2.UnimplementedHandler
}

func NewSecurityService() ISecurityService {
	return &SecurityService{}
}

func (s SecurityService) HandleBearerAuth(
	ctx context.Context,
	operationName api2.OperationName,
	t api2.BearerAuth,
) (context.Context, error) {
	slog.Info("Bearer auth received", "operation", operationName, "token", t.Token)
	// TODO: Validate the token and possibly add user info to the context
	// assignees: rafaelsousa
	return ctx, nil
}
