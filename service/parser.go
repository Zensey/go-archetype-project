package service

type Transaction struct {
	Hash  string
	From  string
	To    string
	Value float32
}

type Parser interface {
	// last parsed block
	GetCurrentBlock() int

	// add address to observer
	Subscribe(address string) bool

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
}
