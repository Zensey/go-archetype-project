package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
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

var version string

func main() {
	l, _ := logger.NewLogger(logger.LogLevelDebug, "demo", logger.BackendConsole)
	l.Infof("Starting up! Version: %s", version)

	domain.SetBalanceUpdateSourcesEnum(strings.Split(os.Getenv("TYPES"), ","))
	n, _ := strconv.Atoi(os.Getenv("N"))
	if n <= 0 {
		n = 5
	}

	db := pg.Connect(&pg.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
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
	jobPeriod := time.Duration(n) * time.Minute
	go utils.RunPeriodically(ctx, wg, jobPeriod, func() {
		userID := 0
		err := dao.CancelLastNOddLedgerRecordsInTx(userID)
		if err != nil {
			l.Error("Callee error>", err)
		}
	})

	h := handler.NewHandler(l, dao)
	r := chi.NewRouter()
	r.Post("/update-balance", h.UpdateBalance)

	// start http
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return
	}
	srv := http.Server{
		Handler: r,
	}
	go srv.Serve(listener)

	// process shutdown
	utils.WaitSigTerm()
	srv.Shutdown(context.Background())
	cancel()
	wg.Wait()

	return
}
