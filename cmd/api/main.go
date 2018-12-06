package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/Zensey/go-archetype-project/pkg/app"
	"github.com/Zensey/go-archetype-project/pkg/logger"
	_ "github.com/lib/pq"
)

var version string

const ctxTimeout = 10 * time.Second

type AppEvHandler struct {
	logger.Logger
	srv http.Server
	hnd *Handler
}

func (a *AppEvHandler) OnStart(app *app.App, conf app.Config) error {
	err := app.ConnectDb()
	if err != nil {
		return err
	}
	err = app.ConnectMq()
	if err != nil {
		return err
	}

	// Ensure socket is open
	listener, err := net.Listen("tcp", conf.ApiAddr)
	if err != nil {
		return err
	}
	a.Info("Listening on", listener.Addr())
	a.srv = http.Server{Handler: NewHandler(app)}
	go a.srv.Serve(listener)
	return err
}

func (a *AppEvHandler) OnStop(app *app.App) {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	a.srv.Shutdown(ctx)
	app.Redis.Close()
	app.Db.Close()
}

func main() {
	l, _ := logger.NewLogger(logger.LogLevelInfo, "server", logger.BackendConsole)
	aa := &AppEvHandler{Logger: l}
	app := app.NewApp(l, app.IAppEventHandler(aa))
	app.Run()
}
