package app

import (
	"os"
	"os/signal"
	"syscall"

	"fmt"
	"github.com/Zensey/go-archetype-project/pkg/logger"
)

type IConfig interface {
	ReadConfig() error
	GetLogger() *logger.Logger
}

type App struct {
	lg   *logger.Logger
	eh   IAppEventHandler
	conf IConfig
}

type IAppEventHandler interface {
	OnStart(a *App) error
	OnStop(a *App)
}

func NewApp(eh IAppEventHandler, conf IConfig) *App {
	return &App{eh: eh, conf: conf}
}

func (app *App) Run() {
	err := func() (err error) {
		if err = app.conf.ReadConfig(); err != nil {
			return
		}
		app.lg = app.conf.GetLogger()
		err = app.eh.OnStart(app)
		return
	}()
	if err != nil {
		if app.lg != nil {
			app.lg.Error(err)
		} else {
			fmt.Println("App.Run err:", err)
		}
		os.Exit(1)
	}

	waitSigTerm()
	app.lg.Info("stopping ..")
	app.eh.OnStop(app)
}

func waitSigTerm() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
