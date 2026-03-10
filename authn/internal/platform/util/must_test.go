package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMust(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		val := Must(123, nil)
		assert.Equal(t, 123, val)
	})

	t.Run("panic", func(t *testing.T) {
		assert.Panics(t, func() {
			Must(0, errors.New("error"))
		})
	})
}
