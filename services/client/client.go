package client

import (
	"context"
	"net"

	"github.com/zensey/go-archetype-project/protocol"
	"github.com/zensey/go-archetype-project/services/pow_service"
	"github.com/zensey/go-archetype-project/transport"

	"go.uber.org/zap"
)

type Client struct {
	logger        *zap.Logger
	listenAddress string
	quotesCount   int
	powService    *pow_service.PoW
}

func New(logger *zap.Logger, listenAddress string, quotesCount int) *Client {
	powService := pow_service.New(0)

	return &Client{
		logger:        logger,
		listenAddress: listenAddress,
		powService:    powService,
		quotesCount:   quotesCount,
	}
}

func (c *Client) Run() error {
	c.logger.Sugar().Debugln("Connect to server ", c.listenAddress)
	conn, err := net.Dial("tcp", c.listenAddress)
	if err != nil {
		c.logger.Sugar().Errorln("Error connecting to server:", err)
		return err
	}

	cw := transport.New(conn)
	defer cw.Close()

	for i := 0; i < c.quotesCount; i++ {
		err := c.requestQuote(c.logger, c.powService, cw)
		if err != nil {
			return err
		}
	}

	// wait for next challenge and say good bye
	_, err = cw.ReadMessage()
	if err != nil {
		c.logger.Sugar().Errorln("Error reading msg:", err)
		return err
	}
	err = cw.WriteMessage("bye")
	if err != nil {
		c.logger.Sugar().Errorln("Error sending message:", err)
		return err
	}
	return nil
}

func (c *Client) requestQuote(logger *zap.Logger, powService *pow_service.PoW, cw *transport.ConnWrapper) error {
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

	responseStr, err := c.powService.ComputeResponse(context.Background(), challenge)
	if err != nil {
		logger.Sugar().Errorln("Error computing PoW response:", err)
		return err
	}
	logger.Sugar().Debugln("Send response to server:", responseStr)
	err = cw.WriteMessage(responseStr)
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
