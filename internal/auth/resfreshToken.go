package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	n, err := rand.Read(key)

	if err != nil {
		return "", fmt.Errorf("reading random bytes: %w", err)
	}
	if n != len(key) {
		return "", fmt.Errorf("short read: got %d, want %d", n, len(key))
	}
	return hex.EncodeToString(key), nil
}
