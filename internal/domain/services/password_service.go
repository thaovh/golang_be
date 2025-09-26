package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// PasswordService handles password-related operations
type PasswordService struct{}

// NewPasswordService creates a new password service
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// HashPassword hashes a password with a random salt
func (ps *PasswordService) HashPassword(password string) (hash, salt string, err error) {
	// Generate random salt
	saltBytes := make([]byte, 16)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate salt: %w", err)
	}
	salt = hex.EncodeToString(saltBytes)

	// Hash password with salt
	hashBytes := sha256.Sum256([]byte(password + salt))
	hash = hex.EncodeToString(hashBytes[:])

	return hash, salt, nil
}

// VerifyPassword verifies a password against its hash and salt
func (ps *PasswordService) VerifyPassword(password, hash, salt string) bool {
	// Hash the provided password with the stored salt
	hashBytes := sha256.Sum256([]byte(password + salt))
	computedHash := hex.EncodeToString(hashBytes[:])

	// Compare hashes
	return computedHash == hash
}

// GenerateRandomPassword generates a random password
func (ps *PasswordService) GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

	password := make([]byte, length)
	for i := range password {
		randomByte := make([]byte, 1)
		if _, err := rand.Read(randomByte); err != nil {
			return "", fmt.Errorf("failed to generate random password: %w", err)
		}
		password[i] = charset[randomByte[0]%byte(len(charset))]
	}

	return string(password), nil
}
