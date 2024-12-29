package protocol

import (
	"errors"

	"github.com/zensey/go-archetype-project/services/pow_service"
	"github.com/zensey/go-archetype-project/services/quotes"
	"github.com/zensey/go-archetype-project/transport"
	"go.uber.org/zap"
)

var (
	ErrEndSession = errors.New("end of session")
)

func ServerProtocolHandler(logger *zap.Logger, quoteRepo *quotes.QuoteRepo, powService *pow_service.PoW, cw transport.IConnWrapper, resource string) error {
	remote := cw.RemoteAddr()
	logger.Sugar().Debugln("Send challenge to", remote)
	challengeStr, err := powService.GenerateChallenge(resource)
	if err != nil {
		logger.Sugar().Errorln("Error:", err)
		return err
	}
	logger.Sugar().Debugln("challenge", challengeStr)
	err = cw.WriteMessage(challengeStr)
	if err != nil {
		logger.Sugar().Errorln("Error sending challenge:", err)
		return err
	}

	responseStr, err := cw.ReadMessage()
	if err != nil {
		logger.Sugar().Errorln("ReadMessage error:", err)
		return err
	}
	logger.Sugar().Debugln("Got response message from", remote)
	if responseStr == MsgQueryEndSession {
		return ErrEndSession
	}
	response, err := Unmarshal(responseStr)
	if err != nil {
		logger.Sugar().Errorln(err)
		return err
	}

	err = powService.VerifyResponse(response, resource)
	if err != nil {
		logger.Sugar().Errorln("PoW verify err:", err)
		return err
	}
	err = cw.WriteMessage(quoteRepo.GetRandomQuote())
	if err != nil {
		logger.Sugar().Errorln("Error sending quote:", err)
		return err
	}
	return nil
}
