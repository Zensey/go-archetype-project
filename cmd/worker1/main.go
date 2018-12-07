package main

import (
	"github.com/Zensey/go-archetype-project/pkg/app"
	_ "github.com/lib/pq"
)

var version string

type AppEvHandler struct {
	handler *Handler
}

func (a *AppEvHandler) OnStart(app *app.App, conf app.Config) error {
	if err := app.ConnectDb(); err != nil {
		return err
	}
	if err := app.ConnectMq(); err != nil {
		return err
	}

	a.handler = newHandler(app)
	go a.handler.receive()
	return nil
}

func (a *AppEvHandler) OnStop(app *app.App) {
	a.handler.terminate() // graceful shutdown
	app.Redis.Close()
	app.Db.Close()
}

func main() {
	eh := &AppEvHandler{}
	app := app.NewApp("worker#1", app.IAppEventHandler(eh))
	app.Run()
}
