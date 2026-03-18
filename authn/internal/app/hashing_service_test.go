package app_test

import (
	"context"
	"testing"

	"github.com/faber-numeris/luciole-auth/authn/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestHashingService(t *testing.T) {
	s := app.NewHashingService()
	ctx := context.Background()

	t.Run("HashPassword success", func(t *testing.T) {
		password := []byte("password123")
		hash, err := s.HashPassword(ctx, password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})
}
