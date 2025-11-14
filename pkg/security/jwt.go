package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateRandomSecret generates a random secret for JWT
func GenerateRandomSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random secret: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// TODO: Implement JWT token generation and validation
// This will be needed for authentication in Week 2
