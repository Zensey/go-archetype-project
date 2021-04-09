package main

import (
	"fmt"
	"net"
	"time"
)

const (
	pingMsg  = "ping"
	helloMsg = "hello"
)

type msg struct {
	addr *net.UDPAddr
	data []byte
}

type peerState int

const (
	DISCONECTED = peerState(0)
	OK          = peerState(1)
	WELCOMED    = peerState(2)

	NOTACTIVE = 1
	ACTIVE    = 2
	HASADRESS = 3
)

type peer struct {
	addr    *net.UDPAddr
	lastRcv time.Time

	state peerState
}

func (p *peer) fsm(cond int) {
	switch cond {

	case ACTIVE:
		if p.state == DISCONECTED {
			fmt.Println("Peer alive:", p.addr)
		}
		p.refreshLastRcv()
		if p.state == DISCONECTED {
			p.setState(OK)
		}

	case NOTACTIVE:
		if p.state > DISCONECTED {
			fmt.Println("Lost connect>", p.addr)
			p.setState(DISCONECTED)
		}

	default:
		panic("unknown condition")
	}
}

//func (p *peer) updateState_(s int) {
//	switch s {
//	case 1:
//		if p.state == 0 {
//			fmt.Println("Client alive:", p.addr)
//		}
//		p.refreshLastRcv()
//		if p.state == 0 {
//			p.setState(1)
//		}
//
//	case 0:
//		if p.state > 0 {
//			fmt.Println("Lost connect>", p.addr)
//			p.setState(0)
//		}
//	}
//}

func (p *peer) setState(s peerState) {
	p.state = s
}

func (p *peer) refreshLastRcv() {
	p.lastRcv = time.Now()
}

func (p *peer) isAlive() bool {
	return p.addr != nil && p.lastRcv.Add(10*time.Second).After(time.Now())
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
