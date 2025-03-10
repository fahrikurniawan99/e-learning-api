package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateActivationToken menghasilkan token acak sepanjang 16 byte yang diubah ke format hex
func GenerateActivationToken() (string, error) {
	bytes := make([]byte, 16) // Panjang 16 byte sesuai dengan randomBytes(16) di JS
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
