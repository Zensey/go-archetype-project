package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Zensey/go-archetype-project/pkg/cfg"
	"github.com/Zensey/go-archetype-project/pkg/domain"
	"github.com/Zensey/go-archetype-project/pkg/handler"
	"github.com/Zensey/go-archetype-project/pkg/svc"
	"github.com/Zensey/go-archetype-project/pkg/utils"
	"github.com/Zensey/slog"
	"github.com/go-chi/chi"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/oklog/run"
)

var version string

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*domain.Customer)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	l := slog.ConsoleLogger()
	l.SetLevel(slog.LevelInfo)
	l.Infof("Starting up! Version: %s", version)

	db := pg.Connect(&pg.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Database: os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	defer db.Close()

	// wait while db is starting
	for {
		_, err := db.Exec("SELECT 1")
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	err := createSchema(db)
	if err != nil {
		l.Error("createSchema err:", err)
		return
	}

	cr := &cfg.ErplyCredentials{
		Username:   os.Getenv("ERPLY_USERNAME"),
		Password:   os.Getenv("ERPLY_PASSWORD"),
		ClientCode: os.Getenv("ERPLY_CLIENTCODE"),
	}
	svc := svc.NewCustomerService(db, cr, l)

	h := handler.NewHandler(l, svc, os.Getenv("API_KEY"))
	r := chi.NewRouter()
	r.Post("/save-customer", h.SaveCustomer)
	r.Get("/customers", h.GetCustomers)
	// workaround for swagger - CORS issue
	r.Options("/save-customer", h.SaveCustomer)
	r.Options("/customers", h.GetCustomers)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "."))
	handler.FileServer(r, "/files", filesDir)

	// start http
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		l.Error("net.Listen err:", err)
		return
	}
	srv := http.Server{
		Handler: r,
	}

	// Shutdown gracefully
	{
		var g run.Group
		s := utils.NewSigTermHandler()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		g.Add(func() error { return svc.SyncCustomersPeriodic(ctx) }, func(err error) { cancel(); svc.WaitSyncCustomersFinish() })
		g.Add(func() error { return s.Wait() }, func(err error) { s.Stop() })
		g.Add(func() error { return srv.Serve(listener) }, func(err error) { srv.Shutdown(ctx) })
		err := g.Run()
		if err != nil {
			l.Error("run.Group error:", err)
		}
	}
	l.Error("Exit>")
}
