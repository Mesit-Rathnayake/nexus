package node

import (
	"github.com/Mesit-Rathnayake/nexus/internal/blockchain"
	"github.com/Mesit-Rathnayake/nexus/internal/transaction"
)

type Node struct {
	ID         string
	Blockchain *blockchain.Blockchain
	Mempool    *transaction.Pool
}

func NewNode(id string) *Node {
	return &Node{
		ID:         id,
		Blockchain: blockchain.NewBlockchain(),
		Mempool:    transaction.NewPool(),
	}
}
