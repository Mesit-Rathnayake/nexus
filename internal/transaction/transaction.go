package transaction

import (
	"bytes"
	"encoding/binary"

	"github.com/Mesit-Rathnayake/nexus/internal/crypto"
)

type Transaction struct {
	ID        [32]byte
	From      string
	To        string
	Amount    uint64
	Nonce     uint64
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

func (tx *Transaction) Sign(
	keyPair *crypto.KeyPair,
) error {
	signature, err := crypto.Sign(
		keyPair.PrivateKey,
		tx.signingData(),
	)
	if err != nil {
		return err
	}

	tx.Signature = signature
	tx.ID = crypto.Hash(tx.signingData())

	return nil
}
