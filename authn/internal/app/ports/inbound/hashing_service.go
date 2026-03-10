package inboundport

import (
	"context"
)

type HashingService interface {
	HashPassword(ctx context.Context, password string) (string, error)
}
