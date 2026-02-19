package service

import "context"

type EncryptionService interface {
	Encrypt(ctx context.Context, textToEncrypt string) (string, error)
	Decrypt(ctx context.Context, encryptedText string) (string, error)
}
