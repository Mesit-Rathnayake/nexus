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
	if p == nil {
		return fmt.Errorf("pool is nil")
	}

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
	if p == nil {
		return nil, false
	}

	tx, exists := p.transactions[id]

	return tx, exists
}

func (p *Pool) Size() int {
	if p == nil {
		return 0
	}

	return len(p.transactions)
}

func (p *Pool) All() []*Transaction {
	if p == nil {
		return nil
	}

	result := make([]*Transaction, 0, len(p.transactions))

	for _, tx := range p.transactions {
		result = append(result, tx)
	}

	return result
}

func (p *Pool) Remove(id [32]byte) {
	if p == nil {
		return
	}

	delete(p.transactions, id)
}

func (p *Pool) Clear() {
	if p == nil {
		return
	}

	p.transactions = make(map[[32]byte]*Transaction)
}
