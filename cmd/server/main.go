package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
	"github.com/zensey/go-archetype-project/protocol"
)

func main() {
	// Define the address and port to listen on
	address := "0.0.0.0:8080"

	// Start listening for incoming connections
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on", address)

	for {
		// Accept a new connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		// Handle the connection in a separate goroutine
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read messages from the client
	reader := bufio.NewReader(conn)
	for {
		// Read the client's message
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed by client:", conn.RemoteAddr())
			return
		}

		// Print the received message
		fmt.Printf("Message from %s: %s", conn.RemoteAddr(), message)

		// Echo the message back to the client
		_, err = conn.Write([]byte("Server received: " + message))
		if err != nil {
			fmt.Println("Error sending response:", err)
			return
		}

		h := protocol.PoW{}

		///
		hasher, err := hash.NewHasher("sha256")
		if err != nil {
			// handle error
			return
		}
		p := pow.New(hasher)
		hash := h.Unmarshal(message)

		err = p.Verify(hash, "127.0.0.1")
		if err != nil {
			// hashcash token failed verification.
			return
		}
		fmt.Println("OK")
		// GetRandomQuote()
	}
}
