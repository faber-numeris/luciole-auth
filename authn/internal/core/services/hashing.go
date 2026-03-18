package services

import (
	"context"

	"crypto/sha256"

	inboundport "github.com/faber-numeris/luciole-auth/authn/internal/ports/inbound"
)

type hashingService struct{}

func NewHashingService() inboundport.HashingService {
	return &hashingService{}
}

func (s *hashingService) HashPassword(ctx context.Context, password []byte) ([]byte, error) {

	h := sha256.New()
	h.Write(password)
	hash := h.Sum(nil)

	return hash, nil
}
