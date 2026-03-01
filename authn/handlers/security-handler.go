package handlers

import (
	"context"
	"errors"
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

type userContextKey string

const UserIDKey userContextKey = "user_id"

func (s SecurityService) HandleBearerAuth(
	ctx context.Context,
	operationName api2.OperationName,
	t api2.BearerAuth,
) (context.Context, error) {
	slog.Info("Bearer auth received", "operation", operationName, "token", t.Token)

	if t.Token == "" {
		slog.Warn("Empty token provided")
		return ctx, errors.New("missing bearer token")
	}

	slog.Debug("Token validated", "token", t.Token)
	ctx = context.WithValue(ctx, UserIDKey, t.Token)
	return ctx, nil
}
