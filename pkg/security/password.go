package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

// HashPassword generates bcrypt hash of password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// CheckPassword compares password with hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
