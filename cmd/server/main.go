package main

import (
	"context"
	"flag"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/zensey/go-archetype-project/services/pow_service"
	"github.com/zensey/go-archetype-project/services/quotes"
	"github.com/zensey/go-archetype-project/services/server"
	"github.com/zensey/go-archetype-project/utils"
	"go.uber.org/zap/zapcore"
)

func main() {
	host := flag.String("host", "0.0.0.0", "host of server")
	port := flag.Int("port", 9999, "server port")
	quotesFile := flag.String("quotes", "quotes.yml", "quotes yml file")
	challengeDifficulty := flag.Int("dif", 5, "PoW challenge difficulty")

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
		logger.Sugar().Errorln("Wrong PoW challenge difficulty:", *challengeDifficulty)
		return
	}
	powService := pow_service.New(*challengeDifficulty)

	srv := server.New(quoteService, logger, listenAddress, powService)	
	if err := srv.Start(ctx); err != nil {
		logger.Sugar().Errorln("Error starting server:", err)
		return
	}

	<-ctx.Done()
	srv.Shutdown()
	logger.Sugar().Info("Exit")
}
