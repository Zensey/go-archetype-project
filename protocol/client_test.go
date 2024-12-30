package protocol_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/zensey/go-archetype-project/protocol"
	"github.com/zensey/go-archetype-project/services/pow_service"
	"github.com/zensey/go-archetype-project/transport/mocks"
	"github.com/zensey/go-archetype-project/utils"
	"go.uber.org/zap/zapcore"
)

// Simulate a sitation with an end of session initiated by client
func TestClientRequestQuoteOnce(t *testing.T) {

	logger := utils.GetLogger(zapcore.DebugLevel)
	defer logger.Sync()

	// mock network connection
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockIConnWrapper(ctrl)

	challengeDifficulty := 3
	powService := pow_service.New(challengeDifficulty, "secret")
	resourse := "res"
	challenge, err := powService.GenerateChallenge(resourse)
	if err != nil {
		t.Fail()
	}

	mock.EXPECT().
		RemoteAddr().
		Return("127.0.0.1:11111").
		AnyTimes()
	mock.EXPECT().
		WriteMessage(gomock.Eq(protocol.MsgQueryQuote)).
		Return(nil)
	mock.EXPECT().
		ReadMessage().
		Return(protocol.MsgChallenge, nil)
	mock.EXPECT().
		ReadMessage().
		Return(challenge, nil)
	mock.EXPECT().
		WriteMessage(gomock.Eq(protocol.MsgResponse)).
		Return(nil)
	mock.EXPECT().
		WriteMessage(gomock.Any()).
		DoAndReturn(func(msg string) error {
			h, err := protocol.Unmarshal(msg)
			if err != nil {
				return err
			}
			err = powService.VerifyResponse(h, resourse)
			if err != nil {
				return err
			}
			return nil
		})
	mock.EXPECT().
		ReadMessage().
		Return(protocol.MsgQueryQuoteResponse, nil)
	mock.EXPECT().
		ReadMessage().
		Return("Word of wisdom quote", nil)
	mock.EXPECT().
		WriteMessage(gomock.Eq(protocol.MsgEndSession)).
		Return(nil)

	err = protocol.ClientRequestQuote(logger, powService, mock)
	if err != nil {
		t.Error(err)
	}
	err = protocol.ClientRequestEndSession(logger, mock)
	if err != nil {
		t.Fail()
	}
}
