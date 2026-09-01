package transaction

import "fmt"

type Pool struct {
	transactions map[[32]byte]*Transaction
}

func NewPool() *Pool {
	return &Pool{
		transactions: make(map[[32]byte]*Transaction),
	}
}

func (p *Pool) Add(tx *Transaction) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}

	if !tx.Verify() {
		return fmt.Errorf("invalid transaction")
	}

	if _, exists := p.transactions[tx.ID]; exists {
		return fmt.Errorf("transaction already exists")
	}

	p.transactions[tx.ID] = tx

	return nil
}

func (p *Pool) Get(id [32]byte) (*Transaction, bool) {
	tx, exists := p.transactions[id]

	return tx, exists
}

func (p *Pool) Size() int {
	return len(p.transactions)
}

func (p *Pool) All() []*Transaction {
	result := make([]*Transaction, 0, len(p.transactions))

	for _, tx := range p.transactions {
		result = append(result, tx)
	}

	return result
}

func (p *Pool) Remove(id [32]byte) {
	delete(p.transactions, id)
}

func (p *Pool) Clear() {
	p.transactions = make(map[[32]byte]*Transaction)
}
