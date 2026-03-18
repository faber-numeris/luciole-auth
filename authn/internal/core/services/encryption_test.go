package services_test

import (
	"context"
	"testing"

	"github.com/faber-numeris/luciole-auth/authn/internal/core/services"
	"github.com/stretchr/testify/assert"
)

func TestEncryptionService(t *testing.T) {
	s := services.NewEncryptionService()
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
