package transaction

import (
	"testing"

	"github.com/Mesit-Rathnayake/nexus/internal/crypto"
)

func createPoolTransaction(t *testing.T) *Transaction {
	t.Helper()

	wallet, err := crypto.NewWallet()
	if err != nil {
		t.Fatalf("NewWallet() failed: %v", err)
	}

	tx := &Transaction{
		To:     "nex-recipient",
		Amount: 10,
		Nonce:  0,
	}

	if err := tx.Sign(wallet); err != nil {
		t.Fatalf("Sign() failed: %v", err)
	}

	return tx
}

func TestPoolAdd(t *testing.T) {
	pool := NewPool()

	tx := createPoolTransaction(t)

	if err := pool.Add(tx); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	if pool.Size() != 1 {
		t.Fatalf("pool size = %d, want 1", pool.Size())
	}
}

func TestPoolRejectsInvalidTransaction(t *testing.T) {
	pool := NewPool()

	tx := createPoolTransaction(t)

	tx.Amount = 999999

	if err := pool.Add(tx); err == nil {
		t.Fatal("pool accepted invalid transaction")
	}
}

func TestPoolRejectsDuplicateTransaction(t *testing.T) {
	pool := NewPool()

	tx := createPoolTransaction(t)

	if err := pool.Add(tx); err != nil {
		t.Fatalf("first Add() failed: %v", err)
	}

	if err := pool.Add(tx); err == nil {
		t.Fatal("pool accepted duplicate transaction")
	}
}

func TestPoolRemove(t *testing.T) {
	pool := NewPool()

	tx := createPoolTransaction(t)

	if err := pool.Add(tx); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	pool.Remove(tx.ID)

	if pool.Size() != 0 {
		t.Fatalf("pool size = %d, want 0", pool.Size())
	}
}
