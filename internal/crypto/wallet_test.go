package crypto

import "testing"

func TestNewWallet(t *testing.T) {
	wallet, err := NewWallet()

	if err != nil {
		t.Fatalf("NewWallet() failed: %v", err)
	}

	if wallet == nil {
		t.Fatal("NewWallet() returned nil wallet")
	}

	if wallet.KeyPair == nil {
		t.Fatal("wallet has no key pair")
	}

	if wallet.KeyPair.PrivateKey == nil {
		t.Fatal("wallet has no private key")
	}

	if wallet.Address == "" {
		t.Fatal("wallet address is empty")
	}

	if len(wallet.Address) != 43 {
		t.Fatalf(
			"address length = %d, want 43",
			len(wallet.Address),
		)
	}
}

func TestWalletAddressesAreUnique(t *testing.T) {
	wallet1, err := NewWallet()

	if err != nil {
		t.Fatalf("first wallet creation failed: %v", err)
	}

	wallet2, err := NewWallet()

	if err != nil {
		t.Fatalf("second wallet creation failed: %v", err)
	}

	if wallet1.Address == wallet2.Address {
		t.Fatal("two wallets generated the same address")
	}
}

func TestWalletSign(t *testing.T) {
	wallet, err := NewWallet()

	if err != nil {
		t.Fatalf("NewWallet() failed: %v", err)
	}

	message := []byte("hello nexus")

	signature, err := wallet.Sign(message)

	if err != nil {
		t.Fatalf("wallet signing failed: %v", err)
	}

	if signature == nil {
		t.Fatal("wallet returned nil signature")
	}

	if !Verify(
		&wallet.KeyPair.PublicKey,
		message,
		signature,
	) {
		t.Fatal("wallet signature could not be verified")
	}
}
