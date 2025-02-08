package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zensey/go-archetype-project/service/eth"
	"github.com/Zensey/go-archetype-project/service/handlers"
)

func main() {
	adr := flag.String("adr", "", "subscribe to a given address (only one value)")
	block := flag.Int("block", 0, "start scan from a given block")
	flag.Parse()

	observer := eth.New()
	if *adr != "" {
		observer.Subscribe(*adr)
	}
	if *block != 0 {
		observer.SetCurrentBlockID(*block)
	}

	observer.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	handlers.SetHttpHandlers(observer)
	go http.ListenAndServe(":8181", nil)

	<-sigs
	observer.Stop()
}
