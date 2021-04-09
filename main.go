package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/pion/stun"
)

var stunServer = "stun.l.google.com:19302"
var gameServer = flag.String("game_server", "", "Game server address")

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
	if !hasAlivePeers(a.peers) {
		if err := sendBindingRequest(a.conn, a.stunAddr); err != nil {
			return err
		}
	}

	for i, p := range a.peers {
		//fmt.Println("p>", p.addr, p.isAlive(), p.state)
		if p.addr != nil {

			switch p.state {
			case 1:
				fmt.Println("Send HELLO >", p.addr)
				sendStr(helloMsg, a.conn, p.addr)
				(&a.peers[i]).setState(2)
			case 0, 2:
				fmt.Println("Send PING >", p.addr)
				sendStr(pingMsg, a.conn, p.addr)
			}

			//(&a.peers[i]).fsm(HASADRESS)
		}

		if !p.isAlive() {
			(&a.peers[i]).fsm(NOTACTIVE)
		}
	}
	return nil
}

func (a *app) clientLogicKeepAlive() error {
	if !hasAlivePeers(a.peers) {
		if err := sendBindingRequest(a.conn, a.stunAddr); err != nil {
			return err
		}
	}

	for _, p := range a.peers {
		sendStr(pingMsg, a.conn, p.addr)
	}
	return nil
}

func main() {
	flag.Parse()
	var err error

	a := app{}
	a.isServer = true
	if *gameServer != "" {
		fmt.Println("game_server addr", *gameServer)
		addr, err := net.ResolveUDPAddr(udp, *gameServer)
		if err != nil {
			log.Println("resolve peeraddr:", err)
			return
		}
		a.peers = append(a.peers, peer{addr: addr})
		a.isServer = false
	}

	a.stunAddr, err = net.ResolveUDPAddr(udp, stunServer)
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

	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				return
			}
			//if !stun.IsMessage(msg.data) {
			//	fmt.Println("m>", string(msg.data), msg.addr)
			//}

			switch {
			case string(msg.data) == pingMsg:
				//fmt.Println("m>", string(msg.data))

				pi := findPeerByAddr(msg.addr, a.peers)
				if pi < 0 {
					p := peer{addr: msg.addr}
					a.peers = append(a.peers, p)
					pi = len(a.peers) - 1
				}
				(&a.peers[pi]).fsm(ACTIVE)

			case string(msg.data) == helloMsg:
				fmt.Println("HELLO")

			case stun.IsMessage(msg.data):
				var xorAddr stun.XORMappedAddress
				err := decodeStunMsg(msg.data, &xorAddr)
				if err != nil {
					log.Println("decodeStunMsg:", err)
					continue
				}
				if a.publicAddr.String() != xorAddr.String() {
					a.publicAddr = xorAddr
					if a.isServer {
						log.Printf("Game started. Players should join with './game -game_server %s'\n", xorAddr)
					} else {
						log.Println("Peer public address:", a.publicAddr)
					}
				}

			default:
				log.Panicln("unknown msg", msg)
			}

		case peerStr := <-peerAddrChan:
			addr, err := net.ResolveUDPAddr(udp, peerStr)
			if err != nil {
				log.Panicln("resolve peeraddr:", err)
			}
			pi := findPeerByAddr(addr, a.peers)
			if pi < 0 {
				a.peers = append(a.peers, peer{addr: addr})
			}

		case <-time.Tick(keepaliveMs * time.Millisecond):
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
