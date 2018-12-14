package main

import (
	"github.com/Zensey/go-archetype-project/pkg/app"
	"github.com/Zensey/go-archetype-project/pkg/config"
	_ "github.com/lib/pq"
)

var version string

type AppEvHandler struct {
	handler *Handler
	cc      config.Config
}

func (a *AppEvHandler) OnStart(app *app.App) error {
	if err := a.cc.ConnectDb(); err != nil {
		return err
	}
	if err := a.cc.ConnectMq(); err != nil {
		return err
	}

	a.handler = newHandler(&a.cc)
	go a.handler.receive()
	return nil
}

func (a *AppEvHandler) OnStop(app *app.App) {
	a.handler.terminate() // graceful shutdown
	a.cc.Redis.Close()
	a.cc.Db.Close()
}

func main() {
	eh := &AppEvHandler{}
	eh.cc.LoggerTag = "worker#2"

	app := app.NewApp(app.IAppEventHandler(eh), app.IConfig(&eh.cc))
	app.Run()
}
