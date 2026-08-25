package blockchain

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/Mesit-Rathnayake/nexus/internal/crypto"
)

type BlockHeader struct {
	Index        uint64
	Timestamp    int64
	PreviousHash [32]byte
	MerkleRoot   [32]byte
	Nonce        uint64
}

type Block struct {
	Header BlockHeader
	Hash   [32]byte
}

func (h BlockHeader) Bytes() []byte {
	var buffer bytes.Buffer

	_ = binary.Write(&buffer, binary.BigEndian, h.Index)
	_ = binary.Write(&buffer, binary.BigEndian, h.Timestamp)
	buffer.Write(h.PreviousHash[:])
	buffer.Write(h.MerkleRoot[:])
	_ = binary.Write(&buffer, binary.BigEndian, h.Nonce)

	return buffer.Bytes()
}

func (b *Block) CalculateHash() [32]byte {
	return crypto.Hash(b.Header.Bytes())
}

func NewBlock(index uint64, previousHash [32]byte) *Block {
	header := BlockHeader{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		PreviousHash: previousHash,
		Nonce:        0,
	}

	block := &Block{
		Header: header,
	}

	block.Hash = block.CalculateHash()

	return block
}
