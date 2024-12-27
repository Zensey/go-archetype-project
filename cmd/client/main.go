package main

import (
	"flag"
	"strconv"

	"github.com/zensey/go-archetype-project/services/client"
	"github.com/zensey/go-archetype-project/utils"
)

func main() {
	// flags
	host := flag.String("host", "localhost", "host of server")
	port := flag.Int("port", 9999, "server port")
	quotesCount := flag.Int("count", 1, "number of quotes to get")
	flag.Parse()

	logger := utils.GetLogger(false)
	defer logger.Sync()

	address := *host + ":" + strconv.Itoa(*port)
	c := client.New(logger, address, *quotesCount)
	c.Run()
}
