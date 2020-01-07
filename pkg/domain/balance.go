package domain

import "errors"

type BalanceUpdate struct {
	TransactionID string  `json:"transactionId" pg:"id,pk"`
	Amount        float64 `json:"amount,string" pg:"amount"`
	State         string  `json:"state"`
	IsCanceled    bool    `json:-`
	Source        string  `json:-`
}

func (b *BalanceUpdate) Validate() error {
	switch b.Source {
	case "game":
	case "server":
	case "payment":

	default:
		return errors.New("unknown source")
	}

	switch b.State {
	case "win":
		if b.Amount < 0 {
			b.Amount = -b.Amount
		}
	case "lost":
		if b.Amount > 0 {
			b.Amount = -b.Amount
		}
	default:
		return errors.New("unknown state")
	}

	return nil
}

type BalanceUpdaterDAO interface {
	UpdateBalanceInTx(r BalanceUpdate) error
}
