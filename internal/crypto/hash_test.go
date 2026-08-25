package crypto

import (
	"encoding/hex"
	"testing"
)

func TestHash(t *testing.T) {
	hash := Hash([]byte("hello world"))

	expected := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	if actual := hex.EncodeToString(hash[:]); actual != expected {
		t.Fatalf("Hash() = %s, want %s", actual, expected)
	}
}

func TestHashDeterministic(t *testing.T) {
	hash1 := Hash([]byte("hello world"))
	hash2 := Hash([]byte("hello world"))

	if hash1 != hash2 {
		t.Fatal("Hash() produced different results for identical input")
	}
}

func TestHashDifferentInput(t *testing.T) {
	hash1 := Hash([]byte("hello world"))
	hash2 := Hash([]byte("Hello world"))

	if hash1 == hash2 {
		t.Fatal("Hash() produced the same result for different inputs")
	}
}

func TestHashEmptyInput(t *testing.T) {
	hash := Hash([]byte{})

	if len(hash) != 32 {
		t.Fatalf("SHA-256 hash length = %d, want 32", len(hash))
	}
}
