package protocol

import (
	"context"

	"github.com/zensey/go-archetype-project/services/pow_service"
	"github.com/zensey/go-archetype-project/transport"
	"go.uber.org/zap"
)

func ClientRequestQuote(logger *zap.Logger, powService *pow_service.PoW, cw transport.IConnWrapper) error {
	logger.Sugar().Debugln("Begin the flow. Wait for challenge")
	challengeStr, err := cw.ReadMessage()
	if err != nil {
		logger.Sugar().Errorln("Error reading response:", err)
		return err
	}
	logger.Sugar().Debugln("Got a challenge:", challengeStr)
	challenge, err := Unmarshal(challengeStr)
	if err != nil {
		logger.Sugar().Errorln("Unmarshal error:", err)
		return err
	}

	responseStr, err := powService.ComputeResponse(context.Background(), challenge)
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

func ClientRequestEndSession(logger *zap.Logger, cw transport.IConnWrapper) error {
	// wait for a next challenge and discard it
	_, err := cw.ReadMessage()
	if err != nil {
		logger.Sugar().Errorln("Error reading msg:", err)
		return err
	}
	// send end of session notification
	err = cw.WriteMessage(MsgQueryEndSession)
	if err != nil {
		logger.Sugar().Errorln("Error sending message:", err)
		return err
	}
	return nil
}