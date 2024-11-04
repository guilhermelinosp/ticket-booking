package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"ticket-booking/configs/logs"
)

type Cryptography interface {
	EncryptPassword(password string) (string, error)
	VerifyPassword(password, hashedPassword string) (bool, error)
}

type cryptography struct{}

func NewCryptography() *cryptography {
	return &cryptography{}
}

func (c *cryptography) EncryptPassword(password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		logs.Error("Failed to generate salt: %w", err)
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := generateHash(password, salt)
	return fmt.Sprintf("%s.%s", salt, hash), nil
}

func (c *cryptography) VerifyPassword(password, hashedPassword string) (bool, error) {
	parts := splitHashedPassword(hashedPassword)
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid hashed password format")
	}

	salt, hash := parts[0], parts[1]
	return generateHash(password, salt) == hash, nil
}

func generateSalt() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		logs.Error("Failed to generate salt: %w", err)
		return "", fmt.Errorf("salt generation error: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func generateHash(password, salt string) string {
	hashBytes := sha256.Sum256([]byte(password + salt))
	return base64.StdEncoding.EncodeToString(hashBytes[:])
}

func splitHashedPassword(hashedPassword string) []string {
	return strings.SplitN(hashedPassword, ".", 2)
}
