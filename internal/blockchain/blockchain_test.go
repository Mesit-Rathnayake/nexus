package blockchain

import (
	"testing"

	"github.com/Mesit-Rathnayake/nexus/internal/crypto"
	"github.com/Mesit-Rathnayake/nexus/internal/transaction"
)

func createTestTransaction(
	t *testing.T,
	from string,
	to string,
	amount uint64,
	nonce uint64,
) *transaction.Transaction {

	t.Helper()

	wallet, err := crypto.NewWallet()
	if err != nil {
		t.Fatalf("NewWallet() failed: %v", err)
	}

	tx := &transaction.Transaction{
		From:   from,
		To:     to,
		Amount: amount,
		Nonce:  nonce,
	}

	if err := tx.Sign(wallet); err != nil {
		t.Fatalf("transaction signing failed: %v", err)
	}

	return tx
}

func TestNewBlockchain(t *testing.T) {
	bc := NewBlockchain()

	if bc == nil {
		t.Fatal("NewBlockchain() returned nil")
	}

	if len(bc.Blocks) != 1 {
		t.Fatalf("block count = %d, want 1", len(bc.Blocks))
	}

	genesis := bc.Blocks[0]

	if genesis.Header.Index != 0 {
		t.Fatalf("genesis index = %d, want 0", genesis.Header.Index)
	}

	if genesis.Header.PreviousHash != [32]byte{} {
		t.Fatal("genesis block should have an empty previous hash")
	}

	if !bc.IsValid() {
		t.Fatal("new blockchain should be valid")
	}
}

func TestAddBlock(t *testing.T) {
	bc := NewBlockchain()

	tx := createTestTransaction(
		t,
		"Alice",
		"Bob",
		10,
		0,
	)

	block := bc.AddBlock([]*transaction.Transaction{tx})

	if len(bc.Blocks) != 2 {
		t.Fatalf("block count = %d, want 2", len(bc.Blocks))
	}

	if block.Header.Index != 1 {
		t.Fatalf("block index = %d, want 1", block.Header.Index)
	}

	if block.Header.PreviousHash != bc.Blocks[0].Hash {
		t.Fatal("new block does not reference previous block")
	}

	if len(block.Transactions) != 1 {
		t.Fatalf(
			"transaction count = %d, want 1",
			len(block.Transactions),
		)
	}

	if !bc.IsValid() {
		t.Fatal("blockchain should be valid after adding a block")
	}
}

func TestBlockchainDetectsTampering(t *testing.T) {
	bc := NewBlockchain()

	tx := createTestTransaction(
		t,
		"Alice",
		"Bob",
		10,
		0,
	)

	bc.AddBlock([]*transaction.Transaction{tx})

	if !bc.IsValid() {
		t.Fatal("blockchain should initially be valid")
	}

	// Tamper with the transaction.
	bc.Blocks[1].Transactions[0].Amount = 1000

	if bc.IsValid() {
		t.Fatal("blockchain accepted a tampered transaction")
	}
}

func TestBlockchainDetectsBrokenLink(t *testing.T) {
	bc := NewBlockchain()

	tx1 := createTestTransaction(
		t,
		"Alice",
		"Bob",
		10,
		0,
	)

	tx2 := createTestTransaction(
		t,
		"Bob",
		"Charlie",
		5,
		0,
	)

	bc.AddBlock([]*transaction.Transaction{tx1})
	bc.AddBlock([]*transaction.Transaction{tx2})

	if !bc.IsValid() {
		t.Fatal("blockchain should initially be valid")
	}

	// Break the chain.
	bc.Blocks[2].Header.PreviousHash = [32]byte{}

	if bc.IsValid() {
		t.Fatal("blockchain accepted a broken previous-hash link")
	}
}

func TestBlockchainDetectsModifiedBlockData(t *testing.T) {
	bc := NewBlockchain()

	tx := createTestTransaction(
		t,
		"Alice",
		"Bob",
		10,
		0,
	)

	bc.AddBlock([]*transaction.Transaction{tx})

	if !bc.IsValid() {
		t.Fatal("blockchain should initially be valid")
	}

	// Modify the block header without recalculating the hash.
	bc.Blocks[1].Header.Nonce++

	if bc.IsValid() {
		t.Fatal("blockchain accepted modified block data")
	}
}
