package main

import (
	"bufio"
	"context"
	"flag"
	"net"
	"strconv"
	"strings"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
	"github.com/zensey/go-archetype-project/consts"
	"github.com/zensey/go-archetype-project/protocol"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	host := flag.String("host", "localhost", "host of server")
	port := flag.Int("port", 9999, "server port")
	quotesCount := flag.Int("count", 1, "number of quotes to get")
	flag.Parse()

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger := zap.Must(loggerConfig.Build())
	defer logger.Sync()

	hasher, err := hash.NewHasher(consts.HashAlgo)
	if err != nil {
		return
	}
	powService := pow.New(hasher)

	address := *host + ":" + strconv.Itoa(*port)
	logger.Sugar().Debugln("Connect to server ", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		logger.Sugar().Errorln("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for i := 0; i < *quotesCount; i++ {
		logger.Sugar().Debugln("Begin the flow. Wait for challenge")
		challengeStr, err := reader.ReadString('\n')
		if err != nil {
			logger.Sugar().Errorln("Error reading response:", err)
			return
		}
		logger.Sugar().Debugln("Got a challenge:", challengeStr)
		challenge, _ := protocol.Unmarshal(challengeStr)

		response, err := powService.Compute(context.Background(), challenge, consts.PoWMaxIterations)
		if err != nil {
			logger.Sugar().Errorln("Error computing PoW response:", err)
			return
		}
		logger.Sugar().Debugln("Send response to server:", response)
		_, err = conn.Write([]byte(response.String() + "\n"))
		if err != nil {
			logger.Sugar().Errorln("Error sending message:", err)
			return
		}

		logger.Sugar().Debugln("Read a quote from server")
		serverResponse, err := reader.ReadString('\n')
		if err != nil {
			logger.Sugar().Errorln("Error reading response:", err)
			return
		}
		serverResponse = strings.TrimSuffix(serverResponse, "\n")
		logger.Sugar().Infoln("The quote:", serverResponse)
	}
}
