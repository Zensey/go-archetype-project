package quotes

import (
	"math/rand"
)

type Quotes struct {
	Data []string
}

func New(collection []string) *Quotes {
	return &Quotes{
		Data: collection,
	}
}

func (q *Quotes) GetRandomQuote() string {
	return q.Data[rand.Intn(len(q.Data))]
}
