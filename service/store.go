package service

type TxStore interface {
	AddTx(address string, tx Transaction)
	GetTransactions(address string) []Transaction
}
