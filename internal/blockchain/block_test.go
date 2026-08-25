package blockchain

import (
	"testing"
)

func TestNewBlock(t *testing.T) {
	var previousHash [32]byte

	block := NewBlock(1, previousHash)

	if block == nil {
		t.Fatal("NewBlock() returned nil")
	}

	if block.Header.Index != 1 {
		t.Fatalf("Index = %d, want 1", block.Header.Index)
	}

	if block.Header.PreviousHash != previousHash {
		t.Fatal("PreviousHash was not stored correctly")
	}

	if block.Hash == [32]byte{} {
		t.Fatal("Block hash should not be empty")
	}
}

func TestBlockHashChangesWhenHeaderChanges(t *testing.T) {
	var previousHash [32]byte

	block := NewBlock(1, previousHash)

	originalHash := block.Hash

	block.Header.Nonce++

	newHash := block.CalculateHash()

	if originalHash == newHash {
		t.Fatal("hash did not change after modifying block data")
	}
}

func TestBlockHashIsDeterministic(t *testing.T) {
	var previousHash [32]byte

	block1 := NewBlock(1, previousHash)
	block2 := block1

	hash1 := block1.CalculateHash()
	hash2 := block2.CalculateHash()

	if hash1 != hash2 {
		t.Fatal("same block data produced different hashes")
	}
}
