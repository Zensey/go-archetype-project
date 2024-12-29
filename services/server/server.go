package server

import (
	"context"
	"errors"
	"net"
	"sync"

	"github.com/zensey/go-archetype-project/protocol"
	"github.com/zensey/go-archetype-project/services/pow_service"
	"github.com/zensey/go-archetype-project/services/quotes"
	"github.com/zensey/go-archetype-project/transport"
	"go.uber.org/zap"
)

type Server struct {
	quoteService  *quotes.QuoteRepo
	logger        *zap.Logger
	listenAddress string

	listener   net.Listener
	wg         sync.WaitGroup
	powService *pow_service.PoW
}

func New(quoteService *quotes.QuoteRepo, logger *zap.Logger, listenAddress string, powService *pow_service.PoW) *Server {
	return &Server{
		quoteService:  quoteService,
		logger:        logger,
		listenAddress: listenAddress,
		powService:    powService,
	}
}

func (s *Server) Shutdown() {
	s.logger.Sugar().Infoln("Shutdown...")
	s.listener.Close()
	s.wg.Wait()
}

func (s *Server) Start(ctx context.Context) error {
	var err error
	lc := net.ListenConfig{}
	s.listener, err = lc.Listen(ctx, "tcp", s.listenAddress)
	if err != nil {
		return err
	}
	s.logger.Sugar().Infoln("Server is listening on", s.listenAddress)
	s.wg.Add(1)
	defer s.wg.Done()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			default:
				conn, err := s.listener.Accept()
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						return
					}
					s.logger.Sugar().Errorln("Error accepting connection:", err)
					continue
				}

				s.logger.Sugar().Debugln("Got connection from", conn.RemoteAddr())
				s.wg.Add(1)
				go func() {
					defer s.wg.Done()

					cw := transport.New(conn)
					s.handleConnection(cw, ctx)
				}()
			}
		}
	}()
	return nil
}

func (s *Server) handleConnection(cw *transport.ConnWrapper, ctx context.Context) {
	defer cw.Close()

	const resource = "quote"

	/* we assume client can do any number of requests */
	for {
		select {
		case <-ctx.Done():
			return

		default:
			err := protocol.ServerProtocolHandler(s.logger, s.quoteService, s.powService, cw, resource)
			if err != nil {
				return
			}
		}
	}
}
