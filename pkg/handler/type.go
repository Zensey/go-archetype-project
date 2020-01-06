package handler

type BalanceUpdate struct {
	TransactionID string  `json:"transactionId"`
	Amount        float64 `json:"amount,string"`
	State         string  `json:"state"`
}
