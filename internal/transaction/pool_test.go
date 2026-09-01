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

func TestNewPool(t *testing.T) {
	pool := NewPool()

	if pool == nil {
		t.Fatal("NewPool() returned nil")
	}

	if pool.Size() != 0 {
		t.Fatalf("new pool size = %d, want 0", pool.Size())
	}
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

func TestPoolGet(t *testing.T) {
	pool := NewPool()
	tx := createPoolTransaction(t)

	if err := pool.Add(tx); err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	got, exists := pool.Get(tx.ID)

	if !exists {
		t.Fatal("transaction was not found")
	}

	if got != tx {
		t.Fatal("Get() returned a different transaction")
	}
}

func TestPoolRejectsInvalidTransaction(t *testing.T) {
	pool := NewPool()
	tx := createPoolTransaction(t)

	// Tamper with the transaction after signing.
	tx.Amount = 999999

	if err := pool.Add(tx); err == nil {
		t.Fatal("pool accepted invalid transaction")
	}

	if pool.Size() != 0 {
		t.Fatal("invalid transaction was added to pool")
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

	if pool.Size() != 1 {
		t.Fatalf("pool size = %d, want 1", pool.Size())
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

	_, exists := pool.Get(tx.ID)

	if exists {
		t.Fatal("removed transaction still exists")
	}
}

func TestPoolClear(t *testing.T) {
	pool := NewPool()

	tx1 := createPoolTransaction(t)
	tx2 := createPoolTransaction(t)

	if err := pool.Add(tx1); err != nil {
		t.Fatalf("adding tx1 failed: %v", err)
	}

	if err := pool.Add(tx2); err != nil {
		t.Fatalf("adding tx2 failed: %v", err)
	}

	if pool.Size() != 2 {
		t.Fatalf("pool size = %d, want 2", pool.Size())
	}

	pool.Clear()

	if pool.Size() != 0 {
		t.Fatalf("pool size after Clear() = %d, want 0", pool.Size())
	}
}

func TestPoolAll(t *testing.T) {
	pool := NewPool()

	tx1 := createPoolTransaction(t)
	tx2 := createPoolTransaction(t)

	if err := pool.Add(tx1); err != nil {
		t.Fatalf("adding tx1 failed: %v", err)
	}

	if err := pool.Add(tx2); err != nil {
		t.Fatalf("adding tx2 failed: %v", err)
	}

	all := pool.All()

	if len(all) != 2 {
		t.Fatalf("All() returned %d transactions, want 2", len(all))
	}
}
