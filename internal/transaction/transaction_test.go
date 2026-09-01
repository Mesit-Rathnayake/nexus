package transaction

import (
	"testing"

	"github.com/Mesit-Rathnayake/nexus/internal/crypto"
)

func TestTransactionSigning(t *testing.T) {
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

	if tx.From != wallet.Address {
		t.Fatalf(
			"transaction From = %s, want %s",
			tx.From,
			wallet.Address,
		)
	}

	if tx.PublicKey == "" {
		t.Fatal("transaction public key is empty")
	}

	if tx.Signature == nil {
		t.Fatal("transaction signature is nil")
	}

	if tx.ID == [32]byte{} {
		t.Fatal("transaction ID is empty")
	}

	if !tx.Verify() {
		t.Fatal("valid transaction failed verification")
	}
}

func TestTransactionRejectsTamperedAmount(t *testing.T) {
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

	if !tx.Verify() {
		t.Fatal("transaction should initially be valid")
	}

	tx.Amount = 1000

	if tx.Verify() {
		t.Fatal("tampered transaction was accepted")
	}
}

func TestTransactionRejectsFakeSender(t *testing.T) {
	alice, err := crypto.NewWallet()
	if err != nil {
		t.Fatalf("Alice wallet creation failed: %v", err)
	}

	bob, err := crypto.NewWallet()
	if err != nil {
		t.Fatalf("Bob wallet creation failed: %v", err)
	}

	tx := &Transaction{
		To:     bob.Address,
		Amount: 10,
		Nonce:  0,
	}

	if err := tx.Sign(alice); err != nil {
		t.Fatalf("Sign() failed: %v", err)
	}

	if !tx.Verify() {
		t.Fatal("transaction should initially be valid")
	}

	// Pretend Alice's transaction came from Bob.
	tx.From = bob.Address

	if tx.Verify() {
		t.Fatal("transaction accepted a fake sender")
	}
}

func TestTransactionRejectsTamperedPublicKey(t *testing.T) {
	wallet, err := crypto.NewWallet()
	if err != nil {
		t.Fatalf("NewWallet() failed: %v", err)
	}

	attacker, err := crypto.NewWallet()
	if err != nil {
		t.Fatalf("attacker wallet creation failed: %v", err)
	}

	tx := &Transaction{
		To:     "nex-recipient",
		Amount: 10,
		Nonce:  0,
	}

	if err := tx.Sign(wallet); err != nil {
		t.Fatalf("Sign() failed: %v", err)
	}

	if !tx.Verify() {
		t.Fatal("transaction should initially be valid")
	}

	tx.PublicKey = crypto.PublicKeyHex(
		&attacker.KeyPair.PublicKey,
	)

	if tx.Verify() {
		t.Fatal("transaction accepted a tampered public key")
	}
}
