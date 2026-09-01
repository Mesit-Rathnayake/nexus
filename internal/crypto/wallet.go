package crypto

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
)

type Wallet struct {
	KeyPair *KeyPair
	Address string
}

func NewWallet() (*Wallet, error) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	address := deriveAddress(&keyPair.PublicKey)

	return &Wallet{
		KeyPair: keyPair,
		Address: address,
	}, nil
}

func deriveAddress(publicKey *ecdsa.PublicKey) string {
	publicKeyBytes := serializePublicKey(publicKey)

	hash := Hash(publicKeyBytes)

	return "nex" + hex.EncodeToString(hash[:20])
}

func serializePublicKey(publicKey *ecdsa.PublicKey) []byte {
	result := make([]byte, 64)

	xBytes := publicKey.X.Bytes()
	yBytes := publicKey.Y.Bytes()

	copy(result[32-len(xBytes):32], xBytes)
	copy(result[64-len(yBytes):], yBytes)

	return result
}

func PublicKeyHex(publicKey *ecdsa.PublicKey) string {
	return hex.EncodeToString(
		serializePublicKey(publicKey),
	)
}

func (w *Wallet) Sign(data []byte) (*Signature, error) {
	if w == nil {
		return nil, fmt.Errorf("wallet is nil")
	}

	if w.KeyPair == nil {
		return nil, fmt.Errorf("wallet has no key pair")
	}

	if w.KeyPair.PrivateKey == nil {
		return nil, fmt.Errorf("wallet has no private key")
	}

	return Sign(w.KeyPair.PrivateKey, data)
}
