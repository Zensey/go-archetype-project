package store

import (
	"sync"

	"github.com/Zensey/go-archetype-project/service"
)

type TxRepo struct {
	mu        sync.Mutex
	addresses []string
	txmap     map[string][]service.Transaction
}

func New() *TxRepo {
	return &TxRepo{
		txmap: make(map[string][]service.Transaction),
	}
}

// add transaction to store
func (r *TxRepo) AddTx(address string, tx service.Transaction) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.txmap[address] = append(r.txmap[address], tx)
}

// list of inbound or outbound transactions for an address
func (r *TxRepo) GetTransactions(address string) []service.Transaction {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := make([]service.Transaction, len(r.txmap[address]))
	copy(res, r.txmap[address])
	return res
}
