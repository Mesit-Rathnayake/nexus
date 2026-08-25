package crypto

import "crypto/sha256"

// Hash returns the SHA-256 hash of the provided data.
func Hash(data []byte) [32]byte {
	return sha256.Sum256(data)
}
