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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	secret       = "secret"
	bits         = 4
	challengeTTL = 30 * time.Second
)

func main() {
	host := flag.String("host", "localhost", "host of server")
	port := flag.Int("port", 9999, "server port")
	quotesFile := flag.String("quotes", "quotes.yml", "quotes yml file")
	flag.Parse()

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger := zap.Must(loggerConfig.Build())
	defer logger.Sync()

	// Start listening for incoming connections
	address := *host + ":" + strconv.Itoa(*port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	logger.Info(fmt.Sprintln("Server is listening on", address))

	hasher, err := hash.NewHasher("sha256")
	if err != nil {
		return
	}
	powService := pow.New(hasher, pow.WithChallengeExpDuration(challengeTTL))

	qoutesCollection, err := utils.LoadFromYaml(*quotesFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	quoteService := quotes.New(qoutesCollection)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Sugar().Errorln("Error accepting connection:", err)
			continue
		}
		logger.Sugar().Debugln("Got connection from", conn.RemoteAddr())
		go handleConnection(conn, hasher, powService, quoteService, logger)
	}
}

func handleConnection(conn net.Conn, hasher hash.Hasher, powService *pow.POW, quoteService *quotes.Quotes, logger *zap.Logger) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	resource := "quote"

	// we assume client can do any number of requests
	for {
		logger.Sugar().Debugln("Send challenge to", conn.RemoteAddr())
		challenge, err := pow.InitHashcash(bits, resource, pow.SignExt(secret, hasher))
		if err != nil {
			logger.Sugar().Errorln("Error:", err)
			return
		}
		challengeStr := challenge.String()
		logger.Sugar().Debugln("challenge", challengeStr)
		_, err = conn.Write([]byte(challengeStr + "\n"))
		if err != nil {
			logger.Sugar().Errorln("Error sending challenge:", err)
			return
		}

		// Read the client's response
		responseStr, err := reader.ReadString('\n')
		if err != nil {
			logger.Sugar().Errorln("Connection closed by client:", conn.RemoteAddr())
			return
		}
		logger.Sugar().Debugln("Got response message from ", conn.RemoteAddr())
		response, err := protocol.Unmarshal(responseStr)
		if err != nil {
			logger.Sugar().Errorln(err)
			return
		}
		err = powService.Verify(response, resource)
		if err != nil {
			logger.Sugar().Errorln("PoW verify err:", err)
			return
		}
		_, err = conn.Write([]byte(quoteService.GetRandomQuote() + "\n"))
		if err != nil {
			logger.Sugar().Errorln("Error sending quote:", err)
			return
		}
	}
}
