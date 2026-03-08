package app

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type HashingService interface {
	HashPassword(ctx context.Context, password string) (string, error)
}

type hashingService struct{}

func NewHashingService() HashingService {
	return &hashingService{}
}

func (s *hashingService) HashPassword(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
