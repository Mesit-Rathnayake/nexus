package crypto

import "testing"

func TestGenerateKeyPair(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	if keyPair.PrivateKey == nil {
		t.Fatal("private key is nil")
	}

	if keyPair.PublicKey.X == nil || keyPair.PublicKey.Y == nil {
		t.Fatal("public key is invalid")
	}
}

func TestSignAndVerify(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	message := []byte("Alice sends 10 NEX to Bob")

	signature, err := Sign(keyPair.PrivateKey, message)
	if err != nil {
		t.Fatalf("Sign() failed: %v", err)
	}

	if !Verify(&keyPair.PublicKey, message, signature) {
		t.Fatal("valid signature was rejected")
	}
}

func TestInvalidSignature(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	message := []byte("Alice sends 10 NEX to Bob")

	signature, err := Sign(keyPair.PrivateKey, message)
	if err != nil {
		t.Fatalf("Sign() failed: %v", err)
	}

	tamperedMessage := []byte("Alice sends 1000 NEX to Bob")

	if Verify(&keyPair.PublicKey, tamperedMessage, signature) {
		t.Fatal("tampered message was accepted")
	}
}
