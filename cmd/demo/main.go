package main

import (
	"net/http"

	"github.com/Zensey/go-archetype-project/pkg/handler"
	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/Zensey/go-archetype-project/pkg/migrations"
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

func main() {
	l.Infof("Starting up! Version: %s", version)

	db := pg.Connect(&pg.Options{
		Addr:     "db:5432",
		Database: "db",
		User:     "db",
		Password: "xxx",
	})
	migrations.Run(db, l)

	h := handler.NewHandler(l)
	r := chi.NewRouter()
	r.Post("/update-balance", h.Handle)

	http.ListenAndServe(":8080", r)
	return
}
