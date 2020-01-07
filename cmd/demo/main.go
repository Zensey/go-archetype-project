package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Zensey/go-archetype-project/pkg/data"
	"github.com/Zensey/go-archetype-project/pkg/handler"
	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/Zensey/go-archetype-project/pkg/migrations"
	"github.com/Zensey/go-archetype-project/pkg/utils"
	"github.com/go-chi/chi"
	"github.com/go-pg/pg/v9"
)

var (
	l       logger.Logger
	version string
)

func init() {
	l, _ = logger.NewLogger(logger.LogLevelInfo, "demo", logger.BackendConsole)
}

const tickerPeriod = 10 * time.Second

func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	l.Infof("Starting up! Version: %s", version)

	db := pg.Connect(&pg.Options{
		Addr:     "db:5432",
		Database: "db",
		User:     "db",
		Password: "xxx",
	})
	for {
		_, err := db.Exec("SELECT 1")
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	migrations.Run(db, l)

	dao := data.NewDAO(l, db)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go utils.RunPeriodically(ctx, wg, tickerPeriod, func() {
		err := dao.CancelLastNOddRecords(10)
		if err != nil {
			fmt.Println("callee >", err)
		}
	})

	h := handler.NewHandler(l, dao)
	r := chi.NewRouter()
	r.Post("/update-balance", h.UpdateBalance)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return
	}
	srv := http.Server{
		Handler: r,
	}
	go srv.Serve(listener)

	utils.WaitSigTerm()
	srv.Shutdown(context.Background())
	cancel()
	wg.Wait()

	return
}
