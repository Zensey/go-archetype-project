package protocol_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/zensey/go-archetype-project/protocol"
	"github.com/zensey/go-archetype-project/services/pow_service"
	"github.com/zensey/go-archetype-project/services/quotes"
	"github.com/zensey/go-archetype-project/transport/mocks"
	"github.com/zensey/go-archetype-project/utils"
	"go.uber.org/zap/zapcore"
)

// Simulate a sitation with an end of session initiated by client
func TestEndOfSessionCase(t *testing.T) {
	logger := utils.GetLogger(zapcore.DebugLevel)
	defer logger.Sync()

	qoutesCollection := []string{"quote 1", "quote 2"}
	quoteService := quotes.New(qoutesCollection)
	challengeDifficulty := 2
	powService := pow_service.New(challengeDifficulty, "secret")

	// mock network connection
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockIConnWrapper(ctrl)

	mock.EXPECT().
		RemoteAddr().
		Return("127.0.0.1:11111").
		AnyTimes()
	mock.EXPECT().
		ReadMessage().
		Return(protocol.MsgEndSession, nil)

	resourse := "res"
	err := protocol.ServerProtocolHandler(logger, quoteService, powService, mock, resourse)
	if err != protocol.ErrEndSessionOK {
		t.Fail()
	}
}

// Simulate a situation with a wrong response from client
func TestWrongResponse(t *testing.T) {
	logger := utils.GetLogger(zapcore.DebugLevel)
	defer logger.Sync()

	qoutesCollection := []string{"quote 1", "quote 2"}
	quoteService := quotes.New(qoutesCollection)
	powService := pow_service.New(2, "secret")

	// mock network connection
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockIConnWrapper(ctrl)

	mock.EXPECT().
		RemoteAddr().
		Return("127.0.0.1:11111").
		AnyTimes()

	mock.EXPECT().
		ReadMessage().
		Return(protocol.MsgQueryQuote, nil)
	mock.EXPECT().
		WriteMessage(gomock.Eq(protocol.MsgChallenge)).
		Return(nil)
	mock.EXPECT().
		WriteMessage(gomock.Any()).
		Return(nil)
	mock.EXPECT().
		ReadMessage().
		Return(protocol.MsgResponse, nil)
	mock.EXPECT().
		ReadMessage().
		Return("1:xxxxxxxxx:yyyyyyyyyyy:zzzzzzzzzzz", nil)

	resourse := "res"
	err := protocol.ServerProtocolHandler(logger, quoteService, powService, mock, resourse)
	if err != protocol.ErrWrongFormat {
		t.Fail()
	}
}

// Simulate a situation with a correct response from client
func TestCorrectResponse(t *testing.T) {

	logger := utils.GetLogger(zapcore.DebugLevel)
	defer logger.Sync()

	qoutesCollection := []string{"quote1"}
	quoteService := quotes.New(qoutesCollection)
	powService := pow_service.New(2, "secret")

	// mock network connection
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mocks.NewMockIConnWrapper(ctrl)
	var challenge string

	mock.EXPECT().
		RemoteAddr().
		Return("127.0.0.1:11111").
		AnyTimes()
	mock.EXPECT().
		ReadMessage().
		Return(protocol.MsgQueryQuote, nil)
	mock.EXPECT().
		WriteMessage(gomock.Eq(protocol.MsgChallenge)).
		Return(nil)
	mock.EXPECT().
		WriteMessage(gomock.Any()).
		DoAndReturn(func(msg string) error {
			/* capture a challenge */
			challenge = msg
			return nil
		})
	mock.EXPECT().
		ReadMessage().
		Return(protocol.MsgResponse, nil)
	mock.EXPECT().
		ReadMessage().
		DoAndReturn(func() (string, error) {
			h, err := protocol.Unmarshal(challenge)
			if err != nil {
				return "", err
			}
			resp, err := powService.ComputeResponse(context.Background(), h)
			if err != nil {
				return "", err
			}
			return resp, nil
		})
	mock.EXPECT().
		WriteMessage(gomock.Eq(protocol.MsgQueryQuoteResponse)).
		Return(nil)
	mock.EXPECT().
		WriteMessage(gomock.Eq(qoutesCollection[0])).
		Return(nil)

	resourse := "res"
	err := protocol.ServerProtocolHandler(logger, quoteService, powService, mock, resourse)
	if err != nil {
		t.Fail()
	}
}
