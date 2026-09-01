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
		From:   "Alice",
		To:     "Bob",
		Amount: 10,
		Nonce:  0,
	}

	if err := tx.Sign(wallet); err != nil {
		t.Fatalf("Sign() failed: %v", err)
	}

	if tx.Signature == nil {
		t.Fatal("transaction signature is nil")
	}

	if tx.ID == [32]byte{} {
		t.Fatal("transaction ID is empty")
	}

	if !crypto.Verify(
		&wallet.KeyPair.PublicKey,
		tx.signingData(),
		tx.Signature,
	) {
		t.Fatal("transaction signature is invalid")
	}
}
