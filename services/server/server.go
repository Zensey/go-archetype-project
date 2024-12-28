package server

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/zensey/go-archetype-project/protocol"
	"github.com/zensey/go-archetype-project/services/pow_service"
	"github.com/zensey/go-archetype-project/services/quotes"
	"github.com/zensey/go-archetype-project/transport"
	"go.uber.org/zap"
)

const (
	secret              = "secret"
	defaultChallengeTTL = 30 * time.Second
)

type Server struct {
	quoteService  *quotes.Quotes
	logger        *zap.Logger
	listenAddress string

	listener            net.Listener
	wg                  sync.WaitGroup
	challengeDifficulty int
	powService          *pow_service.PoW

	// readTimeout         time.Duration
	// collection        QuotesCollection
	// challengeProvider PoWChallengeProvider

}

func New(quoteService *quotes.Quotes, logger *zap.Logger, listenAddress string, challengeDifficulty int) *Server {
	powService := pow_service.New(challengeDifficulty)

	return &Server{
		quoteService:        quoteService,
		logger:              logger,
		listenAddress:       listenAddress,
		powService:          powService,
		challengeDifficulty: challengeDifficulty,
	}
}

func (s *Server) Shutdown() {
	s.logger.Sugar().Infoln("Shutdown...")
	s.listener.Close()
	s.wg.Wait()
}

func (s *Server) Start(ctx context.Context) {
	var err error
	lc := net.ListenConfig{}
	s.listener, err = lc.Listen(ctx, "tcp", s.listenAddress)

	if err != nil {
		s.logger.Sugar().Errorln("Error starting server:", err)
		return
	}
	s.logger.Sugar().Infoln("Server is listening on", s.listenAddress)

	s.wg.Add(1)
	defer s.wg.Done()
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
			s.logger.Sugar().Debugln("Send challenge to", cw.RemoteAddr())
			challengeStr, err := s.powService.GenerateChallenge(resource)
			if err != nil {
				s.logger.Sugar().Errorln("Error:", err)
				return
			}

			s.logger.Sugar().Debugln("challenge", challengeStr)
			err = cw.WriteMessage(challengeStr)
			if err != nil {
				s.logger.Sugar().Errorln("Error sending challenge:", err)
				return
			}

			responseStr, err := cw.ReadMessage()
			if err != nil {
				s.logger.Sugar().Errorln("Connection closed by client:", cw.RemoteAddr())
				return
			}
			s.logger.Sugar().Debugln("Got response message from", cw.RemoteAddr())
			if responseStr == "bye" {
				return
			}
			response, err := protocol.Unmarshal(responseStr)
			if err != nil {
				s.logger.Sugar().Errorln(err)
				return
			}
			err = s.powService.VerifyResponse(response, resource)
			if err != nil {
				s.logger.Sugar().Errorln("PoW verify err:", err)
				return
			}
			err = cw.WriteMessage(s.quoteService.GetRandomQuote())
			if err != nil {
				s.logger.Sugar().Errorln("Error sending quote:", err)
				return
			}
		}
	}
}
