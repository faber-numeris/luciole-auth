package app

import (
	"context"

	inboundport "github.com/faber-numeris/luciole-auth/authn/internal/app/ports/inbound"
	"golang.org/x/crypto/bcrypt"
)

type hashingService struct{}

func NewHashingService() inboundport.HashingService {
	return &hashingService{}
}

func (s *hashingService) HashPassword(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
