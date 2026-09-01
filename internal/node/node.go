package node

import (
	"github.com/Mesit-Rathnayake/nexus/internal/blockchain"
	"github.com/Mesit-Rathnayake/nexus/internal/transaction"
)

type Node struct {
	ID         string
	Address    string
	Blockchain *blockchain.Blockchain
	Mempool    *transaction.Pool
}

func NewNode(id string, address string) *Node {
	return &Node{
		ID:         id,
		Address:    address,
		Blockchain: blockchain.NewBlockchain(),
		Mempool:    transaction.NewPool(),
	}
}
