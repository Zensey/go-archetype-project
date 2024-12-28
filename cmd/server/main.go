package main

import (
	"context"
	"flag"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/zensey/go-archetype-project/services/quotes"
	"github.com/zensey/go-archetype-project/services/server"
	"github.com/zensey/go-archetype-project/utils"
	"go.uber.org/zap/zapcore"
)

func main() {
	host := flag.String("host", "localhost", "host of server")
	port := flag.Int("port", 9999, "server port")
	quotesFile := flag.String("quotes", "quotes.yml", "quotes yml file")
	challengeDifficulty := flag.Int("dif", 4, "PoW challenge difficulty, >=1, default value is 4")

	flag.Parse()

	logger := utils.GetLogger(zapcore.InfoLevel)
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Start listening for incoming connections
	listenAddress := *host + ":" + strconv.Itoa(*port)

	qoutesCollection, err := utils.LoadFromYaml(*quotesFile)
	if err != nil {
		logger.Sugar().Errorln("Error loading quotes:", err)
		return
	}
	quoteService := quotes.New(qoutesCollection)

	if *challengeDifficulty < 1 {
		logger.Sugar().Errorln("Wrong difficulty:")
		return
	}
	srv := server.New(quoteService, logger, listenAddress, *challengeDifficulty)
	go srv.Start(ctx)

	<-ctx.Done()
	srv.Shutdown()
	logger.Sugar().Info("Exit")
}
