package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
	"github.com/zensey/go-archetype-project/protocol"
	"github.com/zensey/go-archetype-project/quotes"
	"github.com/zensey/go-archetype-project/utils"
)

const (
	secret       = "secret"
	bits         = 5
	challengeTTL = 30 * time.Second
)

func main() {
	host := flag.String("host", "localhost", "host of server")
	port := flag.Int("port", 9999, "server port")
	flag.Parse()
	address := *host + ":" + strconv.Itoa(*port)

	// Start listening for incoming connections
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on", address)

	hasher, err := hash.NewHasher("sha256")
	if err != nil {
		return
	}
	powService := pow.New(hasher, pow.WithChallengeExpDuration(challengeTTL))

	qoutesCollection, err := utils.LoadFromYaml("quotes.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	quoteService := quotes.New(qoutesCollection)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, hasher, powService, quoteService)
	}
}

func handleConnection(conn net.Conn, hasher hash.Hasher, powService *pow.POW, quoteService *quotes.Quotes) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	resource := "quote"
	fmt.Printf("Got connection from %s\n", conn.RemoteAddr())

	// Send challenge
	challenge, err := pow.InitHashcash(bits, resource, pow.SignExt(secret, hasher))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	challengeStr := challenge.String()
	// fmt.Println(challenge)
	_, err = conn.Write([]byte(challengeStr + "\n"))
	if err != nil {
		fmt.Println("Error sending challenge:", err)
		return
	}

	// Read the client's response
	responseStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Connection closed by client:", conn.RemoteAddr())
		return
	}
	// fmt.Printf("Got response message from %s\n", conn.RemoteAddr())
	response, err := protocol.Unmarshal(responseStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = powService.Verify(response, resource)
	if err != nil {
		fmt.Println("PoW verify err:", err)
		return
	}
	_, err = conn.Write([]byte(quoteService.GetRandomQuote() + "\n"))
	if err != nil {
		fmt.Println("Error sending quote:", err)
		return
	}
}
