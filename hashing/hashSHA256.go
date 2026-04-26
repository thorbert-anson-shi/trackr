// Package hashing hashes strings into other strings
package hashing

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashSHA256(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
