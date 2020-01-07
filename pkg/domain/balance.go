package domain

import (
	"sort"

	"github.com/Zensey/go-archetype-project/pkg/utils"
)

var balanceUpdateSources []string

func SetBalanceUpdateSourcesEnum(l []string) {
	sort.Strings(l)
	balanceUpdateSources = l
}

type BalanceUpdate struct {
	TransactionID string  `json:"transactionId" pg:"id,pk"`
	Amount        float64 `json:"amount,string" pg:"amount"`
	State         string  `json:"state"`
	IsCanceled    bool    `json:-`
	Source        string  `json:-`
	UserID        int     `json:-`
}

func (b *BalanceUpdate) Validate() error {
	if !utils.HasString(balanceUpdateSources, b.Source) {
		return NewLogicError("unknown source")
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
		return NewLogicError("unknown state")
	}

	return nil
}

type BalanceUpdaterDAO interface {
	UpdateBalanceInTx(r BalanceUpdate) error
}
