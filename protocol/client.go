package protocol

import (
	"context"

	"github.com/zensey/go-archetype-project/services/pow_service"
	"github.com/zensey/go-archetype-project/transport"
	"go.uber.org/zap"
)

func ClientRequestQuote(logger *zap.Logger, powService *pow_service.PoW, cw transport.IConnWrapper) error {
	logger.Sugar().Debugln("Begin the flow")

	err := cw.WriteMessage(MsgQueryQuote)
	if err != nil {
		logger.Sugar().Errorln("Error sending response:", err)
		return err
	}

	// wait for response
	msgType, err := cw.ReadMessage()
	if err != nil {
		logger.Sugar().Errorln("Error reading token:", err)
		return err
	}
	logger.Sugar().Debugln("Token:", msgType)
	if msgType == MsgChallenge {
		// Wait for challenge paiload
		challengeStr, err := cw.ReadMessage()
		if err != nil {
			logger.Sugar().Errorln("Error reading challenge:", err)
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
		err = cw.WriteMessage(MsgResponse)
		if err != nil {
			logger.Sugar().Errorln("Error sending response:", err)
			return err
		}
		err = cw.WriteMessage(responseStr)
		if err != nil {
			logger.Sugar().Errorln("Error sending response:", err)
			return err
		}
	}

	// wait for msg type
	msgType, err = cw.ReadMessage()
	if err != nil {
		logger.Sugar().Errorln("Error reading token:", err)
		return err
	}
	if msgType == MsgQueryQuoteResponse {
		quote, err := cw.ReadMessage()
		if err != nil {
			logger.Sugar().Errorln("Error reading response:", err)
			return err
		}
		logger.Sugar().Infoln("The quote:", quote)
	}
	return nil
}

// Send end of session notification
// The last message from client
func ClientRequestEndSession(logger *zap.Logger, cw transport.IConnWrapper) error {
	err := cw.WriteMessage(MsgEndSession)
	if err != nil {
		logger.Sugar().Errorln("Error sending message:", err)
		return err
	}
	return nil
}
