package client

import (
	"context"
	"net"

	"github.com/PoW-HC/hashcash/pkg/hash"
	"github.com/PoW-HC/hashcash/pkg/pow"
	"github.com/zensey/go-archetype-project/consts"
	"github.com/zensey/go-archetype-project/protocol"
	"github.com/zensey/go-archetype-project/transport"
	"go.uber.org/zap"
)

type Client struct {
	logger        *zap.Logger
	listenAddress string
	quotesCount   int
	hasher        hash.Hasher
	powService    *pow.POW
}

func New(logger *zap.Logger, listenAddress string, quotesCount int) *Client {
	hasher, err := hash.NewHasher("sha256")
	if err != nil {
		return nil
	}
	powService := pow.New(hasher)

	return &Client{
		logger:        logger,
		listenAddress: listenAddress,
		hasher:        hasher,
		powService:    powService,
		quotesCount:   quotesCount,
	}
}

func (c *Client) Run() {
	c.logger.Sugar().Debugln("Connect to server ", c.listenAddress)
	conn, err := net.Dial("tcp", c.listenAddress)
	if err != nil {
		c.logger.Sugar().Errorln("Error connecting to server:", err)
		return
	}

	cw := transport.New(conn)
	defer cw.Close()

	for i := 0; i < c.quotesCount; i++ {
		err := makeRequest(c.logger, c.powService, cw)
		if err != nil {
			break
		}
	}

	// wait for next challenge and say good bye
	_, err = cw.ReadMessage()
	if err != nil {
		c.logger.Sugar().Errorln("Error reading msg:", err)
		return
	}
	err = cw.WriteMessage("bye")
	if err != nil {
		c.logger.Sugar().Errorln("Error sending message:", err)
	}
}

func makeRequest(logger *zap.Logger, powService *pow.POW, cw *transport.ConnWrapper) error {
	logger.Sugar().Debugln("Begin the flow. Wait for challenge")
	challengeStr, err := cw.ReadMessage()
	if err != nil {
		logger.Sugar().Errorln("Error reading response:", err)
		return err
	}
	logger.Sugar().Debugln("Got a challenge:", challengeStr)
	challenge, err := protocol.Unmarshal(challengeStr)
	if err != nil {
		logger.Sugar().Errorln("Unmarshal error:", err)
		return err
	}

	response, err := powService.Compute(context.Background(), challenge, consts.PoWComputeMaxIterations)
	if err != nil {
		logger.Sugar().Errorln("Error computing PoW response:", err)
		return err
	}
	logger.Sugar().Debugln("Send response to server:", response)
	err = cw.WriteMessage(response.String())
	if err != nil {
		logger.Sugar().Errorln("Error sending response:", err)
		return err
	}

	logger.Sugar().Debugln("Read a quote from server")
	serverResponse, err := cw.ReadMessage()
	if err != nil {
		logger.Sugar().Errorln("Error reading response:", err)
		return err
	}
	logger.Sugar().Infoln("The quote:", serverResponse)
	return nil
}
