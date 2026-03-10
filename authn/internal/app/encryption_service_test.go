package app_test

import (
	"context"
	"testing"

	"github.com/faber-numeris/luciole-auth/authn/internal/app"
	"github.com/faber-numeris/luciole-auth/authn/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestEncryptionService(t *testing.T) {
	mockCfg := mocks.NewMockIServiceConfig(t)
	s := app.NewEncryptionService(mockCfg)
	ctx := context.Background()

	t.Run("Encrypt/Decrypt success", func(t *testing.T) {
		text := "plain-text"
		encrypted, err := s.Encrypt(ctx, text)
		assert.NoError(t, err)
		assert.Equal(t, text, encrypted) // currently returns as is in TODO implementation

		decrypted, err := s.Decrypt(ctx, encrypted)
		assert.NoError(t, err)
		assert.Equal(t, text, decrypted)
	})
}
