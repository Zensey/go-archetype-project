package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/Zensey/go-archetype-project/pkg/app"
	_ "github.com/lib/pq"
)

var version string

const ctxTimeout = 10 * time.Second

type AppEvHandler struct {
	srv http.Server
	hnd *Handler
}

func (a *AppEvHandler) OnStart(app *app.App, conf app.Config) error {
	if err := app.ConnectDb(); err != nil {
		return err
	}
	if err := app.ConnectMq(); err != nil {
		return err
	}

	// Ensure socket is open
	listener, err := net.Listen("tcp", conf.ApiAddr)
	if err != nil {
		return err
	}
	app.Info("Listening on", listener.Addr(), "See report on http://localhost"+conf.ApiAddr+"/api/report")
	a.srv = http.Server{Handler: NewHandler(app)}
	go a.srv.Serve(listener)
	return nil
}

func (a *AppEvHandler) OnStop(app *app.App) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	a.srv.Shutdown(ctx)
	app.Redis.Close()
	app.Db.Close()
}

func main() {
	eh := &AppEvHandler{}
	app := app.NewApp("api-service", app.IAppEventHandler(eh))
	app.Run()
}
