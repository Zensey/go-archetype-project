package main

import (
	"flag"
	"strconv"

	"github.com/zensey/go-archetype-project/services/client"
	"github.com/zensey/go-archetype-project/services/pow_service"
	"github.com/zensey/go-archetype-project/utils"
	"go.uber.org/zap/zapcore"
)

func main() {
	// flags
	host := flag.String("host", "localhost", "host of server")
	port := flag.Int("port", 9999, "server port")
	quotesCount := flag.Int("count", 1, "number of quotes to get")
	flag.Parse()

	logger := utils.GetLogger(zapcore.InfoLevel)
	defer logger.Sync()

	powService := pow_service.New(0)

	address := *host + ":" + strconv.Itoa(*port)
	c := client.New(logger, address, *quotesCount, powService)
	c.Run()
}
