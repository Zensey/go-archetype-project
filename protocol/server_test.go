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
	logger := utils.GetLogger(zapcore.InfoLevel)
	defer logger.Sync()

	qoutesCollection := []string{"quote 1", "quote 2"}
	quoteService := quotes.New(qoutesCollection)
	challengeDifficulty := 2
	powService := pow_service.New(challengeDifficulty)

	// mock network connection
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockIConnWrapper(ctrl)

	m.EXPECT().
		RemoteAddr().
		Return("127.0.0.1:11111").
		AnyTimes()
	m.EXPECT().
		WriteMessage(gomock.Any()).
		Return(nil)
	m.EXPECT().
		ReadMessage().
		Return(protocol.MsgQueryEndSession, nil)

	resourse := "res"
	err := protocol.ServerProtocolHandler(logger, quoteService, powService, m, resourse)
	if err != protocol.ErrEndSession {
		t.Fail()
	}
}

// Simulate a situation with a wrong response from client
func TestWrongResponse(t *testing.T) {
	logger, _ := utils.SetupLogsCapture()
	defer logger.Sync()

	qoutesCollection := []string{"quote 1", "quote 2"}
	quoteService := quotes.New(qoutesCollection)
	powService := pow_service.New(2)

	// mock network connection
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockIConnWrapper(ctrl)

	m.EXPECT().
		RemoteAddr().
		Return("127.0.0.1:11111").
		AnyTimes()
	m.EXPECT().
		WriteMessage(gomock.Any()).
		Return(nil)
	m.EXPECT().
		ReadMessage().
		Return("1:xxxxxxxxx:yyyyyyyyyyy:zzzzzzzzzzz", nil)

	resourse := "res"
	err := protocol.ServerProtocolHandler(logger, quoteService, powService, m, resourse)
	if err != protocol.ErrWrongFormat {
		t.Fail()
	}
}

// Simulate a situation with a correct response from client
func TestCorrectResponse(t *testing.T) {

	logger := utils.GetLogger(zapcore.InfoLevel)
	defer logger.Sync()

	qoutesCollection := []string{"quote1"}
	quoteService := quotes.New(qoutesCollection)
	powService := pow_service.New(2)

	// mock network connection
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockIConnWrapper(ctrl)
	var challenge string

	m.EXPECT().
		RemoteAddr().
		Return("127.0.0.1:11111").
		AnyTimes()
	m.EXPECT().
		WriteMessage(gomock.Any()).
		DoAndReturn(func(msg string) error {
			challenge = msg // capture a challenge
			return nil
		})
	m.EXPECT().
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
	m.EXPECT().
		WriteMessage(gomock.Eq(qoutesCollection[0])).
		Return(nil)

	resourse := "res"
	err := protocol.ServerProtocolHandler(logger, quoteService, powService, m, resourse)
	if err != nil {
		t.Fail()
	}
}
