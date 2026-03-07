package application

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type IHashingService interface {
	HashPassword(ctx context.Context, password string) (string, error)
}

type HashingService struct{}

func NewHashingService() IHashingService {
	return &HashingService{}
}

func (s *HashingService) HashPassword(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
