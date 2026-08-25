package blockchain

import (
	"testing"

	"github.com/Mesit-Rathnayake/nexus/internal/crypto"
)

func TestMerkleRoot(t *testing.T) {
	tx1 := crypto.Hash([]byte("transaction 1"))
	tx2 := crypto.Hash([]byte("transaction 2"))
	tx3 := crypto.Hash([]byte("transaction 3"))
	tx4 := crypto.Hash([]byte("transaction 4"))

	root := MerkleRoot([][32]byte{
		tx1,
		tx2,
		tx3,
		tx4,
	})

	if root == [32]byte{} {
		t.Fatal("Merkle root should not be empty")
	}
}

func TestMerkleRootDeterministic(t *testing.T) {
	hashes := [][32]byte{
		crypto.Hash([]byte("transaction 1")),
		crypto.Hash([]byte("transaction 2")),
		crypto.Hash([]byte("transaction 3")),
		crypto.Hash([]byte("transaction 4")),
	}

	root1 := MerkleRoot(hashes)
	root2 := MerkleRoot(hashes)

	if root1 != root2 {
		t.Fatal("same transactions produced different Merkle roots")
	}
}

func TestMerkleRootChangesWhenTransactionChanges(t *testing.T) {
	hashes := [][32]byte{
		crypto.Hash([]byte("transaction 1")),
		crypto.Hash([]byte("transaction 2")),
		crypto.Hash([]byte("transaction 3")),
		crypto.Hash([]byte("transaction 4")),
	}

	originalRoot := MerkleRoot(hashes)

	hashes[2] = crypto.Hash([]byte("tampered transaction"))

	newRoot := MerkleRoot(hashes)

	if originalRoot == newRoot {
		t.Fatal("Merkle root did not change after transaction modification")
	}
}

func TestMerkleRootSingleTransaction(t *testing.T) {
	hash := crypto.Hash([]byte("transaction"))

	root := MerkleRoot([][32]byte{hash})

	if root != hash {
		t.Fatal("single transaction should produce its own hash as the root")
	}
}

func TestMerkleRootOddNumberOfTransactions(t *testing.T) {
	hashes := [][32]byte{
		crypto.Hash([]byte("transaction 1")),
		crypto.Hash([]byte("transaction 2")),
		crypto.Hash([]byte("transaction 3")),
	}

	root := MerkleRoot(hashes)

	if root == [32]byte{} {
		t.Fatal("Merkle root should not be empty")
	}
}
