package jwt

import (
	"crypto/rand"
	"encoding/base64"
)

// NewOpaqueToken generates a random opaque token for early-stage scaffold usage.
func NewOpaqueToken(size int) (string, error) {
	if size <= 0 {
		size = 32
	}
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
