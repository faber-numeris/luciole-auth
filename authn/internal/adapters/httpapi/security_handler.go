package httpapi

import (
	"context"
	"errors"
	"log/slog"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/httpapi/gen"
)

type SecurityHandler struct {
	api.UnimplementedHandler
}

func NewSecurityHandler() *SecurityHandler {
	return &SecurityHandler{}
}

type userContextKey string

const UserIDKey userContextKey = "user_id"

func (s SecurityHandler) HandleBearerAuth(
	ctx context.Context,
	operationName api.OperationName,
	t api.BearerAuth,
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
