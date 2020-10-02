package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Zensey/go-archetype-project/pkg/app"
)

var version string

func main() {
	app, err := app.NewApp()
	if err != nil {
		app.Errorf("Error: %v", err)
		return
	}
	app.Info("Serving..")
	err = app.Start()
	if err != nil {
		app.Errorf("Error: %v", err)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	err = app.Stop()
	if err != nil {
		app.Errorf("Error: %v", err)
		return
	}
}
