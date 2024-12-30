package e2e

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/zensey/go-archetype-project/services/client"
	"github.com/zensey/go-archetype-project/services/pow_service"
	"github.com/zensey/go-archetype-project/services/quotes"
	"github.com/zensey/go-archetype-project/services/server"
	"github.com/zensey/go-archetype-project/utils"
	"go.uber.org/zap/zapcore"
)

// Stress server by running multiple clients sending multiple requests
func TestE2EMultipleClients(t *testing.T) {
	listenAddress := "localhost:9999"
	quotesCount := 100
	challengeDifficulty := 2

	logger := utils.GetLogger(zapcore.InfoLevel)
	defer logger.Sync()
	clientLogger := utils.GetLogger(zapcore.ErrorLevel)
	defer clientLogger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quotesFile := "../assets/quotes.yml"
	qoutesCollection, err := utils.LoadFromYaml(quotesFile)
	if err != nil {
		logger.Sugar().Errorln("Error loading quotes:", err)
		t.Fail()
		return
	}
	quoteService := quotes.New(qoutesCollection)
	powService := pow_service.New(challengeDifficulty, "secret")

	srv := server.New(quoteService, logger, listenAddress, powService)
	if err := srv.Start(ctx); err != nil {
		logger.Sugar().Errorln("Error starting server:", err)
		t.Fail()
		return
	}

	var errCount atomic.Uint32
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		c := client.New(clientLogger, listenAddress, quotesCount, powService)
		wg.Add(1)
		go func() {
			if err := c.Run(); err != nil {
				errCount.Add(1)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	srv.Shutdown()

	if errCount.Load() > 0 {
		t.Fail()
	}
}
