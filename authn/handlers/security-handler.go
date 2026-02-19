package handlers

import (
	"context"
	"log/slog"

	api "github.com/faber-numeris/luciole-auth/api/gen"
)

type ISecurityService = api.SecurityHandler

type SecurityService struct {
	api.UnimplementedHandler
}

func NewSecurityService() ISecurityService {
	return &SecurityService{}
}

func (s SecurityService) HandleBearerAuth(
	ctx context.Context,
	operationName api.OperationName,
	t api.BearerAuth,
) (context.Context, error) {
	slog.Info("Bearer auth received", "operation", operationName, "token", t.Token)
	// TODO: Validate the token and possibly add user info to the context
	// assignees: rafaelsousa
	return ctx, nil
}
