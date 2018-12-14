package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/Zensey/go-archetype-project/pkg/app"
	"github.com/Zensey/go-archetype-project/pkg/config"
	_ "github.com/lib/pq"
)

var version string

const ctxTimeout = 10 * time.Second

type AppEvHandler struct {
	srv http.Server
	hnd *Handler

	cc config.Config
}

func (a *AppEvHandler) OnStart(app *app.App) error {
	if err := a.cc.ConnectDb(); err != nil {
		return err
	}
	if err := a.cc.ConnectMq(); err != nil {
		return err
	}

	// Ensure that socket is open
	listener, err := net.Listen("tcp", a.cc.ApiAddr)
	if err != nil {
		return err
	}
	a.cc.GetLogger().Info("Listening on", listener.Addr(), "See report on http://localhost"+a.cc.ApiAddr+"/api/report")
	a.srv = http.Server{Handler: NewHandler(&a.cc)}
	go a.srv.Serve(listener)
	return nil
}

func (a *AppEvHandler) OnStop(app *app.App) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	a.srv.Shutdown(ctx)
	a.cc.Redis.Close()
	a.cc.Db.Close()
}

func main() {
	eh := new(AppEvHandler)
	eh.cc.LoggerTag = "api-service"

	app := app.NewApp(app.IAppEventHandler(eh), app.IConfig(&eh.cc))
	app.Run()
}
