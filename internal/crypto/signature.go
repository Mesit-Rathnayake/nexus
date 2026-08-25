package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func Sign(privateKey *ecdsa.PrivateKey, data []byte) (*Signature, error) {
	hash := Hash(data)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

func Verify(
	publicKey *ecdsa.PublicKey,
	data []byte,
	signature *Signature,
) bool {
	hash := Hash(data)

	return ecdsa.Verify(
		publicKey,
		hash[:],
		signature.R,
		signature.S,
	)
}
