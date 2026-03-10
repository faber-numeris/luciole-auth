package mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToUUID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		id := "550e8400-e29b-41d4-a716-446655440000"
		u, err := StringToUUID(id)
		assert.NoError(t, err)
		assert.Equal(t, id, u.String())
	})

	t.Run("error", func(t *testing.T) {
		_, err := StringToUUID("invalid")
		assert.Error(t, err)
	})
}
