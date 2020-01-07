package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Zensey/go-archetype-project/pkg/data"
	"github.com/Zensey/go-archetype-project/pkg/domain"
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

const tickerPeriod = 10 * time.Minute

func main() {
	l.Infof("Starting up! Version: %s", version)
	domain.SetBalanceUpdateSources(strings.Split(os.Getenv("TYPES"), ","))

	db := pg.Connect(&pg.Options{
		Addr:     "db:5432",
		Database: "db",
		User:     "db",
		Password: "xxx",
	})
	// wait while db is starting
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
		err := dao.CancelLastNOddLedgerRecordsInTx(10)
		if err != nil {
			l.Error("Callee error>", err)
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

	// wait for process termination
	utils.WaitSigTerm()
	srv.Shutdown(context.Background())
	cancel()
	wg.Wait()

	return
}
