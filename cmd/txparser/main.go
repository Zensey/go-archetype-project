package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zensey/go-archetype-project/service/eth"
	"github.com/Zensey/go-archetype-project/service/handlers"
)

func main() {
	observer := eth.New()
	observer.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	handlers.SetHttpHandlers(observer)
	go http.ListenAndServe(":8181", nil)

	<-sigs
	observer.Stop()
}
