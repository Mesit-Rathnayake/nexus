package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/Mesit-Rathnayake/nexus/internal/crypto"
)

type Transaction struct {
	ID        [32]byte
	From      string
	To        string
	Amount    uint64
	Nonce     uint64
	PublicKey string
	Signature *crypto.Signature
}

func (tx *Transaction) signingData() []byte {
	var buffer bytes.Buffer

	buffer.WriteString(tx.From)
	buffer.WriteString(tx.To)

	_ = binary.Write(&buffer, binary.BigEndian, tx.Amount)
	_ = binary.Write(&buffer, binary.BigEndian, tx.Nonce)

	return buffer.Bytes()
}

func (tx *Transaction) CalculateID() [32]byte {
	return crypto.Hash(tx.signingData())
}

func (tx *Transaction) Sign(wallet *crypto.Wallet) error {
	if wallet == nil {
		return fmt.Errorf("wallet is nil")
	}

	tx.From = wallet.Address
	tx.PublicKey = crypto.PublicKeyHex(&wallet.KeyPair.PublicKey)

	signature, err := wallet.Sign(tx.signingData())
	if err != nil {
		return err
	}

	tx.Signature = signature
	tx.ID = crypto.Hash(tx.signingData())

	return nil
}

func (tx *Transaction) Verify() bool {
	if tx == nil || tx.Signature == nil {
		return false
	}

	publicKey, err := parsePublicKey(tx.PublicKey)
	if err != nil {
		return false
	}

	derivedAddress := deriveAddressFromPublicKey(publicKey)

	if derivedAddress != tx.From {
		return false
	}

	return crypto.Verify(
		publicKey,
		tx.signingData(),
		tx.Signature,
	)
}

func parsePublicKey(encoded string) (*ecdsa.PublicKey, error) {
	data, err := hex.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	if len(data) != 64 {
		return nil, fmt.Errorf("invalid public key length")
	}

	x := new(big.Int).SetBytes(data[:32])
	y := new(big.Int).SetBytes(data[32:])

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}, nil
}

func deriveAddressFromPublicKey(
	publicKey *ecdsa.PublicKey,
) string {
	data := make([]byte, 0, 64)

	data = append(data, publicKey.X.Bytes()...)
	data = append(data, publicKey.Y.Bytes()...)

	hash := crypto.Hash(data)

	return "nex" + hex.EncodeToString(hash[:20])
}
