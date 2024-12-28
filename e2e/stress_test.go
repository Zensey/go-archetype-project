package e2e

import (
	"context"
	"sync"
	"testing"

	"github.com/zensey/go-archetype-project/services/client"
	"github.com/zensey/go-archetype-project/services/quotes"
	"github.com/zensey/go-archetype-project/services/server"
	"github.com/zensey/go-archetype-project/utils"
	"go.uber.org/zap/zapcore"
)

// Stress server by running multiple clients sending multiple requests
func TestMultipleClients(t *testing.T) {
	listenAddress := "localhost:9999"
	quotesCount := 100
	challengeDifficulty := 2

	serverLogger := utils.GetLogger(zapcore.InfoLevel)
	defer serverLogger.Sync()
	clientLogger := utils.GetLogger(zapcore.ErrorLevel)
	defer clientLogger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	quotesFile := "../assets/quotes.yml"
	qoutesCollection, err := utils.LoadFromYaml(quotesFile)
	if err != nil {
		serverLogger.Sugar().Errorln("Error loading quotes:", err)
		return
	}
	quoteService := quotes.New(qoutesCollection)
	srv := server.New(quoteService, serverLogger, listenAddress, challengeDifficulty)
	go srv.Start(ctx)

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		c := client.New(clientLogger, listenAddress, quotesCount)
		wg.Add(1)
		go func() {
			c.Run()
			wg.Done()
		}()
	}
	wg.Wait()
	srv.Shutdown()
}
