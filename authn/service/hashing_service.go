package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

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
