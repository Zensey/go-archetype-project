package protocol

import (
	"errors"

	"github.com/zensey/go-archetype-project/services/pow_service"
	"github.com/zensey/go-archetype-project/services/quotes"
	"github.com/zensey/go-archetype-project/transport"
	"go.uber.org/zap"
)

var (
	ErrEndSessionOK = errors.New("end of session")
	ErrWrongMessage = errors.New("wrong msg")
)

// Returns ErrEndSessionOK if client closes session
func ServerProtocolHandler(logger *zap.Logger, quoteRepo *quotes.QuoteRepo, powService *pow_service.PoW, cw transport.IConnWrapper, resource string) error {
	remote := cw.RemoteAddr()

	msgType, err := cw.ReadMessage()
	if err != nil {
		logger.Sugar().Errorln("ReadMessage error:", err)
		return err
	}
	logger.Sugar().Debugln("Got a message:", msgType)
	if msgType == MsgEndSession {
		return ErrEndSessionOK
	}
	if msgType != MsgQueryQuote {
		return ErrWrongMessage
	}

	logger.Sugar().Debugln("Send challenge to", remote)
	challengeStr, err := powService.GenerateChallenge(resource)
	if err != nil {
		logger.Sugar().Errorln("Error:", err)
		return err
	}
	err = cw.WriteMessage(MsgChallenge)
	if err != nil {
		logger.Sugar().Errorln("Error sending challenge:", err)
		return err
	}
	logger.Sugar().Debugln("challenge", challengeStr)
	err = cw.WriteMessage(challengeStr)
	if err != nil {
		logger.Sugar().Errorln("Error sending challenge payload:", err)
		return err
	}

	msgType, err = cw.ReadMessage()
	if err != nil {
		logger.Sugar().Errorln("ReadMessage error:", err)
		return err
	}
	logger.Sugar().Debugln("Msg type", msgType)
	if msgType != MsgResponse {
		return ErrWrongMessage
	}
	responseStr, err := cw.ReadMessage()
	if err != nil {
		logger.Sugar().Errorln("ReadMessage error:", err)
		return err
	}
	logger.Sugar().Debugln("Got response message from", remote, responseStr)
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

	err = cw.WriteMessage(MsgQueryQuoteResponse)
	if err != nil {
		logger.Sugar().Errorln("Error sending quote:", err)
		return err
	}
	err = cw.WriteMessage(quoteRepo.GetRandomQuote())
	if err != nil {
		logger.Sugar().Errorln("Error sending quote:", err)
		return err
	}
	return nil
}
