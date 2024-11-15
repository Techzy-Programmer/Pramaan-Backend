package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(input string) string {
	hash := sha256.Sum256([]byte(input))
	hashHex := hex.EncodeToString(hash[:])

	return hashHex
}
