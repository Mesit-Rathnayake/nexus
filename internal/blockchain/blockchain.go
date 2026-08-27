package blockchain

import (
	"github.com/Mesit-Rathnayake/nexus/internal/transaction"
)

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	genesis := NewBlock(
		0,
		[32]byte{},
		nil,
	)

	return &Blockchain{
		Blocks: []*Block{genesis},
	}
}

func (bc *Blockchain) LatestBlock() *Block {
	if len(bc.Blocks) == 0 {
		return nil
	}

	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) AddBlock(
	transactions []*transaction.Transaction,
) *Block {

	previousBlock := bc.LatestBlock()

	block := NewBlock(
		previousBlock.Header.Index+1,
		previousBlock.Hash,
		transactions,
	)

	bc.Blocks = append(bc.Blocks, block)

	return block
}
