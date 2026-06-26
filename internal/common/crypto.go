package common

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateSecureKey(prefix string) string {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return prefix + hex.EncodeToString(b)
}
