package httpapi

import (
	"context"
	"testing"

	"github.com/faber-numeris/luciole-auth/authn/internal/adapters/httpapi/gen"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHandler_HandleBearerAuth(t *testing.T) {
	s := NewSecurityHandler()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		token := "token-123"
		newCtx, err := s.HandleBearerAuth(ctx, "GetProfile", api.BearerAuth{Token: token})
		assert.NoError(t, err)
		assert.Equal(t, token, newCtx.Value(UserIDKey))
	})

	t.Run("missing token", func(t *testing.T) {
		_, err := s.HandleBearerAuth(ctx, "GetProfile", api.BearerAuth{Token: ""})
		assert.Error(t, err)
		assert.Equal(t, "missing bearer token", err.Error())
	})
}
