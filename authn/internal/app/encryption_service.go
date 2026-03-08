package app

import (
	"context"

	"github.com/faber-numeris/luciole-auth/authn/internal/infrastructure/config"
)

type EncryptionService interface {
	Encrypt(ctx context.Context, textToEncrypt string) (string, error)
	Decrypt(ctx context.Context, encryptedText string) (string, error)
}

type encryptionService struct {
	cfg config.IServiceConfig
}

func NewEncryptionService(cfg config.IServiceConfig) EncryptionService {
	return &encryptionService{cfg: cfg}
}

func (s *encryptionService) Encrypt(ctx context.Context, textToEncrypt string) (string, error) {
	return textToEncrypt, nil // TODO: implement
}

func (s *encryptionService) Decrypt(ctx context.Context, encryptedText string) (string, error) {
	return encryptedText, nil // TODO: implement
}
