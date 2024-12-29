package quotes

import (
	"math/rand"
)

type QuoteRepo struct {
	Data []string
}

func New(collection []string) *QuoteRepo {
	return &QuoteRepo{
		Data: collection,
	}
}

func (q *QuoteRepo) GetRandomQuote() string {
	return q.Data[rand.Intn(len(q.Data))]
}
