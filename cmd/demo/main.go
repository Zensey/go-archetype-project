package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
)

var version string

const address = ":8080"

func main() {
	rand.Seed(0)

	app, err := newServer()
	if err != nil {
		app.Errorf("Error: %v", err)
		return
	}
	app.Info("Serving..")
	err = app.start()
	if err != nil {
		app.Errorf("Error: %v", err)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	err = app.stop()
	if err != nil {
		app.Errorf("Error: %v", err)
		return
	}
}
