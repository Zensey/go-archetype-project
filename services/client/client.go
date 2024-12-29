package client

import (
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

func New(logger *zap.Logger, listenAddress string, quotesCount int, powService *pow_service.PoW) *Client {
	return &Client{
		logger:        logger,
		listenAddress: listenAddress,
		powService:    powService,
		quotesCount:   quotesCount,
	}
}

func (c *Client) Run() error {
	c.logger.Sugar().Debugln("Connect to server", c.listenAddress)
	conn, err := net.Dial("tcp", c.listenAddress)
	if err != nil {
		c.logger.Sugar().Errorln("Error connecting to server:", err)
		return err
	}

	cw := transport.New(conn)
	defer cw.Close()

	for i := 0; i < c.quotesCount; i++ {
		err := protocol.ClientRequestQuote(c.logger, c.powService, cw)
		if err != nil {
			return err
		}
	}
	return protocol.ClientRequestEndSession(c.logger, cw)
}
