package main

import (
	"log"
	"net"
	"time"
)

const (
	pingMsg  = "ping"
	helloMsg = "hello"

	aliveTimeoutSecs = 10
)

type msg struct {
	addr *net.UDPAddr
	data []byte
}

type peerState int
type condition int

const (
	DISCONECTED = peerState(0)
	OK          = peerState(1)
	HELLO_DONE  = peerState(2)

	NOTACTIVE = condition(1)
	ACTIVE    = condition(2)
	HASADRESS = condition(3)
)

type peer struct {
	addr       *net.UDPAddr
	lastActive time.Time
	state      peerState
}

func (p *peer) fsm(cond condition) {
	switch cond {

	case ACTIVE:
		if p.state == DISCONECTED {
			log.Println("Peer alive:", p.addr)
		}
		p.refreshLastActive()
		if p.state == DISCONECTED {
			p.setState(OK)
		}

	case NOTACTIVE:
		if p.state > DISCONECTED {
			log.Println("Lost connection:", p.addr)
			p.setState(DISCONECTED)
		}

	default:
		panic("unknown condition")
	}
}

func (p *peer) setState(s peerState) {
	p.state = s
}

func (p *peer) refreshLastActive() {
	p.lastActive = time.Now()
}

func (p *peer) isAlive() bool {
	return p.lastActive.Add(aliveTimeoutSecs * time.Second).After(time.Now())
}

func findPeerByAddr(addr *net.UDPAddr, pp []peer) int {
	for i, v := range pp {
		if v.addr.String() == addr.String() {
			return i
		}
	}
	return -1
}

func hasAlivePeers(pp []peer) bool {
	for _, p := range pp {
		if p.isAlive() {
			return true
		}
	}
	return false
}
