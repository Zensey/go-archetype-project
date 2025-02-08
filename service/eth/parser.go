package eth

import (
	"log"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Zensey/go-archetype-project/service"
	"github.com/Zensey/go-archetype-project/service/eth/protocol"
)

type empty struct{}

type Observer struct {
	shutdown chan empty

	currentBlockID atomic.Int32
	mu             sync.Mutex
	addresses      []string

	txStore service.TxStore
}

func New(txStore service.TxStore) *Observer {
	p := &Observer{
		shutdown: make(chan empty),
		txStore:  txStore,
	}
	p.SetCurrentBlockID(0x1b4)
	return p
}

func (p *Observer) SetCurrentBlockID(id int) {
	p.currentBlockID.Store(int32(id))
}

func (p *Observer) Start() {
	go p.startPoller()
}

func (p *Observer) Stop() {
	close(p.shutdown)
}

func (p *Observer) startPoller() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var (
		err         error
		lastBlockID int
	)

	for {
		select {
		case <-ticker.C:
		case <-p.shutdown:
			return
		}

		// if no addresses then sleep
		if p.getAddressesLen() == 0 {
			continue
		}

		currentBlockID := int(p.currentBlockID.Load())
		if lastBlockID <= currentBlockID {
			lastBlockID, err = protocol.GetLastFinalizedBlockID()
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("lastBlockID is 0x%x", lastBlockID)
		}

		block, err := protocol.GetBlock(currentBlockID)
		if err != nil {
			log.Println("GetBlock error:", err)
			continue
		}

		p.mu.Lock()
		for _, t := range block.Transactions {
			for _, adr := range p.addresses {
				if adr == strings.ToLower(t.From) || adr == strings.ToLower(t.To) {

					raw := new(big.Float).SetInt(t.Value.AsBigInt())
					normalized := raw.Quo(raw, new(big.Float).SetInt(protocol.Eth1()))
					value, _ := normalized.Float32()

					p.txStore.AddTx(adr, service.Transaction{
						Hash:  t.Hash,
						From:  t.From,
						To:    t.To,
						Value: value,
					})
				}
			}
		}
		p.mu.Unlock()

		p.currentBlockID.Add(1)
	}
}

func (p *Observer) getAddressesLen() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.addresses)
}

// last parsed block
func (p *Observer) GetCurrentBlock() int {
	return int(p.currentBlockID.Load())
}

// add address to observer
func (p *Observer) Subscribe(address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	has := false
	for _, adr := range p.addresses {
		if adr == address {
			has = true
			break
		}
	}
	if !has {
		p.addresses = append(p.addresses, strings.ToLower(address))
	}
	return !has
}

// list of inbound or outbound transactions for an address
func (p *Observer) GetTransactions(address string) []service.Transaction {
	return p.txStore.GetTransactions(address)
}
