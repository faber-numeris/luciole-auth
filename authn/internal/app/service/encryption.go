package service

import (
	"context"

	inboundport "github.com/faber-numeris/luciole-auth/authn/internal/app/ports/inbound"
)

type encryptionService struct {
}

func NewEncryptionService() inboundport.EncryptionService {
	return &encryptionService{}
}

func (s *encryptionService) Encrypt(ctx context.Context, textToEncrypt string) (string, error) {
	return textToEncrypt, nil // TODO: implement
}

func (s *encryptionService) Decrypt(ctx context.Context, encryptedText string) (string, error) {
	return encryptedText, nil // TODO: implement
}
