package main

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
)

const maxIterations = 1 << 30

func main() {
	hasher, err := hash.NewHasher("sha256")
	if err != nil {
		return
	}
	p := pow.New(hasher)

	hashcash, err := pow.InitHashcash(5, "127.0.0.1", pow.SignExt("secret", hasher))
	if err != nil {
		return
	}

	solution, err := p.Compute(context.Background(), hashcash, maxIterations)
	if err != nil {
		return
	}

	message := solution.String() + "\n"

	// Define the server address and port to connect to
	address := "localhost:8080"

	// Connect to the server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to server at", address)

	// Send the message to the server
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	// Read the server's response
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Print the server's response
	fmt.Printf("Server response: %s", response)
}
