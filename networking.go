package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/pion/stun"
)

func listen(conn *net.UDPConn) <-chan msg {
	messages := make(chan msg)
	go func() {
		for {
			buf := make([]byte, 1024)

			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				close(messages)
				return
			}
			buf = buf[:n]

			messages <- msg{
				addr: addr,
				data: buf,
			}
		}
	}()
	return messages
}

func sendBindingRequest(conn *net.UDPConn, addr *net.UDPAddr) error {
	m := stun.MustBuild(stun.TransactionID, stun.BindingRequest)

	err := send(m.Raw, conn, addr)
	if err != nil {
		return fmt.Errorf("binding: %w", err)
	}

	return nil
}

func send(msg []byte, conn *net.UDPConn, addr *net.UDPAddr) error {
	_, err := conn.WriteToUDP(msg, addr)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

func sendStr(msg string, conn *net.UDPConn, addr *net.UDPAddr) error {
	return send([]byte(msg), conn, addr)
}

func decodeStunMsg(message []byte, xorAddr *stun.XORMappedAddress) error {
	m := new(stun.Message)
	m.Raw = message
	decErr := m.Decode()
	if decErr != nil {
		return decErr
	}
	if getErr := xorAddr.GetFrom(m); getErr != nil {
		return getErr
	}
	return nil
}

func helperGetPeerAddr() <-chan string {
	result := make(chan string)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			log.Println("Enter remote peer address:")
			peer, _ := reader.ReadString('\n')
			result <- strings.Trim(peer, " \r\n")
		}
	}()

	return result
}
