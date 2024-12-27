package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/zensey/go-archetype-project/services/quotes"
	"github.com/zensey/go-archetype-project/services/server"
	"github.com/zensey/go-archetype-project/utils"
)

func main() {
	host := flag.String("host", "localhost", "host of server")
	port := flag.Int("port", 9999, "server port")
	quotesFile := flag.String("quotes", "quotes.yml", "quotes yml file")
	flag.Parse()

	logger := utils.GetLogger(true)
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
	srv := server.New(quoteService, logger, listenAddress)
	go srv.Start(ctx)

	<-ctx.Done()
	srv.Shutdown()
	logger.Sugar().Info("Exit")
}
