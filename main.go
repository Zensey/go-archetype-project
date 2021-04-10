package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/pion/stun"
)

var gameServer = flag.String("game_server", "", "Game server address")
var stunServer = flag.String("stun", "stun.l.google.com:19302", "Stun server address")

const (
	udp         = "udp4"
	keepaliveMs = 2000
)

type app struct {
	publicAddr stun.XORMappedAddress
	peers      []peer

	isServer bool
	conn     *net.UDPConn
	stunAddr *net.UDPAddr
}

func (a *app) serverLogicKeepAlive() error {
	if !hasAlivePeers(&a.peers) {
		if err := sendBindingRequest(a.conn, a.stunAddr); err != nil {
			return err
		}
	}

	for i, p := range a.peers {
		if p.addr != nil {
			switch p.state {
			case OK:
				log.Println("Send hello:", p.addr)
				sendStr(helloMsg, a.conn, p.addr)
				(&a.peers[i]).setState(HELLO_DONE)

			case DISCONECTED, HELLO_DONE:
				sendStr(pingMsg, a.conn, p.addr)
			}
		}

		if !p.isAlive() {
			(&a.peers[i]).fsm(NOTACTIVE)
		}
	}
	return nil
}

func (a *app) clientLogicKeepAlive() error {
	if !hasAlivePeers(&a.peers) {
		if err := sendBindingRequest(a.conn, a.stunAddr); err != nil {
			return err
		}
	}

	for _, p := range a.peers {
		sendStr(pingMsg, a.conn, p.addr)
	}
	return nil
}

func (a *app) onStunMsg(msg msg) {
	xorAddr := stun.XORMappedAddress{}

	err := decodeStunMsg(msg.data, &xorAddr)
	if err != nil {
		log.Println("decodeStunMsg:", err)
	}
	if a.publicAddr.String() != xorAddr.String() {
		a.publicAddr = xorAddr
		if a.isServer {
			log.Printf("Game started. Players should join with './game -game_server %s'\n", xorAddr)
		} else {
			log.Println("Peer public address:", a.publicAddr)
		}
	}
}

func main() {
	flag.Parse()
	var err error

	a := app{}
	a.isServer = true
	if *gameServer != "" {
		addr, err := net.ResolveUDPAddr(udp, *gameServer)
		if err != nil {
			log.Println("resolve peeraddr:", err)
			return
		}
		a.peers = append(a.peers, peer{addr: addr})
		a.isServer = false
	}

	a.stunAddr, err = net.ResolveUDPAddr(udp, *stunServer)
	if err != nil {
		log.Fatalln("resolve serveraddr:", err)
	}

	a.conn, err = net.ListenUDP(udp, nil)
	if err != nil {
		log.Fatalln("dial:", err)
	}
	defer a.conn.Close()
	log.Printf("Listening on %s\n", a.conn.LocalAddr())

	messageChan := listen(a.conn)
	peerAddrChan := helperGetPeerAddr()
	timer := time.Tick(keepaliveMs * time.Millisecond)

	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				return
			}

			switch {
			case string(msg.data) == pingMsg:
				i := findPeerByAddr(msg.addr, &a.peers)
				if i < 0 {
					p := peer{addr: msg.addr}
					a.peers = append(a.peers, p)
					i = len(a.peers) - 1
				}
				(&a.peers[i]).fsm(ACTIVE)

			case string(msg.data) == helloMsg:
				log.Println("Hello")

			case stun.IsMessage(msg.data):
				a.onStunMsg(msg)

			default:
				log.Panicln("unknown msg", msg)
			}

		case peerStr := <-peerAddrChan:
			addr, err := net.ResolveUDPAddr(udp, peerStr)
			if err != nil {
				log.Panicln("resolve peeraddr:", err)
			}
			i := findPeerByAddr(addr, &a.peers)
			if i < 0 {
				a.peers = append(a.peers, peer{addr: addr})
			}

		case <-timer:
			if a.isServer {
				err = a.serverLogicKeepAlive()
			} else {
				err = a.clientLogicKeepAlive()
			}
			if err != nil {
				log.Panicln("keepalive:", err)
			}
		}
	}
}
