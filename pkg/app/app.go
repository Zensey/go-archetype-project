package app

import (
	"os"
	"os/signal"
	"syscall"

	"fmt"
	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

type App struct {
	loggerTag string
	logger.Logger
	conf Config

	Db    *sqlx.DB
	Redis *redis.Client

	eh IAppEventHandler
}

type IAppEventHandler interface {
	OnStart(a *App, conf Config) error
	OnStop(a *App)
}

func NewApp(loggerTag string, eh IAppEventHandler) *App {
	return &App{eh: eh, loggerTag: loggerTag}
}

func waitSigTerm() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func (app *App) Run() {
	conf, err := GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	app.conf = conf
	app.Logger, err = logger.NewLogger(logger.LookupLogLevel(app.conf.LogLevel), app.loggerTag, logger.LookupLogBackend(app.conf.LogBackend))
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := app.eh.OnStart(app, conf); err != nil {
		app.Error("OnStart", err)
		return
	}
	waitSigTerm()
	app.Info("stopping ..")
	app.eh.OnStop(app)
}

func (app *App) ConnectDb() (err error) {
	app.Db, err = sqlx.Connect("postgres", app.conf.PgDsn)
	return
}

func (app *App) ConnectMq() error {
	app.Redis = redis.NewClient(&redis.Options{
		Addr:     app.conf.RedisAddr,
		Password: app.conf.RedisPass,
		DB:       0,
	})
	_, err := app.Redis.Ping().Result()
	return err
}

func (app *App) GetBadWords() []string {
	return app.conf.BadWords
}
