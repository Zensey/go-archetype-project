package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"strconv"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
	"github.com/zensey/go-archetype-project/protocol"
)

const (
	maxIterations = 1 << 30
)

func main() {
	host := flag.String("host", "localhost", "host of server")
	port := flag.Int("port", 9999, "server port")
	flag.Parse()
	address := *host + ":" + strconv.Itoa(*port)

	// Connect to the server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to server at", address)

	// Read the server's challengeStr
	challengeStr, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	// fmt.Printf("Server challenge: %s", challengeStr)
	challenge, _ := protocol.Unmarshal(challengeStr)

	hasher, err := hash.NewHasher("sha256")
	if err != nil {
		return
	}
	powService := pow.New(hasher)

	resp, err := powService.Compute(context.Background(), challenge, maxIterations)
	if err != nil {
		return
	}
	// fmt.Println("resp:", resp)

	// Send the message to the server
	_, err = conn.Write([]byte(resp.String() + "\n"))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	serverResponse, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	fmt.Println("Quote:", serverResponse)
}
